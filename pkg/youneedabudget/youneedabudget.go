package youneedabudget

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"go.bmvs.io/ynab"
	"go.bmvs.io/ynab/api"
	"go.bmvs.io/ynab/api/transaction"

	"github.com/esselius/ybs"
)

type YNAB struct {
	client ynab.ClientServicer
}

type Budget struct {
	ID   string
	Name string
}

type Account struct {
	ID   string
	Name string
	Note string
}

func New(token string) YNAB {
	client := ynab.NewClient(token)
	return YNAB{client}
}

func (y YNAB) Budgets() ([]Budget, error) {
	budgets, err := y.client.Budget().GetBudgets()
	if err != nil {
		return nil, err
	}

	var result []Budget
	for _, b := range budgets {
		result = append(result, Budget{
			ID:   b.ID,
			Name: b.Name,
		})
	}

	return result, nil
}

func (y YNAB) Accounts(budget Budget) ([]Account, error) {
	accounts, err := y.client.Account().GetAccounts(budget.ID)
	if err != nil {
		return nil, err
	}

	var result []Account
	for _, a := range accounts {
		if a.Note == nil || a.Deleted {
			continue
		}
		result = append(result, Account{
			ID:   a.ID,
			Name: a.Name,
			Note: *a.Note,
		})
	}

	return result, nil
}

func (y YNAB) AppendTransactions(budget Budget, bankAccount ybs.BankAccount, tranactions []ybs.Transaction) ([]ybs.Transaction, error) {
	var payloadTransactions []transaction.PayloadTransaction

	var account Account
	accounts, err := y.Accounts(budget)
	if err != nil {
		return nil, err
	}

	for _, a := range accounts {
		if a.Name == bankAccount.Name {
			account = a
			break
		}
	}
	if account.ID == "" {
		return nil, errors.New("no account found")
	}

	var importedPayloadTransactions []transaction.PayloadTransaction
	var importIds []string
	for _, t := range tranactions {
		var importIdPtr *string
		var importId string
		if t.Status == "cleared" {
			importId = generateImportId(importIds, t)
			importIdPtr = &importId
			importIds = append(importIds, importId)
		}
		desc := t.Description
		importedPayloadTransactions = append(importedPayloadTransactions, transaction.PayloadTransaction{
			AccountID: account.ID,
			Date:      api.Date{Time: t.Date},
			Amount:    int64(t.Amount * 1000),
			Cleared:   transaction.ClearingStatus(t.Status),
			Approved:  false,
			PayeeName: &desc,
			ImportID:  importIdPtr,
		})
	}

	existingTransactions, err := y.client.Transaction().GetTransactionsByAccount(
		budget.ID,
		account.ID,
		&transaction.Filter{
			Since: &api.Date{
				time.Now().AddDate(
					0,
					-1,
					0,
					),
			},
		})
	if err != nil {
		return nil, err
	}

	for _, pt := range importedPayloadTransactions {
		for _, et := range existingTransactions {
			if pt.Date == et.Date && pt.Amount == et.Amount && pt.Cleared == et.Cleared {
				break
			}
			payloadTransactions = append(payloadTransactions, importedPayloadTransactions...)
		}

	}

	createdTransactions, err := y.client.Transaction().CreateTransactions(budget.ID, payloadTransactions)
	if err != nil {
		return nil, err
	}

	var transactions []ybs.Transaction
	for _, t := range createdTransactions.Transactions {
		transactions = append(transactions, ybs.Transaction{
			Date:        t.Date.Time,
			Description: *t.PayeeName,
			Amount:      float64(t.Amount),
		})
	}

	return transactions, nil
}

func (y YNAB) BankImport(bank ybs.BankService, ui ybs.UserInterface) error {
	budgets, err := y.Budgets()
	if err != nil {
		return err
	}

	budget, err := chooseBudget(budgets, ui)
	if err != nil {
		return err
	}

	account, err := chooseAccount(y, budget, ui)
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

	transactions, err = y.AppendTransactions(budget, account, transactions)
	if err != nil {
		return err
	}

	return ui.ShowTransactions(transactions)
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

func chooseAccount(y YNAB, budget Budget, userInterface ybs.UserInterface) (ybs.BankAccount, error) {
	accounts, err := y.Accounts(budget)
	if err != nil {
		return ybs.BankAccount{}, err
	}

	var accountNames []string
	for _, a := range accounts {
		accountNames = append(accountNames, a.Name)
	}

	accountName, err := userInterface.Choose("Choose account", accountNames)
	if err != nil {
		return ybs.BankAccount{}, err
	}

	bankNumberRegex := regexp.MustCompile(`^Skandia: (.*)`)

	var account ybs.BankAccount
	for _, a := range accounts {
		if a.Name == accountName {
			account = ybs.BankAccount{
				Name:   a.Name,
				Number: bankNumberRegex.FindStringSubmatch(a.Note)[1],
			}
			break
		}
	}
	return account, nil
}

func chooseBudget(budgets []Budget, userInterface ybs.UserInterface) (Budget, error) {
	var budgetNames []string
	for _, budget := range budgets {
		budgetNames = append(budgetNames, budget.Name)
	}

	budgetName, err := userInterface.Choose("Choose budget", budgetNames)
	if err != nil {
		return Budget{}, err
	}

	var budget Budget
	for _, b := range budgets {
		if b.Name == budgetName {
			budget = b
			break
		}
	}

	return budget, nil
}
