package skandia_export

import (
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/esselius/ybs"
	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
)

type SEK int64

type ExportTransaction struct {
	Date        time.Time
	Description string
	Amount      SEK
	Balance     SEK
}

type Export struct {
	AccountNumber string
	Period        string
	Transactions  []ExportTransaction
}

type AccountTransactionType int

const (
	Generic AccountTransactionType = iota
	Klarna
	iZettle
	Zettle
	Autogiro
	Swish
	Transfer
)

type AccountTransaction struct {
	AccountNumber   string
	BookingDate     time.Time
	TransactionDate time.Time
	Description     string
	Amount          SEK
	Balance         SEK
	Type            AccountTransactionType
}

func ReadExport(file string) (Export, error) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		return Export{}, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	sheet := "Kontoutdrag"

	accountNumber, err := f.GetCellValue(sheet, "B1")
	if err != nil {
		return Export{}, err
	}

	period, err := f.GetCellValue(sheet, "B2")
	if err != nil {
		return Export{}, err
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return Export{}, err
	}

	var transactions []ExportTransaction
	for _, row := range rows[4:] {
		date, err := time.Parse("2006-01-02", row[0])
		if err != nil {
			return Export{}, err
		}

		amount, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return Export{}, err
		}

		balance, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return Export{}, err
		}

		transactions = append(transactions, ExportTransaction{
			Date:        date,
			Description: row[1],
			Amount:      SEK(math.Round(amount * 100)),
			Balance:     SEK(math.Round(balance * 100)),
		})
	}

	return Export{
		AccountNumber: accountNumber,
		Period:        period,
		Transactions:  transactions,
	}, nil
}

func ExportToAccountTransactions(export Export) []AccountTransaction {
	return lo.Map(export.Transactions, func(t ExportTransaction, _ int) AccountTransaction {
		return AccountTransaction{
			AccountNumber: export.AccountNumber,
			BookingDate:   t.Date,
			Description:   t.Description,
			Amount:        t.Amount,
			Balance:       t.Balance,
		}
	})
}

var transactionDate = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})( kontaktlös|) (.*)`)
var klarna1 = regexp.MustCompile(`^Autogiro K\*(.*)`)
var klarna2 = regexp.MustCompile(`^Klarna \* (.*)`)
var klarna3 = regexp.MustCompile(`^Utbetalning autogiro K\*(.*)`)
var autogiro = regexp.MustCompile(`^Autogiro (.*)`)
var swish = regexp.MustCompile(`^Swish (till|från) (.*)`)
var zettle = regexp.MustCompile(`^ZTL\*(.*)`)
var izettle = regexp.MustCompile(`^IZ \*(.*)`)
var transfer = regexp.MustCompile(`^Överf.*`)

func setTransactionDate(at AccountTransaction, _ int) AccountTransaction {
	if transactionDate.MatchString(at.Description) {
		transactionDate, err := time.Parse("2006-01-02", transactionDate.FindStringSubmatch(at.Description)[1])
		if err != nil {
			panic(err)
		}

		at.TransactionDate = transactionDate
	} else {
		at.TransactionDate = at.BookingDate
	}

	return at
}

func setTransactionType(at AccountTransaction, _ int) AccountTransaction {
	switch {
	case klarna1.MatchString(at.Description):
		at.Type = Klarna
	case klarna2.MatchString(at.Description):
		at.Type = Klarna
	case klarna3.MatchString(at.Description):
		at.Type = Klarna
	case autogiro.MatchString(at.Description):
		at.Type = Autogiro
	case swish.MatchString(at.Description):
		at.Type = Swish
	case zettle.MatchString(at.Description):
		at.Type = Zettle
	case izettle.MatchString(at.Description):
		at.Type = iZettle
	case transfer.MatchString(at.Description):
		at.Type = Transfer
	}

	return at
}

func cleanTransactionDescription(at AccountTransaction, _ int) AccountTransaction {
	if transactionDate.MatchString(at.Description) {
		at.Description = transactionDate.FindStringSubmatch(at.Description)[3]
	}
	if klarna1.MatchString(at.Description) {
		at.Description = klarna1.FindStringSubmatch(at.Description)[1]
	}
	if klarna2.MatchString(at.Description) {
		at.Description = klarna2.FindStringSubmatch(at.Description)[1]
	}
	if klarna3.MatchString(at.Description) {
		at.Description = klarna3.FindStringSubmatch(at.Description)[1]
	}
	if autogiro.MatchString(at.Description) {
		at.Description = autogiro.FindStringSubmatch(at.Description)[1]
	}
	if swish.MatchString(at.Description) {
		at.Description = swish.FindStringSubmatch(at.Description)[2]
	}
	if zettle.MatchString(at.Description) {
		at.Description = zettle.FindStringSubmatch(at.Description)[1]
	}
	if izettle.MatchString(at.Description) {
		at.Description = izettle.FindStringSubmatch(at.Description)[1]
	}

	return at
}

func filterExcelExports(file os.FileInfo, _ int) bool {
	pattern1 := regexp.MustCompile(`\d{11}_\d{4}-\d{2}-\d{2}-\d{4}-\d{2}-\d{2}\.xlsx`)
	pattern2 := regexp.MustCompile(`\d{4}-\d{3}\.\d{3}-\d_\d{4}-\d{2}-\d{2}-\d{4}-\d{2}-\d{2}\.xlsx`)

	switch {
	case pattern1.MatchString(file.Name()):
		return true
	case pattern2.MatchString(file.Name()):
		return true
	}

	return false
}

func ReadAccountTransactions(path string) ([]AccountTransaction, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return lo.Uniq(lo.FlatMap(lo.Filter(files, filterExcelExports), func(file os.FileInfo, _ int) []AccountTransaction {
		export, err := ReadExport(filepath.Join(path, file.Name()))
		if err != nil {
			panic(err)
		}
		return lo.Map(lo.Map(lo.Map(
			ExportToAccountTransactions(export),
			setTransactionDate),
			setTransactionType),
			cleanTransactionDescription)
	})), nil
}

type SkandiaExport struct {
	Path string
}

func (s SkandiaExport) Login(ui ybs.UserInterface) error { panic("not implemented") }
func (s SkandiaExport) Logout() error                    { panic("not implemented") }

func (s SkandiaExport) Transactions(account ybs.BankAccount, ui ybs.UserInterface) ([]ybs.Transaction, error) {
	ats, err := ReadAccountTransactions(s.Path)
	if err != nil {
		return nil, err
	}

	var tranactions []ybs.Transaction

	for _, at := range lo.Filter(ats, func(at AccountTransaction, _ int) bool { return at.AccountNumber == account.Number }) {
		tranactions = append(tranactions, ybs.Transaction{
			Date:        at.BookingDate,
			Description: at.Description,
			Amount:      float64(at.Amount) / 100,
			Status:      "cleared",
		})
	}

	return tranactions, nil
}
