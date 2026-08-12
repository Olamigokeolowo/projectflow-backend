package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Olamigokeolowo/projectflow-backend/internal/decision"
)

func main() {
	r := gin.Default()

	repo := decision.NewInMemoryRepository()
	service := decision.NewService(repo)
	decisionHandler := decision.NewHandler(service)

	v1 := r.Group("/api/v1")
	{
		decisions := v1.Group("/decisions")
		decisions.GET("", decisionHandler.List)
		decisions.GET("/:id", decisionHandler.Get)
		decisions.POST("", decisionHandler.Create)
		decisions.GET("/:id/tasks", decisionHandler.ListTasks)
	}

	r.Run(":8080")
}