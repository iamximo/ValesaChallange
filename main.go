package main

import (
	"Valesa/Challange/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/accounts", handlers.HandlerCreateAccount)
	r.GET("/accounts", handlers.HandlerGetAccounts)
	r.GET("/accounts/:id", handlers.HandlerGetAccountByID)
	r.GET("/accounts/:id/transactions", handlers.HandlerGetTransactionByID)
	r.POST("/accounts/:id/transactions", handlers.HandlerCreateTransaction)
	r.POST("/transfer", handlers.HandlerTransfer)

	r.Run("localhost:8080")
}
