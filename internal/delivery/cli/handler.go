package cli

import (
	"fmt"

	"github.com/icoder-new/installment-cli/internal/domain"
	"github.com/icoder-new/installment-cli/internal/usecase"
)

type Handler struct {
	calculator *usecase.InstallmentCalculator
	flagParser *FlagParser
	prompter   *UserPrompter
	printer    *ResultPrinter
}

func NewHandler(calculator *usecase.InstallmentCalculator) *Handler {
	return &Handler{
		calculator: calculator,
		flagParser: NewFlagParser(),
		prompter:   NewUserPrompter(),
		printer:    NewResultPrinter(),
	}
}

func (h *Handler) Run() error {
	product, err := h.parseAndCollectInput()
	if err != nil {
		return err
	}

	totalPayment, err := h.calculator.CalculateInstallment(product)
	if err != nil {
		return fmt.Errorf("ошибка при расчете рассрочки: %w", err)
	}

	h.printer.PrintInstallmentResult(product, totalPayment)
	return nil
}

func (h *Handler) parseAndCollectInput() (domain.Product, error) {
	flags, err := h.flagParser.Parse()
	if err != nil {
		return domain.Product{}, err
	}

	if flags.Interactive {
		return h.handleInteractiveMode(flags)
	}

	return flags.ToProduct(), nil
}

func (h *Handler) handleInteractiveMode(flags *Flags) (domain.Product, error) {
	if flags.HasPartialData() {
		return h.collectInteractiveInputWithDefaults(flags), nil
	}
	return h.collectInteractiveInput(), nil
}

func (h *Handler) collectInteractiveInputWithDefaults(flags *Flags) domain.Product {
	product := domain.Product{
		Type:        h.prompter.PromptProductType(flags.ProductType),
		Price:       h.prompter.PromptPrice(flags.Price),
		PhoneNumber: h.prompter.PromptPhoneNumber(flags.PhoneNumber),
	}

	product.PeriodMonths = h.prompter.PromptInstallmentPeriod(flags.Months, product.Type)

	return product
}

func (h *Handler) collectInteractiveInput() domain.Product {
	return h.collectInteractiveInputWithDefaults(&Flags{})
}
