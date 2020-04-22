package mock

import "github.com/esselius/ybs"

type BudgetService struct {
	BankImportFn func(bankService ybs.BankService, userInterface ybs.UserInterface) error
	BankImportInvoked bool
}

func (bs *BudgetService) BankImport(bankService ybs.BankService, userInterface ybs.UserInterface) error {
	bs.BankImportInvoked = true
	return bs.BankImportFn(bankService, userInterface)
}