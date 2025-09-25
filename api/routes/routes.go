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
	MovementHandler *handler.MovementHandler
	MessageHandler  *handler.MessageHandler
}

func (r *RouteHandler) Routes() []models.Route {
	routes := []models.Route{}

	routes = append(routes, gmailRoutes...)
	routes = append(routes, AccountRoutes(r.AccountHandler)...)
	routes = append(routes, MovementsRoutes(r.MovementHandler)...)
	routes = append(routes, MessagesRoutes(r.MessageHandler)...)

	return routes
}
