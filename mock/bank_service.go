package mock

import "github.com/esselius/ybs"

type BankService struct {
	LoginFn func(userInterface ybs.UserInterface) error
	LoginInvoked bool

	LogoutFn func() error
	LogoutInvoked bool

	TransactionsFn func(account ybs.BankAccount) ([]ybs.Transaction, error)
	TransactionsInvoked bool
}

func (bs *BankService) Login(userInterface ybs.UserInterface) error {
	bs.LoginInvoked = true
	return bs.LoginFn(userInterface)
}

func (bs *BankService) Logout() error {
	bs.LogoutInvoked = true
	return bs.LogoutFn()
}

func (bs *BankService) Transactions(account ybs.BankAccount) ([]ybs.Transaction, error) {
	bs.TransactionsInvoked = true
	return bs.TransactionsFn(account)
}
