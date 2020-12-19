package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/tag"
	"regexp"
	"testing"
)

type Suite struct {
	suite.Suite
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	tags       []models.Tag
	repository tag.Repository
	bdError    error
}

func (s *Suite) SetupSuite() {
	var db *sql.DB
	var err error

	t1 := models.Tag{
		TagId: 1,
		Name:  "LOL",
	}
	s.tags = []models.Tag{t1}

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

	s.repository = NewTagGormRepo(s.DB)
}

func (s *Suite) TestGetTagById() {
	testId := 1
	testTag := s.tags[0]

	s.mock.ExpectQuery("SELECT").
		WithArgs(testId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(testTag.TagId, testTag.Name))

	res, err := s.repository.GetTagById(testId)

	require.NoError(s.T(), err)

	require.Nil(s.T(), deep.Equal(models.Tag{
		TagId: 1,
		Name:  "LOL",
	}, res))
}

func (s *Suite) TestGetTagByIdError() {
	testId := 1

	s.mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnError(s.bdError)

	_, err := s.repository.GetTagById(testId)
	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestGetTagByName() {
	testId := "LOL"
	testTag := s.tags[0]

	s.mock.ExpectQuery("SELECT").
		WithArgs(testId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(testTag.TagId, testTag.Name))

	res, err := s.repository.GetTagByName(testId)

	require.NoError(s.T(), err)

	require.Nil(s.T(), deep.Equal(models.Tag{
		TagId: 1,
		Name:  "LOL",
	}, res))
}

func (s *Suite) TestGetTagByNameError() {
	testId := "LOL"

	s.mock.ExpectQuery("SELECT").
		WithArgs(testId).
		WillReturnError(s.bdError)

	_, err := s.repository.GetTagByName(testId)

	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestCreateTag() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").
		WithArgs("LOL").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	_, err := s.repository.CreateTag("LOL")
	require.NoError(s.T(), err)
}

func (s *Suite) TestCreateTagError() {
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").
		WithArgs("LOL").
		WillReturnError(s.bdError)
	s.mock.ExpectRollback()

	_, err := s.repository.CreateTag("LOL")

	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestGetOrCreateTag() {
	testId := "LOL"
	testTag := s.tags[0]

	s.mock.ExpectQuery("SELECT").
		WithArgs(testId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(testTag.TagId, testTag.Name))

	res, err := s.repository.GetOrCreateTag(testId)

	require.NoError(s.T(), err)

	require.Nil(s.T(), deep.Equal(models.Tag{
		TagId: 1,
		Name:  "LOL",
	}, res))
}

func (s *Suite) TestGetOrCreateTagError() {
	testId := "LOL"

	s.mock.ExpectQuery("SELECT").
		WithArgs(testId).
		WillReturnError(s.bdError)

	_, err := s.repository.GetOrCreateTag(testId)

	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) TestGetOrCreateTagExists() {
	testId := "LOL"

	s.mock.ExpectQuery("SELECT").
		WithArgs(testId).
		WillReturnError(gorm.ErrRecordNotFound)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").
		WithArgs("LOL").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	_, err := s.repository.GetOrCreateTag("LOL")
	require.NoError(s.T(), err)
}

func (s *Suite) TestFilterTags() {
	testId := "LO"
	testTag := s.tags[0]

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tags`)).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(testTag.TagId, testTag.Name))

	_, err := s.repository.FilterTags(testId)

	require.NoError(s.T(), err)
}

func (s *Suite) TestFilterTagsError() {
	testId := "LO"

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tags`)).
		WithArgs().
		WillReturnError(s.bdError)

	_, err := s.repository.FilterTags(testId)

	require.Error(s.T(), err)
	require.Equal(s.T(), err, s.bdError)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestMeetings(t *testing.T) {
	suite.Run(t, new(Suite))
}

/*
Testing method
GetTagById(id int) (models.Tag, error)
GetTagByName(name string) (models.Tag, error)
CreateTag(name string) (models.Tag, error)
GetOrCreateTag(name string) (models.Tag, error)
FilterTags(startsWith string) ([]models.Tag, error)
*/
