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
func NewNotificationService(cfg *config.Config) (NotificationService, error) {
	ns := &notificationService{
		cfg:  cfg,
		done: make(chan bool),
	}

	if err := ns.connect(); err != nil {
		log.Printf("Initial RabbitMQ connection failed: %v. Will retry in background.", err)
		go ns.handleReconnect() // Start background reconnection attempts
	}

	return ns, nil
}

func (s *notificationService) connect() error {
	s.mu.Lock() // Lock for initial state check and modification

	// Case 1: Fully connected and channel is open
	if s.isConnected && s.channel != nil {
		s.mu.Unlock()
		return nil
	}

	// Case 2: Connection is believed to be active (s.isConnected == true), but channel is nil.
	if s.isConnected && s.channel == nil {
		log.Println("connect: Connection believed active, but channel is nil. Attempting to re-open channel.")
		if s.conn == nil {
			log.Println("connect: Inconsistent state - s.isConnected is true but s.conn is nil. Forcing full reconnect.")
			s.isConnected = false // Mark for full reconnect
		} else {
			// Attempt to open a new channel on the existing connection.
			// openChannelAndDeclareExchange assumes s.mu is held by caller.
			err := s.openChannelAndDeclareExchange() // This uses s.conn and sets s.channel
			if err == nil {
				log.Println("connect: Channel re-opened successfully on existing connection.")
				// s.isConnected is still true. s.channel is now set.
				// Setup NotifyClose for the new channel.
				s.setupChannelNotifyCloseHandler(s.channel) // Pass the newly created channel
				s.mu.Unlock()
				return nil
			}
			log.Printf("connect: Failed to re-open channel on existing connection: %v. Proceeding to full reconnect.", err)
			s.isConnected = false // Mark for full reconnect
			if s.conn != nil {    // s.conn should not be nil here, but defensive
				s.conn.Close() // Best effort close
				s.conn = nil
			}
			// s.channel is already nil or failed to open.
		}
	}
	// s.mu is still locked if we are here.
	// This point is reached if:
	// 1. s.isConnected was initially false.
	// 2. s.isConnected was true, s.channel was nil, but s.conn was also nil (inconsistent state).
	// 3. s.isConnected was true, s.channel was nil, s.conn was not nil, but opening channel failed.

	// Case 3: Full reconnect needed (s.isConnected is false at this point)
	// Ensure previous connection resources are cleaned up.
	if s.conn != nil {
		s.conn.Close() // Idempotent
		s.conn = nil
	}
	s.channel = nil // Ensure channel is also nil

	s.mu.Unlock() // Unlock before the blocking Dial call

	log.Println("Attempting to connect to RabbitMQ (Dial)...")
	connAttempt, err := amqp.Dial(s.cfg.RabbitMQURL)
	if err != nil {
		// No lock needed here as we are returning. State remains 'not connected'.
		return fmt.Errorf("failed to dial RabbitMQ: %w", err)
	}

	s.mu.Lock() // Re-acquire lock to update shared state
	s.conn = connAttempt
	log.Println("RabbitMQ connection established.")

	// Now open channel and declare exchange
	if err := s.openChannelAndDeclareExchange(); err != nil {
		if s.conn != nil {
			s.conn.Close()
			s.conn = nil
		}
		s.isConnected = false // Ensure state reflects failure
		s.mu.Unlock()
		return err
	}

	s.isConnected = true
	log.Println("RabbitMQ channel opened and default exchange declared.")

	// Setup NotifyClose handlers for the new connection and the new channel
	s.setupConnectionNotifyCloseHandler(s.conn) // Pass the newly created connection
	s.setupChannelNotifyCloseHandler(s.channel) // Pass the newly created channel

	s.mu.Unlock()
	return nil
}

// openChannelAndDeclareExchange attempts to open a new channel and declare the default exchange.
// It assumes s.mu is held by the caller and s.conn is valid.
func (s *notificationService) openChannelAndDeclareExchange() error {
	if s.conn == nil {
		return fmt.Errorf("cannot open channel, RabbitMQ connection is nil")
	}

	ch, err := s.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}
	s.channel = ch // Set new channel

	// Declare a default exchange
	err = s.channel.ExchangeDeclare(
		s.cfg.RabbitMQEventExchangeName, // name
		ExchangeTopic,                   // type
		true,                            // durable
		false,                           // auto-deleted
		false,                           // internal
		false,                           // no-wait
		nil,                             // arguments
	)
	if err != nil {
		if s.channel != nil {
			s.channel.Close() // Close the newly opened channel if exchange declaration fails
		}
		s.channel = nil // Nullify s.channel as it's unusable
		return fmt.Errorf("failed to declare exchange '%s': %w", s.cfg.RabbitMQEventExchangeName, err)
	}
	log.Printf("RabbitMQ exchange '%s' declared successfully.", s.cfg.RabbitMQEventExchangeName)
	return nil
}

