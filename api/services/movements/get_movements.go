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

func GetMovements(gClient *googleapi.GoogleClient) gin.HandlerFunc {
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

		page, err := strconv.ParseInt(c.Query("page"), 10, 64)
		if err != nil || page == 0 {
			page = 1
		}

		movements, totalPages, err := movementsService.GetMovements(c, page)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "get_movements_failed",
				Error: err,
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: "failed to get movements",
			})

			return
		}

		models.NewResponseOK(c, models.Response{
			Message: "movements retrieved successfully",
			Data:    &models.MovementsListResponse{TotalPages: totalPages, Page: page, Movements: movements},
		})
	}
}
