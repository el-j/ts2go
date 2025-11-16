package background

import (
	"context"
	"time"

	"github.com/el-j/ts2go/saas/backend/logger"
)

// Job represents a background job
type Job func(ctx context.Context) error

// Scheduler manages background jobs
type Scheduler struct {
	jobs     map[string]*ScheduledJob
	stopChan chan struct{}
}

// ScheduledJob represents a job with schedule
type ScheduledJob struct {
	Name     string
	Job      Job
	Interval time.Duration
	LastRun  time.Time
}

// NewScheduler creates a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		jobs:     make(map[string]*ScheduledJob),
		stopChan: make(chan struct{}),
	}
}

// AddJob adds a job to the scheduler
func (s *Scheduler) AddJob(name string, job Job, interval time.Duration) {
	s.jobs[name] = &ScheduledJob{
		Name:     name,
		Job:      job,
		Interval: interval,
	}
	logger.Log.Info().
		Str("job", name).
		Dur("interval", interval).
		Msg("Background job scheduled")
}

// Start starts the scheduler
func (s *Scheduler) Start(ctx context.Context) {
	logger.Log.Info().
		Int("jobs", len(s.jobs)).
		Msg("Starting background job scheduler")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info().Msg("Scheduler stopping")
			return

		case <-s.stopChan:
			logger.Log.Info().Msg("Scheduler stopped")
			return

		case now := <-ticker.C:
			for _, job := range s.jobs {
				if now.Sub(job.LastRun) >= job.Interval {
					go s.runJob(ctx, job, now)
				}
			}
		}
	}
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	close(s.stopChan)
}

func (s *Scheduler) runJob(ctx context.Context, job *ScheduledJob, now time.Time) {
	logger.Log.Debug().
		Str("job", job.Name).
		Msg("Running background job")

	start := time.Now()
	err := job.Job(ctx)
	duration := time.Since(start)

	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("job", job.Name).
			Dur("duration", duration).
			Msg("Background job failed")
	} else {
		logger.Log.Info().
			Str("job", job.Name).
			Dur("duration", duration).
			Msg("Background job completed")
	}

	job.LastRun = now
}