// setupConnectionNotifyCloseHandler sets up a goroutine to listen for connection close events.
// conn parameter is the specific connection to monitor.
func (s *notificationService) setupConnectionNotifyCloseHandler(conn *amqp.Connection) {
	go func() {
		select {
		case <-s.done:
			return
		case closeErr := <-conn.NotifyClose(make(chan *amqp.Error)): // Monitor the passed 'conn'
			s.mu.Lock()
			// Only act if the connection that closed is still the active one.
			if s.conn == conn {
				s.isConnected = false
				s.channel = nil // Channel is implicitly gone
				s.conn = nil    // Clear the dead connection reference
				log.Printf("RabbitMQ connection closed: %v. Attempting to reconnect...", closeErr)
				s.mu.Unlock() // Unlock before calling handleReconnect
				go s.handleReconnect()
			} else {
				log.Printf("Notification for closure of an old/stale RabbitMQ connection: %v.", closeErr)
				s.mu.Unlock()
			}
		}
	}()
}

// setupChannelNotifyCloseHandler sets up a goroutine to listen for channel close events.
// ch parameter is the specific channel to monitor.
func (s *notificationService) setupChannelNotifyCloseHandler(ch *amqp.Channel) {
	go func() {
		select {
		case <-s.done:
			return
		// The 'err' variable here corresponds to the one mentioned for line 180 in the original error report.
		case err := <-ch.NotifyClose(make(chan *amqp.Error)): // Monitor the passed 'ch'
			s.mu.Lock()
			// Only act if the channel that closed is still the active one.
			if s.channel == ch {
				s.channel = nil // Mark active channel as nil
				// This log corresponds to the original line 180 context.
				log.Printf("Currently active RabbitMQ channel closed: %v. It will be reopened on the next operation or reconnect cycle.", err)
				// Do NOT set s.isConnected = false here. If the connection is alive,
				// next PublishEvent -> connect() will see s.channel == nil and try to reopen it.
			} else {
				log.Printf("Notification for closure of an old/stale RabbitMQ channel: %v.", err)
			}
			s.mu.Unlock()
		}
	}()
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
			} else {
				log.Printf("RabbitMQ reconnection failed: %v. Retrying in 10s...", err)
			}
		}
	}
}

// PublishEvent publishes a generic event.
func (s *notificationService) PublishEvent(exchange, routingKey string, payload interface{}) error {
	s.mu.Lock()
	if !s.isConnected || s.channel == nil {
		s.mu.Unlock() // Unlock before calling connect to avoid deadlock
		log.Println("PublishEvent: RabbitMQ not connected or channel nil. Attempting to connect/reopen.")
		if err := s.connect(); err != nil {
			log.Printf("PublishEvent: Failed to connect/reopen for RabbitMQ: %v", err)
			return fmt.Errorf("RabbitMQ not connected and failed to reconnect: %w", err)
		}
		s.mu.Lock() // Re-lock after connect attempt
		// Check again after connect attempt
		if !s.isConnected || s.channel == nil {
			s.mu.Unlock()
			log.Printf("PublishEvent: Still not connected or channel nil after reconnect attempt for exchange %s, key %s.", exchange, routingKey)
			return fmt.Errorf("RabbitMQ is not connected or channel is nil")
		}
	}
	// Hold the lock until publish is done.
	defer s.mu.Unlock()

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	currentChannel := s.channel // Capture channel to use for this publish operation

	err = currentChannel.Publish(
		exchange,   // exchange name
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		})

	if err != nil {
		// If publishing fails, it might be due to a closed channel/connection.
		// The NotifyClose handlers should eventually catch this, but we can also react here.
		log.Printf("Failed to publish event to RabbitMQ (exchange: %s, key: %s): %v. Marking for reconnect check.", exchange, routingKey, err)

		// It's possible the channel or connection died just before or during Publish.
		// Setting s.isConnected = false or s.channel = nil here directly can be racy
		// if NotifyClose handlers are also trying to modify state.
		// Rely on NotifyClose handlers to update s.isConnected and s.channel.
		// The error will be returned to the caller. Next call will trigger connect() if needed.
		// If err is amqp.ErrClosed, that's a strong indicator.
		if err == amqp.ErrClosed {
			// This indicates the channel or connection was closed.
			// Ensure state reflects this possibility for next operation.
			// No need to unlock and relock s.mu here as we are already under lock.
			if s.channel == currentChannel { // If the channel we used is still the main one
				s.channel = nil // It's definitely gone or problematic
			}
			// s.isConnected might also need to be false if it was a connection issue.
			// For simplicity, let NotifyClose manage s.isConnected.
			// The next call to PublishEvent will run connect() if channel is nil.
		}
		return fmt.Errorf("failed to publish event to RabbitMQ: %w", err)
	}

	log.Printf("Event published to RabbitMQ: Exchange=%s, RoutingKey=%s", exchange, routingKey) // Payload removed for brevity in logs
	return nil
}

// PublishUserRegisteredEvent sends a user registration event.
func (s *notificationService) PublishUserRegisteredEvent(payload map[string]interface{}) error {
	if s == nil {
		log.Println("NotificationService is nil, cannot publish UserRegisteredEvent.")
		return fmt.Errorf("NotificationService is nil") // Return an error
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

	log.Println("Closing NotificationService...")
	close(s.done) // Signal reconnect and NotifyClose goroutines to stop

	if s.channel != nil {
		log.Println("Closing RabbitMQ channel...")
		if err := s.channel.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ channel: %v", err)
		} else {
			log.Println("RabbitMQ channel closed.")
		}
		s.channel = nil
	}
	if s.conn != nil {
		log.Println("Closing RabbitMQ connection...")
		if err := s.conn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
		} else {
			log.Println("RabbitMQ connection closed.")
		}
		s.conn = nil
	}
	s.isConnected = false
	log.Println("NotificationService closed.")
}
