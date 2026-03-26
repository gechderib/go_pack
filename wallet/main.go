package main

import (
	"fmt"
	"wallet/account"
)

func main() {

	acc := account.Account{Balance: 100.0, Logger: account.ConsoleLogger{}}

	acc.Deposit(50.0)

	err := acc.Withdraw(30.0)
	if err != nil {
		fmt.Println(err)
	}

	balance := acc.GetBalance()

	fmt.Printf("Balance: %.2f\n", balance)
}
