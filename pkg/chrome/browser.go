package chrome

import (
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/sclevine/agouti"

	"github.com/esselius/ybs"
)

type Browser struct {
	driver         *agouti.WebDriver
	page           *agouti.Page
	downloadFolder ybs.DownloadFolder
}

func New(headless bool) (Browser, error) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return Browser{}, err
	}
	downloadFolder := DownloadFolder{Directory: tmpDir}

	options := []agouti.Option{
		agouti.ChromeOptions("prefs", map[string]map[string]string{
			"download": {
				"default_directory": downloadFolder.Directory,
			},
		}),
	}
	if headless {
		options = append(options, agouti.ChromeOptions("args", []string{"--headless"}))

	}
	driver := agouti.ChromeDriver(options...)

	err = driver.Start()
	if err != nil {
		return Browser{}, err
	}

	page, err := driver.NewPage(agouti.Browser("chrome"))
	if err != nil {
		return Browser{}, err
	}

	err = page.SetImplicitWait(20000)
	if err != nil {
		return Browser{}, err
	}

	return Browser{driver, page, downloadFolder}, nil
}

func (b *Browser) Close() error {
	return b.driver.Stop()
}

func (b Browser) Get(url string) error {
	return b.page.Navigate(url)
}

func (b Browser) ClickLink(text string) error {
	selection := b.page.FindByLink(text)

	return selection.Click()
}

func (b Browser) ClickButton(text string) error {
	selection := b.page.FindByButton(text)

	return selection.Click()
}

func (b Browser) ClickDiv(class string) error {
	selection := b.page.FindByClass(class)

	return selection.Click()
}

func (b Browser) TextField(name, text string) error {
	selection := b.page.FindByName(name)
	return selection.SendKeys(text)
}

func (b Browser) LookFor(text string) (bool, error) {
	selector := "#he-main-wrapper > main > header > section > div > h1"
	content, err := b.page.Find(selector).Text()
	if err != nil && err.Error() != fmt.Sprintf("failed to select element from selection 'CSS: %s [single]': element not found", selector) {
		return false, err
	}
	return content == text, nil
}

func (b Browser) ScanQrCode() (string, error) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	err = file.Close()
	if err != nil {
		return "", err
	}

	err = b.page.Screenshot(file.Name())
	if err != nil {
		return "", err
	}
	imgdata, err := os.Open(file.Name())
	if err != nil {
		return "", err
	}

	img, _, err := image.Decode(imgdata)
	if err != nil {
		return "", err
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.DecodeWithoutHints(bmp)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

func (b Browser) DownloadFolder() ybs.DownloadFolder {
	return b.downloadFolder
}
