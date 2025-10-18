package routes

import (
	"transaction-tracker/api/handler"
	"transaction-tracker/api/models"
)

func MessagesRoutes(h *handler.MessageHandler) []models.Route {
	return []models.Route{
		{
			Method:      models.POST,
			Endpoint:    "/messages",
			HandlerFunc: h.ProcessMessage,
			ApiVersion:  API_VERSION,
		},
		{
			Method:      models.GET,
			Endpoint:    "/messages/:id",
			HandlerFunc: h.GetMessageByID,
			ApiVersion:  API_VERSION,
		},
		{
			Method:      models.GET,
			Endpoint:    "/messages/history/:id",
			HandlerFunc: h.GetMessagesByHistory,
			ApiVersion:  API_VERSION,
		},
	}
}
