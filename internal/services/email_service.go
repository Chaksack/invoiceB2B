package services

import (
	"fmt"
	"invoiceB2B/internal/config"
	"log"
	"strings"
	"time"

	"gopkg.in/mail.v2"
)

type EmailService interface {
	SendEmail(to, subject, body string) error
	SendEmailWithAttachment(to, subject, body, attachmentPath, attachmentName string) error
}

type emailService struct {
	cfg    *config.Config
	dialer *mail.Dialer
}

// formatEmailHTML creates an HTML email with a logo header, title, and content
func formatEmailHTML(subject, body string) string {
	// Replace newlines with HTML line breaks
	bodyWithBreaks := strings.ReplaceAll(body, "\n", "<br/>")

	// Create HTML template with logo header, title, and content
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            text-align: center;
            margin-bottom: 20px;
            padding-bottom: 20px;
            border-bottom: 1px solid #eee;
        }
        .logo {
            max-width: 150px;
            height: auto;
        }
        .title {
            font-size: 24px;
            font-weight: bold;
            margin: 20px 0;
            color: #4F46E5;
        }
        .content {
            margin-bottom: 30px;
        }
        .footer {
            font-size: 12px;
            color: #777;
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
        }
    </style>
</head>
<body>
    <div class="header">
        <img class="logo" src="https://tailwindui.com/plus/img/logos/mark.svg?color=indigo&shade=600" alt="Company Logo">
    </div>
    <div class="title">%s</div>
    <div class="content">%s</div>
    <div class="footer">
        &copy; %d Invoice B2B. All rights reserved.
    </div>
</body>
</html>
`
	// Format the HTML template with the subject, body, and current year
	return fmt.Sprintf(htmlTemplate, subject, subject, bodyWithBreaks, time.Now().Year())
}

func NewEmailService(cfg *config.Config) EmailService {
	port := cfg.SMTPPort
	if port == 0 {
		port = 465
	}

	d := mail.NewDialer(cfg.SMTPHost, port, cfg.SMTPUser, cfg.SMTPPassword)

	return &emailService{
		cfg:    cfg,
		dialer: d,
	}
}

func (s *emailService) SendEmail(to, subject, body string) error {
	if s.cfg.AppEnv != "production" && (s.cfg.SMTPHost == "smtp.gmail.com" || s.cfg.SMTPHost == "") {
		log.Printf("DEV MODE: Email not sent. To: %s, Subject: %s, Body: %s\n", to, subject, body)
		return nil
	}

	// Format the email body with HTML template
	htmlBody := formatEmailHTML(subject, body)

	m := mail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPSenderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	log.Printf("Attempting to send email to %s from %s via %s:%d", to, s.cfg.SMTPSenderEmail, s.cfg.SMTPHost, s.cfg.SMTPPort)

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s. Error: %v", to, err)
		return fmt.Errorf("could not send email: %w", err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

func (s *emailService) SendEmailWithAttachment(to, subject, body, attachmentPath, attachmentName string) error {
	if s.cfg.AppEnv != "production" && (s.cfg.SMTPHost == "smtp.gmail.com") {
		log.Printf("DEV MODE: Email with attachment not sent. To: %s, Subject: %s, Attachment: %s, Body: %s\n", to, subject, attachmentName, body)
		return nil
	}

	// Format the email body with HTML template
	htmlBody := formatEmailHTML(subject, body)

	m := mail.NewMessage()
	m.SetHeader("From", s.cfg.SMTPSenderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)
	if attachmentPath != "" {
		m.Attach(attachmentPath, mail.Rename(attachmentName))
	}

	log.Printf("Attempting to send email with attachment '%s' to %s from %s via %s:%d", attachmentName, to, s.cfg.SMTPSenderEmail, s.cfg.SMTPHost, s.cfg.SMTPPort)

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send email with attachment to %s. Error: %v", to, err)
		return fmt.Errorf("could not send email with attachment: %w", err)
	}

	log.Printf("Email with attachment sent successfully to %s", to)
	return nil
}
