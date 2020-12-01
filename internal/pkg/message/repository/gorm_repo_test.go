package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"konami_backend/internal/pkg/message"
	"konami_backend/internal/pkg/models"
	"testing"
)

type Suite struct {
	suite.Suite
	DB      *gorm.DB
	mock    sqlmock.Sqlmock

	bdError error
	repository message.Repository
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

	s.repository = NewMeetingGormRepo(s.DB)
}

func (s *Suite) TestGetSessions() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err := s.repository.GetMessages(0)
	require.NoError(s.T(), err)
}

func (s *Suite) TestGetSessionsError() {
	s.mock.ExpectQuery("SELECT").
		WillReturnError(s.bdError)

	_, err := s.repository.GetMessages(0)
	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestSaveMes() {
	s.mock.ExpectQuery("INSERT INTO").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	gg := models.Message{
		Id:        1,
		AuthorId:  2,
		MeetingId: 3,
		Text:      "qwer",
		Timestamp: "2006-01-02T15:04:05.000Z",
	}

	_, err := s.repository.SaveMessage(gg)
	require.NoError(s.T(), err)
}

func (s *Suite) TestSaveMesErr() {
	s.mock.ExpectQuery("INSERT INTO").
		WillReturnError(s.bdError)

	gg := models.Message{
		Id:        1,
		AuthorId:  2,
		MeetingId: 3,
		Text:      "qwer",
		Timestamp: "2006-01-02T15:04:05.000Z",
	}

	_, err := s.repository.SaveMessage(gg)
	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestSaveMesErrSec() {
	gg := models.Message{
		Id:        1,
		AuthorId:  2,
		MeetingId: 3,
		Text:      "qwer",
		Timestamp: "2006-01-02T15:04:05.000Zlsnfln",
	}

	_, err := s.repository.SaveMessage(gg)
	require.Error(s.T(), err)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestMeetings(t *testing.T) {
	suite.Run(t, new(Suite))
}
