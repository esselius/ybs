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
)

type Browser struct {
	driver            *agouti.WebDriver
	page              *agouti.Page
	downloadDirectory string
}

func New(headless bool) (Browser, error) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return Browser{}, err
	}

	options := []agouti.Option{
		agouti.ChromeOptions("prefs", map[string]map[string]string{
			"download": {
				"default_directory": tmpDir,
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

	return Browser{driver, page, tmpDir}, nil
}

func (b Browser) Close() error {
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

func (b Browser) ClickDivByAttribute(attribute string, value string) error {
	selection := b.page.Find(fmt.Sprintf("div[%s=\"%s\"]", attribute, value))

	return selection.Click()
}

func (b Browser) TextField(name, text string) error {
	selection := b.page.FindByName(name)
	return selection.SendKeys(text)
}

func (b Browser) Find(selector, text string) (bool, error) {
	content, err := b.page.Find(selector).Text()
	if err != nil && err.Error() != fmt.Sprintf("failed to select element from selection 'CSS: %s [single]': element not found", selector) {
		return false, err
	}
	return content == text, nil
}

func (b Browser) Table(selector string) ([][]string, error) {
	var result [][]string
	table := b.page.Find(selector)

	rows := table.All("tbody tr")
	rowCount, err := rows.Count()
	if err != nil {
		return nil, err
	}

	for i := 0; i < rowCount; i++ {
		cells := rows.At(i).All("td")
		cellCount, err := cells.Count()
		if err != nil {
			return nil, err
		}

		var rowCells []string
		for j := 0; j < cellCount; j++ {
			cell, err := cells.At(j).Text()
			if err != nil {
				return nil, err
			}
			rowCells = append(rowCells, cell)

		}
		result = append(result, rowCells)
	}

	return result, nil
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

func (b Browser) DownloadDirectory() string {
	return b.downloadDirectory
}
