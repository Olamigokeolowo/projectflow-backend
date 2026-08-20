package main

import (
	"context"

	"github.com/Olamigokeolowo/projectflow-backend/internal/decision"
	"github.com/Olamigokeolowo/projectflow-backend/internal/events"
	"github.com/Olamigokeolowo/projectflow-backend/internal/middleware"
	"github.com/Olamigokeolowo/projectflow-backend/internal/user"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())

	ctx := context.Background()
	queue := events.NewInMemoryQueue(100)
	events.StartWorker(ctx, queue)

	decisionRepo := decision.NewInMemoryRepository()
	decisionService := decision.NewService(decisionRepo, queue)
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
		decisions.Use(middleware.AuthRequired())
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
