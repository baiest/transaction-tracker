package services

import (
	"strconv"
	"transaction-tracker/api/models"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

func GetMovements() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := c.MustGet("logger").(*loggerModels.Logger)

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

			models.NewResponseInternalServerError(c)

			return
		}

		models.NewResponseOK(c, models.Response{
			Data: &models.MovementsListResponse{TotalPages: totalPages, Page: page, Movements: movements},
		})
	}
}
