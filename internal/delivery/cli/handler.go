package cli

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/icoder-new/installment-cli/internal/domain"
	"github.com/icoder-new/installment-cli/internal/usecase"
)

type Handler struct {
	calculator *usecase.InstallmentCalculator
}

func NewHandler(calculator *usecase.InstallmentCalculator) *Handler {
	return &Handler{
		calculator: calculator,
	}
}

func (h *Handler) parseFlags() (domain.Product, error) {
	flags := ParseFlags()

	if err := flags.Validate(); err != nil {
		return domain.Product{}, err
	}
	if flags.Interactive {
		if flags.ProductType != "" || flags.Price > 0 || flags.PhoneNumber != "" || flags.Months > 0 {
			return h.collectInteractiveInputWithDefaults(flags)
		}
		return h.collectInteractiveInput()
	}

	return flags.ToProduct(), nil
}

func (h *Handler) collectInteractiveInputWithDefaults(flags *Flags) (domain.Product, error) {
	reader := bufio.NewReader(os.Stdin)
	var product domain.Product

	product.Type = h.promptProductType(reader, flags.ProductType)

	product.Price = h.promptPrice(reader, flags.Price)

	product.PhoneNumber = h.promptPhoneNumber(reader, flags.PhoneNumber)

	product.PeriodMonths = h.promptInstallmentPeriod(reader, flags.Months, product.Type)

	return product, nil
}

func (h *Handler) collectInteractiveInput() (domain.Product, error) {
	return h.collectInteractiveInputWithDefaults(&Flags{})
}

func (h *Handler) promptProductType(reader *bufio.Reader, defaultValue string) domain.ProductType {
	for {
		defaultChoice := ""
		switch strings.ToLower(defaultValue) {
		case "смартфон":
			defaultChoice = "1"
		case "компьютер":
			defaultChoice = "2"
		case "телевизор":
			defaultChoice = "3"
		}

		prompt := "Выберите тип товара (1-Смартфон, 2-Компьютер, 3-Телевизор)"
		if defaultChoice != "" {
			prompt += fmt.Sprintf(" [%s]", defaultChoice)
		}
		prompt += ": "

		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" && defaultChoice != "" {
			input = defaultChoice
		}

		if input == "" {
			fmt.Println("Ошибка: необходимо выбрать тип товара")
			continue
		}

		switch input {
		case "1":
			return domain.Smartphone
		case "2":
			return domain.Computer
		case "3":
			return domain.TV
		default:
			fmt.Println("Ошибка: выберите число от 1 до 3")
		}
	}
}

func (h *Handler) promptPrice(reader *bufio.Reader, defaultValue float64) float64 {
	for {
		prompt := "Введите цену товара (сомони)"
		if defaultValue > 0 {
			prompt += fmt.Sprintf(" [%s]", h.formatPrice(defaultValue))
		}
		prompt += ": "

		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" && defaultValue > 0 {
			return defaultValue
		}

		parsedPrice, err := strconv.ParseFloat(input, 64)
		if err != nil || parsedPrice <= 0 {
			fmt.Println("Ошибка: введите корректную цену (положительное число)")
			continue
		}

		return math.Round(parsedPrice*100) / 100
	}
}

func (h *Handler) promptPhoneNumber(reader *bufio.Reader, defaultValue string) string {
	for {
		prompt := "Введите номер телефона (в формате 992XXXXXXXXX)"
		if defaultValue != "" {
			prompt += fmt.Sprintf(" [%s]", defaultValue)
		}
		prompt += ": "

		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" && defaultValue != "" {
			return defaultValue
		} else if input == "" {
			fmt.Println("Ошибка: номер телефона не может быть пустым")
			continue
		}

		re := regexp.MustCompile(`[^0-9]`)
		cleanInput := re.ReplaceAllString(input, "")

		if len(cleanInput) == 9 {
			cleanInput = "992" + cleanInput
		}

		if !h.isValidPhoneNumber(cleanInput) {
			fmt.Println("Ошибка: неверный формат номера телефона. Используйте формат 992XXXXXXXXX")
			continue
		}

		return cleanInput
	}
}

func (h *Handler) promptInstallmentPeriod(reader *bufio.Reader, defaultValue int, productType domain.ProductType) int {
	minPeriod, maxPeriod := 0, 0
	switch productType {
	case domain.Smartphone:
		minPeriod, maxPeriod = 3, 9
	case domain.Computer:
		minPeriod, maxPeriod = 3, 12
	case domain.TV:
		minPeriod, maxPeriod = 3, 18
	}

	for {
		prompt := fmt.Sprintf("Введите срок рассрочки (от %d до %d месяцев)", minPeriod, maxPeriod)
		if defaultValue > 0 {
			prompt += fmt.Sprintf(" [%d]", defaultValue)
		}
		prompt += ": "

		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" && defaultValue > 0 {
			return defaultValue
		}

		parsedMonths, err := strconv.Atoi(input)
		if err != nil || parsedMonths < minPeriod || parsedMonths > maxPeriod {
			fmt.Printf("Ошибка: введите число от %d до %d\n", minPeriod, maxPeriod)
			continue
		}

		return parsedMonths
	}
}

func (h *Handler) Run() error {
	product, err := h.parseFlags()
	if err != nil {
		return err
	}

	totalPayment, err := h.calculator.CalculateInstallment(product)
	if err != nil {
		return fmt.Errorf("ошибка при расчете рассрочки: %w", err)
	}

	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║            РАССРОЧКА                   ║")
	fmt.Println("╠════════════════════════════════════════╣")
	fmt.Printf("║ %-15s %15.2f сомони ║\n", "Цена товара:", product.Price)
	fmt.Printf("║ %-16s %15s %-5d ║\n", "Срок:", "", product.PeriodMonths)
	fmt.Println("╠════════════════════════════════════════╣")
	fmt.Printf("║ %s %15.2f сомони ║\n", "Итоговая сумма:", totalPayment)
	fmt.Printf("║ %-15s %15.2f сомони ║\n", "Переплата:", totalPayment-product.Price)
	fmt.Println("╚════════════════════════════════════════╝")

	return nil
}

func (h *Handler) formatPrice(price float64) string {
	return strconv.FormatFloat(price, 'f', 2, 64)
}

func (h *Handler) isValidPhoneNumber(phone string) bool {
	matched, _ := regexp.MatchString(`^\+?992\d{9}$`, phone)
	return matched
}
