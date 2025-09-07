package routes

import (
	"transaction-tracker/api/models"
	services "transaction-tracker/api/services/google"
)

var (
	googleRoutes = []models.Route{
		{
			Endpoint:       "/google/auth/generate",
			Method:         models.POST,
			HandlerFunc:    services.GoogleGenerateAuthLink(),
			ApiVersion:     API_VERSION,
			NoRequiresAuth: true,
		},
		{
			Endpoint:       "/google/auth/save",
			Method:         models.GET,
			HandlerFunc:    services.GoogleLogin(),
			ApiVersion:     API_VERSION,
			NoRequiresAuth: true,
		},
		{
			Endpoint:    "/google/auth/refresh",
			Method:      models.POST,
			HandlerFunc: services.Refresh(),
			ApiVersion:  API_VERSION,
		},
	}
)
