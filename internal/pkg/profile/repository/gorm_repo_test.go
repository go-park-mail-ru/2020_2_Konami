package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"konami_backend/internal/pkg/profile"
	"testing"
)

type Suite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository profile.Repository
	bdError    error
}

func (s *Suite) SetupSuite() {
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

	s.repository = NewProfileGormRepo(s.DB)
}

func (s *Suite) TestGetAll() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, err := s.repository.GetAll(profile.FilterParams{})

	require.NoError(s.T(), err)
}

func (s *Suite) TestGetAllError() {
	s.mock.ExpectQuery("SELECT").
		WillReturnError(s.bdError)

	_, err := s.repository.GetAll(profile.FilterParams{})

	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestGetCred() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, _, err := s.repository.GetCredentials("qwerty")

	require.NoError(s.T(), err)
}

func (s *Suite) TestGetLabel() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err := s.repository.GetLabel(1)

	require.NoError(s.T(), err)
}

func (s *Suite) TestGetSub() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err := s.repository.GetTagSubscriptions(1)

	require.NoError(s.T(), err)
}

func (s *Suite) TestGetLabelErr() {
	s.mock.ExpectQuery("SELECT").
		WillReturnError(s.bdError)

	_, err := s.repository.GetLabel(1)

	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestGetCredNo() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, _, err := s.repository.GetCredentials("qwerty")

	require.Error(s.T(), err)
	require.Equal(s.T(), err, profile.ErrUserNonExistent)
}

func (s *Suite) TestGetCredError() {
	s.mock.ExpectQuery("SELECT").
		WillReturnError(s.bdError)

	_, _, err := s.repository.GetCredentials("qwerty")

	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestGetSubs() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, err := s.repository.GetUserSubscriptions(profile.FilterParams{
		PrevId:      1,
		CountLimit:  1,
		ReqAuthorId: -1,
	})

	require.NoError(s.T(), err)
}

func (s *Suite) TestCheckUserSubscription() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, err := s.repository.CheckUserSubscription(1, 2)

	require.NoError(s.T(), err)
}

func (s *Suite) TestCreateSubscription() {
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("INSERT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err := s.repository.CreateSubscription(1, 2)

	require.NoError(s.T(), err)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestMeetings(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGetProfile() {
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

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	_, err := s.repository.GetProfile(-1, 1337)

	require.NoError(s.T(), err)
}
