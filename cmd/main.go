package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Olamigokeolowo/projectflow-backend/internal/cache"
	"github.com/Olamigokeolowo/projectflow-backend/internal/decision"
	"github.com/Olamigokeolowo/projectflow-backend/internal/events"
	"github.com/Olamigokeolowo/projectflow-backend/internal/metrics"
	"github.com/Olamigokeolowo/projectflow-backend/internal/middleware"
	"github.com/Olamigokeolowo/projectflow-backend/internal/user"
	"github.com/Olamigokeolowo/projectflow-backend/internal/workspace"
)

func main() {
	r := gin.New()

	metricsCollector := metrics.New()

	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(metricsCollector))

	workerCtx, cancelWorker := context.WithCancel(context.Background())

	queue := events.NewInMemoryQueue(100)
	events.StartWorker(workerCtx, queue)

	userRepo := user.NewInMemoryRepository()
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	workspaceRepo := workspace.NewInMemoryRepository()
	workspaceService := workspace.NewService(workspaceRepo, userService)
	workspaceHandler := workspace.NewHandler(workspaceService)

	decisionCache := cache.NewInMemoryCache()
	decisionRepo := decision.NewInMemoryRepository()
	decisionService := decision.NewService(decisionRepo, queue, decisionCache, workspaceService)
	decisionHandler := decision.NewHandler(decisionService)

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

		workspaces := v1.Group("/workspaces")
		workspaces.Use(middleware.AuthRequired())
		{
			workspaces.POST("", workspaceHandler.Create)
			workspaces.GET("", workspaceHandler.List)
			workspaces.POST("/:id/members", workspaceHandler.AddMember)
			workspaces.GET("/:id/members", workspaceHandler.ListMembers)
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutdown signal received, starting graceful shutdown")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shut down: %v", err)
	}

	cancelWorker()

	log.Println("server exited cleanly")
}