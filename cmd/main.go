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
	"github.com/Olamigokeolowo/projectflow-backend/internal/activity"
	"github.com/Olamigokeolowo/projectflow-backend/internal/cache"
	"github.com/Olamigokeolowo/projectflow-backend/internal/comment"
	"github.com/Olamigokeolowo/projectflow-backend/internal/database"
	"github.com/Olamigokeolowo/projectflow-backend/internal/decision"
	"github.com/Olamigokeolowo/projectflow-backend/internal/events"
	"github.com/Olamigokeolowo/projectflow-backend/internal/metrics"
	"github.com/Olamigokeolowo/projectflow-backend/internal/middleware"
	"github.com/Olamigokeolowo/projectflow-backend/internal/task"
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

	dbConnString := os.Getenv("DATABASE_URL")
	if dbConnString == "" {
		dbConnString = "postgres://projectflow:devpassword@localhost:5432/projectflow"
	}
	dbPool, err := database.Connect(context.Background(), dbConnString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	userRepo := user.NewPostgresRepository(dbPool)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	workspaceRepo := workspace.NewInMemoryRepository()

	activityRepo := activity.NewInMemoryRepository()
	activityService := activity.NewService(activityRepo, workspaceRepo)

	workspaceService := workspace.NewService(workspaceRepo, userService, activityService)
	workspaceHandler := workspace.NewHandler(workspaceService)

	decisionCache := cache.NewInMemoryCache()
	decisionRepo := decision.NewInMemoryRepository()
	decisionService := decision.NewService(decisionRepo, queue, decisionCache, workspaceService, activityService)
	decisionHandler := decision.NewHandler(decisionService)

	taskRepo := task.NewInMemoryRepository()
	taskService := task.NewService(taskRepo, decisionService, workspaceService, activityService)
	taskHandler := task.NewHandler(taskService)

	commentResolver := comment.NewResolver(decisionService, taskService)
	commentRepo := comment.NewInMemoryRepository()
	commentService := comment.NewService(commentRepo, commentResolver, workspaceService, activityService)
	commentHandler := comment.NewHandler(commentService)

	activityHandler := activity.NewHandler(activityService)

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
			workspaces.GET("/:id/activity", activityHandler.List)
		}

		decisions := v1.Group("/decisions")
		decisions.Use(middleware.AuthRequired())
		{
			decisions.GET("", decisionHandler.List)
			decisions.GET("/:id", decisionHandler.Get)
			decisions.POST("", decisionHandler.Create)
			decisions.PATCH("/:id", decisionHandler.Update)
			decisions.DELETE("/:id", decisionHandler.Delete)
			decisions.GET("/slow", decisionHandler.SlowOperation)
			decisions.GET("/:id/tasks", taskHandler.List)
			decisions.POST("/:id/tasks", taskHandler.Create)
			decisions.GET("/:id/comments", commentHandler.ListForDecision)
			decisions.POST("/:id/comments", commentHandler.CreateForDecision)
		}

		tasks := v1.Group("/tasks")
		tasks.Use(middleware.AuthRequired())
		{
			tasks.PATCH("/:id", taskHandler.Update)
			tasks.DELETE("/:id", taskHandler.Delete)
			tasks.GET("/:id/comments", commentHandler.ListForTask)
			tasks.POST("/:id/comments", commentHandler.CreateForTask)
		}

		comments := v1.Group("/comments")
		comments.Use(middleware.AuthRequired())
		{
			comments.DELETE("/:id", commentHandler.Delete)
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