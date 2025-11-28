package handler

import (
	"transaction-tracker/api/models"
	"transaction-tracker/internal/accounts/usecase"
	"transaction-tracker/logger"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AccountHandler handles HTTP requests for the accounts domain.
type AccountHandler struct {
	accountsUsecase usecase.AccountsUsecase
}

// NewAccountHandler creates a new instance of AccountsHandler.
func NewAccountHandler(ucm usecase.AccountsUsecase) *AccountHandler {
	return &AccountHandler{
		accountsUsecase: ucm,
	}
}

// GoogleGenerateAuthLink handles the GET /accounts/login request.
func (h *AccountHandler) GoogleGenerateAuthLink(c *gin.Context) {
	models.NewResponseOK(c, models.Response{
		Message: h.accountsUsecase.GetAuthURL(),
	})
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	_, account, err := getContextDependencies(c)
	if err != nil {
		return
	}

	models.NewResponseOK(c, models.Response{
		Data: account,
	})
}

func (h *AccountHandler) SaveLogin(c *gin.Context) {
	l, ok := c.Get("logger")
	if !ok {
		models.NewResponseInternalServerError(c)
		return
	}

	log := l.(*loggerModels.Logger)

	code := c.Query("code")
	if code == "" {
		models.NewResponseInvalidRequest(c, models.Response{Message: "code is required"})
		return
	}

	account, err := h.accountsUsecase.SaveGoogleAccount(c.Request.Context(), code)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "save_token_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)

		return
	}

	err = h.accountsUsecase.UpdateAccount(c.Request.Context(), account)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "update_account_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)

		return
	}

	token, refreshToken, redirectURL, err := h.accountsUsecase.GenerateTokens(c.Request.Context(), account)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "create_token_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)

		return
	}

	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/", "localhost", false, true) // 7 días
	c.Redirect(302, redirectURL)
}

func (h *AccountHandler) Refresh(c *gin.Context) {
	l, ok := c.Get("logger")
	if !ok {
		models.NewResponseInternalServerError(c)
		return
	}

	log := l.(*loggerModels.Logger)

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "missing_refresh_token",
			Error: err,
		})

		models.NewResponseUnauthorized(c, models.Response{
			Message: "missing refresh token",
		})

		return
	}

	token, err := h.accountsUsecase.VerifyToken(refreshToken)
	if err != nil {
		models.NewResponseUnauthorized(c, models.Response{
			Message: "invalid refresh token",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		models.NewResponseUnauthorized(c, models.Response{
			Message: "invalid claims",
		})
		return
	}

	id, ok := claims["id"].(string)
	if !ok || id == "" {
		models.NewResponseUnauthorized(c, models.Response{
			Message: "invalid refresh token payload",
		})
		return
	}

	account, err := h.accountsUsecase.GetAccount(c.Request.Context(), id)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "get_account_failed",
			Error: err,
			AdditionalParams: []loggerModels.Properties{
				logger.MapToProperties(map[string]string{
					"account_id": id,
				}),
			},
		})

		models.NewResponseInternalServerError(c)

		return
	}

	err = h.accountsUsecase.RefreshGoogleToken(c.Request.Context(), account)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "refresh_google_token_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)

		return
	}

	newToken, newRefreshToken, _, err := h.accountsUsecase.GenerateTokens(c.Request.Context(), account)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "create_token_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)

		return
	}

	c.SetCookie("token", newToken, 3600, "/", "localhost", false, true)
	c.SetCookie("refresh_token", newRefreshToken, 604800, "/", "localhost", false, true)

	models.NewResponseOK(c, models.Response{
		Message: "tokens refreshed successfully",
	})
}

func (h *AccountHandler) CreateWatcher(c *gin.Context) {
	log, account, err := getContextDependencies(c)
	if err != nil {
		return
	}

	err = h.accountsUsecase.CreateWatcher(c.Request.Context(), account)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "create_watcher_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)

		return
	}

	models.NewResponseOK(c, models.Response{
		Message: "watcher created successfully",
	})
}

func (h *AccountHandler) DeleteWatcher(c *gin.Context) {
	log, account, err := getContextDependencies(c)
	if err != nil {
		return
	}

	err = h.accountsUsecase.DeleteWatcher(c.Request.Context(), account)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "delete_watcher_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)

		return
	}

	models.NewResponseOK(c, models.Response{
		Message: "watcher deleted successfully",
	})
}
