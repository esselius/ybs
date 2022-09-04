package skandia

import (
	"github.com/esselius/ybs"
)

type SkandiaFile struct {
	Path string
}

func (s SkandiaFile) Login(ui ybs.UserInterface) error { panic("not implemented") }
func (s SkandiaFile) Logout() error                    { panic("not implemented") }

func (s SkandiaFile) Transactions(account ybs.BankAccount) ([]ybs.Transaction, error) {
	return ExcelToTransactions(s.Path)
}
