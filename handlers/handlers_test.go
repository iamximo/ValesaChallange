package handlers

import (
	"Valesa/Challange/models"
	"Valesa/Challange/storage"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandlerCreateAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/accounts", HandlerCreateAccount)

	t.Run("Valid account creation", func(t *testing.T) {
		storage.Reset()
		initialBalance := float32(100)
		body := models.AccountRequest{
			Owner:          "Joaquin Gonzalez",
			InitialBalance: &initialBalance,
		}

		bodyJSON, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(bodyJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Account
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Joaquin Gonzalez", response.Owner)
		assert.Equal(t, initialBalance, response.Balance)
		assert.NotEmpty(t, response.ID)
	})

	t.Run("Invalid account creation (negative balance)", func(t *testing.T) {
		storage.Reset()
		initialBalance := float32(-50)
		body := models.AccountRequest{
			Owner:          "Joaq 2",
			InitialBalance: &initialBalance,
		}

		bodyJSON, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(bodyJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		expectedResponse := `{"error":"Initial balance cannot be less than 0"}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})

	t.Run("Empty Owner", func(t *testing.T) {
		storage.Reset()
		initialBalance := float32(50)
		body := models.AccountRequest{
			Owner:          "",
			InitialBalance: &initialBalance,
		}

		bodyJSON, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(bodyJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		expectedResponse := `{"error":"Owner cannot be empty"}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})
}

func TestHandlerGetAccountByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/accounts/:id", HandlerGetAccountByID)

	t.Run("Valid account retrieval", func(t *testing.T) {
		storage.Reset()
		account := &models.Account{
			ID:      "12345",
			Owner:   "Joaquin Gonzalez",
			Balance: 100,
		}
		storage.SaveAccount(account)

		req, _ := http.NewRequest(http.MethodGet, "/accounts/12345", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Account
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, account.ID, response.ID)
		assert.Equal(t, account.Owner, response.Owner)
		assert.Equal(t, account.Balance, response.Balance)
	})

	t.Run("Invalid account retrieval", func(t *testing.T) {
		storage.Reset()
		req, _ := http.NewRequest(http.MethodGet, "/accounts/456_not_exist", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		expectedResponse := `{"error":"Invalid account ID"}`
		assert.JSONEq(t, expectedResponse, w.Body.String())
	})
}

func TestHandlerGetAccounts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/accounts", HandlerGetAccounts)

	t.Run("Retrieve all accounts", func(t *testing.T) {
		storage.Reset()
		account := &models.Account{
			ID:      "12345",
			Owner:   "Joaquin Gonzalez",
			Balance: 100,
		}
		storage.SaveAccount(account)

		req, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.Account
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, account.ID, response[0].ID)
	})
}

func TestHandlerCreateTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/accounts/:id/transactions", HandlerCreateTransaction)

	t.Run("Valid deposit", func(t *testing.T) {
		storage.Reset()
		account := &models.Account{
			ID:      "12345",
			Owner:   "Joaquin Gonzalez",
			Balance: 100,
		}
		storage.SaveAccount(account)

		amount := float32(50)
		body := models.TransactionRequest{
			Type:   "deposit",
			Amount: &amount,
		}

		bodyJSON, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/accounts/12345/transactions", bytes.NewBuffer(bodyJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Transaction
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "deposit", response.Type)
		assert.Equal(t, amount, response.Amount)
		assert.Equal(t, account.ID, response.AccountID)

		updatedAccount, _ := storage.GetAccount(account.ID)
		assert.Equal(t, float32(150), updatedAccount.Balance)
	})
}

func TestHandlerTransfer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/transfer", HandlerTransfer)

	t.Run("Valid transfer", func(t *testing.T) {
		storage.Reset()
		accountFrom := &models.Account{
			ID:      "12345",
			Owner:   "Joaquin Gonzalez",
			Balance: 100,
		}
		accountTo := &models.Account{
			ID:      "67890",
			Owner:   "Joaq 2",
			Balance: 50,
		}
		storage.SaveAccount(accountFrom)
		storage.SaveAccount(accountTo)

		amount := float32(50)
		body := models.TransferRequest{
			AccountIDfrom: "12345",
			AccountIDto:   "67890",
			Amount:        &amount,
		}

		bodyJSON, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/transfer", bytes.NewBuffer(bodyJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		updatedFrom, _ := storage.GetAccount(accountFrom.ID)
		updatedTo, _ := storage.GetAccount(accountTo.ID)
		assert.Equal(t, float32(50), updatedFrom.Balance)
		assert.Equal(t, float32(100), updatedTo.Balance)
	})
}
