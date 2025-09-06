package movements

import (
	"strconv"
	"time"
	"transaction-tracker/api/models"
	"transaction-tracker/database/mongo/schemas"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

func CreateMovement() gin.HandlerFunc {
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

		account := c.MustGet("account").(*models.Account)

		value, err := strconv.ParseFloat(c.PostForm("value"), 64)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "invalid_value",
				Error: err,
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: "invalid value",
			})

			return
		}

		movementType := c.PostForm("type")

		isNegative := movementType == "expense"

		movementTopic := c.PostForm("topic")
		if movementTopic == "" {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: "invalid topic",
			})

			return
		}

		movement := schemas.NewMovement(
			account.Email,
			"",
			time.Now(),
			value,
			isNegative,
			movementTopic,
			c.PostForm("detail"),
		)

		err = movementsService.CreateMovement(c, movement)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "create_movement_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		models.NewResponseOK(c, models.Response{
			Data: movement,
		})
	}
}
