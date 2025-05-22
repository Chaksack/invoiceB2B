package services

import (
	"fmt"
	"invoiceB2B/internal/config"
	"log"
	// "strconv" // No longer needed here

	"gopkg.in/mail.v2"
)

type EmailService interface {
	SendEmail(to, subject, body string) error
}

type emailService struct {
	cfg    *config.Config
	dialer *mail.Dialer
}

func NewEmailService(cfg *config.Config) EmailService {
	port := cfg.SMTPPort
	if port == 0 {
		port = 587
	}

	d := mail.NewDialer(cfg.SMTPHost, port, cfg.SMTPUser, cfg.SMTPPassword)

	return &emailService{
		cfg:    cfg,
		dialer: d,
	}
}

func (s *emailService) SendEmail(to, subject, body string) error {
	if s.cfg.AppEnv != "production" && (s.cfg.SMTPHost == "smtp.example.com" || s.cfg.SMTPHost == "") {
		log.Printf("DEV MODE: Email not sent. To: %s, Subject: %s, Body: %s\n", to, subject, body)
		return nil
	}

	m := mail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPSenderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	log.Printf("Attempting to send email to %s from %s via %s:%d", to, s.cfg.SMTPSenderEmail, s.cfg.SMTPHost, s.cfg.SMTPPort)

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s. Error: %v", to, err)
		return fmt.Errorf("could not send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}
