package routes

import (
	"transaction-tracker/api/models"
	services "transaction-tracker/api/services/google"
)

var (
	gmailRoutes = []models.Route{
		{
			Endppoint:   "/gmail/watchers",
			Method:      models.DELETE,
			HandlerFunc: services.GoogleDeleteWath,
			ApiVersion:  API_VERSION,
		},
		{
			Endppoint:   "/gmail/emails/:historyID",
			Method:      models.GET,
			HandlerFunc: services.GetEmailByHistoryID,
			ApiVersion:  API_VERSION,
		},
		{
			Endppoint:   "/gmail/emails/:historyID/save",
			Method:      models.POST,
			HandlerFunc: services.StoreEmailByFilters,
			ApiVersion:  API_VERSION,
		},
	}
)
