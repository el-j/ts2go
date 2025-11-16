package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ts2go_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ts2go_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Job metrics
	JobsEnqueued = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ts2go_jobs_enqueued_total",
			Help: "Total number of jobs enqueued",
		},
		[]string{"priority"},
	)

	JobsProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ts2go_jobs_processed_total",
			Help: "Total number of jobs processed",
		},
		[]string{"status"}, // completed, failed
	)

	JobProcessingDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "ts2go_job_processing_duration_seconds",
			Help:    "Job processing duration in seconds",
			Buckets: []float64{1, 5, 10, 30, 60, 120, 300, 600},
		},
	)

	QueueDepth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ts2go_queue_depth",
			Help: "Current depth of job queue",
		},
		[]string{"priority"},
	)

	// Worker metrics
	WorkerActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "ts2go_worker_active",
			Help: "Number of active workers",
		},
	)

	WorkerJobsProcessing = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "ts2go_worker_jobs_processing",
			Help: "Number of jobs currently being processed",
		},
	)

	// File metrics
	FilesUploaded = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ts2go_files_uploaded_total",
			Help: "Total number of files uploaded",
		},
		[]string{"type"}, // input, output
	)

	FilesUploadedBytes = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "ts2go_files_uploaded_bytes_total",
			Help: "Total bytes of files uploaded",
		},
	)

	StorageUsageBytes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ts2go_storage_usage_bytes",
			Help: "Current storage usage in bytes",
		},
		[]string{"user_id"},
	)

	// Database metrics
	DatabaseConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ts2go_database_connections",
			Help: "Number of database connections",
		},
		[]string{"state"}, // open, idle, in_use
	)

	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ts2go_database_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"operation"},
	)

	// Cache metrics
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ts2go_cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type"},
	)

	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ts2go_cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type"},
	)

	// WebSocket metrics
	WebSocketConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "ts2go_websocket_connections",
			Help: "Number of active WebSocket connections",
		},
	)

	WebSocketMessagesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ts2go_websocket_messages_total",
			Help: "Total number of WebSocket messages",
		},
		[]string{"direction"}, // sent, received
	)

	// Authentication metrics
	AuthAttemptsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ts2go_auth_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"method", "status"}, // jwt/apikey, success/failure
	)

	// Rate limit metrics
	RateLimitExceeded = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ts2go_rate_limit_exceeded_total",
			Help: "Total number of rate limit violations",
		},
		[]string{"endpoint"},
	)
)

// RecordHTTPRequest records HTTP request metrics
func RecordHTTPRequest(method, path string, status int, duration float64) {
	HTTPRequestsTotal.WithLabelValues(method, path, string(rune(status))).Inc()
	HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
}

// RecordJobEnqueued records job enqueue metrics
func RecordJobEnqueued(priority string) {
	JobsEnqueued.WithLabelValues(priority).Inc()
}

// RecordJobProcessed records job processing completion
func RecordJobProcessed(status string, duration float64) {
	JobsProcessed.WithLabelValues(status).Inc()
	JobProcessingDuration.Observe(duration)
}

// UpdateQueueDepth updates queue depth gauge
func UpdateQueueDepth(priority string, depth int) {
	QueueDepth.WithLabelValues(priority).Set(float64(depth))
}

// RecordFileUpload records file upload metrics
func RecordFileUpload(fileType string, sizeBytes int64) {
	FilesUploaded.WithLabelValues(fileType).Inc()
	FilesUploadedBytes.Add(float64(sizeBytes))
}

// RecordCacheAccess records cache hit/miss
func RecordCacheAccess(cacheType string, hit bool) {
	if hit {
		CacheHits.WithLabelValues(cacheType).Inc()
	} else {
		CacheMisses.WithLabelValues(cacheType).Inc()
	}
}

// RecordAuthAttempt records authentication attempt
func RecordAuthAttempt(method string, success bool) {
	status := "failure"
	if success {
		status = "success"
	}
	AuthAttemptsTotal.WithLabelValues(method, status).Inc()
}
