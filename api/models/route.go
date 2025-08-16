package models

import (
	"transaction-tracker/googleapi"

	"github.com/gin-gonic/gin"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	DELETE Method = "DELETE"
)

type Route struct {
	Method      Method
	Endppoint   string
	HandlerFunc func(*googleapi.GoogleClient) gin.HandlerFunc
	ApiVersion  string
}
