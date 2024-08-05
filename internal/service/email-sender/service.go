package emailsender

import "context"

type EmailSender interface {
	SendEmail(ctx context.Context, recipient string, message string) error
}

type Service struct {
	sender EmailSender
}

func New(emailSender EmailSender) *Service {
	return &Service{
		sender: emailSender,
	}
}
