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
	Login() error
	Logout() error
	Transactions(account Account) ([]Transaction, error)
}

type BudgetService interface {
	Budgets() ([]Budget, error)
	Accounts(budget Budget) ([]Account, error)
	AppendTransactions(Budget, Account, []Transaction) error
}

type UserInterface interface {
	Ask(message string) (string, error)
	Choose(message string, options []string) (string, error)
	ShowQrCode(message string) error
}

type Browser interface {
	Get(url string) error
	ClickButton(text string) error
	ClickLink(text string) error
	ClickDiv(class string) error
	TextField(name, text string) error
	ScanQrCode() (string, error)
	LookFor(text string) (bool, error)
	DownloadFolder() DownloadFolder
}

type DownloadFolder interface {
	LatestFileWithPrefix(string) (string, error)
}
