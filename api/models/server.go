package models

import (
	"fmt"
	"runtime/debug"
	"transaction-tracker/logger"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Port   string
	engine *gin.Engine
}

func NewServer(port int) *Server {
	engine := gin.Default()

	engine.Use(InitLogger())
	engine.Use(RecoveryWithJSON())

	return &Server{
		Port:   fmt.Sprintf(":%d", port),
		engine: engine,
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token := c.GetHeader("Authorization")
		// if token == "" {
		// 	NewResponseUnauthorized(c, Response{
		// 		Message: "missing authorization token",
		// 	})

		// 	return
		// }

		account, err := func() (*Account, error) {
			return &Account{Email: "juanballesteros2001@gmail.com"}, nil
		}()

		if err != nil {
			NewResponseUnauthorized(c, Response{
				Message: "invalid or expired token",
			})

			return
		}

		c.Set("account", account)

		c.Next()
	}
}

func InitLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log, err := logger.GetLogger(c, "transaction-tracker")
		if err != nil {
			NewResponseInternalServerError(c)
			return
		}

		c.Set("logger", log)
		c.Next()
	}
}

func RecoveryWithJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := c.MustGet("logger").(*loggerModels.Logger)

		defer func() {
			if r := recover(); r != nil {
				log.Panic(loggerModels.LogProperties{
					Event: "panic_recovered",
					Error: fmt.Errorf("%v\n%s", r, debug.Stack()),
				})

				NewResponseInternalServerError(c)
			}
		}()
		c.Next()
	}
}

func (s *Server) AddRoutes(routes []Route) {
	api := s.engine.Group("/api")

	for _, r := range routes {
		groupPublic := api.Group(r.ApiVersion)
		groupPrivate := api.Group(r.ApiVersion, AuthMiddleware())

		if r.NoRequiresAuth {
			groupPublic.Handle(string(r.Method), r.Endpoint, r.HandlerFunc)
		} else {
			groupPrivate.Handle(string(r.Method), r.Endpoint, r.HandlerFunc)
		}
	}
}

func (s *Server) Run() {
	s.engine.Run(s.Port)
}
