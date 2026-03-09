package projects

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestCheckAccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	projectID := uuid.New()
	userID := uuid.New()

	// Test 1: Owner (returns 1)
	mock.ExpectQuery(`SELECT COUNT\(\*\)`).WithArgs(projectID, userID).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	hasAccess, err := repo.CheckAccess(context.Background(), projectID, userID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !hasAccess {
		t.Errorf("Expected owner to have access")
	}

	// Test 2: Not Owner/Team Member (returns 0)
	mock.ExpectQuery(`SELECT COUNT\(\*\)`).WithArgs(projectID, userID).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	hasAccess, err = repo.CheckAccess(context.Background(), projectID, userID)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if hasAccess {
		t.Errorf("Expected unauthorized user to not have access")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
