package youneedabudget

import (
	"fmt"
	"strings"

	"go.bmvs.io/ynab"
	"go.bmvs.io/ynab/api"
	"go.bmvs.io/ynab/api/transaction"

	"github.com/esselius/ybs"
)

type YNAB struct {
	client ynab.ClientServicer
}

func New(token string) YNAB {
	client := ynab.NewClient(token)
	return YNAB{client}
}

func (y YNAB) Transactions(budget ybs.Budget, account ybs.Account) ([]ybs.Transaction, error) {
	panic("implement me")
}

func (y YNAB) Budgets() ([]ybs.Budget, error) {
	budgets, err := y.client.Budget().GetBudgets()
	if err != nil {
		return nil, err
	}

	var result []ybs.Budget
	for _, b := range budgets {
		result = append(result, ybs.Budget{
			ID:   b.ID,
			Name: b.Name,
		})
	}

	return result, nil
}

func (y YNAB) Accounts(budget ybs.Budget) ([]ybs.Account, error) {
	accounts, err := y.client.Account().GetAccounts(budget.ID)
	if err != nil {
		return nil, err
	}

	var result []ybs.Account
	for _, a := range accounts {
		if a.Note == nil || a.Deleted {
			continue
		}
		result = append(result, ybs.Account{
			ID:   a.ID,
			Name: a.Name,
			Note: *a.Note,
		})
	}

	return result, nil
}

func (y YNAB) AppendTransactions(budget ybs.Budget, account ybs.Account, tranactions []ybs.Transaction) error {
	var payloadTransactions []transaction.PayloadTransaction

	var importIds []string
	for _, t := range tranactions {
		importId := generateImportId(importIds, t)
		importIds = append(importIds, importId)
		desc := t.Description
		payloadTransactions = append(payloadTransactions, transaction.PayloadTransaction{
			AccountID: account.ID,
			Date:      api.Date{Time: t.Date},
			Amount:    int64(t.Amount * 1000),
			Cleared:   "cleared",
			Approved:  false,
			PayeeName: &desc,
			ImportID:  &importId,
		})
	}

	_, err := y.client.Transaction().CreateTransactions(budget.ID, payloadTransactions)
	if err != nil {
		return err
	}

	return nil
}

func (y YNAB) BankImport(bank ybs.BankService, tty ybs.UserInterface) error {
	budgets, err := y.Budgets()
	if err != nil {
		return err
	}

	budget, err := chooseBudget(budgets, tty)
	if err != nil {
		return err
	}

	account, err := chooseAccount(y, budget, tty)
	if err != nil {
		return err
	}

	transactions, err := bank.Transactions(account)
	if err != nil {
		return err
	}

	err = bank.Logout()
	if err != nil {
		return err
	}

	err = y.AppendTransactions(budget, account, transactions)
	if err != nil {
		return err
	}

	return nil
}

func generateImportId(importIds []string, transaction ybs.Transaction) string {
	var importId string
	for i := 1; i < 999; i++ {
		importId = importIdFormat(transaction, i)
		if !contains(importIds, importId) {
			break
		}
	}
	return importId
}

func importIdFormat(transaction ybs.Transaction, i int) string {
	return fmt.Sprintf("YNAB:%d:%s:%d", int64(transaction.Amount*1000), transaction.Date.Format("2006-01-02"), i)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
