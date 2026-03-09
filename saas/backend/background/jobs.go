package background

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/el-j/ts2go/saas/backend/logger"
	"github.com/el-j/ts2go/saas/backend/redis"
	"github.com/el-j/ts2go/saas/backend/storage"
)

// CleanupOldJobsJob removes jobs older than 30 days
func CleanupOldJobsJob(db *sql.DB, redisClient *redis.Client) Job {
	return func(ctx context.Context) error {
		cutoff := time.Now().AddDate(0, 0, -30)

		// Delete old jobs from database and get their IDs
		query := `DELETE FROM transpilations WHERE completed_at < $1 RETURNING id`
		rows, err := db.QueryContext(ctx, query, cutoff)
		if err != nil {
			return fmt.Errorf("failed to delete old jobs: %w", err)
		}
		defer rows.Close()

		var deletedJobs []string
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err == nil {
				deletedJobs = append(deletedJobs, id)
			}
		}

		logger.Log.Info().
			Int("deleted", len(deletedJobs)).
			Msg("Cleaned up old jobs from database")

		// Clean up Redis job data as well
		if len(deletedJobs) > 0 {
			redisKeys := make([]string, len(deletedJobs))
			for i, id := range deletedJobs {
				redisKeys[i] = "job:" + id
			}
			err = redisClient.Del(ctx, redisKeys...).Err()
			if err != nil {
				logger.Log.Warn().Err(err).Msg("Failed to clean up Redis job data")
			} else {
				logger.Log.Info().
					Int("deleted_keys", len(redisKeys)).
					Msg("Cleaned up old job data from Redis")
			}
		}

		return nil
	}
}

// CleanupOldFilesJob removes files older than 90 days
func CleanupOldFilesJob(db *sql.DB, storageClient *storage.Client) Job {
	return func(ctx context.Context) error {
		cutoff := time.Now().AddDate(0, 0, -90)

		// Get old files
		query := `SELECT id, storage_path FROM files WHERE created_at < $1 LIMIT 1000`
		rows, err := db.QueryContext(ctx, query, cutoff)
		if err != nil {
			return fmt.Errorf("failed to query old files: %w", err)
		}
		defer rows.Close()

		deleted := 0
		for rows.Next() {
			var id, storagePath string
			if err := rows.Scan(&id, &storagePath); err != nil {
				continue
			}

			// Delete from storage
			if err := storageClient.DeleteFile(ctx, storagePath); err != nil {
				logger.Log.Warn().
					Err(err).
					Str("file_id", id).
					Msg("Failed to delete file from storage")
				continue
			}

			// Delete from database
			_, err := db.ExecContext(ctx, `DELETE FROM files WHERE id = $1`, id)
			if err != nil {
				logger.Log.Warn().
					Err(err).
					Str("file_id", id).
					Msg("Failed to delete file record")
				continue
			}

			deleted++
		}

		logger.Log.Info().
			Int("deleted", deleted).
			Msg("Cleaned up old files")

		return nil
	}
}

// UpdateQueueMetricsJob updates queue depth metrics
func UpdateQueueMetricsJob(redisClient *redis.Client) Job {
	return func(ctx context.Context) error {
		priorities := []string{"low", "normal", "high", "urgent"}

		for _, priority := range priorities {
			queueKey := fmt.Sprintf("queue:transpilation:%s", priority)
			count, err := redisClient.ZCard(ctx, queueKey).Result()
			if err != nil {
				continue
			}

			// Update metrics
			// metrics.UpdateQueueDepth(priority, int(count))

			logger.Log.Debug().
				Str("priority", priority).
				Int64("depth", count).
				Msg("Queue depth")
		}

		return nil
	}
}

// GenerateUsageReportsJob generates daily usage reports
func GenerateUsageReportsJob(db *sql.DB) Job {
	return func(ctx context.Context) error {
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

		// Aggregate usage by user
		query := `
			SELECT 
				user_id,
				COUNT(*) as jobs_count,
				SUM(COALESCE(input_size_bytes, 0)) as total_input_bytes,
				SUM(COALESCE(output_size_bytes, 0)) as total_output_bytes,
				AVG(COALESCE(processing_time_ms, 0)) as avg_processing_time
			FROM transpilations
			WHERE DATE(created_at) = $1
			GROUP BY user_id
		`

		rows, err := db.QueryContext(ctx, query, yesterday)
		if err != nil {
			return fmt.Errorf("failed to generate usage report: %w", err)
		}
		defer rows.Close()

		userCount := 0
		for rows.Next() {
			var userID string
			var jobsCount, totalInputBytes, totalOutputBytes int64
			var avgProcessingTime float64

			if err := rows.Scan(&userID, &jobsCount, &totalInputBytes, &totalOutputBytes, &avgProcessingTime); err != nil {
				continue
			}

			// Store usage record
			insertQuery := `
				INSERT INTO usage_records (id, user_id, resource_type, quantity, metadata, created_at)
				VALUES (gen_random_uuid(), $1, 'daily_summary', $2, $3, NOW())
			`

			metadata := fmt.Sprintf(`{
				"date": "%s",
				"jobs": %d,
				"input_bytes": %d,
				"output_bytes": %d,
				"avg_processing_ms": %.2f
			}`, yesterday, jobsCount, totalInputBytes, totalOutputBytes, avgProcessingTime)

			_, err := db.ExecContext(ctx, insertQuery, userID, jobsCount, metadata)
			if err != nil {
				logger.Log.Warn().
					Err(err).
					Str("user_id", userID).
					Msg("Failed to store usage record")
			}

			userCount++
		}

		logger.Log.Info().
			Str("date", yesterday).
			Int("users", userCount).
			Msg("Generated usage reports")

		return nil
	}
}

// CleanupExpiredSessionsJob removes expired Redis sessions
func CleanupExpiredSessionsJob(redisClient *redis.Client) Job {
	return func(ctx context.Context) error {
		// Redis handles TTL automatically, but we can clean up orphaned keys
		logger.Log.Debug().Msg("Session cleanup (Redis handles TTL)")
		return nil
	}
}
