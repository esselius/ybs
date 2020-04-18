package sync

import (
	"strings"

	"github.com/esselius/ybs"
)

func BankToBudget(bankService ybs.BankService, budgetService ybs.BudgetService, userInterface ybs.UserInterface) error {
	budgets, err := budgetService.Budgets()
	if err != nil {
		return err
	}

	budget, err := chooseBudget(budgets, userInterface)
	if err != nil {
		return err
	}

	account, err := chooseAccount(budgetService, budget, userInterface)
	if err != nil {
		return err
	}

	err = bankService.Login()
	if err != nil {
		return err
	}

	transactions, err := bankService.Transactions(account)
	if err != nil {
		return err
	}

	err = bankService.Logout()
	if err != nil {
		return err
	}

	err = budgetService.AppendTransactions(budget, account, transactions)
	if err != nil {
		return err
	}

	return nil
}

func chooseAccount(y ybs.BudgetService, budget ybs.Budget, userInterface ybs.UserInterface) (ybs.Account, error) {
	accounts, err := y.Accounts(budget)
	if err != nil {
		return ybs.Account{}, err
	}

	var accountNames []string
	for _, a := range accounts {
		accountNames = append(accountNames, a.Name)
	}

	accountName, err := userInterface.Choose("Choose account", accountNames)
	if err != nil {
		return ybs.Account{}, err
	}

	var account ybs.Account
	for _, a := range accounts {
		if strings.Contains(a.Name, accountName) {
			account = a
			break
		}
	}
	return account, nil
}

func chooseBudget(budgets []ybs.Budget, userInterface ybs.UserInterface) (ybs.Budget, error) {
	var budgetNames []string
	for _, budget := range budgets {
		budgetNames = append(budgetNames, budget.Name)
	}

	budgetName, err := userInterface.Choose("Choose budget", budgetNames)
	if err != nil {
		return ybs.Budget{}, err
	}

	var budget ybs.Budget
	for _, b := range budgets {
		if b.Name == budgetName {
			budget = b
			break
		}
	}

	return budget, nil
}
