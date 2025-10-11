package routes

import (
	"transaction-tracker/api/handler"
	"transaction-tracker/api/models"
)

func ExtractsRoutes(h *handler.ExtractsHandler) []models.Route {
	return []models.Route{
		{
			Endpoint:    "/extracts",
			Method:      models.POST,
			HandlerFunc: h.GetAllExtracts,
			ApiVersion:  API_VERSION,
		},
	}
}
