package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/session"
	"regexp"
	"testing"
)


type Suite struct {
	suite.Suite
	DB      *gorm.DB
	mock    sqlmock.Sqlmock
	sessions []models.Session

	repository session.Repository
}

func (s *Suite) SetupSuite() {
	var db  *sql.DB
	var err error

	s1 := models.Session{
		UserId: 1,
		Token:  "tokkken",
	}
	s.sessions = []models.Session{s1}

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	s.DB.LogMode(true)

	s.repository = NewSessionGormRepo(s.DB)
}

func (s *Suite) TestGetSessions() {
	testId := "tokkken"
	testSession := s.sessions[0]

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "sessions" WHERE (Token = $1)`)).
		WithArgs(testId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "userid", "token"}).
			AddRow(1, testSession.UserId, testSession.Token))

	res, err := s.repository.GetUserId(testSession.Token)

	require.NoError(s.T(), err)

	require.Nil(s.T(), deep.Equal(0, res))
}


func (s *Suite) TestDeleteSessions() {
	testId := "tokkken"
	testSession := s.sessions[0]

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "sessions" WHERE (token = $1)`)).
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
