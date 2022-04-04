package skandia

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"

	"github.com/esselius/ybs"
)

type Skandia struct {
	Browser ybs.Browser
}

func (s Skandia) Login(ui ybs.UserInterface) error {
	personNumber, err := ui.Ask("person number please")
	if err != nil {
		return err
	}

	err = s.Browser.Get("http://skandia.se")
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)


	err = s.Browser.ClickButton("Acceptera alla cookies")
	if err != nil {
		return err
	}

	err = s.Browser.ClickLink("Logga in")
	if err != nil {
		return err
	}

	err = s.Browser.ClickLink("Mobilt BankID")
	if err != nil {
		return err
	}

	err = s.Browser.TextField("NationalIdentificationNumber", personNumber)
	if err != nil {
		return err
	}

	err = s.Browser.ClickButton("Logga in med Mobilt BankID")
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	qrCode, err := s.Browser.ScanQrCode()
	if err != nil {
		return err
	}
	err = ui.ShowQrCode(qrCode)
	if err != nil {
		return err
	}

	found, err := s.Browser.Find("#he-main-wrapper > main > header > section > div > h1", "Du är utloggad")
	if err != nil {
		return err
	}
	if found {
		return errors.New("logged out")
	}

	return nil
}

func (s Skandia) Logout() error {
	err := s.Browser.ClickDiv("he-avatar")
	if err != nil {
		return err
	}

	err = s.Browser.ClickLink("Logga ut")
	if err != nil {
		return err
	}

	_, err = s.Browser.Find("#he-main-wrapper > main > header > section > div > h1", "Nu är du utloggad")
	return err
}

func (s Skandia) Transactions(bankAccount ybs.BankAccount) ([]ybs.Transaction, error) {
	err := s.Browser.ClickButton("Konton")
	if err != nil {
		return nil, err
	}

	err = s.Browser.ClickLink("Kontoöversikt")
	if err != nil {
		return nil, err
	}

	err = s.Browser.ClickLink(fmt.Sprintf("%s (%s)", bankAccount.Name, bankAccount.Number))
	if err != nil {
		return nil, err
	}

	// TODO Figure out how to wait for elements showing up
	time.Sleep(5 * time.Second)

	err = s.Browser.ClickLink("Exportera till Excel")
	if err != nil {
		return nil, err
	}

	// TODO Check file with exponential backoff instead
	time.Sleep(10 * time.Second)

	filename, err := latestFileWithPrefix(s.Browser.DownloadDirectory(), bankAccount.Number)
	if err != nil {
		return nil, err
	}

	return ExcelToTransactions(filename)
}

func latestFileWithPrefix(path, prefix string) (string, error) {
	fileList, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	var files []os.FileInfo
	for _, f := range fileList {
		if strings.HasPrefix(f.Name(), prefix) {
			files = append(files, f)
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	return fmt.Sprintf("%s/%s", path, files[len(files)-1].Name()), nil
}

func ExcelToTransactions(filename string) ([]ybs.Transaction, error) {
	file, err := xlsx.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	sheet := file.Sheet["Kontoutdrag"]

	var transactions []ybs.Transaction
	for i := 4; i <= sheet.MaxRow; i++ {
		if sheet.Row(i).Cells == nil {
			break
		}
		var transaction ybs.Transaction
		row := sheet.Row(i)

		date, err := time.Parse("2006-01-02", row.Cells[0].Value)
		if err != nil {
			return nil, err
		}
		amount, err := strconv.ParseFloat(row.Cells[2].Value, 64)
		if err != nil {
			return nil, err
		}
		transaction = ybs.Transaction{
			Date:        date,
			Description: row.Cells[1].Value,
			Amount:      amount,
			Status:      "cleared",
		}
		transactions = append(transactions, transaction)
	}

	return scrub(transactions), nil
}

func scrub(transactions []ybs.Transaction) []ybs.Transaction {
	var result []ybs.Transaction
	for _, t := range transactions {
		r := cleanPayee(t)
		result = append(result, r)
	}

	return result
}

func cleanPayee(t ybs.Transaction) ybs.Transaction {
	result := t
	transactionDateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}( kontaktlös|) (.*)`)
	klarnaAutogiroRegex := regexp.MustCompile(`^Autogiro K\*(.*)`)
	klarnaRegex := regexp.MustCompile(`^Klarna \* (.*)`)
	genericAutogiroRegex := regexp.MustCompile(`^Autogiro (.*)`)
	swishRegex := regexp.MustCompile(`^Swish (till|från) (.*)`)

	switch {
	case transactionDateRegex.MatchString(t.Description):
		result.Description = transactionDateRegex.FindStringSubmatch(t.Description)[2]
	case klarnaAutogiroRegex.MatchString(t.Description):
		result.Description = klarnaAutogiroRegex.FindStringSubmatch(t.Description)[1]
	case klarnaRegex.MatchString(t.Description):
		result.Description = klarnaRegex.FindStringSubmatch(t.Description)[1]
	case genericAutogiroRegex.MatchString(t.Description):
		result.Description = genericAutogiroRegex.FindStringSubmatch(t.Description)[1]
	case swishRegex.MatchString(t.Description):
		result.Description = swishRegex.FindStringSubmatch(t.Description)[2]
	}

	return result
}
