package main

import (
	"log"
	"os"

	"github.com/esselius/ybs/pkg/chrome"
	"github.com/esselius/ybs/pkg/skandia"
	"github.com/esselius/ybs/pkg/sync"
	"github.com/esselius/ybs/pkg/terminal"
	"github.com/esselius/ybs/pkg/youneedabudget"
)

func main() {
	userInterface := terminal.New()

	budget := youneedabudget.New(os.Getenv("YNAB_TOKEN"))

	chromeBrowser, err := chrome.New(true)
	if err != nil {
		log.Fatal(err)
	}
	defer chromeBrowser.Close()

	bank := skandia.Skandia{
		Browser: chromeBrowser,
		UserInterface: userInterface,
	}

	err = sync.BankToBudget(bank, budget, userInterface)
	if err != nil {
		log.Fatal(err)
	}
}
