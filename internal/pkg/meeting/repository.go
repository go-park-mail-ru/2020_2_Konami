package meeting

import (
	"errors"
	"konami_backend/internal/pkg/models"
	"time"
)

var ErrMeetingNotFound = errors.New("meeting not found")
var ErrNoSeatsLeft = errors.New("no meeting seats left")

type FilterParams struct {
	StartDate  time.Time
	EndDate    time.Time
	PrevId     int
	CountLimit int
	UserId     int
}

type Repository interface {
	CreateMeeting(meeting models.Meeting) (meetingId int, err error)
	GetMeeting(meetingId, userId int, authorized bool) (models.MeetingDetails, error)
	UpdateMeeting(userId int, update models.MeetingUpdate) error
	GetNextMeetings(params FilterParams) ([]models.Meeting, error)
	GetTopMeetings(params FilterParams) ([]models.Meeting, error)
	FilterLiked(params FilterParams) ([]models.Meeting, error)
	FilterRegistered(params FilterParams) ([]models.Meeting, error)
	FilterRecommended(params FilterParams) ([]models.Meeting, error)
	FilterTagged(params FilterParams, tagId int) ([]models.Meeting, error)
	FilterSimilar(params FilterParams, meetingId int) ([]models.Meeting, error)

	SearchMeetings(params FilterParams, meetingName string) ([]models.Meeting, error)
}
