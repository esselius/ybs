package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/esselius/ybs"
	"github.com/esselius/ybs/pkg/skandia_export"
	"github.com/esselius/ybs/pkg/youneedabudget"
)

func main() {
	path := os.Args[1]

	bankService := skandia_export.SkandiaExport{Path: path}
	budgetService := youneedabudget.New(os.Getenv("YNAB_TOKEN"))

	budgets, err := budgetService.Budgets()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Date,Payee,Amount")
	for _, b := range budgets {
		accounts, err := budgetService.Accounts(b)
		if err != nil {
			log.Fatal(err)
		}

		for _, a := range accounts {
			transactions, err := bankService.Transactions(ybs.BankAccount{
				Number: regexp.MustCompile(`^Skandia: (.*)`).FindStringSubmatch(a.Note)[1],
			}, nil)
			if err != nil {
				log.Fatal(err)
			}

			for _, t := range transactions {
				// fmt.Println(t)
				fmt.Printf("%s,\"%s\",%.2f\n", t.Date.Format("2006-01-02"), t.Description, t.Amount)
			}
		}
	}
	// bankService := skandia.SkandiaFile{
	// 	Path: path,
	// }

	// budgets, err := budgetService.Budgets()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// budget, err := youneedabudget.ChooseBudget(budgets, userInterface)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// account, err := budgetService.ChooseAccount(budget, userInterface)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// transactions, err := bankService.Transactions(account, userInterface)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// transactions, err = budgetService.AppendTransactions(budget, account, transactions)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// userInterface.ShowTransactions(transactions)
}
