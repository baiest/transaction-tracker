package routes

import (
	"transaction-tracker/api/models"
	services "transaction-tracker/api/services/movements"
)

var (
	movementsRoutes = []models.Route{
		{
			Endpoint:    "/movements",
			Method:      models.GET,
			HandlerFunc: services.GetMovements(),
			ApiVersion:  API_VERSION,
		},
		{
			Endpoint:    "/movements/years/:year",
			Method:      models.GET,
			HandlerFunc: services.GetMovementsByYear(),
			ApiVersion:  API_VERSION,
		},
		{
			Endpoint:    "/movements/years/:year/months/:month",
			Method:      models.GET,
			HandlerFunc: services.GetMovementsByMonth(),
			ApiVersion:  API_VERSION,
		},
	}
)
