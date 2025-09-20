package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (r *Response) ContainsMessage() bool {
	return r.Message != ""
}

func (r *Response) DataOrMessage() any {
	if r.ContainsMessage() {
		return gin.H{"message": r.Message}
	}

	return r.Data
}

func NewResponseOK(c *gin.Context, response Response) {
	c.JSON(http.StatusOK, response.DataOrMessage())
}

func NewResponseCreated(c *gin.Context, response Response) {
	c.JSON(http.StatusCreated, response.DataOrMessage())
}

func NewResponseInternalServerError(c *gin.Context) {
	response := Response{Message: "something was wrong, please try again"}

	c.AbortWithStatusJSON(http.StatusInternalServerError, response.DataOrMessage())
}

func NewResponseNotFound(c *gin.Context, response Response) {
	c.JSON(http.StatusNotFound, response.DataOrMessage())
}

func NewResponseInvalidRequest(c *gin.Context, response Response) {
	c.JSON(http.StatusBadRequest, response.DataOrMessage())
}

func NewResponseUnauthorized(c *gin.Context, response Response) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, response.DataOrMessage())
}
