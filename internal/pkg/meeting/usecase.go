package meeting

import "konami_backend/internal/pkg/models"

type UseCase interface {
	CreateMeeting(authorId int, data models.MeetingData) (meetingId int, err error)
	GetMeeting(meetingId, userId int) (models.Meeting, error)
	UpdateMeeting(userId int, update models.MeetingUpdate) error
	GetAll() ([]models.MeetingCard, error)
	FilterToday() ([]models.MeetingCard, error)
	FilterTomorrow() ([]models.MeetingCard, error)
	FilterFuture() ([]models.MeetingCard, error)
	FilterLiked(userId int) ([]models.MeetingCard, error)
	FilterRegistered(userId int) ([]models.MeetingCard, error)
}
