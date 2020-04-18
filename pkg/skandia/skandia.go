package skandia

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"

	"github.com/esselius/ybs"
)

type Skandia struct {
	Browser       ybs.Browser
	UserInterface ybs.UserInterface
}

func (s Skandia) Login() error {
	personNumber, err := s.UserInterface.Ask("person number please")
	if err != nil {
		return err
	}

	err = s.Browser.Get("http://skandia.se")
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

	qrCode, err := s.Browser.ScanQrCode()
	if err != nil {
		return err
	}
	err = s.UserInterface.ShowQrCode(qrCode)
	if err != nil {
		return err
	}

	found, err := s.Browser.LookFor("Du är utloggad")
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
	return s.Browser.ClickLink("Logga ut")
}

func (s Skandia) Transactions(account ybs.Account) ([]ybs.Transaction, error) {
	bankAccount := s.BudgetAccountToBankAccount(account)
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

	time.Sleep(2 * time.Second)

	err = s.Browser.ClickLink("Exportera till Excel")
	if err != nil {
		return nil, err
	}

	time.Sleep(2 * time.Second)

	filename, err := s.Browser.DownloadFolder().LatestFileWithPrefix(bankAccount.Number)
	if err != nil {
		return nil, err
	}

	return ExcelToTransactions(filename)
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

func (s Skandia) BudgetAccountToBankAccount(account ybs.Account) ybs.BankAccount {
	bankNumberRegex := regexp.MustCompile(`^Skandia: (.*)`)
	bankAccount := ybs.BankAccount{
		Name:   account.Name,
		Number: bankNumberRegex.FindStringSubmatch(account.Note)[1],
	}
	return bankAccount
}

func cleanPayee(t ybs.Transaction) ybs.Transaction {
	result := t
	transactionDateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}( kontaktlös|) (.*)`)
	klarnaAutogiroRegex := regexp.MustCompile(`^Autogiro K\*(.*)`)
	genericAutogiroRegex := regexp.MustCompile(`^Autogiro (.*)`)
	swishRegex := regexp.MustCompile(`^Swish (till|från) (.*)`)

	switch {
	case transactionDateRegex.MatchString(t.Description):
		result.Description = transactionDateRegex.FindStringSubmatch(t.Description)[2]
	case klarnaAutogiroRegex.MatchString(t.Description):
		result.Description = klarnaAutogiroRegex.FindStringSubmatch(t.Description)[1]
	case genericAutogiroRegex.MatchString(t.Description):
		result.Description = genericAutogiroRegex.FindStringSubmatch(t.Description)[1]
	case swishRegex.MatchString(t.Description):
		result.Description = swishRegex.FindStringSubmatch(t.Description)[2]
	}

	return result
}
