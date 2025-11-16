package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/el-j/ts2go/saas/backend/config"
	"github.com/el-j/ts2go/saas/backend/db"
	"github.com/el-j/ts2go/saas/backend/logger"
	"github.com/el-j/ts2go/saas/backend/queue"
	"github.com/el-j/ts2go/saas/backend/redis"
	"github.com/el-j/ts2go/saas/backend/storage"
	"github.com/el-j/ts2go/saas/backend/worker"
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
		ServiceName: "ts2go-worker",
	})

	logger.Log.Info().Msg("Starting ts2go worker service")

	// Connect to database
	database, err := db.Connect(
		cfg.Database.URL,
		cfg.Database.MaxOpenConns,
		cfg.Database.MaxIdleConns,
		cfg.Database.ConnMaxLifetime,
	)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	logger.Log.Info().Msg("✅ Connected to PostgreSQL database")

	// Connect to Redis
	redisClient, err := redis.Connect(cfg.Redis.URL, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer redisClient.Close()

	logger.Log.Info().Msg("✅ Connected to Redis")

	// Connect to storage
	storageClient, err := storage.NewClient(storage.Config{
		Endpoint:   getEnv("MINIO_ENDPOINT", "localhost:9000"),
		AccessKey:  getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		SecretKey:  getEnv("MINIO_SECRET_KEY", "minioadmin"),
		BucketName: getEnv("MINIO_BUCKET", "ts2go-files"),
		UseSSL:     getEnv("MINIO_USE_SSL", "false") == "true",
	})
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to storage")
	}

	logger.Log.Info().Msg("✅ Connected to MinIO storage")

	// Initialize services
	jobQueue := queue.NewQueue(redisClient)
	storageRepo := storage.NewRepository(database.DB)

	// Create worker
	w := worker.NewWorker(jobQueue, storageClient, storageRepo, worker.Config{
		MaxConcurrent: getEnvInt("WORKER_MAX_CONCURRENT", 5),
		WorkDir:       getEnv("WORKER_WORK_DIR", "/tmp/ts2go-worker"),
	})

	// Start worker in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		if err := w.Start(ctx); err != nil {
			errChan <- err
		}
	}()

	logger.Log.Info().Msg("🚀 Worker started and processing jobs")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		logger.Log.Info().Msg("Shutdown signal received")
	case err := <-errChan:
		logger.Log.Error().Err(err).Msg("Worker error")
	}

	// Graceful shutdown
	logger.Log.Info().Msg("🛑 Shutting down worker...")
	cancel()
	w.Stop()

	logger.Log.Info().Msg("Worker stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}
