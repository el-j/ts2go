package worker

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/el-j/ts2go/saas/backend/email"
	"github.com/el-j/ts2go/saas/backend/logger"
	"github.com/el-j/ts2go/saas/backend/metrics"
	"github.com/el-j/ts2go/saas/backend/queue"
	"github.com/el-j/ts2go/saas/backend/storage"
	"github.com/google/uuid"
)

// Worker processes transpilation jobs from the queue
type Worker struct {
	id            string
	queue         *queue.Queue
	storage       *storage.Client
	storageRepo   *storage.Repository
	emailService  *email.Service
	db            *sql.DB
	maxConcurrent int
	workDir       string
	stopChan      chan struct{}
}

// Config holds worker configuration
type Config struct {
	MaxConcurrent int    // Maximum concurrent jobs
	WorkDir       string // Working directory for transpilation
	PollInterval  time.Duration
}

// NewWorker creates a new worker instance
func NewWorker(queue *queue.Queue, storage *storage.Client, storageRepo *storage.Repository, emailService *email.Service, db *sql.DB, cfg Config) *Worker {
	if cfg.MaxConcurrent == 0 {
		cfg.MaxConcurrent = 5
	}
	if cfg.WorkDir == "" {
		cfg.WorkDir = "/tmp/ts2go-worker"
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = 2 * time.Second
	}

	// Create work directory
	os.MkdirAll(cfg.WorkDir, 0755)

	return &Worker{
		id:            uuid.New().String(),
		queue:         queue,
		storage:       storage,
		storageRepo:   storageRepo,
		emailService:  emailService,
		db:            db,
		maxConcurrent: cfg.MaxConcurrent,
		workDir:       cfg.WorkDir,
		stopChan:      make(chan struct{}),
	}
}

// Start begins processing jobs
func (w *Worker) Start(ctx context.Context) error {
	logger.Log.Info().
		Str("worker_id", w.id).
		Int("max_concurrent", w.maxConcurrent).
		Str("work_dir", w.workDir).
		Msg("Worker starting")

	metrics.WorkerActive.Inc()
	defer metrics.WorkerActive.Dec()

	sem := make(chan struct{}, w.maxConcurrent)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info().Str("worker_id", w.id).Msg("Worker stopping (context cancelled)")
			return ctx.Err()

		case <-w.stopChan:
			logger.Log.Info().Str("worker_id", w.id).Msg("Worker stopping (stop signal)")
			return nil

		case <-ticker.C:
			// Try to dequeue a job
			select {
			case sem <- struct{}{}: // Acquire semaphore
				go func() {
					defer func() { <-sem }() // Release semaphore

					job, err := w.queue.Dequeue(ctx, w.id)
					if err != nil {
						if err.Error() != "no jobs available" {
							logger.Log.Error().Err(err).Msg("Failed to dequeue job")
						}
						return
					}

					if job == nil {
						return
					}

					logger.Log.Info().
						Str("worker_id", w.id).
						Str("job_id", job.ID.String()).
						Str("user_id", job.UserID.String()).
						Msg("Processing job")

					metrics.WorkerJobsProcessing.Inc()
					defer metrics.WorkerJobsProcessing.Dec()

					if err := w.processJob(ctx, job); err != nil {
						logger.Log.Error().
							Err(err).
							Str("job_id", job.ID.String()).
							Msg("Job processing failed")
					}
				}()
			default:
				// All workers busy, skip this tick
			}
		}
	}
}

// Stop stops the worker gracefully
func (w *Worker) Stop() {
	close(w.stopChan)
}

