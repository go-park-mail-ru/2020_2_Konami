//go:generate mockgen -source=usecase.go -destination=./usecase_mock.go -package=meeting
package meeting

import "konami_backend/internal/pkg/models"

type UseCase interface {
	CreateMeeting(authorId int, data models.MeetingData) (meetingId int, err error)
	GetMeeting(meetingId, userId int, authorized bool) (models.Meeting, error)
	UpdateMeeting(userId int, update models.MeetingUpdate) error
	GetAll(userId int) ([]models.Meeting, error)
	FilterToday(userId int) ([]models.Meeting, error)
	FilterTomorrow(userId int) ([]models.Meeting, error)
	FilterFuture(userId int) ([]models.Meeting, error)
	FilterLiked(userId int) ([]models.Meeting, error)
	FilterRegistered(userId int) ([]models.Meeting, error)
}
