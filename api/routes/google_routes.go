package routes

import (
	"transaction-tracker/api/models"
	"transaction-tracker/api/services"
)

var (
	googleRoutes = []models.Route{
		{Endppoint: "/google/auth/generate",
			Method:      models.POST,
			HandlerFunc: services.GoogleGenerateAuthLink,
			ApiVersion:  API_VERSION,
		},
		{
			Endppoint:   "/google/auth/save",
			Method:      models.GET,
			HandlerFunc: services.GoogleLogin,
			ApiVersion:  API_VERSION,
		},
	}
)
