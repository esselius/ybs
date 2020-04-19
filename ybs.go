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

type Budget struct {
	ID       string
	Name     string
}

type Account struct {
	ID   string
	Name string
	Note string
}

type BankService interface {
	Login(UserInterface) error
	Logout() error
	Transactions(Account) ([]Transaction, error)
}

type BudgetService interface {
	Budgets() ([]Budget, error)
	Accounts(Budget) ([]Account, error)
	AppendTransactions(Budget, Account, []Transaction) error
	BankImport(BankService, UserInterface) error
}

type UserInterface interface {
	Ask(string) (string, error)
	Choose(string, []string) (string, error)
	ShowQrCode(string) error
}

type Browser interface {
	Get(string) error
	ClickButton(string) error
	ClickLink(string) error
	ClickDiv(string) error
	TextField(string, string) error
	ScanQrCode() (string, error)
	LookFor(string) (bool, error)
	DownloadFolder() DownloadFolder
}

type DownloadFolder interface {
	LatestFileWithPrefix(string) (string, error)
}