// processJob handles a single transpilation job
func (w *Worker) processJob(ctx context.Context, job *queue.TranspilationJob) error {
	startTime := time.Now()

	// Create job-specific directory
	jobDir := filepath.Join(w.workDir, job.ID.String())
	if err := os.MkdirAll(jobDir, 0755); err != nil {
		w.queue.FailJob(ctx, job.ID, fmt.Sprintf("Failed to create job directory: %v", err))
		return err
	}
	defer os.RemoveAll(jobDir)

	// Parse input files
	if len(job.InputFiles) == 0 {
		w.queue.FailJob(ctx, job.ID, "No input files provided")
		return fmt.Errorf("no input files")
	}

	// Download input files from storage
	inputDir := filepath.Join(jobDir, "input")
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		w.queue.FailJob(ctx, job.ID, fmt.Sprintf("Failed to create input directory: %v", err))
		return err
	}

	logger.Log.Debug().
		Str("job_id", job.ID.String()).
		Int("file_count", len(job.InputFiles)).
		Msg("Downloading input files")

	localInputFiles := []string{}
	for _, fileMeta := range job.InputFiles {
		fileUUID, err := uuid.Parse(fileMeta.Path)
		if err != nil {
			continue
		}

		// Get file metadata
		fileRecord, err := w.storageRepo.GetFile(ctx, fileUUID)
		if err != nil {
			logger.Log.Warn().Err(err).Str("file_id", fileMeta.Path).Msg("Failed to get file metadata")
			continue
		}

		// Download file
		reader, err := w.storage.DownloadFile(ctx, fileRecord.StoragePath)
		if err != nil {
			logger.Log.Warn().Err(err).Str("file_id", fileMeta.Path).Msg("Failed to download file")
			continue
		}

		// Save to local file
		localPath := filepath.Join(inputDir, fileRecord.OriginalName)
		localFile, err := os.Create(localPath)
		if err != nil {
			reader.Close()
			continue
		}

		if _, err := localFile.ReadFrom(reader); err != nil {
			reader.Close()
			localFile.Close()
			continue
		}

		reader.Close()
		localFile.Close()
		localInputFiles = append(localInputFiles, localPath)
	}

	if len(localInputFiles) == 0 {
		w.queue.FailJob(ctx, job.ID, "Failed to download any input files")
		return fmt.Errorf("no input files downloaded")
	}

	// Create output directory
	outputDir := filepath.Join(jobDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		w.queue.FailJob(ctx, job.ID, fmt.Sprintf("Failed to create output directory: %v", err))
		return err
	}

	// Run transpilation
	logger.Log.Info().
		Str("job_id", job.ID.String()).
		Msg("Running transpilation")

	if err := w.runTranspilation(ctx, localInputFiles, outputDir, job.Settings); err != nil {
		errorMsg := fmt.Sprintf("Transpilation failed: %v", err)
		w.queue.FailJob(ctx, job.ID, errorMsg)
		metrics.RecordJobProcessed("failed", time.Since(startTime).Seconds())

		// Send failure email
		w.sendFailureEmail(ctx, job, errorMsg)

		return err
	}

	// Upload output files
	logger.Log.Debug().
		Str("job_id", job.ID.String()).
		Msg("Uploading output files")

	outputFiles, err := w.uploadOutputFiles(ctx, outputDir, *job.ProjectID)
	if err != nil {
		w.queue.FailJob(ctx, job.ID, fmt.Sprintf("Failed to upload output: %v", err))
		return err
	}

	// Complete job
	if err := w.queue.CompleteJob(ctx, job.ID, outputFiles); err != nil {
		logger.Log.Error().Err(err).Str("job_id", job.ID.String()).Msg("Failed to mark job complete")
		return err
	}

	processingTime := time.Since(startTime)
	metrics.RecordJobProcessed("completed", processingTime.Seconds())

	// Send success email
	w.sendSuccessEmail(ctx, job, outputFiles)

	logger.Log.Info().
		Str("job_id", job.ID.String()).
		Dur("processing_time", processingTime).
		Int("output_files", len(outputFiles)).
		Msg("Job completed successfully")

	return nil
}

