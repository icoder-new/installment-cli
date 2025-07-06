package cli

import (
	"fmt"

	"github.com/icoder-new/installment-cli/internal/domain"
)

type ResultPrinter struct{}

func NewResultPrinter() *ResultPrinter {
	return &ResultPrinter{}
}

func (rp *ResultPrinter) PrintInstallmentResult(product domain.Product, totalPayment float64) {
	rp.printHeader()
	rp.printProductInfo(product)
	rp.printSeparator()
	rp.printTotalInfo(product, totalPayment)
	rp.printFooter()
}

func (rp *ResultPrinter) printHeader() {
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║            РАССРОЧКА                   ║")
}

func (rp *ResultPrinter) printProductInfo(product domain.Product) {
	fmt.Println("╠════════════════════════════════════════╣")
	fmt.Printf("║ %-15s %15.2f сомони ║\n", "Цена товара:", product.Price)
	fmt.Printf("║ %-16s %15s %-5d ║\n", "Срок:", "", product.PeriodMonths)
}

func (rp *ResultPrinter) printSeparator() {
	fmt.Println("╠════════════════════════════════════════╣")
}

func (rp *ResultPrinter) printTotalInfo(product domain.Product, totalPayment float64) {
	fmt.Printf("║ %s %15.2f сомони ║\n", "Итоговая сумма:", totalPayment)
	fmt.Printf("║ %-15s %15.2f сомони ║\n", "Переплата:", totalPayment-product.Price)
}

func (rp *ResultPrinter) printFooter() {
	fmt.Println("╚════════════════════════════════════════╝")
}
