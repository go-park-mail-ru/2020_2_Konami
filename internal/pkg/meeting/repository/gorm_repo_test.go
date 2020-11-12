package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"konami_backend/internal/pkg/meeting"
	"konami_backend/internal/pkg/models"
	"regexp"
	"testing"
)

type Suite struct {
	suite.Suite
	DB      *gorm.DB
	mock    sqlmock.Sqlmock
	meetings []models.Meeting
	meetingUpd models.MeetingUpdate

	repository meeting.Repository
}

func (s *Suite) SetupSuite() {
	var db  *sql.DB
	var err error

	t1 := models.Tag{
		TagId: 1,
		Name:  "LOL",
	}
	tags1 := []*models.Tag{&t1}

	ml1 := models.MeetingLabel{
		Id:    2,
		Title: "123",
		Cover: "123",
	}

	c1 := models.MeetingCard{
		Label:     &ml1,
		AuthorId:  1,
		Text:      "123",
		Tags:      tags1,
		Address:   "123",
		City:      "123",
		StartDate: "2020-11-11T23:19:00.000Z",
		EndDate:   "2020-11-12T23:19:00.000Z",
		Seats:     10,
		SeatsLeft: 20,
	}

	m1 := models.Meeting{
		Card: &c1,
		Like: false,
		Reg:  false,
	}

	s.meetings = []models.Meeting{m1}

	b := true
	f1 := models.MeetUpdateFields{
		Reg:  &b,
		Like: &b,
	}

	s.meetingUpd = models.MeetingUpdate{
		MeetId: 0,
		Fields: &f1,
	}

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	s.DB.LogMode(false)

	s.repository = NewMeetingGormRepo(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *Suite) TestCreateMeeting() {
	/*	s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "meetings" ("author_id","title","text","img_src","city","address","start_date","end_date","seats","seats_left") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "meetings"."id"`)).
			WithArgs(0, 0).
			WillReturnResult(driver.RowsAffected(1))
		s.mock.ExpectRollback()

		_, err := s.repository.CreateMeeting(s.meetings[0])
		require.NoError(s.T(), err)
	*/
}

func (s *Suite) TestGetMeetings() {
	superTest := ""

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meetings`)).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"Title"}).
			AddRow(superTest))

	_, err := s.repository.GetAll(1)

	require.NoError(s.T(), err)
}

func (s *Suite) TestGetMeeting() {
	testMeetingId := 1
	var m Meeting

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "meetings" WHERE (id = $1)`)).
		WithArgs(testMeetingId).
		WillReturnRows(sqlmock.NewRows([]string{"Id"}).
			AddRow(&m.Id))

	res, err := s.repository.GetMeeting(testMeetingId, 1, true)

	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(models.Meeting{
		Card: res.Card,
		Like: false,
		Reg:  false,
	}, res))
}

func (s *Suite) TestUpdateMeeting() {

}

func (s *Suite) TestFilterToday() {

}

func (s *Suite) TestFilterTomorrow() {

}

func (s *Suite) TestFilterFuture() {

}

func (s *Suite) TestFilterLiked() {

}

func (s *Suite) TestFilterRegistered() {

}

func TestMeetings(t *testing.T) {
	suite.Run(t, new(Suite))
}
/*
Testing method
CreateMeeting(meeting models.Meeting) (meetingId int, err error)
GetMeeting(meetingId, userId int, authorized bool) (models.Meeting, error)
UpdateMeeting(userId int, update models.MeetingUpdate) error
GetAll(userId int) ([]models.Meeting, error)
FilterToday(userId int) ([]models.Meeting, error)
FilterTomorrow(userId int) ([]models.Meeting, error)
FilterFuture(userId int) ([]models.Meeting, error)
FilterLiked(userId int) ([]models.Meeting, error)
FilterRegistered(userId int) ([]models.Meeting, error)
*/
