package cli

import (
	"flag"
	"strings"

	"github.com/icoder-new/installment-cli/internal/domain"
)

type Flags struct {
	Help        bool
	Interactive bool
	ProductType string
	Price       float64
	PhoneNumber string
	Months      int
}

type FlagParser struct {
	validator *FlagValidator
}

func NewFlagParser() *FlagParser {
	return &FlagParser{
		validator: NewFlagValidator(),
	}
}

func (fp *FlagParser) Parse() (*Flags, error) {
	flags := &Flags{}
	fp.defineFlags(flags)
	flag.Parse()

	return flags, fp.validator.Validate(flags)
}

func (fp *FlagParser) defineFlags(flags *Flags) {
	flag.BoolVar(&flags.Help, "h", false, "Показать помощь")
	flag.BoolVar(&flags.Help, "help", false, "Показать помощь (длинная форма)")

	flag.BoolVar(&flags.Interactive, "i", false, "Интерактивный режим")
	flag.BoolVar(&flags.Interactive, "interactive", false, "Интерактивный режим (длинная форма)")

	flag.StringVar(&flags.ProductType, "p", "", "Тип товара (Смартфон, Компьютер, Телевизор)")
	flag.StringVar(&flags.ProductType, "product", "", "Тип товара (длинная форма)")

	flag.Float64Var(&flags.Price, "c", 0, "Цена товара")
	flag.Float64Var(&flags.Price, "cost", 0, "Цена товара (длинная форма)")

	flag.StringVar(&flags.PhoneNumber, "n", "", "Номер телефона")
	flag.StringVar(&flags.PhoneNumber, "number", "", "Номер телефона (длинная форма)")

	flag.IntVar(&flags.Months, "m", 0, "Срок рассрочки")
	flag.IntVar(&flags.Months, "months", 0, "Срок рассрочки (длинная форма)")
}

func (f *Flags) ToProduct() domain.Product {
	var productType domain.ProductType
	switch strings.ToLower(f.ProductType) {
	case "1", "смартфон":
		productType = domain.Smartphone
	case "2", "компьютер":
		productType = domain.Computer
	case "3", "телевизор":
		productType = domain.TV
	default:
		productType = domain.ProductType(f.ProductType)
	}

	return domain.Product{
		Type:         productType,
		Price:        f.Price,
		PhoneNumber:  f.PhoneNumber,
		PeriodMonths: f.Months,
	}
}

func (f *Flags) HasPartialData() bool {
	return f.ProductType != "" || f.Price > 0 || f.PhoneNumber != "" || f.Months > 0
}

func (f *Flags) IsComplete() bool {
	return f.ProductType != "" && f.Price > 0 && f.PhoneNumber != "" && f.Months > 0
}
