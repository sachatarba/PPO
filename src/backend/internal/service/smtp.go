package service

import (
	"log"
	"net/smtp"

	"github.com/sachatarba/course-db/internal/config"
)

type SmtpService struct {
	config *config.SmtpConfig
}

func NewSmtpService(config *config.SmtpConfig) *SmtpService {
	return &SmtpService{
		config: config,
	}
}

func (s *SmtpService) SendMail(text string, receiver string, subject string) error {
	message := []byte("From: " + s.config.FromAddres + "\r\n" +
		"To: " + receiver + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		text + "\r\n")

	auth := smtp.PlainAuth("", s.config.FromAddres, s.config.Password, s.config.SmtpHost)
	err := smtp.SendMail(s.config.SmtpHost+":"+s.config.SmtpPort, auth, s.config.FromAddres, []string{receiver}, message)
	if err != nil {
		log.Printf("Failed to send email: %v\n", err)
		return err
	}
	log.Println("Email sent successfully!")
	return nil
}
