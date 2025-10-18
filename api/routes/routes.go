package routes

import (
	"transaction-tracker/api/handler"
	"transaction-tracker/api/models"
)

const (
	API_VERSION = "v1"
)

// Routes holds all the application handlers.
type RouteHandler struct {
	AccountHandler  *handler.AccountHandler
	MessageHandler  *handler.MessageHandler
	ExtractHandler  *handler.ExtractsHandler
	MovementHandler *handler.MovementHandler
}

func (r *RouteHandler) Routes() []models.Route {
	routes := []models.Route{}

	routes = append(routes, AccountRoutes(r.AccountHandler)...)
	routes = append(routes, MessagesRoutes(r.MessageHandler)...)
	routes = append(routes, ExtractsRoutes(r.ExtractHandler)...)
	routes = append(routes, MovementsRoutes(r.MovementHandler)...)

	return routes
}
