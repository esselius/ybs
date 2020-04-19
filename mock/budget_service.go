package mock

import "github.com/esselius/ybs"

var transactions []ybs.Transaction

type BudgetService struct {}

var budgets []ybs.Budget = []ybs.Budget{
	{
		ID:   "uuid",
		Name: "home budget",
	},
}
var accounts []ybs.Account = []ybs.Account{
	{
		ID:   "uuid",
		Name: "main account",
		Note: "BudgetService: 123",
	},
}

func (m BudgetService) Budgets() ([]ybs.Budget, error) {
	return budgets, nil
}

func (m BudgetService) Accounts(budget ybs.Budget) ([]ybs.Account, error) {
	return accounts, nil
}

func (m BudgetService) Transactions(budget ybs.Budget, account ybs.Account) ([]ybs.Transaction, error) {
	return transactions, nil
}

func (m BudgetService) AppendTransactions(_ ybs.Budget, _ ybs.Account, t []ybs.Transaction) error {
	transactions = append(transactions, t...)
	return nil
}
