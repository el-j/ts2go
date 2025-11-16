package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/el-j/ts2go/saas/backend/redis"
	"github.com/google/uuid"
)

// JobPriority defines job priority levels
type JobPriority int

const (
	PriorityLow JobPriority = iota
	PriorityNormal
	PriorityHigh
	PriorityUrgent
)

// TranspilationJob represents a transpilation job
type TranspilationJob struct {
	ID           uuid.UUID              `json:"id"`
	UserID       uuid.UUID              `json:"user_id"`
	ProjectID    *uuid.UUID             `json:"project_id,omitempty"`
	Status       string                 `json:"status"` // pending, processing, completed, failed
	InputFiles   []FileMetadata         `json:"input_files"`
	OutputFiles  []FileMetadata         `json:"output_files,omitempty"`
	Settings     map[string]interface{} `json:"settings"`
	Priority     JobPriority            `json:"priority"`
	WorkerID     string                 `json:"worker_id,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	StartedAt    *time.Time             `json:"started_at,omitempty"`
	CompletedAt  *time.Time             `json:"completed_at,omitempty"`
}

// FileMetadata represents file metadata
type FileMetadata struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	SizeBytes   int64  `json:"size_bytes"`
	ContentType string `json:"content_type"`
}

// Queue manages transpilation job queue
type Queue struct {
	redis *redis.Client
}

// NewQueue creates a new job queue
func NewQueue(redisClient *redis.Client) *Queue {
	return &Queue{redis: redisClient}
}

// Enqueue adds a job to the queue
func (q *Queue) Enqueue(ctx context.Context, job *TranspilationJob) error {
	job.Status = "pending"
	job.CreatedAt = time.Now()

	jobJSON, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	// Store job data
	jobKey := fmt.Sprintf("job:%s", job.ID)
	if err := q.redis.SetWithExpiry(ctx, jobKey, string(jobJSON), 24*time.Hour); err != nil {
		return fmt.Errorf("failed to store job: %w", err)
	}

	// Add to priority queue
	queueKey := q.getQueueKey(job.Priority)
	score := float64(time.Now().Unix())

	// Use ZADD (sorted set) for priority queue
	if err := q.redis.ZAdd(ctx, queueKey, string(jobJSON), score).Err(); err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}

	// Add to user's job list
	userJobsKey := fmt.Sprintf("user:%s:jobs", job.UserID)
	if err := q.redis.SAdd(ctx, userJobsKey, job.ID.String()).Err(); err != nil {
		return fmt.Errorf("failed to add to user jobs: %w", err)
	}

	return nil
}

// Dequeue retrieves the next job from the queue
func (q *Queue) Dequeue(ctx context.Context, workerID string) (*TranspilationJob, error) {
	// Check queues in priority order
	priorities := []JobPriority{PriorityUrgent, PriorityHigh, PriorityNormal, PriorityLow}

	for _, priority := range priorities {
		queueKey := q.getQueueKey(priority)

		// Get first job from sorted set
		results, err := q.redis.ZRange(ctx, queueKey, 0, 0).Result()
		if err != nil || len(results) == 0 {
			continue
		}

		jobJSON := results[0]

		// Remove from queue
		if err := q.redis.ZRem(ctx, queueKey, jobJSON).Err(); err != nil {
			continue
		}

		// Parse job
		var job TranspilationJob
		if err := json.Unmarshal([]byte(jobJSON), &job); err != nil {
			continue
		}

		// Update job status
		job.Status = "processing"
		job.WorkerID = workerID
		now := time.Now()
		job.StartedAt = &now

		// Save updated job
		if err := q.UpdateJob(ctx, &job); err != nil {
			return nil, err
		}

		return &job, nil
	}

	return nil, fmt.Errorf("no jobs available")
}

// UpdateJob updates a job's status and data
func (q *Queue) UpdateJob(ctx context.Context, job *TranspilationJob) error {
	jobJSON, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	jobKey := fmt.Sprintf("job:%s", job.ID)
	return q.redis.SetWithExpiry(ctx, jobKey, string(jobJSON), 24*time.Hour)
}

// GetJob retrieves a job by ID
func (q *Queue) GetJob(ctx context.Context, jobID uuid.UUID) (*TranspilationJob, error) {
	jobKey := fmt.Sprintf("job:%s", jobID)
	jobJSON, err := q.redis.GetSession(ctx, jobKey)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	var job TranspilationJob
	if err := json.Unmarshal([]byte(jobJSON), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	return &job, nil
}

// ListUserJobs lists all jobs for a user
func (q *Queue) ListUserJobs(ctx context.Context, userID uuid.UUID, limit int) ([]*TranspilationJob, error) {
	userJobsKey := fmt.Sprintf("user:%s:jobs", userID)

	// Get job IDs
	jobIDs, err := q.redis.SMembers(ctx, userJobsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user jobs: %w", err)
	}

	var jobs []*TranspilationJob
	count := 0
	for _, jobIDStr := range jobIDs {
		if limit > 0 && count >= limit {
			break
		}

		jobID, err := uuid.Parse(jobIDStr)
		if err != nil {
			continue
		}

		job, err := q.GetJob(ctx, jobID)
		if err != nil {
			continue
		}

		jobs = append(jobs, job)
		count++
	}

	return jobs, nil
}

// CompleteJob marks a job as completed
func (q *Queue) CompleteJob(ctx context.Context, jobID uuid.UUID, outputFiles []FileMetadata) error {
	job, err := q.GetJob(ctx, jobID)
	if err != nil {
		return err
	}

	job.Status = "completed"
	job.OutputFiles = outputFiles
	now := time.Now()
	job.CompletedAt = &now

	return q.UpdateJob(ctx, job)
}

// FailJob marks a job as failed
func (q *Queue) FailJob(ctx context.Context, jobID uuid.UUID, errorMsg string) error {
	job, err := q.GetJob(ctx, jobID)
	if err != nil {
		return err
	}

	job.Status = "failed"
	job.ErrorMessage = errorMsg
	now := time.Now()
	job.CompletedAt = &now

	return q.UpdateJob(ctx, job)
}

// GetQueueStats returns queue statistics
func (q *Queue) GetQueueStats(ctx context.Context) (map[string]int, error) {
	stats := make(map[string]int)

	priorities := []JobPriority{PriorityLow, PriorityNormal, PriorityHigh, PriorityUrgent}
	for _, priority := range priorities {
		queueKey := q.getQueueKey(priority)
		count, err := q.redis.ZCard(ctx, queueKey).Result()
		if err != nil {
			continue
		}
		stats[priority.String()] = int(count)
	}

	return stats, nil
}

func (q *Queue) getQueueKey(priority JobPriority) string {
	return fmt.Sprintf("queue:transpilation:%s", priority.String())
}

func (p JobPriority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityNormal:
		return "normal"
	case PriorityHigh:
		return "high"
	case PriorityUrgent:
		return "urgent"
	default:
		return "normal"
	}
}
