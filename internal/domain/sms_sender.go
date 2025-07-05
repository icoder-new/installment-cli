package domain

type SMSSender interface {
	SendSMS(phoneNumber string, message string) error
}
