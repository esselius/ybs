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
}
