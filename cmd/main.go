package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Olamigokeolowo/projectflow-backend/internal/decision"
	"github.com/Olamigokeolowo/projectflow-backend/internal/middleware"
	"github.com/Olamigokeolowo/projectflow-backend/internal/user"
)

func main() {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())

	decisionRepo := decision.NewInMemoryRepository()
	decisionService := decision.NewService(decisionRepo)
	decisionHandler := decision.NewHandler(decisionService)

	userRepo := user.NewInMemoryRepository()
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		decisions := v1.Group("/decisions")
		decisions.Use(middleware.AuthRequired()) // everything below this line requires a valid token
		{
			decisions.GET("", decisionHandler.List)
			decisions.GET("/:id", decisionHandler.Get)
			decisions.POST("", decisionHandler.Create)
			decisions.GET("/:id/tasks", decisionHandler.ListTasks)
			decisions.GET("/slow", decisionHandler.SlowOperation)
		}
	}

	r.Run(":8080")
}