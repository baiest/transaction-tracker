package movements

import (
	"strconv"
	"transaction-tracker/api/models"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

func GetMovementsByMonth() gin.HandlerFunc {
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

			models.NewResponseInternalServerError(c)

			return
		}

		models.NewResponseOK(c, models.Response{
			Data: movements,
		})
	}
}
