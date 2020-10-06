package mock

type Browser struct {
	GetFn      func(url string) error
	GetInvoked bool

	ClickButtonFn      func(text string) error
	ClickButtonInvoked bool

	ClickLinkFn      func(text string) error
	ClickLinkInvoked bool

	ClickDivFn      func(class string) error
	ClickDivInvoked bool

	TextFieldFn      func(name, text string) error
	TextFieldInvoked bool

	ScanQrCodeFn      func() (string, error)
	ScanQrCodeInvoked bool

	FindFn      func(selector, text string) (bool, error)
	FindInvoked bool

	DownloadDirectoryFn      func() string
	DownloadDirectoryInvoked bool

	CloseFn      func() error
	CloseInvoked bool
}

func (b *Browser) Get(url string) error {
	b.GetInvoked = true
	return b.GetFn(url)
}

func (b *Browser) ClickButton(text string) error {
	b.ClickButtonInvoked = true
	return b.ClickButtonFn(text)
}

func (b *Browser) ClickLink(text string) error {
	b.ClickLinkInvoked = true
	return nil
}

func (b *Browser) ClickDiv(class string) error {
	b.ClickDivInvoked = true
	return b.ClickDivFn(class)
}

func (b *Browser) TextField(name, text string) error {
	b.TextFieldInvoked = true
	return b.TextFieldFn(name, text)
}

func (b *Browser) ScanQrCode() (string, error) {
	b.ScanQrCodeInvoked = true
	return b.ScanQrCodeFn()
}

func (b *Browser) Find(selector, text string) (bool, error) {
	b.FindInvoked = true
	return b.FindFn(selector, text)
}

func (b *Browser) DownloadDirectory() string {
	b.DownloadDirectoryInvoked = true
	return b.DownloadDirectoryFn()
}

func (b *Browser) Close() error {
	b.CloseInvoked = true
	return b.CloseFn()
}
