package skandia_export_test

import (
	"testing"
	"time"

	"github.com/esselius/ybs/pkg/skandia_export"
	"github.com/stretchr/testify/assert"
)

var accountTransactions = []skandia_export.AccountTransaction{
	{
		AccountNumber:   "1234-567.890-1",
		BookingDate:     time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
		TransactionDate: time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
		Description:     "WorldOfGolf",
		Amount:          -19500,
		Balance:         1943182,
		Type:            skandia_export.Autogiro,
	},
	{
		AccountNumber:   "1234-567.890-1",
		BookingDate:     time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
		TransactionDate: time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
		Description:     "DINEL",
		Amount:          -53600,
		Balance:         1962682,
		Type:            skandia_export.Autogiro,
	},
	{
		AccountNumber:   "1234-567.890-1",
		BookingDate:     time.Date(2022, time.August, 30, 0, 0, 0, 0, time.UTC),
		TransactionDate: time.Date(2022, time.August, 29, 0, 0, 0, 0, time.UTC),
		Description:     "JoeAndTheJuice, GÃ¶teborg",
		Amount:          -15800,
		Balance:         2016282,
		Type:            skandia_export.Generic,
	},
	{
		AccountNumber:   "1234-567.890-1",
		BookingDate:     time.Date(2022, time.August, 30, 0, 0, 0, 0, time.UTC),
		TransactionDate: time.Date(2022, time.August, 29, 0, 0, 0, 0, time.UTC),
		Description:     "VOI TECHNOLOGY, STOCKHOLM",
		Amount:          -2250,
		Balance:         2018532,
		Type:            skandia_export.Generic,
	},
}
var fixtures = map[string]skandia_export.Export{
	"fixtures/12345678901_2022-08-30-2022-08-31.xlsx": {
		AccountNumber: "1234-567.890-1",
		Period:        "2022-08-30 - 2022-08-31",
		Transactions: []skandia_export.ExportTransaction{
			{
				Date:        time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
				Description: "Autogiro WorldOfGolf",
				Amount:      -19500,
				Balance:     1943182,
			},
			{
				Date:        time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
				Description: "Autogiro DINEL",
				Amount:      -53600,
				Balance:     1962682,
			},
			{
				Date:        time.Date(2022, time.August, 30, 0, 0, 0, 0, time.UTC),
				Description: "2022-08-29 JoeAndTheJuice, GÃ¶teborg",
				Amount:      -15800,
				Balance:     2016282,
			},
			{
				Date:        time.Date(2022, time.August, 30, 0, 0, 0, 0, time.UTC),
				Description: "2022-08-29 VOI TECHNOLOGY, STOCKHOLM",
				Amount:      -2250,
				Balance:     2018532,
			},
		},
	},
	"fixtures/12345678901_2022-08-31-2022-08-31.xlsx": {
		AccountNumber: "1234-567.890-1",
		Period:        "2022-08-31 - 2022-08-31",
		Transactions: []skandia_export.ExportTransaction{
			{
				Date:        time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
				Description: "Autogiro WorldOfGolf",
				Amount:      -19500,
				Balance:     1943182,
			},
			{
				Date:        time.Date(2022, time.August, 31, 0, 0, 0, 0, time.UTC),
				Description: "Autogiro DINEL",
				Amount:      -53600,
				Balance:     1962682,
			},
		},
	},
}

func TestSingleFile(t *testing.T) {
	file := "fixtures/12345678901_2022-08-31-2022-08-31.xlsx"

	export, err := skandia_export.ReadExport(file)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fixtures[file], export)
}

func TestReadAccountTransactions(t *testing.T) {
	at, err := skandia_export.ReadAccountTransactions("fixtures")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, accountTransactions, at)
}
