package main

import (
	"fmt"
	"os"

	"github.com/icoder-new/installment-cli/internal/domain"
	"github.com/icoder-new/installment-cli/internal/infra/sms"
	"github.com/icoder-new/installment-cli/internal/usecase"
)

func main() {
	smsSender := sms.NewConsoleSender()
	calculator := usecase.NewInstallmentCalculator(smsSender)

	res, err := calculator.CalculateInstallment(domain.Product{
		Type:         domain.Computer,
		Price:        1000,
		PhoneNumber:  "+992909010101",
		PeriodMonths: 6,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Итого к оплате: %.2f сомони\n", res)
}
