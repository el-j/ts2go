package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user account
type User struct {
	ID                     uuid.UUID  `json:"id" db:"id"`
	Email                  string     `json:"email" db:"email"`
	Username               string     `json:"username" db:"username"`
	PasswordHash           string     `json:"-" db:"password_hash"`
	FullName               *string    `json:"full_name,omitempty" db:"full_name"`
	AvatarURL              *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	EmailVerified          bool       `json:"email_verified" db:"email_verified"`
	EmailVerificationToken *string    `json:"-" db:"email_verification_token"`
	PasswordResetToken     *string    `json:"-" db:"password_reset_token"`
	PasswordResetExpires   *time.Time `json:"-" db:"password_reset_expires"`
	IsActive               bool       `json:"is_active" db:"is_active"`
	LastLoginAt            *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`
}

// PublicUser returns a user object safe for public display
type PublicUser struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FullName  *string   `json:"full_name,omitempty"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ToPublic converts a User to PublicUser
func (u *User) ToPublic() PublicUser {
	return PublicUser{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		FullName:  u.FullName,
		AvatarURL: u.AvatarURL,
		CreatedAt: u.CreatedAt,
	}
}

// APIKey represents an API key for programmatic access
type APIKey struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	Name       string     `json:"name" db:"name"`
	KeyHash    string     `json:"-" db:"key_hash"`
	KeyPrefix  string     `json:"key_prefix" db:"key_prefix"`
	Scopes     []string   `json:"scopes" db:"scopes"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

// Team represents an organization/team
type Team struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description *string   `json:"description,omitempty" db:"description"`
	AvatarURL   *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	OwnerID     uuid.UUID `json:"owner_id" db:"owner_id"`
	PlanType    string    `json:"plan_type" db:"plan_type"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// TeamMember represents a team membership
type TeamMember struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	TeamID    uuid.UUID  `json:"team_id" db:"team_id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Role      string     `json:"role" db:"role"`
	InvitedBy *uuid.UUID `json:"invited_by,omitempty" db:"invited_by"`
	InvitedAt time.Time  `json:"invited_at" db:"invited_at"`
	JoinedAt  *time.Time `json:"joined_at,omitempty" db:"joined_at"`
	IsActive  bool       `json:"is_active" db:"is_active"`
}

// Project represents a code project
type Project struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	TeamID           *uuid.UUID `json:"team_id,omitempty" db:"team_id"`
	Name             string     `json:"name" db:"name"`
	Slug             string     `json:"slug" db:"slug"`
	Description      *string    `json:"description,omitempty" db:"description"`
	Visibility       string     `json:"visibility" db:"visibility"`
	Settings         string     `json:"settings" db:"settings"` // JSONB as string
	TotalFiles       int        `json:"total_files" db:"total_files"`
	TotalSizeBytes   int64      `json:"total_size_bytes" db:"total_size_bytes"`
	LastTranspiledAt *time.Time `json:"last_transpiled_at,omitempty" db:"last_transpiled_at"`
	IsArchived       bool       `json:"is_archived" db:"is_archived"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// Transpilation represents a transpilation job
type Transpilation struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	ProjectID        *uuid.UUID `json:"project_id,omitempty" db:"project_id"`
	JobID            string     `json:"job_id" db:"job_id"`
	Status           string     `json:"status" db:"status"`
	InputFiles       string     `json:"input_files" db:"input_files"`             // JSONB as string
	OutputFiles      *string    `json:"output_files,omitempty" db:"output_files"` // JSONB as string
	Settings         string     `json:"settings" db:"settings"`                   // JSONB as string
	ErrorMessage     *string    `json:"error_message,omitempty" db:"error_message"`
	InputSizeBytes   *int64     `json:"input_size_bytes,omitempty" db:"input_size_bytes"`
	OutputSizeBytes  *int64     `json:"output_size_bytes,omitempty" db:"output_size_bytes"`
	ProcessingTimeMs *int       `json:"processing_time_ms,omitempty" db:"processing_time_ms"`
	WorkerID         *string    `json:"worker_id,omitempty" db:"worker_id"`
	StartedAt        *time.Time `json:"started_at,omitempty" db:"started_at"`
	CompletedAt      *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// UsageRecord represents a usage event
type UsageRecord struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	TeamID       *uuid.UUID `json:"team_id,omitempty" db:"team_id"`
	ResourceType string     `json:"resource_type" db:"resource_type"`
	Quantity     int        `json:"quantity" db:"quantity"`
	Metadata     string     `json:"metadata" db:"metadata"` // JSONB as string
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// File represents a file stored in the system
type File struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	ProjectID       uuid.UUID  `json:"project_id" db:"project_id"`
	TranspilationID *uuid.UUID `json:"transpilation_id,omitempty" db:"transpilation_id"`
	FileType        string     `json:"file_type" db:"file_type"` // input, output, archive
	OriginalName    string     `json:"original_name" db:"original_name"`
	StoragePath     string     `json:"storage_path" db:"storage_path"`
	MimeType        *string    `json:"mime_type,omitempty" db:"mime_type"`
	SizeBytes       int64      `json:"size_bytes" db:"size_bytes"`
	Checksum        *string    `json:"checksum,omitempty" db:"checksum"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	Action       string     `json:"action" db:"action"`
	ResourceType *string    `json:"resource_type,omitempty" db:"resource_type"`
	ResourceID   *uuid.UUID `json:"resource_id,omitempty" db:"resource_id"`
	IPAddress    *string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent    *string    `json:"user_agent,omitempty" db:"user_agent"`
	Metadata     string     `json:"metadata" db:"metadata"` // JSONB as string
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}
