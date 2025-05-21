package services

import (
	"encoding/json"
	"fmt"
	"invoiceB2B/internal/config"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

const (
	// Exchange types
	ExchangeDirect = "direct"
	ExchangeTopic  = "topic"
	ExchangeFanout = "fanout"
)

// NotificationService handles publishing events to RabbitMQ
type NotificationService interface {
	PublishUserRegisteredEvent(payload map[string]interface{}) error
	PublishEvent(exchange, routingKey string, payload interface{}) error
	Close()
}

type notificationService struct {
	cfg         *config.Config
	conn        *amqp.Connection
	channel     *amqp.Channel
	isConnected bool
	mu          sync.Mutex // To protect connection and channel state
	done        chan bool  // To signal goroutines to stop
}

// NewNotificationService initializes the notification service and connects to RabbitMQ.
// It also declares a default exchange.
func NewNotificationService(cfg *config.Config) (NotificationService, error) {
	ns := &notificationService{
		cfg:  cfg,
		done: make(chan bool),
	}

	if err := ns.connect(); err != nil {
		// Log error but allow service to be created. It will attempt to reconnect.
		log.Printf("Initial RabbitMQ connection failed: %v. Will retry in background.", err)
		// Start a goroutine to attempt reconnection periodically if initial connection fails
		go ns.handleReconnect()
	}

	return ns, nil // Return service even if initial connection fails
}

func (s *notificationService) connect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isConnected {
		// If already connected (e.g. by another goroutine), ensure channel is open
		if s.channel == nil || s.channel.IsClosed() {
			if err := s.openChannelAndDeclareExchange(); err != nil {
				s.isConnected = false // Mark as not connected if channel opening fails
				return fmt.Errorf("failed to open channel after reconnecting: %w", err)
			}
		}
		return nil // Already connected and channel is fine
	}

	log.Println("Attempting to connect to RabbitMQ...")
	conn, err := amqp.Dial(s.cfg.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to dial RabbitMQ: %w", err)
	}
	s.conn = conn
	log.Println("RabbitMQ connection established.")

	if err := s.openChannelAndDeclareExchange(); err != nil {
		// If opening channel or declaring exchange fails, close connection and mark as not connected
		if s.conn != nil {
			s.conn.Close()
		}
		return err // Error already includes context
	}

	s.isConnected = true
	log.Println("RabbitMQ channel opened and default exchange declared.")

	// Watch for connection close
	go func() {
		select {
		case <-s.done:
			return
		case err := <-s.conn.NotifyClose(make(chan *amqp.Error)):
			s.mu.Lock()
			s.isConnected = false
			s.channel = nil // Channel is implicitly closed with connection
			s.conn = nil
			s.mu.Unlock()
			log.Printf("RabbitMQ connection closed: %v. Attempting to reconnect...", err)
			go s.handleReconnect() // Non-blocking reconnect attempt
		}
	}()

	// Watch for channel close (less common if connection is stable, but good practice)
	go func() {
		// Need to ensure channel is not nil before setting NotifyClose
		s.mu.Lock()
		ch := s.channel
		s.mu.Unlock()

		if ch != nil {
			select {
			case <-s.done:
				return
			case err := <-ch.NotifyClose(make(chan *amqp.Error)):
				s.mu.Lock()
				// Only mark as not connected if the connection itself isn't already known to be dead
				if s.isConnected {
					s.isConnected = false // Or just try to reopen channel if conn is alive
					log.Printf("RabbitMQ channel closed: %v. Will attempt to reopen on next publish or reconnect.", err)
					// Potentially try to reopen just the channel if connection is still good.
					// For simplicity here, full reconnect might be triggered by isConnected = false.
				}
				s.channel = nil
				s.mu.Unlock()
				// If channel closes but connection is fine, next publish attempt should try to reopen channel.
				// Or, trigger a specific channel reopen attempt here.
			}
		}
	}()

	return nil
}

