package terminal

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mdp/qrterminal"

	"github.com/esselius/ybs"
)

type UserInterface struct {
	endpoints []Endpoint
}

type Endpoint struct {
	name string
	fn   func() error
}

func New() UserInterface {
	return UserInterface{}
}

func (ui UserInterface) Ask(message string) (string, error) {
	var response string
	prompt := &survey.Input{
		Message: message,
	}
	err := survey.AskOne(prompt, &response)
	return response, err
}

func (ui UserInterface) Choose(message string, options []string) (string, error) {
	var response string
	prompt := &survey.Select{
		Message: message,
		Options: options,
	}
	err := survey.AskOne(prompt, &response)
	return response, err
}

func (ui UserInterface) ChooseMultiple(message string, options []string) ([]string, error) {
	var response []string
	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
	}
	err := survey.AskOne(prompt, &response)
	return response, err
}

func (ui UserInterface) ShowQrCode(message string) error {
	qrterminal.GenerateHalfBlock(message, qrterminal.L, os.Stderr)
	return nil
}

func (ui UserInterface) ShowTransactions(transactions []ybs.Transaction) error {
	fmt.Println("date,description,amount")
	for _, t := range transactions {
		fmt.Printf("%s,%s,%f\n", t.Date.Format("2006-01-02"), t.Description, t.Amount)
	}
	return nil
}
