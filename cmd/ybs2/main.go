package main

import (
	"log"
	"os"

	"github.com/esselius/ybs"
	"github.com/esselius/ybs/pkg/skandia"
	"github.com/esselius/ybs/pkg/terminal"
	"github.com/esselius/ybs/pkg/youneedabudget"
)

var (
	userInterface ybs.UserInterface
	bankService   ybs.BankService
	budgetService ybs.BudgetService
	browser       ybs.Browser
)

func main() {
	path := os.Args[1]

	userInterface = terminal.New()

	budgetService = youneedabudget.New(os.Getenv("YNAB_TOKEN"))

	bankService := skandia.SkandiaFile{
		Path: path,
	}

	err := budgetService.BankImport(bankService, userInterface)
	if err != nil {
		log.Fatal(err)
	}
}
