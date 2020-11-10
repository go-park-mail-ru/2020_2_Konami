package meeting

import "konami_backend/internal/pkg/models"

type Repository interface {
	CreateMeeting(meeting models.Meeting) (meetingId int, err error)
	GetMeeting(meetingId, userId int) (models.Meeting, error)
	UpdateMeeting(userId int, update models.MeetingUpdate) error
	GetAll() ([]models.MeetingCard, error)
	FilterToday() ([]models.MeetingCard, error)
	FilterTomorrow() ([]models.MeetingCard, error)
	FilterFuture() ([]models.MeetingCard, error)
	FilterLiked(userId int) ([]models.MeetingCard, error)
	FilterRegistered(userId int) ([]models.MeetingCard, error)
}
