package mock

import "github.com/esselius/ybs"

type UserInterface struct {
	AskFn func(message string) (string, error)
	AskInvoked bool

	ChooseFn func(message string, options []string) (string, error)
	ChooseInvoked bool

	ShowQrCodeFn func(message string) error
	ShowQrCodeInvoked bool

	ShowTransactionsFn func(transactions []ybs.Transaction) error
	ShowTransactionsInvoked bool

	RegisterFn func(string, func() error)
	RegisterInvoked bool

	StartFn func() error
	StartInvoked bool
}

func (ui *UserInterface) Ask(message string) (string, error) {
	ui.AskInvoked = true
	return ui.AskFn(message)
}

func (ui *UserInterface) Choose(message string, options []string) (string, error) {
	ui.ChooseInvoked = true
	return ui.ChooseFn(message, options)
}

func (ui *UserInterface) ShowQrCode(message string) error {
	ui.ShowQrCodeInvoked = true
	return ui.ShowQrCodeFn(message)
}

func (ui *UserInterface) ShowTransactions(transactions []ybs.Transaction) error {
	ui.ShowTransactionsInvoked = true
	return ui.ShowTransactionsFn(transactions)
}

func (ui *UserInterface) Register(name string, fn func() error) {
	ui.RegisterInvoked = true
	ui.RegisterFn(name, fn)
}

func (ui *UserInterface) Start() error {
	ui.StartInvoked = true
	return ui.StartFn()
}
