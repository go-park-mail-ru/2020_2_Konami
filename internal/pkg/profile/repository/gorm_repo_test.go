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

	_, err := s.repository.GetAll()

	require.NoError(s.T(), err)
}

func (s *Suite) TestGetAllError() {
	s.mock.ExpectQuery("SELECT").
		WillReturnError(s.bdError)

	_, err := s.repository.GetAll()

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

func (s *Suite) TestEditPhoto() {
	/*	s.mock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		s.mock.ExpectQuery("UPDATE").
			WithArgs(1,"","PICTURE","","0001-01-01 00:00:00 +0000 UTC","","","","","","","","","","",1)
		s.mock.ExpectCommit()

		err := s.repository.EditProfilePic(1, "PICTURE")

		require.NoError(s.T(), err)
	*/
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

	_, err := s.repository.GetProfile(1337)

	require.NoError(s.T(), err)
}

func (s *Suite) TestCreateProfile() {
	/*	testTags := []*models.Tag{
			{

			},
		}

		testMeeting := []*models.MeetingLabel{
			{
				Id:    1,
				Title: "gg",
				Cover: "gg",
			},
		}

		gg := models.Profile{
			Card:        &models.ProfileCard{
				Label:        &models.ProfileLabel{
					Id:     0,
					Name:   "",
					ImgSrc: "",
				},
				Job:          "",
				InterestTags: nil,
				SkillTags:    nil,
			},
			Gender:      "",
			Birthday:    "",
			City:        "",
			Login:       "",
			PwdHash:     "",
			Telegram:    "",
			Vk:          "",
			Education:   "",
			MeetingTags: testTags,
			Aims:        "",
			Interests:   "",
			Skills:      "",
			Meetings:    testMeeting,
		}

		s.mock.ExpectBegin()
		s.mock.ExpectQuery("INSERT INTO").
			WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(1))

		s.mock.ExpectQuery("INSERT INTO").
			WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(1))

		s.mock.ExpectQuery("INSERT INTO").
			WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(1))


		s.mock.ExpectCommit()


		_, err := s.repository.Create(gg)

		require.NoError(s.T(), err)*/
}

/*
☨☨☨ EditProfile(update models.Profile) error
☨☨☨ EditProfilePic(userId int, imgSrc string) error
☨☨☨ Create(p models.Profile) (userId int, err error)
*/
