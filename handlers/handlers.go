package handlers

import (
	"Valesa/Challange/models"
	"Valesa/Challange/storage"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ACT_DEPOSTIT   = "deposit"
	ACT_WITHDRAWAL = "withdrawal"
)

func getAccountByID(id string) (*models.Account, error) {
	acc, exists := storage.GetAccount(id)
	if !exists {
		return nil, errors.New("account not found")
	}
	return acc, nil
}

func HandlerCreateAccount(ctx *gin.Context) {
	var req models.AccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if *req.InitialBalance < 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Initial balance cannot be less than 0"})
		return
	}

	if req.Owner == "" {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Owner cannot be empty"})
		return
	}

	newAccount := &models.Account{
		ID:      uuid.New().String(),
		Owner:   req.Owner,
		Balance: *req.InitialBalance,
	}

	storage.SaveAccount(newAccount)

	ctx.IndentedJSON(http.StatusCreated, newAccount)
}

func HandlerGetAccountByID(ctx *gin.Context) {
	id := ctx.Param("id")
	account, err := getAccountByID(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, account)
}

func HandlerGetAccounts(ctx *gin.Context) {
	accounts := storage.GetAllAccounts()
	ctx.IndentedJSON(http.StatusOK, accounts)
}

func HandlerGetTransactionByID(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := getAccountByID(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	transactions := storage.GetTransactionsByAccountID(id)

	ctx.IndentedJSON(http.StatusOK, transactions)
}

func HandlerCreateTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	account, err := getAccountByID(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var req models.TransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if *req.Amount <= 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Amount cannot be 0 or negative"})
		return
	}

	if req.Type != ACT_DEPOSTIT && req.Type != ACT_WITHDRAWAL {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Transaction type must be 'deposit' or 'withdrawal'"})
		return
	}

	transaction, err := storage.DoTransaction(account, *req.Amount, req.Type)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, transaction)
}

func HandlerTransfer(ctx *gin.Context) {
	var req models.TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if req.AccountIDfrom == req.AccountIDto {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Account IDs must be different"})
		return
	}

	if *req.Amount <= 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Amount cannot be 0 or negative"})
		return
	}

	fromAccount, err := getAccountByID(req.AccountIDfrom)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid 'from' account ID"})
		return
	}

	toAccount, err := getAccountByID(req.AccountIDto)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid 'to' account ID"})
		return
	}

	withdrawal, err := storage.DoTransaction(fromAccount, *req.Amount, ACT_WITHDRAWAL)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds in 'from' account"})
		return
	}

	deposit, depositErr := storage.DoTransaction(toAccount, *req.Amount, ACT_DEPOSTIT)
	if depositErr != nil {
		_, rollbackErr := storage.DoTransaction(fromAccount, *req.Amount, ACT_DEPOSTIT)
		if rollbackErr != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Critical error: failed to rollback withdrawal"})
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete transfer"})
		return
	}
	ctx.IndentedJSON(http.StatusCreated, []models.Transaction{*withdrawal, *deposit})
}
