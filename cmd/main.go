package main

import (
	"github.com/Olamigokeolowo/projectflow-demo/internal/decision"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	decisionHandler := decision.NewHandler()

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
