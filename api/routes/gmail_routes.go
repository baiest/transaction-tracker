package routes

import (
	"transaction-tracker/api/models"
	services "transaction-tracker/api/services/google"
)

var (
	gmailRoutes = []models.Route{
		{
			Endpoint:    "/gmail/watchers",
			Method:      models.DELETE,
			HandlerFunc: services.GoogleDeleteWath,
			ApiVersion:  API_VERSION,
		},
		{
			Endpoint:    "/gmail/emails/histories/:historyID",
			Method:      models.GET,
			HandlerFunc: services.GetEmailByHistoryID,
			ApiVersion:  API_VERSION,
		},
		{
			Endpoint:    "/gmail/emails/histories/:historyID/save",
			Method:      models.POST,
			HandlerFunc: services.StoreEmailByFilters,
			ApiVersion:  API_VERSION,
		},
		{
			Endpoint:    "/gmail/emails/extracts",
			Method:      models.POST,
			HandlerFunc: services.StoreBankExtracts,
			ApiVersion:  API_VERSION,
		},
	}
)
