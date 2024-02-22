package mail

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	sendgridMailHelper "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/suraboy/upload-file-worksho/app/core/config"
	"log"
)

type SendgridMailRepository interface {
	SendEmail(emailFrom string, emailTo string) error
}

// sendgridMail implements the SendgridMailRepository interface.
type sendgridMail struct {
	AppConfig  *config.AppConfig
	MailClient *sendgrid.Client
}

type Config struct {
	AppConfig  *config.AppConfig
	MailClient *sendgrid.Client
}

// NewSendGridClient creates a new SendGrid email client.
func NewSendGridClient(cfg Config) SendgridMailRepository {
	return &sendgridMail{
		AppConfig:  cfg.AppConfig,
		MailClient: cfg.MailClient,
	}
}

// SendEmail sends an email using SendGrid.
func (s *sendgridMail) SendEmail(emailFrom string, emailTo string) error {
	// Compose the email
	from := sendgridMailHelper.NewEmail("Sirichai Janpan", emailFrom)
	to := sendgridMailHelper.NewEmail("Recipient Name", emailTo)
	subject := "Test Email"
	content := sendgridMailHelper.NewContent("text/plain", "Upload File Successfully.")
	message := sendgridMailHelper.NewV3MailInit(from, subject, to, content)

	// Send the email
	response, err := s.MailClient.Send(message)
	if err != nil {
		log.Printf("Failed to send email: %v\n", err)
		return err
	}

	fmt.Printf("Email sent successfully : %v", response)

	return nil
}
