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
	bankService   ybs.BankService
	budgetService ybs.BudgetService
	browser       ybs.Browser
)

func main() {
	userInterface = terminal.New()

	budgetService = youneedabudget.New(os.Getenv("YNAB_TOKEN"))

	var err error
	browser, err = chrome.New(false)
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

	budgets, err := budgetService.Budgets()
	if err != nil {
		log.Fatal(err)
	}

	budget, err := youneedabudget.ChooseBudget(budgets, userInterface)
	if err != nil {
		log.Fatal(err)
	}

	account, err := budgetService.ChooseAccount(budget, userInterface)
	if err != nil {
		log.Fatal(err)
	}

	transactions, err := bankService.Transactions(account, userInterface)
	if err != nil {
		log.Fatal(err)
	}

	transactions, err = budgetService.AppendTransactions(budget, account, transactions)
	if err != nil {
		log.Fatal(err)
	}

	userInterface.ShowTransactions(transactions)

	err = bankService.Logout()
	if err != nil {
		log.Fatal(err)
	}
}
