package cli

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/icoder-new/installment-cli/internal/domain"
)

const (
	phonePrefix      = "992"
	localPhoneLength = 9
	fullPhoneLength  = 12
)

var (
	phoneRegex     = regexp.MustCompile(`[^0-9]`)
	phoneValidator = regexp.MustCompile(`^\+?992\d{9}$`)

	validProductTypes = map[string]bool{
		string(domain.Smartphone): true,
		string(domain.Computer):   true,
		string(domain.TV):         true,
	}
)

type FlagValidator struct{}

func NewFlagValidator() *FlagValidator {
	return &FlagValidator{}
}

func (fv *FlagValidator) Validate(flags *Flags) error {
	if flags.Help {
		flag.Usage()
		os.Exit(0)
	}

	if err := fv.validateNonInteractiveMode(flags); err != nil {
		return err
	}

	return fv.validateFields(flags)
}

func (fv *FlagValidator) validateNonInteractiveMode(flags *Flags) error {
	if !flags.Interactive && !flags.IsComplete() {
		flag.Usage()
		return fmt.Errorf("все флаги обязательны в нон-интерактивном режиме или используйте интерактивный режим (-i/--interactive)")
	}
	return nil
}

func (fv *FlagValidator) validateFields(flags *Flags) error {
	if err := fv.validateProductType(flags.ProductType); err != nil {
		return err
	}

	if err := fv.validatePrice(flags.Price); err != nil {
		return err
	}

	if err := fv.validatePhoneNumber(flags.PhoneNumber); err != nil {
		return err
	}

	if flags.ProductType != "" && flags.Months > 0 {
		if err := fv.validateMonthsForProduct(flags.Months, flags.ProductType); err != nil {
			return err
		}
	}

	return fv.validateMonths(flags.Months)
}

func (fv *FlagValidator) validateProductType(productType string) error {
	if productType != "" {
		switch strings.ToLower(productType) {
		case "1", "смартфон":
			return nil
		case "2", "компьютер":
			return nil
		case "3", "телевизор":
			return nil
		default:
			return fmt.Errorf("неверный тип товара: %s. Допустимые значения: 1/Смартфон, 2/Компьютер, 3/Телевизор", productType)
		}
	}
	return nil
}

func (fv *FlagValidator) validatePrice(price float64) error {
	if price < 0 {
		return fmt.Errorf("цена товара не может быть отрицательной")
	}
	return nil
}

func (fv *FlagValidator) validatePhoneNumber(phoneNumber string) error {
	if phoneNumber == "" {
		return nil
	}

	cleanPhone := phoneRegex.ReplaceAllString(phoneNumber, "")
	
	switch len(cleanPhone) {
	case 9:
		return nil
	case 12:
		if cleanPhone[:3] == "992" {
			return nil
		}
	}

	return fmt.Errorf("неверный формат номера телефона. Используйте формат: 992XXXXXXXXX, 9XXXXXXXX или +992XXXXXXXXX")
}

func (fv *FlagValidator) validateMonthsForProduct(months int, productType string) error {
	if err := fv.validateMonths(months); err != nil {
		return err
	}

	var maxPeriod int
	switch strings.ToLower(productType) {
	case "1", "смартфон":
		maxPeriod = 9
	case "2", "компьютер":
		maxPeriod = 12
	case "3", "телевизор":
		maxPeriod = 18
	default:
		return fmt.Errorf("неизвестный тип товара: %s", productType)
	}

	if months > maxPeriod {
		return fmt.Errorf("для %s максимальный срок рассрочки %d месяцев", productType, maxPeriod)
	}

	return nil
}

func (fv *FlagValidator) validateMonths(months int) error {
	if months < 0 {
		return fmt.Errorf("срок рассрочки не может быть отрицательным")
	}

	validPeriods := []int{3, 6, 9, 12, 18, 24}
	valid := false
	for _, p := range validPeriods {
		if p == months {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("неверный срок рассрочки. Допустимые значения: %v", validPeriods)
	}

	return nil
}

func (fv *FlagValidator) isValidProductType(productType string) bool {
	return validProductTypes[strings.ToLower(productType)]
}

func (fv *FlagValidator) isValidPhoneNumber(phone string) bool {
	cleanPhone := phoneRegex.ReplaceAllString(phone, "")

	switch len(cleanPhone) {
	case localPhoneLength:
		return true
	case fullPhoneLength:
		return strings.HasPrefix(cleanPhone, phonePrefix)
	default:
		return false
	}
}
