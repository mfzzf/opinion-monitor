package main

import (
	"log"
	"opinion-monitor/internal/api"
	"opinion-monitor/internal/config"
	"opinion-monitor/internal/models"
	"opinion-monitor/internal/worker"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := models.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations
	if err := models.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize job queue
	jobQueue := worker.NewJobQueue()

	// Start worker pool
	workerPool := worker.NewWorkerPool(cfg, db, jobQueue)
	workerPool.Start()

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Serve static files (uploads)
	r.Static("/uploads", cfg.Server.UploadPath)

	// Initialize handlers
	authHandler := api.NewAuthHandler(db, cfg)
	videoHandler := api.NewVideoHandler(db, cfg, jobQueue)
	reportHandler := api.NewReportHandler(db)
	jobHandler := api.NewJobHandler(db)

	// Auth routes
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.GET("/me", authHandler.Me)
	}

	// Protected routes
	apiGroup := r.Group("/api")
	apiGroup.Use(api.AuthMiddleware(cfg))
	{
		// Video routes
		apiGroup.POST("/videos/upload", videoHandler.Upload)
		apiGroup.GET("/videos", videoHandler.List)
		apiGroup.GET("/videos/:id", videoHandler.Get)
		apiGroup.DELETE("/videos/:id", videoHandler.Delete)

		// Report routes
		apiGroup.GET("/reports/:video_id", reportHandler.GetByVideoID)
		apiGroup.GET("/reports", reportHandler.List)

		// Job routes
		apiGroup.GET("/jobs/:id/status", jobHandler.GetStatus)
		apiGroup.GET("/jobs", jobHandler.List)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
