package cli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/icoder-new/installment-cli/internal/domain"
)

type InputValidator interface {
	ValidateProductType(input string) (domain.ProductType, error)
	ValidatePrice(input string) (float64, error)
	ValidatePhoneNumber(phone string) (string, error)
	ValidateInstallmentPeriod(input string, productType domain.ProductType) (int, error)
	GetInstallmentPeriodRange(productType domain.ProductType) (min, max int)
}

type inputValidator struct{}

func NewInputValidator() InputValidator {
	return &inputValidator{}
}

func (v *inputValidator) ValidateProductType(input string) (domain.ProductType, error) {
	switch strings.ToLower(input) {
	case "смартфон":
		return domain.Smartphone, nil
	case "компьютер":
		return domain.Computer, nil
	case "телевизор":
		return domain.TV, nil
	default:
		return "", fmt.Errorf("неверный тип товара. Допустимые значения: Смартфон, Компьютер, Телевизор")
	}
}

func (v *inputValidator) ValidatePrice(input string) (float64, error) {
	price, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("введите корректное число")
	}

	if price <= 0 {
		return 0, fmt.Errorf("цена должна быть положительным числом")
	}

	return price, nil
}

func (v *inputValidator) ValidatePhoneNumber(phone string) (string, error) {
	cleanPhone := regexp.MustCompile(`[^0-9]`).ReplaceAllString(phone, "")

	switch {
	case strings.HasPrefix(cleanPhone, "992") && len(cleanPhone) == 12:
	case len(cleanPhone) == 9:
		cleanPhone = "992" + cleanPhone
	case strings.HasPrefix(phone, "+") && len(cleanPhone) >= 10:
		cleanPhone = cleanPhone[len(cleanPhone)-9:]
		cleanPhone = "992" + cleanPhone
	default:
		return "", fmt.Errorf("неверный формат номера телефона. Используйте формат: 992XXXXXXXXX, 9XXXXXXXX, или +992XXXXXXXXX")
	}

	if len(cleanPhone) != 12 || !strings.HasPrefix(cleanPhone, "992") {
		return "", fmt.Errorf("неверный формат номера телефона. Проверьте введенный номер")
	}

	return cleanPhone, nil
}

func (v *inputValidator) ValidateInstallmentPeriod(input string, productType domain.ProductType) (int, error) {
	months, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("введите корректное число")
	}

	min, max := v.GetInstallmentPeriodRange(productType)
	if months < min || months > max {
		return 0, fmt.Errorf("срок рассрочки должен быть от %d до %d месяцев", min, max)
	}

	return months, nil
}

func (v *inputValidator) GetInstallmentPeriodRange(productType domain.ProductType) (min, max int) {
	switch productType {
	case domain.Smartphone:
		return 3, 9
	case domain.Computer:
		return 3, 24
	case domain.TV:
		return 3, 12
	default:
		return 3, 24
	}
}
