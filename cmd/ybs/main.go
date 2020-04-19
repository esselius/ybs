package main

import (
	"log"
	"os"

	"github.com/esselius/ybs/pkg/browser"
	"github.com/esselius/ybs/pkg/skandia"
	"github.com/esselius/ybs/pkg/terminal"
	"github.com/esselius/ybs/pkg/youneedabudget"
)

func main() {
	tty := terminal.New()

	ynab := youneedabudget.New(os.Getenv("YNAB_TOKEN"))

	chrome, err := browser.New(true)
	if err != nil {
		log.Fatal(err)
	}
	defer chrome.Close()

	bank := skandia.Skandia{
		Browser: chrome,
	}

	err = bank.Login(tty)
	if err != nil {
		log.Fatal(err)
	}

	err = ynab.BankImport(bank, tty)
	if err != nil {
		log.Fatal(err)
	}
}
