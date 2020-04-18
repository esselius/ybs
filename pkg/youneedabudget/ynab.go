package youneedabudget

import (
	"fmt"

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

func (y YNAB) Budgets() ([]ybs.Budget, error) {
	budgets, err := y.client.Budget().GetBudgets()
	if err != nil {
		return nil, err
	}

	var result []ybs.Budget
	for _, b := range budgets {
		result = append(result, ybs.Budget{
			ID: b.ID,
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
			ID:a.ID,
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
		importId :=  generateImportId(importIds, t)
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
