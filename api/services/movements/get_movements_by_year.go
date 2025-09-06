package movements

import (
	"strconv"
	"transaction-tracker/api/models"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

func GetMovementsByYear() gin.HandlerFunc {
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

		movements, err := movementsService.GetMovementsByYear(c, year)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "get_movements_by_year_failed",
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
