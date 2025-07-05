package sms

import (
	"fmt"

	"github.com/icoder-new/installment-cli/internal/domain"
)

type ConsoleSender struct{}

func NewConsoleSender() *ConsoleSender {
	return &ConsoleSender{}
}

func (s *ConsoleSender) SendSMS(phoneNumber string, message string) error {
	fmt.Printf("Уведомление отправлено на номер %s:\n%s\n", phoneNumber, message)
	return nil
}

var _ domain.SMSSender = (*ConsoleSender)(nil)
