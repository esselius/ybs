package skandia

import (
	"os"

	"github.com/esselius/ybs"
)

type SkandiaFile struct {
	Path string
}

func (s SkandiaFile) Login(ui ybs.UserInterface) error { panic("not implemented") }
func (s SkandiaFile) Logout() error                    { panic("not implemented") }

func (s SkandiaFile) Transactions(account ybs.BankAccount, ui ybs.UserInterface) ([]ybs.Transaction, error) {
	fileInfo, err := os.Stat(s.Path)
	if err != nil {
		return nil, err
	}

	transactions := []ybs.Transaction{}

	if fileInfo.IsDir() {
		files, err := MatchExcelExportFiles(s.Path, account.Number)
		if err != nil {
			return nil, err
		}
		files, err = ui.ChooseMultiple("Choose export files", files)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			ts, err := ExcelToTransactions(f)
			if err != nil {
				return nil, err
			}

			transactions = append(transactions, ts...)
		}
	} else {
		if IsExcelExportFile(fileInfo, account.Number) {
			transactions, err = ExcelToTransactions(s.Path)
			if err != nil {
				return nil, err
			}
		}
	}

	return transactions, nil
}
