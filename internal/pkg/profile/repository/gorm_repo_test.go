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
	DB   *gorm.DB
	mock sqlmock.Sqlmock
	repository profile.Repository
	bdError error
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

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestMeetings(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGetProfile() {
	/* GOTCHAS
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "profiles" WHERE id = $1 ORDER BY "profiles"."id" LIMIT 1`)).
		WithArgs(1337).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1337))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "profile_interest_tags" WHERE "profile_interest_tags"."profile_id" = $1`)).
		WithArgs(1337).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	s.mock.ExpectQuery("SELECT").
		WithArgs(sql.NullInt64{}).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "profile_meeting_tags" WHERE "profile_meeting_tags"."profile_id" = $1`)).
		WithArgs(1337).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	/*s.mock.ExpectQuery("SELECT")
	s.mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	_, err := s.repository.GetProfile(1337)

	require.NoError(s.T(), err)
*/
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

	s.mock.ExpectCommit()


	_, err := s.repository.Create(gg)

	require.NoError(s.T(), err)
*/
}

/*
GetProfile(userId int) (models.Profile, error) GOTCHAS
EditProfile(update models.Profile) error
EditProfilePic(userId int, imgSrc string) error
Create(p models.Profile) (userId int, err error)
GetCredentials(login string) (userId int, pwdHash string, err error)
GetLabel(userId int) (models.ProfileLabel, error)
GetSubscriptions(userId int) (tagIds []int, err error)
*/