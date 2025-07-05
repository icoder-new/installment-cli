package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/icoder-new/installment-cli/internal/delivery/cli"
	"github.com/icoder-new/installment-cli/internal/infra/sms"
	"github.com/icoder-new/installment-cli/internal/usecase"
)

func main() {
	smsSender := sms.NewConsoleSender()
	calculator := usecase.NewInstallmentCalculator(smsSender)
	handler := cli.NewHandler(calculator)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Использование: %s [ПАРАМЕТРЫ]

Параметры:
  -h, --help             Показать эту справку
  -i, --interactive      Включить интерактивный режим
  -p, --product ТОВАР    Тип товара (Смартфон, Компьютер, Телевизор)
  -c, --cost ЦЕНА       Цена товара в сомони
  -n, --number НОМЕР    Номер телефона клиента
  -m, --months МЕСЯЦЫ   Срок рассрочки в месяцах

Примеры:
  %[1]s -p Смартфон -c 1000 -n +992001234567 -m 6
  %[1]s --product=Компьютер --cost=2000 --number=+992001234567 --months=12
  %[1]s -i
  %[1]s --interactive

Для интерактивного режима можно указать часть параметров, 
а остальные ввести в диалоговом режиме:
  %[1]s -p Телевизор -i

`, os.Args[0])
	}

	if err := handler.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}
