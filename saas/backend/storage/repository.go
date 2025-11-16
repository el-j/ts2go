package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/el-j/ts2go/saas/backend/db/models"
	"github.com/google/uuid"
)

// Repository handles database operations for files
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new files repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateFile stores file metadata in database
func (r *Repository) CreateFile(ctx context.Context, projectID, transpilationID *uuid.UUID, fileType, originalName, storagePath, mimeType string, sizeBytes int64, checksum string) (*models.File, error) {
	file := &models.File{
		ID:              uuid.New(),
		ProjectID:       *projectID,
		TranspilationID: transpilationID,
		FileType:        fileType,
		OriginalName:    originalName,
		StoragePath:     storagePath,
		MimeType:        &mimeType,
		SizeBytes:       sizeBytes,
		Checksum:        &checksum,
		CreatedAt:       time.Now(),
	}

	query := `
		INSERT INTO files (id, project_id, transpilation_id, file_type, original_name, storage_path, mime_type, size_bytes, checksum, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, project_id, transpilation_id, file_type, original_name, storage_path, mime_type, size_bytes, checksum, created_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		file.ID, file.ProjectID, file.TranspilationID, file.FileType,
		file.OriginalName, file.StoragePath, file.MimeType, file.SizeBytes,
		file.Checksum, file.CreatedAt,
	).Scan(
		&file.ID, &file.ProjectID, &file.TranspilationID, &file.FileType,
		&file.OriginalName, &file.StoragePath, &file.MimeType, &file.SizeBytes,
		&file.Checksum, &file.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	return file, nil
}

// GetFile retrieves a file by ID
func (r *Repository) GetFile(ctx context.Context, fileID uuid.UUID) (*models.File, error) {
	file := &models.File{}
	query := `
		SELECT id, project_id, transpilation_id, file_type, original_name, storage_path, mime_type, size_bytes, checksum, created_at
		FROM files
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, fileID).Scan(
		&file.ID, &file.ProjectID, &file.TranspilationID, &file.FileType,
		&file.OriginalName, &file.StoragePath, &file.MimeType, &file.SizeBytes,
		&file.Checksum, &file.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("file not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return file, nil
}

// ListProjectFiles lists all files for a project
func (r *Repository) ListProjectFiles(ctx context.Context, projectID uuid.UUID) ([]*models.File, error) {
	query := `
		SELECT id, project_id, transpilation_id, file_type, original_name, storage_path, mime_type, size_bytes, checksum, created_at
		FROM files
		WHERE project_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	defer rows.Close()

	var files []*models.File
	for rows.Next() {
		file := &models.File{}
		err := rows.Scan(
			&file.ID, &file.ProjectID, &file.TranspilationID, &file.FileType,
			&file.OriginalName, &file.StoragePath, &file.MimeType, &file.SizeBytes,
			&file.Checksum, &file.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	return files, nil
}

// DeleteFile deletes a file record
func (r *Repository) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	query := `DELETE FROM files WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, fileID)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("file not found")
	}

	return nil
}
