package auth

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/el-j/ts2go/saas/backend/db/models"
	"github.com/google/uuid"
)

// Repository handles database operations for authentication
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateUser creates a new user in the database
func (r *Repository) CreateUser(ctx context.Context, email, username, passwordHash string) (*models.User, error) {
	user := &models.User{
		ID:            uuid.New(),
		Email:         email,
		Username:      username,
		PasswordHash:  passwordHash,
		EmailVerified: false,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	query := `
		INSERT INTO users (id, email, username, password_hash, email_verified, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, email, username, email_verified, is_active, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		user.ID, user.Email, user.Username, user.PasswordHash,
		user.EmailVerified, user.IsActive, user.CreatedAt, user.UpdatedAt,
	).Scan(
		&user.ID, &user.Email, &user.Username,
		&user.EmailVerified, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, full_name, avatar_url,
		       email_verified, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true
	`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.FullName, &user.AvatarURL, &user.EmailVerified, &user.IsActive,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, full_name, avatar_url,
		       email_verified, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE username = $1 AND is_active = true
	`

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.FullName, &user.AvatarURL, &user.EmailVerified, &user.IsActive,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, full_name, avatar_url,
		       email_verified, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.FullName, &user.AvatarURL, &user.EmailVerified, &user.IsActive,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateLastLogin updates the last login timestamp
func (r *Repository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET last_login_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

// CreateAPIKey creates a new API key
func (r *Repository) CreateAPIKey(ctx context.Context, userID uuid.UUID, name, keyHash, keyPrefix string, scopes []string, expiresAt *time.Time) (*models.APIKey, error) {
	apiKey := &models.APIKey{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		Scopes:    scopes,
		ExpiresAt: expiresAt,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, scopes, expires_at, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, user_id, name, key_prefix, scopes, expires_at, is_active, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		apiKey.ID, apiKey.UserID, apiKey.Name, apiKey.KeyHash, apiKey.KeyPrefix,
		apiKey.Scopes, apiKey.ExpiresAt, apiKey.IsActive, apiKey.CreatedAt, apiKey.UpdatedAt,
	).Scan(
		&apiKey.ID, &apiKey.UserID, &apiKey.Name, &apiKey.KeyPrefix,
		&apiKey.Scopes, &apiKey.ExpiresAt, &apiKey.IsActive, &apiKey.CreatedAt, &apiKey.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return apiKey, nil
}

// GetAPIKeyByHash retrieves an API key by its hash
func (r *Repository) GetAPIKeyByHash(ctx context.Context, keyHash string) (*models.APIKey, error) {
	apiKey := &models.APIKey{}
	query := `
		SELECT id, user_id, name, key_hash, key_prefix, scopes, last_used_at, expires_at, is_active, created_at, updated_at
		FROM api_keys
		WHERE key_hash = $1 AND is_active = true
	`

	err := r.db.QueryRowContext(ctx, query, keyHash).Scan(
		&apiKey.ID, &apiKey.UserID, &apiKey.Name, &apiKey.KeyHash, &apiKey.KeyPrefix,
		&apiKey.Scopes, &apiKey.LastUsedAt, &apiKey.ExpiresAt, &apiKey.IsActive,
		&apiKey.CreatedAt, &apiKey.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("API key not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	return apiKey, nil
}

// UpdateAPIKeyLastUsed updates the last used timestamp
func (r *Repository) UpdateAPIKeyLastUsed(ctx context.Context, keyID uuid.UUID) error {
	query := `UPDATE api_keys SET last_used_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), keyID)
	if err != nil {
		return fmt.Errorf("failed to update API key last used: %w", err)
	}
	return nil
}

// ListAPIKeysByUserID lists all API keys for a user
func (r *Repository) ListAPIKeysByUserID(ctx context.Context, userID uuid.UUID) ([]*models.APIKey, error) {
	query := `
		SELECT id, user_id, name, key_prefix, scopes, last_used_at, expires_at, is_active, created_at, updated_at
		FROM api_keys
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}
	defer rows.Close()

	var keys []*models.APIKey
	for rows.Next() {
		key := &models.APIKey{}
		err := rows.Scan(
			&key.ID, &key.UserID, &key.Name, &key.KeyPrefix,
			&key.Scopes, &key.LastUsedAt, &key.ExpiresAt, &key.IsActive,
			&key.CreatedAt, &key.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key: %w", err)
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// DeleteAPIKey soft deletes an API key
func (r *Repository) DeleteAPIKey(ctx context.Context, keyID uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE api_keys SET is_active = false WHERE id = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, keyID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("API key not found or unauthorized")
	}

	return nil
}
