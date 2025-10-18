package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	accountsDomain "transaction-tracker/internal/accounts/domain"
	"transaction-tracker/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupTestContext(method, target string, body io.Reader) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(method, target, body)

	c.Request = req

	mockLogger, _ := logger.GetLogger(c, "test")
	c.Set("logger", mockLogger)

	c.Set("account", &accountsDomain.Account{ID: "accountID"})

	return c, w
}

func TestGetContextDependencies(t *testing.T) {
	t.Run("logger not in context", func(t *testing.T) {
		c := require.New(t)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

		log, acc, err := getContextDependencies(ctx)
		c.Nil(log)
		c.Nil(acc)
		c.EqualError(err, "logger not found in context")
	})

	t.Run("account not in context", func(t *testing.T) {
		c := require.New(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		logger, err := logger.GetLogger(ctx, "test")
		c.NoError(err)
		ctx.Set("logger", logger)

		_, _, err = getContextDependencies(ctx)
		c.ErrorIs(err, errMissingAccountID)
		c.Equal(500, w.Code)
	})

	t.Run("logger and account in context", func(t *testing.T) {
		c := require.New(t)
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

		logger, err := logger.GetLogger(ctx, "test")
		c.NoError(err)
		account := &accountsDomain.Account{ID: "acc123"}

		ctx.Set("logger", logger)
		ctx.Set("account", account)

		log, acc, err := getContextDependencies(ctx)
		c.NoError(err)
		c.Equal(logger, log)
		c.Equal(account, acc)
	})
}
