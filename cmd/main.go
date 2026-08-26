package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/Olamigokeolowo/projectflow-backend/internal/cache"
	"github.com/Olamigokeolowo/projectflow-backend/internal/decision"
	"github.com/Olamigokeolowo/projectflow-backend/internal/events"
	"github.com/Olamigokeolowo/projectflow-backend/internal/metrics"
	"github.com/Olamigokeolowo/projectflow-backend/internal/middleware"
	"github.com/Olamigokeolowo/projectflow-backend/internal/user"
)

func main() {
	r := gin.New()

	metricsCollector := metrics.New()

	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(metricsCollector))

	ctx := context.Background()
	queue := events.NewInMemoryQueue(100)
	events.StartWorker(ctx, queue)

	decisionCache := cache.NewInMemoryCache()
	decisionRepo := decision.NewInMemoryRepository()
	decisionService := decision.NewService(decisionRepo, queue, decisionCache)
	decisionHandler := decision.NewHandler(decisionService)

	userRepo := user.NewInMemoryRepository()
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(200, metricsCollector.Snapshot())
	})

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