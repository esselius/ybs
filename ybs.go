package ybs

import "time"

type Transaction struct {
	Date        time.Time `xlsx:"0"`
	Description string    `xlsx:"1"`
	Amount      float64   `xlsx:"1"`
}

type BankAccount struct {
	Name   string
	Number string
}

type BankService interface {
	Login(UserInterface) error
	Logout() error
	Transactions(BankAccount) ([]Transaction, error)
}

type BudgetService interface {
	BankImport(BankService, UserInterface) error
}

type UserInterface interface {
	Ask(string) (string, error)
	Choose(string, []string) (string, error)
	ShowQrCode(string) error
	ShowTransactions([]Transaction) error
}

type Browser interface {
	Get(string) error
	ClickButton(string) error
	ClickLink(string) error
	ClickDiv(string) error
	TextField(string, string) error
	ScanQrCode() (string, error)
	Find(string, string) (bool, error)
	DownloadDirectory() string
	Close() error
}
