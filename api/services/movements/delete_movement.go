package services

import (
	"transaction-tracker/api/models"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

func DeleteMovement() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := c.MustGet("logger").(*loggerModels.Logger)

		movementsService, err := NewMovementsService(c)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "init_movements_service_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		movementID := c.Param("movementID")
		if movementID == "" {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: "invalid movement id",
			})

			return
		}

		err = movementsService.DeleteMovement(c, movementID)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "delete_movement_failed",
				Error: err,
			})
		}

		models.NewResponseOK(c, models.Response{
			Message: "movement deleted successfully",
		})
	}
}
