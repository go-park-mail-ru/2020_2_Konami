package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/tag"
	"regexp"
	"testing"
)

type Suite struct {
	suite.Suite
	DB      *gorm.DB
	mock    sqlmock.Sqlmock
	tags []models.Tag

	repository tag.Repository
}

func (s *Suite) SetupSuite() {
	var db  *sql.DB
	var err error

	t1 := models.Tag{
		TagId: 1,
		Name:  "LOL",
	}
	s.tags = []models.Tag{t1}

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	s.DB.LogMode(true)

	s.repository = NewTagGormRepo(s.DB)
}

func (s *Suite) TestGetTagById() {
	testId := 1
	testTag := s.tags[0]

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tags" WHERE (id = $1)`)).
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

func (s *Suite) TestGetTagByName() {
	testId := "LOL"
	testTag := s.tags[0]

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tags" WHERE (UPPER(name) = $1)`)).
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

func (s *Suite) TestCreateTag() {
/*	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "tags" ("name") VALUES ($1) RETURNING "tags"."id""`)).
		WithArgs("LOL").
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectRollback()

	_, _ = s.repository.CreateTag("LOL")*/
	//require.NoError(s.T(), err)
}

func (s *Suite) TestGetOrCreateTag() {
	testId := "LOL"
	testTag := s.tags[0]

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tags" WHERE (UPPER(name) = $1)`)).
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
