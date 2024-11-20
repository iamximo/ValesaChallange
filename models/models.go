package models

import "time"

type Account struct {
	ID      string  `json:"id"`
	Owner   string  `json:"owner"`
	Balance float32 `json:"balance"`
}

type Transaction struct {
	ID        string    `json:"id"`
	AccountID string    `json:"accountId"`
	Type      string    `json:"type"`
	Amount    float32   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type AccountRequest struct {
	Owner          string   `json:"owner"`
	InitialBalance *float32 `json:"initial_balance" binding:"required"`
}

type TransactionRequest struct {
	Type   string   `json:"type"`
	Amount *float32 `json:"amount" binding:"required"`
}

type TransferRequest struct {
	AccountIDfrom string   `json:"from_account_id"`
	AccountIDto   string   `json:"to_account_id"`
	Amount        *float32 `json:"amount" binding:"required"`
}
