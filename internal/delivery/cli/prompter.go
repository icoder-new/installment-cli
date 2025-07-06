package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/icoder-new/installment-cli/internal/domain"
)

const (
	productTypePrompt       = "Выберите тип товара (1-Смартфон, 2-Компьютер, 3-Телевизор)"
	pricePrompt             = "Введите цену товара (сомони)"
	phonePrompt             = "Введите номер телефона (в формате 992XXXXXXXXX)"
	installmentPeriodPrompt = "Введите срок рассрочки (доступно: %s)"
)

var (
	commandMap = map[domain.ProductType]string{
		domain.Smartphone: "1",
		domain.Computer:   "2",
		domain.TV:         "3",
	}
)

type PromptBuilder struct {
	basePrompt   string
	defaultValue string
}

func NewPromptBuilder(basePrompt string) *PromptBuilder {
	return &PromptBuilder{basePrompt: basePrompt}
}

func (pb *PromptBuilder) WithDefault(defaultValue string) *PromptBuilder {
	pb.defaultValue = defaultValue
	return pb
}

func (pb *PromptBuilder) Build() string {
	prompt := pb.basePrompt
	if pb.defaultValue != "" {
		prompt += fmt.Sprintf(" [%s]", pb.defaultValue)
	}
	return prompt + ": "
}

type UserPrompter struct {
	reader    *bufio.Reader
	validator InputValidator
}

func NewUserPrompter() *UserPrompter {
	return &UserPrompter{
		reader:    bufio.NewReader(os.Stdin),
		validator: NewInputValidator(),
	}
}

func (p *UserPrompter) PromptProductType(defaultValue string) domain.ProductType {
	defaultChoice := p.getDefaultProductTypeChoice(defaultValue)
	promptBuilder := NewPromptBuilder(productTypePrompt).WithDefault(defaultChoice)

	for {
		fmt.Print(promptBuilder.Build())
		input := p.readInput()
		input = p.handleDefaultValue(input, defaultChoice)

		switch input {
		case "1":
			return domain.Smartphone
		case "2":
			return domain.Computer
		case "3":
			return domain.TV
		}

		productType, err := p.validator.ValidateProductType(input)
		if err == nil {
			return productType
		}

		fmt.Println("Ошибка: выберите 1, 2 или 3, либо введите название товара")
	}
}

func (p *UserPrompter) PromptPrice(defaultValue float64) float64 {
	defaultChoice := p.getDefaultPriceChoice(defaultValue)
	promptBuilder := NewPromptBuilder(pricePrompt).WithDefault(defaultChoice)

	return p.promptFloat64WithValidation(promptBuilder, defaultChoice,
		func(input string) (float64, error) {
			return p.validator.ValidatePrice(input)
		})
}

func (p *UserPrompter) PromptPhoneNumber(defaultValue string) string {
	promptBuilder := NewPromptBuilder(phonePrompt).WithDefault(defaultValue)

	return p.promptStringWithValidation(promptBuilder, defaultValue,
		func(input string) (string, error) {
			return p.validator.ValidatePhoneNumber(input)
		})
}

func (p *UserPrompter) PromptInstallmentPeriod(defaultValue int, productType domain.ProductType) int {
	allowedPeriods := []int{3, 6, 9, 12, 18, 24}
	var availablePeriods []string
	var maxPeriod int

	switch productType {
	case domain.Smartphone:
		maxPeriod = 9
	case domain.Computer:
		maxPeriod = 12
	case domain.TV:
		maxPeriod = 18
	default:
		maxPeriod = 24
	}

	for _, period := range allowedPeriods {
		if period <= maxPeriod {
			availablePeriods = append(availablePeriods, strconv.Itoa(period))
		}
	}

	basePrompt := fmt.Sprintf("Выберите срок рассрочки (доступно: %s)",
		strings.Join(availablePeriods, ", "))

	defaultChoice := ""
	if defaultValue > 0 {
		defaultChoice = strconv.Itoa(defaultValue)
	}

	promptBuilder := NewPromptBuilder(basePrompt).WithDefault(defaultChoice)

	return p.promptIntWithValidation(promptBuilder, defaultChoice,
		func(input string) (int, error) {
			period, err := strconv.Atoi(input)
			if err != nil {
				return 0, fmt.Errorf("введите корректное число")
			}

			valid := false
			for _, p := range allowedPeriods {
				if p == period && period <= maxPeriod {
					valid = true
					break
				}
			}

			if !valid {
				return 0, fmt.Errorf("неверный срок рассрочки. Доступные значения: %v", availablePeriods)
			}

			return period, nil
		})
}

func (p *UserPrompter) promptStringWithValidation(
	promptBuilder *PromptBuilder,
	defaultValue string,
	validator func(string) (string, error),
) string {
	for {
		fmt.Print(promptBuilder.Build())

		input := p.readInput()
		input = p.handleDefaultValue(input, defaultValue)

		if input == "" {
			fmt.Println("Ошибка: поле не может быть пустым")
			continue
		}

		result, err := validator(input)
		if err != nil {
			fmt.Printf("Ошибка: %s\n", err.Error())
			continue
		}

		return result
	}
}

func (p *UserPrompter) promptIntWithValidation(
	promptBuilder *PromptBuilder,
	defaultValue string,
	validator func(string) (int, error),
) int {
	for {
		fmt.Print(promptBuilder.Build())

		input := p.readInput()
		input = p.handleDefaultValue(input, defaultValue)

		if input == "" {
			fmt.Println("Ошибка: поле не может быть пустым")
			continue
		}

		result, err := validator(input)
		if err != nil {
			fmt.Printf("Ошибка: %s\n", err.Error())
			continue
		}

		return result
	}
}

func (p *UserPrompter) promptFloat64WithValidation(
	promptBuilder *PromptBuilder,
	defaultValue string,
	validator func(string) (float64, error),
) float64 {
	for {
		fmt.Print(promptBuilder.Build())

		input := p.readInput()
		input = p.handleDefaultValue(input, defaultValue)

		if input == "" {
			fmt.Println("Ошибка: поле не может быть пустым")
			continue
		}

		result, err := validator(input)
		if err != nil {
			fmt.Printf("Ошибка: %s\n", err.Error())
			continue
		}

		return result
	}
}

func (p *UserPrompter) readInput() string {
	input, _ := p.reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func (p *UserPrompter) handleDefaultValue(input, defaultValue string) string {
	if input == "" && defaultValue != "" {
		return defaultValue
	}
	return input
}

func (p *UserPrompter) getDefaultProductTypeChoice(defaultValue string) string {
	if defaultValue == "" {
		return ""
	}

	productType := domain.ProductType(strings.ToLower(defaultValue))
	if command, exists := commandMap[productType]; exists {
		return command
	}

	return ""
}

func (p *UserPrompter) getDefaultPriceChoice(defaultValue float64) string {
	if defaultValue <= 0 {
		return ""
	}
	return formatPrice(defaultValue)
}
