package models

import (
	"fmt"
	"runtime/debug"
	"transaction-tracker/api/services/accounts"
	"transaction-tracker/internal/accounts/usecase"
	"transaction-tracker/logger"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Server struct {
	Port           string
	accountUsecase usecase.AccountsUseCase
	engine         *gin.Engine
}

func NewServer(accountUsecase usecase.AccountsUseCase, port int) *Server {
	engine := gin.Default()

	engine.Use(InitLogger())
	engine.Use(RecoveryWithJSON())

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	return &Server{
		Port:           fmt.Sprintf(":%d", port),
		engine:         engine,
		accountUsecase: accountUsecase,
	}
}

func (s *Server) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {

			NewResponseUnauthorized(c, Response{
				Message: "missing token authorization",
			})

			return
		}

		token, err := accounts.VerifyToken(tokenString)
		if err != nil {
			NewResponseUnauthorized(c, Response{
				Message: "missing authorization",
			})

			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			NewResponseUnauthorized(c, Response{
				Message: "invalid token",
			})

			return
		}

		id, ok := claims["id"].(string)
		if id == "" || !ok {
			NewResponseUnauthorized(c, Response{
				Message: "email not found",
			})

			return
		}

		account, err := s.accountUsecase.GetAccount(c.Request.Context(), id)
		if err != nil {
			NewResponseInternalServerError(c)

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
		groupPrivate := api.Group(r.ApiVersion, s.AuthMiddleware())

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
