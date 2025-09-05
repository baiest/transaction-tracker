package models

import (
	"fmt"
	"transaction-tracker/googleapi"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Port   string
	engine *gin.Engine
}

func NewServer(port int) *Server {
	return &Server{
		Port:   fmt.Sprintf(":%d", port),
		engine: gin.Default(),
	}
}

func (s *Server) AddRoutes(routes []Route, gClient *googleapi.GoogleClient) {
	api := s.engine.Group("/api")

	for _, r := range routes {
		api.Group(r.ApiVersion).Handle(string(r.Method), r.Endpoint, r.HandlerFunc(gClient))
	}
}

func (s *Server) Run() {
	s.engine.Run(s.Port)
}
