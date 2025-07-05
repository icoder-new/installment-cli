package usecase

import (
	"fmt"

	"github.com/icoder-new/installment-cli/internal/domain"
)

type InstallmentCalculator struct {
	smsSender domain.SMSSender
}

func NewInstallmentCalculator(smsSender domain.SMSSender) *InstallmentCalculator {
	return &InstallmentCalculator{
		smsSender: smsSender,
	}
}

func (uc *InstallmentCalculator) CalculateInstallment(product domain.Product) (float64, error) {
	totalPayment := product.CalculateTotalPayment()
	overpayment := totalPayment - product.Price

	message := fmt.Sprintf(
		"Уважаемый клиент!\n"+
			"Детали вашей покупки:\n"+
			"Товар: %s\n"+
			"Сумма: %.2f сомони\n"+
			"Срок рассрочки: %d мес.\n"+
			"Переплата: %.2f сомони\n"+
			"Итого к оплате: %.2f сомони",
		product.Type,
		product.Price,
		product.PeriodMonths,
		overpayment,
		totalPayment,
	)

	if err := uc.smsSender.SendSMS(product.PhoneNumber, message); err != nil {
		return 0, fmt.Errorf("не удалось отправить SMS: %w", err)
	}

	return totalPayment, nil
}
