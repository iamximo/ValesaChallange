package storage

import (
	"Valesa/Challange/models"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	accountStore     = make(map[string]*models.Account)
	transactions     []models.Transaction
	accountStoreLock sync.Mutex
	transactionsLock sync.Mutex
)

func Reset() {
	accountStoreLock.Lock()
	defer accountStoreLock.Unlock()

	transactionsLock.Lock()
	defer transactionsLock.Unlock()

	accountStore = make(map[string]*models.Account)
	transactions = []models.Transaction{}
}

func SaveAccount(account *models.Account) {
	accountStoreLock.Lock()
	defer accountStoreLock.Unlock()

	accountStore[account.ID] = account
}

func GetAccount(id string) (*models.Account, bool) {
	accountStoreLock.Lock()
	defer accountStoreLock.Unlock()

	account, exists := accountStore[id]
	return account, exists
}

func GetAllAccounts() []models.Account {
	accountStoreLock.Lock()
	defer accountStoreLock.Unlock()

	var accounts []models.Account
	for _, acc := range accountStore {
		accounts = append(accounts, *acc)
	}
	return accounts
}

func GetTransactionsByAccountID(accountID string) []models.Transaction {
	transactionsLock.Lock()
	defer transactionsLock.Unlock()

	var accountTransactions []models.Transaction
	for _, tr := range transactions {
		if tr.AccountID == accountID {
			accountTransactions = append(accountTransactions, tr)
		}
	}
	return accountTransactions
}

func DoTransaction(act *models.Account, amount float32, typ string) (*models.Transaction, error) {
	accountStoreLock.Lock()
	defer accountStoreLock.Unlock()

	transactionsLock.Lock()
	defer transactionsLock.Unlock()

	transaction := &models.Transaction{
		ID:        uuid.New().String(),
		AccountID: act.ID,
		Type:      typ,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	if typ == "deposit" {
		act.Balance += amount
		transactions = append(transactions, *transaction)
		return transaction, nil
	}

	if amount > act.Balance {
		return nil, errors.New("insufficient credit")
	}

	act.Balance -= amount
	transactions = append(transactions, *transaction)
	return transaction, nil
}
