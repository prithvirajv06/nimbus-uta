package service

import (
	"encoding/json"
	"fmt"
	"net/smtp"

	"github.com/prithvirajv06/nimbus-uta/go/notification/config"
	"github.com/prithvirajv06/nimbus-uta/go/notification/internal"
	"github.com/prithvirajv06/nimbus-uta/go/notification/internal/models"
)

type CustomerService struct {
	Cfg *config.Config
}

/*
*

	Will be used during Message Q

*
*/
func (s *CustomerService) HandleUserEvent(message []byte) error {
	var event models.UserMessageEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return err
	}

	switch event.EventType {
	case "USER_CREATED":
		sendNotification(s.Cfg, []string{event.Payload.Email},
			//Subject
			"🚀 Welcome to Nimbus UTA! Your Adventure Begins 🌟",
			internal.GetWelcomeTemplate(event.Payload.FullName(), s.Cfg.FrontEndURL+"/dashboard"),
		)
		// Additional processing logic here
	case "USER_UPDATED":
		fmt.Printf("Handling user updated event for user: %s\n", event.Payload)
	case "USER_DELETED":
		fmt.Printf("Handling user deleted event for user: %s\n", event.Payload)
	}

	return nil
}

func sendNotification(cfg *config.Config, to []string, subject string, body string) error {
	// Configuration
	from := cfg.SMTP.User
	password := cfg.SMTP.Password // Use an App Password, not your real password
	smtpHost := cfg.SMTP.Host
	smtpPort := cfg.SMTP.Port

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email
	// Set MIME headers for HTML email
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	fullMessage := []byte("Subject: " + subject + "\n" + headers + "\n" + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, fullMessage)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	fmt.Println("Email sent successfully!")
	return nil
}
