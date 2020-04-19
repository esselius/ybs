package mock

type UserInterface struct {}

func (m UserInterface) Ask(message string) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (m UserInterface) Choose(message string, options []string) (string, error) {
	return options[0], nil
}

func (m UserInterface) ShowQrCode(message string) error {
	panic("not implemented") // TODO: Implement
}
