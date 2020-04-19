package mock

import (
	"time"

	"github.com/esselius/ybs"
)

var t1 []ybs.Transaction = []ybs.Transaction{
	{
		Date:        time.Time{},
		Description: "Grandma",
		Amount:      100,
	},
}

type BankService struct {}

func (m BankService) Login() error {
	return nil
}

func (m BankService) Logout() error {
	return nil
}

func (m BankService) Transactions(account ybs.Account) ([]ybs.Transaction, error) {
	return t1, nil
}
