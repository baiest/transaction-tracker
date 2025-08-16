package models

import "github.com/gin-gonic/gin"

type Response struct {
	Message string `json:"message"`
}

func NewResponseOK(c *gin.Context, response Response) {
	c.JSON(200, response)
}

func NewResponseInternalServerError(c *gin.Context) {
	c.JSON(500, Response{Message: "something was wrong, please try again"})
}

func NewResponseInvalidRequest(c *gin.Context, response Response) {
	c.JSON(400, response)
}
