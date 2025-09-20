package handler

import (
	"errors"
	"net/http"

	"transaction-tracker/api/services/accounts"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

// getContextDependencies extrae el logger y el accountID del contexto de Gin.
func getContextDependencies(c *gin.Context) (*loggerModels.Logger, *accounts.Account, error) {
	l, ok := c.Get("logger")
	if !ok {
		return nil, nil, errors.New("logger not found in context")
	}

	log := l.(*loggerModels.Logger)

	acc, ok := c.Get("account")
	if !ok {
		log.Error(loggerModels.LogProperties{
			Event: "account_id_not_found",
			Error: errors.New("accountID not found in context"),
		})

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return nil, nil, errors.New("accountID not found")
	}

	account := acc.(*accounts.Account)

	return log, account, nil
}
