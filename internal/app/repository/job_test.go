package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type JobRepositoryTestSuite struct {
	suite.Suite
	DB         *sql.DB
	mock       sqlmock.Sqlmock
	repository *JobRepository
}

func (j *JobRepositoryTestSuite) SetupTest() {
	var err error
	j.DB, j.mock, err = sqlmock.New()

	j.mock.ExpectPrepare("select")
	j.mock.ExpectPrepare("insert")
	j.mock.ExpectPrepare("update")
	j.mock.ExpectPrepare("delete")
	j.repository = NewJobRepository(j.DB)
	require.NoError(j.T(), err)
}

func TestJobRepositoryTestSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(JobRepositoryTestSuite))
}

func (j *JobRepositoryTestSuite) TestGetJobForRun() {

	j.mock.ExpectBegin()
	j.mock.ExpectQuery("select").WillReturnRows(sqlmock.NewRows([]string{"id", "order_number", "created_at", "updated_at", "next_run", "run_cnt"}).
		AddRow("1", "101", time.Now(), time.Now(), time.Now(), "2"))

	j.mock.ExpectQuery("update").WithArgs(int64(1)).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
	j.mock.ExpectCommit()
	jobs, err := j.repository.GetJobForRun(context.Background())
	require.NoError(j.T(), err)

	if len(*jobs) != 1 {
		j.T().Error("expected 1 job")
	}
	job := (*jobs)[0]
	require.NotNil(j.T(), job)
	require.Equal(j.T(), "101", job.OrderNumber)
}

func (j *JobRepositoryTestSuite) TestCreateJobByOrderNumber() {
	orderNumber := "100"
	expectedID := int64(89)
	j.mock.ExpectQuery("insert into").WithArgs(orderNumber).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
	actualID, err := j.repository.CreateJobByOrderNumber(context.Background(), orderNumber)
	require.NoError(j.T(), err)
	require.Equal(j.T(), expectedID, actualID)
}

func (j *JobRepositoryTestSuite) TestUpdateJobByOrderNumber() {
	orderNumber := "100"
	j.mock.ExpectQuery("update").WithArgs(orderNumber).WillReturnRows(sqlmock.NewRows([]string{""}))
	err := j.repository.UpdateJobByOrderNumber(context.Background(), orderNumber)
	require.NoError(j.T(), err)
}

func (j *JobRepositoryTestSuite) TestDeleteJobByOrderNumber() {
	orderNumber := "100"
	j.mock.ExpectQuery("delete").WithArgs(orderNumber).WillReturnRows(sqlmock.NewRows([]string{""}))
	err := j.repository.DeleteJobByOrderNumber(context.Background(), orderNumber)
	require.NoError(j.T(), err)
}
