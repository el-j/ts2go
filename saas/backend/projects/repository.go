package projects

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/el-j/ts2go/saas/backend/db/models"
	"github.com/google/uuid"
)

// Repository handles database operations for projects
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new projects repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new project
func (r *Repository) Create(ctx context.Context, userID uuid.UUID, name, slug, description string, visibility string, teamID *uuid.UUID) (*models.Project, error) {
	project := &models.Project{
		ID:         uuid.New(),
		UserID:     userID,
		TeamID:     teamID,
		Name:       name,
		Slug:       slug,
		Visibility: visibility,
		Settings:   "{}",
		TotalFiles: 0,
		IsArchived: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if description != "" {
		project.Description = &description
	}

	query := `
		INSERT INTO projects (id, user_id, team_id, name, slug, description, visibility, settings, total_files, total_size_bytes, is_archived, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, user_id, team_id, name, slug, description, visibility, settings, total_files, total_size_bytes, is_archived, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		project.ID, project.UserID, project.TeamID, project.Name, project.Slug,
		project.Description, project.Visibility, project.Settings, project.TotalFiles,
		project.TotalSizeBytes, project.IsArchived, project.CreatedAt, project.UpdatedAt,
	).Scan(
		&project.ID, &project.UserID, &project.TeamID, &project.Name, &project.Slug,
		&project.Description, &project.Visibility, &project.Settings, &project.TotalFiles,
		&project.TotalSizeBytes, &project.IsArchived, &project.CreatedAt, &project.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

// GetByID retrieves a project by ID
func (r *Repository) GetByID(ctx context.Context, projectID uuid.UUID) (*models.Project, error) {
	project := &models.Project{}
	query := `
		SELECT id, user_id, team_id, name, slug, description, visibility, settings,
		       total_files, total_size_bytes, last_transpiled_at, is_archived, created_at, updated_at
		FROM projects
		WHERE id = $1 AND is_archived = false
	`

	err := r.db.QueryRowContext(ctx, query, projectID).Scan(
		&project.ID, &project.UserID, &project.TeamID, &project.Name, &project.Slug,
		&project.Description, &project.Visibility, &project.Settings, &project.TotalFiles,
		&project.TotalSizeBytes, &project.LastTranspiledAt, &project.IsArchived,
		&project.CreatedAt, &project.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

// GetBySlug retrieves a project by user ID and slug
func (r *Repository) GetBySlug(ctx context.Context, userID uuid.UUID, slug string) (*models.Project, error) {
	project := &models.Project{}
	query := `
		SELECT id, user_id, team_id, name, slug, description, visibility, settings,
		       total_files, total_size_bytes, last_transpiled_at, is_archived, created_at, updated_at
		FROM projects
		WHERE user_id = $1 AND slug = $2 AND is_archived = false
	`

	err := r.db.QueryRowContext(ctx, query, userID, slug).Scan(
		&project.ID, &project.UserID, &project.TeamID, &project.Name, &project.Slug,
		&project.Description, &project.Visibility, &project.Settings, &project.TotalFiles,
		&project.TotalSizeBytes, &project.LastTranspiledAt, &project.IsArchived,
		&project.CreatedAt, &project.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

// ListByUser lists all projects for a user
func (r *Repository) ListByUser(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]*models.Project, error) {
	query := `
		SELECT id, user_id, team_id, name, slug, description, visibility, settings,
		       total_files, total_size_bytes, last_transpiled_at, is_archived, created_at, updated_at
		FROM projects
		WHERE user_id = $1
	`

	if !includeArchived {
		query += " AND is_archived = false"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err := rows.Scan(
			&project.ID, &project.UserID, &project.TeamID, &project.Name, &project.Slug,
			&project.Description, &project.Visibility, &project.Settings, &project.TotalFiles,
			&project.TotalSizeBytes, &project.LastTranspiledAt, &project.IsArchived,
			&project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// Update updates a project
func (r *Repository) Update(ctx context.Context, projectID uuid.UUID, name, description *string, visibility *string) error {
	query := `
		UPDATE projects
		SET name = COALESCE($2, name),
		    description = COALESCE($3, description),
		    visibility = COALESCE($4, visibility),
		    updated_at = NOW()
		WHERE id = $1 AND is_archived = false
	`

	result, err := r.db.ExecContext(ctx, query, projectID, name, description, visibility)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("project not found or already archived")
	}

	return nil
}

// Delete soft deletes a project (archives it)
func (r *Repository) Delete(ctx context.Context, projectID, userID uuid.UUID) error {
	query := `
		UPDATE projects
		SET is_archived = true, updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND is_archived = false
	`

	result, err := r.db.ExecContext(ctx, query, projectID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("project not found or unauthorized")
	}

	return nil
}

// UpdateStats updates project statistics
func (r *Repository) UpdateStats(ctx context.Context, projectID uuid.UUID, totalFiles int, totalSizeBytes int64) error {
	query := `
		UPDATE projects
		SET total_files = $2,
		    total_size_bytes = $3,
		    last_transpiled_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, projectID, totalFiles, totalSizeBytes)
	if err != nil {
		return fmt.Errorf("failed to update project stats: %w", err)
	}

	return nil
}

// CheckAccess checks if a user has access to a project
func (r *Repository) CheckAccess(ctx context.Context, projectID, userID uuid.UUID) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM projects
		WHERE id = $1 AND (
			user_id = $2 OR
			team_id IN (
				SELECT team_id FROM team_members WHERE user_id = $2 AND is_active = true
			) OR
			visibility = 'public'
		) AND is_archived = false
	`

	err := r.db.QueryRowContext(ctx, query, projectID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check access: %w", err)
	}

	return count > 0, nil
}
