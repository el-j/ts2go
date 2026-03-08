package background

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/el-j/ts2go/saas/backend/redis"
	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
)

func TestCleanupOldJobsJob(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rdb, mockRedis := redismock.NewClientMock()
	redisClient := &redis.Client{Client: rdb}

	job1 := uuid.New().String()
	job2 := uuid.New().String()

	// Mock DB returning 2 deleted rows
	rows := sqlmock.NewRows([]string{"id"}).AddRow(job1).AddRow(job2)
	mockDB.ExpectQuery(regexp.QuoteMeta("DELETE FROM transpilations WHERE completed_at < $1 RETURNING id")).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	// Mock Redis deleting those 2 keys
	mockRedis.ExpectDel("job:"+job1, "job:"+job2).SetVal(2)

	jobFunc := CleanupOldJobsJob(db, redisClient)
	err = jobFunc(context.Background())

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("DB expectations not met: %v", err)
	}
	if err := mockRedis.ExpectationsWereMet(); err != nil {
		t.Errorf("Redis expectations not met: %v", err)
	}
}

func TestUpdateQueueMetricsJob(t *testing.T) {
	rdb, mockRedis := redismock.NewClientMock()
	redisClient := &redis.Client{Client: rdb}

	mockRedis.ExpectZCard("queue:transpilation:low").SetVal(1)
	mockRedis.ExpectZCard("queue:transpilation:normal").SetVal(2)
	mockRedis.ExpectZCard("queue:transpilation:high").SetVal(0)
	mockRedis.ExpectZCard("queue:transpilation:urgent").SetVal(5)

	jobFunc := UpdateQueueMetricsJob(redisClient)
	err := jobFunc(context.Background())

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if err := mockRedis.ExpectationsWereMet(); err != nil {
		t.Errorf("Redis expectations not met: %v", err)
	}
}
