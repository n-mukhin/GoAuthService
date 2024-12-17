package service

import (
	"context"

	"github.com/rs/zerolog/log"
)

type EmailSender interface {
	SendWarningEmail(ctx context.Context, to, oldIP, newIP string) error
}

type EmailService struct {
	SenderAddress string
}

func NewEmailService(sender string) *EmailService {
	return &EmailService{SenderAddress: sender}
}

func (e *EmailService) SendWarningEmail(ctx context.Context, to, oldIP, newIP string) error {
	log.Info().
		Str("event", "ip_change_warning").
		Str("sender", e.SenderAddress).
		Str("to", to).
		Str("old_ip", oldIP).
		Str("new_ip", newIP).
		Msg("Sending IP change warning email")

	return nil
}
