package routes

import (
	"transaction-tracker/api/models"
	gmailServices "transaction-tracker/api/services/gmail"
)

var (
	gmailRoutes = []models.Route{
		{
			Endpoint:    "/gmail/emails/histories/:historyID",
			Method:      models.GET,
			HandlerFunc: gmailServices.GetEmailByHistoryID(),
			ApiVersion:  API_VERSION,
		},
		{
			Endpoint:    "/gmail/emails/histories/:historyID/save",
			Method:      models.POST,
			HandlerFunc: gmailServices.StoreEmailByFilters(),
			ApiVersion:  API_VERSION,
		},
		{
			Endpoint:    "/gmail/emails/extracts",
			Method:      models.POST,
			HandlerFunc: gmailServices.StoreBankExtracts(),
			ApiVersion:  API_VERSION,
		},
	}
)
