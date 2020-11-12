//go:generate mockgen -source=repository.go -destination=./repositoty_mock.go -package=meeting
package meeting

import (
	"errors"
	"konami_backend/internal/pkg/models"
)

var ErrMeetingNotFound = errors.New("meeting not found")

type Repository interface {
	CreateMeeting(meeting models.Meeting) (meetingId int, err error)
	GetMeeting(meetingId, userId int, authorized bool) (models.Meeting, error)
	UpdateMeeting(userId int, update models.MeetingUpdate) error
	GetAll(userId int) ([]models.Meeting, error)
	FilterToday(userId int) ([]models.Meeting, error)
	FilterTomorrow(userId int) ([]models.Meeting, error)
	FilterFuture(userId int) ([]models.Meeting, error)
	FilterLiked(userId int) ([]models.Meeting, error)
	FilterRegistered(userId int) ([]models.Meeting, error)
}
