package cli

import (
	"bufio"
	"flag"
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
	calculator      *usecase.InstallmentCalculator
	interactiveMode bool
}

func NewHandler(calculator *usecase.InstallmentCalculator) *Handler {
	return &Handler{
		calculator:      calculator,
		interactiveMode: false,
	}
}

func (h *Handler) parseFlags() (domain.Product, error) {
	help := flag.Bool("h", false, "Показать справку")
	helpLong := flag.Bool("help", false, "Показать справку (полная форма)")

	interactive := flag.Bool("i", false, "Интерактивный режим")
	interactiveLong := flag.Bool("interactive", false, "Интерактивный режим (полная форма)")

	productType := flag.String("p", "", "Тип товара (Смартфон, Компьютер, Телевизор)")
	productTypeLong := flag.String("product", "", "Тип товара (полная форма)")

	price := flag.Float64("c", 0, "Цена товара в сомони")
	priceLong := flag.Float64("cost", 0, "Цена товара в сомони (полная форма)")

	phone := flag.String("n", "", "Номер телефона клиента")
	phoneLong := flag.String("number", "", "Номер телефона клиента (полная форма)")

	months := flag.Int("m", 0, "Срок рассрочки в месяцах")
	monthsLong := flag.Int("months", 0, "Срок рассрочки в месяцах (полная форма)")

	flag.Parse()

	if *help || *helpLong {
		flag.Usage()
		os.Exit(0)
	}

	h.interactiveMode = *interactive || *interactiveLong

	if *productType == "" && *productTypeLong != "" {
		productType = productTypeLong
	}

	if *price == 0 && *priceLong != 0 {
		price = priceLong
	}

	if *phone == "" && *phoneLong != "" {
		phone = phoneLong
	}

	if *months == 0 && *monthsLong != 0 {
		months = monthsLong
	}

	if h.interactiveMode {
		if *productType != "" || *price > 0 || *phone != "" || *months > 0 {
			return h.collectInteractiveInputWithDefaults(*productType, *price, *phone, *months)
		}
		return h.collectInteractiveInput()
	}

	if *productType == "" || *price <= 0 || *phone == "" || *months == 0 {
		flag.Usage()
		return domain.Product{}, fmt.Errorf("необходимо указать все обязательные параметры или использовать интерактивный режим (-i/--interactive)")
	}

	return domain.Product{
		Type:         domain.ProductType(*productType),
		Price:        *price,
		PhoneNumber:  *phone,
		PeriodMonths: *months,
	}, nil
}

func (h *Handler) collectInteractiveInputWithDefaults(productType string, price float64, phone string, months int) (domain.Product, error) {
	reader := bufio.NewReader(os.Stdin)

	var product domain.Product

	for {
		defaultChoice := ""
		switch strings.ToLower(productType) {
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
			product.Type = domain.Smartphone
		case "2":
			product.Type = domain.Computer
		case "3":
			product.Type = domain.TV
		default:
			fmt.Println("Ошибка: выберите число от 1 до 3")
			continue
		}
		break
	}

	for {
		prompt := "Введите цену товара (сомони)"
		if price > 0 {
			prompt += fmt.Sprintf(" [%s]", h.formatPrice(price))
		}
		prompt += ": "

		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" && price > 0 {
			product.Price = price
			break
		}

		parsedPrice, err := strconv.ParseFloat(input, 64)
		if err != nil || parsedPrice <= 0 {
			fmt.Println("Ошибка: введите корректную цену (положительное число)")
			continue
		}

		product.Price = math.Round(parsedPrice*100) / 100
		break
	}

	for {
		prompt := "Введите номер телефона (в формате 992XXXXXXXXX)"
		if phone != "" {
			prompt += fmt.Sprintf(" [%s]", phone)
		}
		prompt += ": "

		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" && phone != "" {
			product.PhoneNumber = phone
			break
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

		product.PhoneNumber = cleanInput
		break
	}

	minPeriod, maxPeriod := 0, 0
	switch product.Type {
	case domain.Smartphone:
		minPeriod, maxPeriod = 3, 9
	case domain.Computer:
		minPeriod, maxPeriod = 3, 12
	case domain.TV:
		minPeriod, maxPeriod = 3, 18
	}

	for {
		prompt := fmt.Sprintf("Введите срок рассрочки (от %d до %d месяцев)", minPeriod, maxPeriod)
		if months >= minPeriod && months <= maxPeriod {
			prompt += fmt.Sprintf(" [%d]", months)
		}
		prompt += ": "
		fmt.Print(prompt)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" && months >= minPeriod && months <= maxPeriod {
			product.PeriodMonths = months
			break
		} else if input == "" {
			fmt.Printf("Ошибка: введите число от %d до %d\n", minPeriod, maxPeriod)
			continue
		}

		parsedMonths, err := strconv.Atoi(input)
		if err != nil || parsedMonths < minPeriod || parsedMonths > maxPeriod {
			fmt.Printf("Ошибка: введите число от %d до %d\n", minPeriod, maxPeriod)
			continue
		}
		product.PeriodMonths = parsedMonths
		break
	}

	fmt.Println()

	return product, nil
}

func (h *Handler) collectInteractiveInput() (domain.Product, error) {
	return h.collectInteractiveInputWithDefaults("", 0, "", 0)
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
