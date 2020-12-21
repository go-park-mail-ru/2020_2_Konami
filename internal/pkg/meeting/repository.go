//go:generate mockgen -source=repository.go -destination=./repositoty_mock.go -package=meeting
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
	PrevLikes  int
	PrevStart  time.Time
	CountLimit int
	UserId     int
}

type Repository interface {
	CreateMeeting(meeting models.Meeting) (meetingId int, err error)
	GetMeeting(meetingId, userId int, authorized bool) (models.MeetingDetails, error)
	SetLike(meetId int, userId int) error
	RemoveLike(meetId int, userId int) error
	SetReg(meetId int, userId int) error
	RemoveReg(meetId int, userId int) error
	UpdateMeeting(update models.MeetingCard) error
	GetNextMeetings(params FilterParams) ([]models.Meeting, error)
	GetTopMeetings(params FilterParams) ([]models.Meeting, error)
	FilterLiked(params FilterParams) ([]models.Meeting, error)
	FilterSubsLiked(params FilterParams) ([]models.Meeting, error)
	FilterRegistered(params FilterParams) ([]models.Meeting, error)
	FilterSubsRegistered(params FilterParams) ([]models.Meeting, error)
	FilterRecommended(params FilterParams) ([]models.Meeting, error)
	FilterTagged(params FilterParams, tags []string) ([]models.Meeting, error)
	FilterSimilar(params FilterParams, meetingId int) ([]models.Meeting, error)
	SearchMeetings(params FilterParams, meetingName string, limit int) ([]models.Meeting, error)
}
