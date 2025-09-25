package handler

import (
	"errors"
	"net/http"

	"transaction-tracker/internal/accounts/domain"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

var (
	errMissingAccountID = errors.New("accountID not found in context")
)

// getContextDependencies extrae el logger y el accountID del contexto de Gin.
func getContextDependencies(c *gin.Context) (*loggerModels.Logger, *domain.Account, error) {
	l, ok := c.Get("logger")
	if !ok {
		return nil, nil, errors.New("logger not found in context")
	}

	log := l.(*loggerModels.Logger)

	acc, ok := c.Get("account")
	if !ok {
		log.Error(loggerModels.LogProperties{
			Event: "account_id_not_found",
			Error: errMissingAccountID,
		})

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return nil, nil, errMissingAccountID
	}

	account := acc.(*domain.Account)

	return log, account, nil
}
