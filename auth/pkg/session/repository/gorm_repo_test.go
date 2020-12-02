package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"konami_backend/auth/pkg/models"
	"konami_backend/auth/pkg/session"
	"testing"
)

type Suite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	sessions   []models.Session
	bdError    error
	repository session.Repository
}

func (s *Suite) SetupSuite() {
	s1 := models.Session{
		UserId: 1,
		Token:  "tokkken",
	}
	s.sessions = []models.Session{s1}

	var db *sql.DB
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

	s.repository = NewSessionGormRepo(s.DB)
}

func (s *Suite) TestGetSessions() {
	testId := "tokkken"
	testSession := s.sessions[0]

	s.mock.ExpectQuery("SELECT").
		WithArgs(testId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "userid", "token"}).
			AddRow(1, testSession.UserId, testSession.Token))

	_, err := s.repository.GetUserId(testSession.Token)

	require.NoError(s.T(), err)
}

func (s *Suite) TestCreateTagError() {
	s.mock.ExpectQuery("SELECT").
		WithArgs("LOL").
		WillReturnError(s.bdError)

	_, err := s.repository.GetUserId("LOL")
	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestCreateSession() {
	// CANT TEST uuid.New().String()
}

func (s *Suite) TestDeleteSessions() {
	testId := "tokkken"
	testSession := s.sessions[0]

	s.mock.ExpectBegin()
	s.mock.ExpectExec("DELETE FROM").
		WithArgs(testId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	err := s.repository.RemoveSession(testSession.Token)
	require.NoError(s.T(), err)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestMeetings(t *testing.T) {
	suite.Run(t, new(Suite))
}

/*GetUserId(token string) (userId int, err error)
CreateSession(userId int) (token string, err error)
RemoveSession(token string) error
*/
