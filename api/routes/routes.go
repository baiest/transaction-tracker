package routes

import (
	"transaction-tracker/api/models"
)

const (
	API_VERSION = "v1"
)

func Routes() []models.Route {
	routes := []models.Route{}

	routes = append(routes, googleRoutes...)
	routes = append(routes, gmailRoutes...)

	return routes
}
