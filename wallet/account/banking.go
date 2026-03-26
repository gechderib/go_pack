package account

import "fmt"

type Account struct {
	Balance float64
	Logger  TransactionLogger
}

func (a *Account) Deposit(amount float64) {
	a.Balance += amount
	if a.Logger != nil {
		a.Logger.Log("deposit", amount, a.Balance)
	}
}

func (a *Account) Withdraw(amount float64) error {
	if amount > a.Balance {
		// a.Logger.LogTransaction("withdrawal", amount)
		return fmt.Errorf("insufficient funds")
	}
	a.Balance -= amount
	// a.Logger.Log("withdrawal", amount, a.Balance)
	return nil
}

func (a *Account) GetBalance() float64 {
	// a.Logger.LogTransaction("balance inquiry", 0)
	return a.Balance
}
