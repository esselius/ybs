package main

import (
	"log"
	"os"

	"github.com/esselius/ybs"
	"github.com/esselius/ybs/pkg/chrome"
	"github.com/esselius/ybs/pkg/skandia"
	"github.com/esselius/ybs/pkg/terminal"
	"github.com/esselius/ybs/pkg/youneedabudget"
)

var (
	userInterface ybs.UserInterface
	bankService ybs.BankService
	budgetService ybs.BudgetService
	browser ybs.Browser
)

func main() {
	userInterface = terminal.New()

	budgetService = youneedabudget.New(os.Getenv("YNAB_TOKEN"))

	var err error
	browser, err = chrome.New(true)
	if err != nil {
		log.Fatal(err)
	}
	defer browser.Close()

	bankService = skandia.Skandia{
		Browser: browser,
	}

	err = bankService.Login(userInterface)
	if err != nil {
		log.Fatal(err)
	}

	err = budgetService.BankImport(bankService, userInterface)
	if err != nil {
		log.Fatal(err)
	}

	err = bankService.Logout()
	if err != nil {
		log.Fatal(err)
	}
}
