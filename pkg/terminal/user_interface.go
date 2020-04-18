package terminal

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mdp/qrterminal"
)

type UserInterface struct {}

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

func (ui UserInterface) ShowQrCode(message string) error {
	qrterminal.GenerateHalfBlock(message, qrterminal.L, os.Stderr)
	return nil
}