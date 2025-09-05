package services

import (
	"fmt"
	"strconv"
	"transaction-tracker/api/models"
	"transaction-tracker/googleapi"
	"transaction-tracker/logger"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

func GetMovementsByMonth(gClient *googleapi.GoogleClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		log, err := logger.GetLogger(c, "transaction-tracker")
		if err != nil {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: fmt.Sprintf("logger not init: %s", err.Error()),
			})

			return
		}

		email := c.PostForm("email")
		if email == "" {
			log.Info(loggerModels.LogProperties{
				Event: "missing_email",
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: "email is required in x-www-form-urlencoded body",
			})

			return
		}

		gClient.SetEmail(email)

		movementsService, err := NewMovementsService(c)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "init_movements_service_failed",
				Error: err,
			})

			return
		}

		year, err := strconv.Atoi(c.Param("year"))
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "invalid_year",
				Error: err,
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: "invalid year",
			})

			return
		}

		month, err := strconv.Atoi(c.Param("month"))
		if err != nil || month < 1 || month > 12 {
			log.Error(loggerModels.LogProperties{
				Event: "invalid_month",
				Error: err,
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: "invalid month",
			})

			return
		}

		movements, err := movementsService.GetMovementsByMonth(c, year, month)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "get_movements_by_month_failed",
				Error: err,
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: "failed to get movements by month",
			})

			return
		}

		models.NewResponseOK(c, models.Response{
			Data: movements,
		})
	}
}
