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
	MovementHandler *handler.MovementHandler
}

func (r *RouteHandler) Routes() []models.Route {
	routes := []models.Route{}

	routes = append(routes, googleRoutes...)
	routes = append(routes, gmailRoutes...)
	routes = append(routes, MovementsRoutes(r.MovementHandler)...)

	return routes
}
