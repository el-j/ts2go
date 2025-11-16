package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"net/http"

	"github.com/el-j/ts2go/saas/backend/logger"
	"github.com/el-j/ts2go/saas/backend/queue"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// Hub maintains active WebSocket connections
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	queue      *queue.Queue
}

// Client represents a WebSocket client connection
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID uuid.UUID
	jobIDs map[string]bool // Jobs this client is subscribed to
	mu     sync.RWMutex
}

// Message represents a WebSocket message
type Message struct {
	Type   string                 `json:"type"`
	JobID  string                 `json:"job_id,omitempty"`
	Status string                 `json:"status,omitempty"`
	Data   map[string]interface{} `json:"data,omitempty"`
	UserID string                 `json:"user_id,omitempty"`
	Error  string                 `json:"error,omitempty"`
}

// NewHub creates a new WebSocket hub
func NewHub(queue *queue.Queue) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		queue:      queue,
	}
}

// Run starts the hub
func (h *Hub) Run(ctx context.Context) {
	// Start job status monitor
	go h.monitorJobStatus(ctx)

	for {
		select {
		case <-ctx.Done():
			return

		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			logger.Log.Info().
				Str("user_id", client.userID.String()).
				Msg("WebSocket client connected")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			logger.Log.Info().
				Str("user_id", client.userID.String()).
				Msg("WebSocket client disconnected")

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				// Only send to clients interested in this job
				if message.JobID != "" {
					client.mu.RLock()
					interested := client.jobIDs[message.JobID]
					client.mu.RUnlock()

					if !interested {
						continue
					}
				}

				// Only send to owner of the job
				if message.UserID != "" && message.UserID != client.userID.String() {
					continue
				}

				select {
				case client.send <- mustMarshal(message):
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// monitorJobStatus polls for job status changes and broadcasts updates
func (h *Hub) monitorJobStatus(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	lastStatus := make(map[string]string)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Get all subscribed job IDs
			jobIDs := h.getSubscribedJobIDs()

			for jobID := range jobIDs {
				jobUUID, err := uuid.Parse(jobID)
				if err != nil {
					continue
				}

				job, err := h.queue.GetJob(ctx, jobUUID)
				if err != nil {
					continue
				}

				// Check if status changed
				if lastStatus[jobID] != job.Status {
					lastStatus[jobID] = job.Status

					// Broadcast update
					h.broadcast <- &Message{
						Type:   "job_update",
						JobID:  jobID,
						Status: job.Status,
						Data: map[string]interface{}{
							"started_at":   job.StartedAt,
							"completed_at": job.CompletedAt,
							"error":        job.ErrorMessage,
						},
						UserID: job.UserID.String(),
					}

					logger.Log.Debug().
						Str("job_id", jobID).
						Str("status", job.Status).
						Msg("Job status changed")
				}
			}
		}
	}
}

// getSubscribedJobIDs returns all job IDs that clients are subscribed to
func (h *Hub) getSubscribedJobIDs() map[string]bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	jobIDs := make(map[string]bool)
	for client := range h.clients {
		client.mu.RLock()
		for jobID := range client.jobIDs {
			jobIDs[jobID] = true
		}
		client.mu.RUnlock()
	}
	return jobIDs
}

// ServeWS handles WebSocket requests from clients
func (h *Hub) ServeWS(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to upgrade WebSocket")
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))
	client := &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userUUID,
		jobIDs: make(map[string]bool),
	}

	h.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// readPump reads messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Log.Error().Err(err).Msg("WebSocket read error")
			}
			break
		}

		// Handle client messages (e.g., subscribe to job)
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		if msg.Type == "subscribe" && msg.JobID != "" {
			c.mu.Lock()
			c.jobIDs[msg.JobID] = true
			c.mu.Unlock()

			logger.Log.Debug().
				Str("user_id", c.userID.String()).
				Str("job_id", msg.JobID).
				Msg("Client subscribed to job")
		}
	}
}

// writePump writes messages to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Write queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// mustMarshal marshals a message or panics
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
