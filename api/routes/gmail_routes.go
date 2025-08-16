package routes

import (
	"transaction-tracker/api/models"
	"transaction-tracker/api/services"
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
	}
)
