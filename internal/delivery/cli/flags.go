package cli

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/icoder-new/installment-cli/internal/domain"
)

type Flags struct {
	Help          bool
	Interactive   bool
	ProductType   string
	Price         float64
	PhoneNumber   string
	Months        int
}

func ParseFlags() *Flags {
	flags := &Flags{}

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

	flag.Parse()

	return flags
}

func (f *Flags) Validate() error {
	if f.Help {
		flag.Usage()
		os.Exit(0)
	}

	if !f.Interactive && (f.ProductType == "" || f.Price <= 0 || f.PhoneNumber == "" || f.Months == 0) {
		flag.Usage()
		return fmt.Errorf("все флаги обязательны в нон-интерактивном режиме или используйте интерактивный режим (-i/--interactive)")
	}

	if f.ProductType != "" {
		if !isValidProductType(f.ProductType) {
			return fmt.Errorf("неверный тип товара: %s. Допустимые значения: Смартфон, Компьютер, Телевизор", f.ProductType)
		}
	}

	if f.Price < 0 {
		return fmt.Errorf("цена товара не может быть отрицательной")
	}
	if f.PhoneNumber != "" && !isValidPhoneNumber(f.PhoneNumber) {
		return fmt.Errorf("неверный формат номера телефона")
	}

	if f.Months < 0 {
		return fmt.Errorf("срок рассрочки не может быть отрицательным")
	}

	return nil
}

func (f *Flags) ToProduct() domain.Product {
	return domain.Product{
		Type:         domain.ProductType(f.ProductType),
		Price:        f.Price,
		PhoneNumber:  f.PhoneNumber,
		PeriodMonths: f.Months,
	}
}

func isValidProductType(productType string) bool {
	switch strings.ToLower(productType) {
	case "смартфон", "компьютер", "телевизор":
		return true
	default:
		return false
	}
}

func isValidPhoneNumber(phone string) bool {
	re := regexp.MustCompile(`[^0-9]`)
	cleanPhone := re.ReplaceAllString(phone, "")

	switch len(cleanPhone) {
	case 9:
		return true
	case 12:
		return strings.HasPrefix(cleanPhone, "992")
	default:
		return false
	}
}
