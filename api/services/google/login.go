package google

import (
	"errors"
	"fmt"
	"transaction-tracker/api/models"
	"transaction-tracker/api/services/accounts"
	"transaction-tracker/googleapi"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

func GoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := c.MustGet("logger").(*loggerModels.Logger)

		gClient, err := googleapi.NewGoogleClient(c)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "init_google_client_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)
			return
		}

		err = gClient.SaveTokenAndInitServices(c, c.Query("code"))
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "save_token_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		projectID := "transaction-tracker-2473"
		topicName := fmt.Sprintf("projects/%s/topics/gmail-notifications", projectID)

		gmailService, err := gClient.GmailService(c)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "init_gmail_failed",
				Error: err,
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: err.Error(),
			})

			return
		}

		_, _, err = gmailService.CreateWatch(c, topicName)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "create_watch_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		accountsService, err := accounts.NewAccountService(c)
		if err != nil {
			models.NewResponseInternalServerError(c)

			return
		}

		email := gClient.Email()

		account, err := accountsService.GetAccountByEmail(c, email)
		if err != nil {
			if !errors.Is(err, accounts.ErrAccountNotFound) {
				log.Error(loggerModels.LogProperties{
					Event: "get_account_failed",
					Error: err,
				})

				models.NewResponseInternalServerError(c)

				return
			}

			account, err = accounts.NewAccount(email)
			if err != nil {
				log.Error(loggerModels.LogProperties{
					Event: "create_new_account_failed",
					Error: err,
				})

				models.NewResponseInternalServerError(c)

				return
			}

			err = accountsService.CreateAccount(c, account)
			if err != nil {
				models.NewResponseInternalServerError(c)

				return
			}
		}

		token, refreshToken, err := account.GenerateTokens()
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
		c.Redirect(302, "http://localhost:3000")
	}
}
