package domain

import (
	"errors"
	"fmt"
)

type ProductType string

const (
	Smartphone ProductType = "Смартфон"
	Computer   ProductType = "Компьютер"
	TV         ProductType = "Телевизор"
)

var (
	ErrInvalidPrice       = errors.New("цена должна быть больше 0")
	ErrInvalidPhoneNumber = errors.New("необходимо указать номер телефона")
	ErrInvalidProductType = errors.New("неверный тип продукта")
	ErrInvalidPeriod      = errors.New("неверный срок рассрочки")
)

var validPeriods = []int{3, 6, 9, 12, 18, 24}

type Product struct {
	Type         ProductType
	Price        float64
	PhoneNumber  string
	PeriodMonths int
}

func (p *Product) Validate() error {
	if p.Price <= 0 {
		return ErrInvalidPrice
	}

	if p.PhoneNumber == "" {
		return ErrInvalidPhoneNumber
	}

	switch p.Type {
	case Smartphone, Computer, TV:
	default:
		return fmt.Errorf("%w: %s", ErrInvalidProductType, p.Type)
	}

	minPeriod, maxPeriod := p.getValidPeriods()

	if p.PeriodMonths < minPeriod || p.PeriodMonths > maxPeriod {
		return fmt.Errorf("%w: для %s допустимый срок от %d до %d месяцев",
			ErrInvalidPeriod, p.Type, minPeriod, maxPeriod)
	}

	if (p.PeriodMonths-3)%3 != 0 || p.PeriodMonths < 3 || p.PeriodMonths > 24 {
		return fmt.Errorf("%w: допустимые значения: %v", ErrInvalidPeriod, validPeriods)
	}

	return nil
}

func (p *Product) getValidPeriods() (int, int) {
	switch p.Type {
	case Smartphone:
		return 3, 9
	case Computer:
		return 3, 12
	case TV:
		return 3, 18
	default:
		return 0, 0
	}
}

func (p *Product) GetInterestRate() float64 {
	switch p.Type {
	case Smartphone:
		return 0.03
	case Computer:
		return 0.04
	case TV:
		return 0.05
	default:
		return 0
	}
}

func (p *Product) CalculateTotalPayment() float64 {
	if p.PeriodMonths <= 3 {
		return p.Price
	}

	extraPeriods := (p.PeriodMonths - 3) / 3
	interestRate := p.GetInterestRate()
	return p.Price * (1 + float64(extraPeriods)*interestRate)
}