func (s *notificationService) openChannelAndDeclareExchange() error {
	if s.conn == nil || s.conn.IsClosed() {
		return fmt.Errorf("cannot open channel, RabbitMQ connection is not established or closed")
	}

	ch, err := s.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}
	s.channel = ch

	// Declare a default exchange (e.g., topic exchange for events)
	// This should be configurable and idempotent.
	err = s.channel.ExchangeDeclare(
		s.cfg.RabbitMQEventExchangeName, // name
		ExchangeTopic,                   // type (topic is flexible for routing keys)
		true,                            // durable
		false,                           // auto-deleted
		false,                           // internal
		false,                           // no-wait
		nil,                             // arguments
	)
	if err != nil {
		// If channel was opened, try to close it before returning error
		s.channel.Close()
		s.channel = nil
		return fmt.Errorf("failed to declare exchange '%s': %w", s.cfg.RabbitMQEventExchangeName, err)
	}
	log.Printf("RabbitMQ exchange '%s' declared successfully.", s.cfg.RabbitMQEventExchangeName)
	return nil
}

func (s *notificationService) handleReconnect() {
	for {
		select {
		case <-s.done:
			log.Println("Stopping RabbitMQ reconnection attempts.")
			return
		case <-time.After(10 * time.Second): // Retry interval
			log.Println("Attempting RabbitMQ reconnection...")
			if err := s.connect(); err == nil {
				log.Println("RabbitMQ reconnected successfully.")
				return // Exit goroutine on successful reconnect
			}
			log.Printf("RabbitMQ reconnection failed: %v. Retrying in 10s...", err)
		}
	}
}

// PublishEvent publishes a generic event to a specified exchange and routing key.
func (s *notificationService) PublishEvent(exchange, routingKey string, payload interface{}) error {
	s.mu.Lock()
	if !s.isConnected || s.channel == nil || s.channel.IsClosed() {
		// Attempt to reconnect/reopen channel before giving up
		s.mu.Unlock() // Unlock before calling connect to avoid deadlock
		log.Println("PublishEvent: RabbitMQ not connected or channel closed. Attempting to connect/reopen channel.")
		if err := s.connect(); err != nil {
			log.Printf("PublishEvent: Failed to connect/reopen channel for RabbitMQ: %v", err)
			return fmt.Errorf("RabbitMQ not connected and failed to reconnect: %w", err)
		}
		s.mu.Lock() // Re-lock after connect attempt
		// Check again after connect attempt
		if !s.isConnected || s.channel == nil || s.channel.IsClosed() {
			s.mu.Unlock()
			log.Printf("PublishEvent: Still not connected after reconnect attempt for exchange %s, key %s.", exchange, routingKey)
			return fmt.Errorf("RabbitMQ is not connected")
		}
	}
	// Hold the lock until publish is done to prevent channel being closed by another goroutine.
	defer s.mu.Unlock()

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	err = s.channel.Publish(
		exchange,   // exchange name
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // Make messages persistent
			Body:         body,
		})
	if err != nil {
		// If publishing fails, it might be due to a closed channel/connection.
		// Mark as not connected to trigger reconnect on next attempt.
		s.isConnected = false
		s.channel = nil // Channel might be invalid
		// Consider closing connection too, or let NotifyClose handle it.
		log.Printf("Failed to publish event to RabbitMQ (exchange: %s, key: %s): %v. Marking for reconnect.", exchange, routingKey, err)
		return fmt.Errorf("failed to publish event to RabbitMQ: %w", err)
	}
	log.Printf("Event published to RabbitMQ: Exchange=%s, RoutingKey=%s, Payload=%+v", exchange, routingKey, payload)
	return nil
}

// PublishUserRegisteredEvent sends a user registration event.
func (s *notificationService) PublishUserRegisteredEvent(payload map[string]interface{}) error {
	if s == nil { // Gracefully handle if notificationService itself is nil (e.g. failed init in main)
		log.Println("NotificationService is nil, cannot publish UserRegisteredEvent.")
		return nil // Or return an error
	}
	return s.PublishEvent(
		s.cfg.RabbitMQEventExchangeName,
		s.cfg.RabbitMQUserRegisteredRoutingKey,
		payload,
	)
}

// Close gracefully closes the RabbitMQ channel and connection.
func (s *notificationService) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	close(s.done) // Signal reconnect goroutines to stop

	if s.channel != nil && !s.channel.IsClosed() {
		log.Println("Closing RabbitMQ channel...")
		if err := s.channel.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ channel: %v", err)
		} else {
			log.Println("RabbitMQ channel closed.")
		}
	}
	if s.conn != nil && !s.conn.IsClosed() {
		log.Println("Closing RabbitMQ connection...")
		if err := s.conn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
		} else {
			log.Println("RabbitMQ connection closed.")
		}
	}
	s.isConnected = false
}
