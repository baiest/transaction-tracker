package google

import (
	"transaction-tracker/api/models"
	"transaction-tracker/api/services/accounts"
	"transaction-tracker/googleapi"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := c.MustGet("logger").(*loggerModels.Logger)

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			models.NewResponseUnauthorized(c, models.Response{
				Message: "missing refresh token",
			})
			return
		}

		token, err := accounts.VerifyToken(refreshToken)
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

		accountsService, err := accounts.NewAccountService(c)
		if err != nil {
			models.NewResponseInternalServerError(c)
			return
		}

		account, err := accountsService.GetAccount(c, id)
		if err != nil {
			models.NewResponseUnauthorized(c, models.Response{
				Message: "account not found",
			})
			return
		}

		newAccess, newRefresh, err := account.GenerateTokens()
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "generate_tokens_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		c.SetCookie("token", newAccess, 900, "/", "localhost", false, true)             // 15 min
		c.SetCookie("refresh_token", newRefresh, 604800, "/", "localhost", false, true) // 7 días

		googleService, err := googleapi.NewGoogleClient(c)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "google_init_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		googleService.SetEmail(account.Email)

		_, err = googleService.RefreshToken(c)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "refresh_google_token_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		models.NewResponseOK(c, models.Response{
			Data: map[string]string{
				"token":         newAccess,
				"refresh_token": newRefresh,
			},
		})
	}
}
