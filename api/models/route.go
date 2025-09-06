package models

import (
	"github.com/gin-gonic/gin"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	DELETE Method = "DELETE"
)

type Route struct {
	Method         Method
	Endpoint       string
	HandlerFunc    gin.HandlerFunc
	ApiVersion     string
	NoRequiresAuth bool
}
