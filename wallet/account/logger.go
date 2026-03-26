package account

import "fmt"

type TransactionLogger interface {
	Log(transactionType string, amount float64, balance float64)
}

type ConsoleLogger struct{}

func (c ConsoleLogger) Log(transactionType string, amount float64, balance float64) {
	fmt.Printf("[LOG] %s: %.2f | Balance: %.2f\n", transactionType, amount, balance)
}