// runTranspilation executes the ts2go transpiler
func (w *Worker) runTranspilation(ctx context.Context, inputFiles []string, outputDir string, settings map[string]interface{}) error {
	// Build command
	// Assuming ts2go CLI is in PATH or at a known location
	ts2goBin := os.Getenv("TS2GO_BIN")
	if ts2goBin == "" {
		ts2goBin = "/usr/local/bin/ts2go" // Default location
	}

	args := []string{"transpile"}
	args = append(args, "--output", outputDir)

	// Add settings as flags
	if settings != nil {
		if optimizeGoFmt, ok := settings["optimize_go_fmt"].(bool); ok && optimizeGoFmt {
			args = append(args, "--optimize-go-fmt")
		}
		if preserveComments, ok := settings["preserve_comments"].(bool); ok && preserveComments {
			args = append(args, "--preserve-comments")
		}
	}

	// Add input files
	args = append(args, inputFiles...)

	// Create command
	cmd := exec.CommandContext(ctx, ts2goBin, args...)
	cmd.Dir = filepath.Dir(inputFiles[0])

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("transpilation command failed: %w, output: %s", err, string(output))
	}

	logger.Log.Debug().
		Str("output", string(output)).
		Msg("Transpilation output")

	return nil
}

// uploadOutputFiles uploads all files from output directory to storage
func (w *Worker) uploadOutputFiles(ctx context.Context, outputDir string, projectID uuid.UUID) ([]queue.FileMetadata, error) {
	var outputFiles []queue.FileMetadata

	// Walk output directory
	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Open file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Generate storage path
		relPath, _ := filepath.Rel(outputDir, path)
		storagePath := fmt.Sprintf("projects/%s/output/%s/%s", projectID, uuid.New(), relPath)

		// Upload to storage
		_, err = w.storage.UploadFile(ctx, storagePath, file, info.Size(), "text/plain")
		if err != nil {
			return err
		}

		// Save metadata to database
		fileRecord, err := w.storageRepo.CreateFile(
			ctx,
			&projectID,
			nil,
			"output",
			filepath.Base(path),
			storagePath,
			"text/plain",
			info.Size(),
			"",
		)
		if err != nil {
			return err
		}

		outputFiles = append(outputFiles, queue.FileMetadata{
			Name:        fileRecord.OriginalName,
			Path:        fileRecord.ID.String(),
			SizeBytes:   info.Size(),
			ContentType: "text/plain",
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return outputFiles, nil
}

// sendSuccessEmail sends job completion email
func (w *Worker) sendSuccessEmail(ctx context.Context, job *queue.TranspilationJob, outputFiles []queue.FileMetadata) {
	if w.emailService == nil {
		return
	}

	// Get user email
	var email string
	err := w.db.QueryRowContext(ctx, `SELECT email FROM users WHERE id = $1`, job.UserID).Scan(&email)
	if err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to get user email")
		return
	}

	// Get project name
	var projectName string
	if job.ProjectID != nil {
		w.db.QueryRowContext(ctx, `SELECT name FROM projects WHERE id = $1`, job.ProjectID).Scan(&projectName)
	}
	if projectName == "" {
		projectName = "Untitled Project"
	}

	// Send email (non-blocking)
	go func() {
		err := w.emailService.SendJobComplete(email, job.ID.String(), projectName, len(outputFiles))
		if err != nil {
			logger.Log.Warn().Err(err).Msg("Failed to send completion email")
		}
	}()
}

// sendFailureEmail sends job failure email
func (w *Worker) sendFailureEmail(ctx context.Context, job *queue.TranspilationJob, errorMsg string) {
	if w.emailService == nil {
		return
	}

	// Get user email
	var email string
	err := w.db.QueryRowContext(ctx, `SELECT email FROM users WHERE id = $1`, job.UserID).Scan(&email)
	if err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to get user email")
		return
	}

	// Get project name
	var projectName string
	if job.ProjectID != nil {
		w.db.QueryRowContext(ctx, `SELECT name FROM projects WHERE id = $1`, job.ProjectID).Scan(&projectName)
	}
	if projectName == "" {
		projectName = "Untitled Project"
	}

	// Send email (non-blocking)
	go func() {
		err := w.emailService.SendJobFailed(email, job.ID.String(), projectName, errorMsg)
		if err != nil {
			logger.Log.Warn().Err(err).Msg("Failed to send failure email")
		}
	}()
}
