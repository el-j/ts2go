package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/el-j/ts2go/saas/backend/auth"
	"github.com/el-j/ts2go/saas/backend/background"
	"github.com/el-j/ts2go/saas/backend/cache"
	"github.com/el-j/ts2go/saas/backend/config"
	"github.com/el-j/ts2go/saas/backend/db"
	"github.com/el-j/ts2go/saas/backend/email"
	"github.com/el-j/ts2go/saas/backend/logger"
	"github.com/el-j/ts2go/saas/backend/middleware"
	"github.com/el-j/ts2go/saas/backend/projects"
	"github.com/el-j/ts2go/saas/backend/queue"
	"github.com/el-j/ts2go/saas/backend/redis"
	"github.com/el-j/ts2go/saas/backend/storage"
	"github.com/el-j/ts2go/saas/backend/transpilation"
	"github.com/el-j/ts2go/saas/backend/websocket"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger.InitLogger(logger.Config{
		Level:       "info",
		Pretty:      cfg.Server.Environment != "production",
		ServiceName: "ts2go-api",
	})

	// Set Gin mode based on environment
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	database, err := db.Connect(
		cfg.Database.URL,
		cfg.Database.MaxOpenConns,
		cfg.Database.MaxIdleConns,
		cfg.Database.ConnMaxLifetime,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("✅ Connected to PostgreSQL database")

	// Connect to Redis
	redisClient, err := redis.Connect(cfg.Redis.URL, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	log.Println("✅ Connected to Redis")

	// Connect to MinIO/S3
	storageClient, err := storage.NewClient(storage.Config{
		Endpoint:   "localhost:9000",
		AccessKey:  "minioadmin",
		SecretKey:  "minioadmin",
		BucketName: "ts2go-files",
		UseSSL:     false,
	})
	if err != nil {
		log.Fatalf("Failed to connect to storage: %v", err)
	}

	log.Println("✅ Connected to MinIO storage")

	// Initialize services
	authService := auth.NewService(cfg.JWT.Secret, cfg.JWT.Expiry, cfg.JWT.RefreshExp)
	authRepo := auth.NewRepository(database.DB)
	projectsRepo := projects.NewRepository(database.DB)
	storageRepo := storage.NewRepository(database.DB)
	jobQueue := queue.NewQueue(redisClient)

	// Initialize cache
	_ = cache.NewCache(redisClient) // Available for future use
	log.Println("✅ Cache service initialized")

	// Initialize email service
	_ = email.NewService(email.Config{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}) // Initialized, used by worker
	log.Println("✅ Email service initialized")

	// Initialize background job scheduler
	scheduler := background.NewScheduler()
	scheduler.AddJob("cleanup_old_jobs", background.CleanupOldJobsJob(database.DB, redisClient), 24*time.Hour)
	scheduler.AddJob("cleanup_old_files", background.CleanupOldFilesJob(database.DB, storageClient), 24*time.Hour)
	scheduler.AddJob("update_queue_metrics", background.UpdateQueueMetricsJob(redisClient), 1*time.Minute)
	scheduler.AddJob("generate_usage_reports", background.GenerateUsageReportsJob(database.DB), 24*time.Hour)
	scheduler.AddJob("cleanup_expired_sessions", background.CleanupExpiredSessionsJob(redisClient), 1*time.Hour)
	log.Println("✅ Background job scheduler initialized")

	// Initialize WebSocket hub
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wsHub := websocket.NewHub(jobQueue)
	go wsHub.Run(ctx)

	// Start background job scheduler
	go scheduler.Start(ctx)
	defer scheduler.Stop()

	// Initialize middleware
	authMiddleware := middleware.NewMiddleware(authService, authRepo)
	rateLimiter := middleware.NewRateLimiter(redisClient, middleware.DefaultRateLimitConfig())

	// Initialize handlers
	authHandler := NewAuthHandler(authService, authRepo)
	projectsHandler := projects.NewHandler(projectsRepo)
	storageHandler := storage.NewHandler(storageClient, storageRepo)
	transpilationHandler := transpilation.NewHandler(jobQueue)

	// Setup Gin router
	router := gin.Default()

	// Global middleware
	router.Use(corsMiddleware())
	router.Use(middleware.MetricsMiddleware()) // Prometheus metrics
	router.Use(middleware.ErrorHandler())      // Error handling with request ID
	router.Use(middleware.Recovery())          // Panic recovery
	router.Use(middleware.RequestLogger())     // Structured logging
	router.Use(rateLimiter.Limit())            // Rate limiting

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		status := gin.H{"status": "healthy"}

		// Check database
		if err := database.Health(); err != nil {
			status["database"] = "unhealthy: " + err.Error()
			c.JSON(500, status)
			return
		}
		status["database"] = "healthy"

		// Check Redis
		if err := redisClient.Health(c.Request.Context()); err != nil {
			status["redis"] = "unhealthy: " + err.Error()
			c.JSON(500, status)
			return
		}
		status["redis"] = "healthy"

		c.JSON(200, status)
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
		}

		// Protected routes (require JWT)
		authenticated := v1.Group("")
		authenticated.Use(authMiddleware.RequireAuth())
		{
			// Current user
			authenticated.GET("/auth/me", authHandler.GetCurrentUser)

			// API keys
			authenticated.POST("/api-keys", authHandler.CreateAPIKey)
			authenticated.GET("/api-keys", authHandler.ListAPIKeys)
			authenticated.DELETE("/api-keys/:id", authHandler.DeleteAPIKey)

			// Projects
			authenticated.POST("/projects", projectsHandler.Create)
			authenticated.GET("/projects", projectsHandler.List)
			authenticated.GET("/projects/:id", projectsHandler.Get)
			authenticated.PUT("/projects/:id", projectsHandler.Update)
			authenticated.DELETE("/projects/:id", projectsHandler.Delete)

			// File storage
			authenticated.POST("/files/upload", storageHandler.UploadFile)
			authenticated.GET("/files/:id/download", storageHandler.DownloadFile)
			authenticated.GET("/files/:id/url", storageHandler.GetPresignedURL)
			authenticated.DELETE("/files/:id", storageHandler.DeleteFile)

			// Transpilation
			authenticated.POST("/transpile", transpilationHandler.Transpile)
			authenticated.GET("/transpile", transpilationHandler.ListUserJobs)
			authenticated.GET("/transpile/:job_id", transpilationHandler.GetJobStatus)
			authenticated.POST("/transpile/:job_id/cancel", transpilationHandler.CancelJob)

			// WebSocket for real-time updates
			authenticated.GET("/ws", wsHub.ServeWS)
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("🚀 Server starting on %s (environment: %s)", addr, cfg.Server.Environment)
	log.Printf("📚 API docs will be available at http://localhost:%d/swagger", cfg.Server.Port)

	// Graceful shutdown
	go func() {
		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Server shutting down...")
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
