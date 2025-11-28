package routes

import (
	"transaction-tracker/api/handler"
	"transaction-tracker/api/models"
)

func AccountRoutes(h *handler.AccountHandler) []models.Route {
	return []models.Route{
		{
			Endpoint:       "/accounts/login",
			Method:         models.POST,
			HandlerFunc:    h.GoogleGenerateAuthLink,
			ApiVersion:     API_VERSION,
			NoRequiresAuth: true,
		},
		{
			Endpoint:       "/accounts/save",
			Method:         models.GET,
			HandlerFunc:    h.SaveLogin,
			ApiVersion:     API_VERSION,
			NoRequiresAuth: true,
		},
		{
			Endpoint:       "/accounts/refresh",
			Method:         models.POST,
			HandlerFunc:    h.Refresh,
			ApiVersion:     API_VERSION,
			NoRequiresAuth: true,
		},
		{
			Endpoint:    "/accounts/watchers",
			Method:      models.POST,
			HandlerFunc: h.CreateWatcher,
			ApiVersion:  API_VERSION,
		},
		{
			Endpoint:    "/accounts/watchers",
			Method:      models.DELETE,
			HandlerFunc: h.DeleteWatcher,
			ApiVersion:  API_VERSION,
		},
		{
			Endpoint:    "/accounts/me",
			Method:      models.GET,
			HandlerFunc: h.GetAccount,
			ApiVersion:  API_VERSION,
		},
	}
}
