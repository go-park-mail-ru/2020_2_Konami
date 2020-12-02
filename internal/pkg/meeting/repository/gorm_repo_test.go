package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"konami_backend/internal/pkg/meeting"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	DB      *gorm.DB
	mock    sqlmock.Sqlmock

	bdError error
	repository meeting.Repository
}

func (s *Suite) SetupSuite() {
	var db  *sql.DB
	var err error

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{
		DriverName:           "postgres",
		DSN:                  "sqlmock_db_0",
		PreferSimpleProtocol: true,
		Conn:                 db,
	}), &gorm.Config{})

	s.bdError = errors.New("some bd error")
	require.NoError(s.T(), err)

	s.repository = NewMeetingGormRepoLite(s.DB)
}

func (s *Suite) TestGetSessions() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err := s.repository.GetMeeting(1, 2, true)
	require.NoError(s.T(), err)
}

func (s *Suite) TestSearchMeetErr() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err := s.repository.SearchMeetings(meeting.FilterParams{
		StartDate:  time.Time{},
		EndDate:    time.Time{},
		PrevId:     0,
		CountLimit: 0,
		UserId:     0,
	}, "testmeet", 100)
	require.NoError(s.T(), err)
}

func (s *Suite) TestSearchMeet() {
	s.mock.ExpectQuery("SELECT").
		WillReturnError(s.bdError)

	_, err := s.repository.SearchMeetings(meeting.FilterParams{
		StartDate:  time.Time{},
		EndDate:    time.Time{},
		PrevId:     0,
		CountLimit: 0,
		UserId:     0,
	}, "testmeet", 100)
	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestFilterMeet() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err := s.repository.FilterTagged(meeting.FilterParams{
		StartDate:  time.Time{},
		EndDate:    time.Time{},
		PrevId:     0,
		CountLimit: 0,
		UserId:     0,
	}, 100)
	require.NoError(s.T(), err)
}

func (s *Suite) TestFilterTopMeet() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err := s.repository.GetTopMeetings(meeting.FilterParams{
		StartDate:  time.Time{},
		EndDate:    time.Time{},
		PrevId:     0,
		CountLimit: 0,
		UserId:     0,
	})
	require.NoError(s.T(), err)
}

func (s *Suite) TestFilterReg() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err := s.repository.FilterRegistered(meeting.FilterParams{
		StartDate:  time.Time{},
		EndDate:    time.Time{},
		PrevId:     0,
		CountLimit: 0,
		UserId:     0,
	})
	require.NoError(s.T(), err)
}

func (s *Suite) TestFilterTopMeetErr() {
	s.mock.ExpectQuery("SELECT").
		WillReturnError(s.bdError)

	_, err := s.repository.GetTopMeetings(meeting.FilterParams{
		StartDate:  time.Time{},
		EndDate:    time.Time{},
		PrevId:     0,
		CountLimit: 0,
		UserId:     0,
	})
	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestFilterMeetErr() {
	s.mock.ExpectQuery("SELECT").
		WillReturnError(s.bdError)

	_, err := s.repository.FilterTagged(meeting.FilterParams{
		StartDate:  time.Time{},
		EndDate:    time.Time{},
		PrevId:     0,
		CountLimit: 0,
		UserId:     0,
	}, 100)
	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestGetSessionsError() {
	s.mock.ExpectQuery("SELECT").
		WillReturnError(s.bdError)

	_, err := s.repository.GetMeeting(1, 2, true)
	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestGetNextMeetings() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err := s.repository.GetNextMeetings(meeting.FilterParams{
		StartDate:  time.Time{},
		EndDate:    time.Time{},
		PrevId:     0,
		CountLimit: 0,
		UserId:     0,
	})

	require.NoError(s.T(), err)
}


func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestMeetings(t *testing.T) {
	suite.Run(t, new(Suite))
}