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
	Endpoint    string
	HandlerFunc func(*googleapi.GoogleClient) gin.HandlerFunc
	ApiVersion  string
}
