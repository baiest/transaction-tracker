package routes

import (
	"transaction-tracker/api/handler"
	"transaction-tracker/api/models"
)

func MovementsRoutes(h *handler.MovementHandler) []models.Route {
	return []models.Route{
		{
			Endpoint:    "/movements",
			Method:      models.GET,
			HandlerFunc: h.GetMovements,
			ApiVersion:  API_VERSION,
		},
		{
			Endpoint:    "/movements/:id",
			Method:      models.GET,
			HandlerFunc: h.GetMovementByID,
			ApiVersion:  API_VERSION,
		},
		{
			Endpoint:    "/movements",
			Method:      models.POST,
			HandlerFunc: h.CreateMovement,
			ApiVersion:  API_VERSION,
		},
		{
			Endpoint:    "/movements/:id",
			Method:      models.DELETE,
			HandlerFunc: h.DeleteMovement,
			ApiVersion:  API_VERSION,
		},
		{
			Endpoint:    "/movements/years/:year",
			Method:      models.GET,
			HandlerFunc: h.GetMovementsByYear,
			ApiVersion:  API_VERSION,
		},
		{
			Endpoint:    "/movements/years/:year/months/:month",
			Method:      models.GET,
			HandlerFunc: h.GetMovementsByMonth,
			ApiVersion:  API_VERSION,
		},
	}
}
