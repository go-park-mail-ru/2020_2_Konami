package meeting

import "konami_backend/internal/pkg/models"

type UseCase interface {
	CreateMeeting(authorId int, data models.MeetingData) (meetingId int, err error)
	GetMeeting(meetingId, userId int, authorized bool) (models.MeetingDetails, error)
	UpdateMeeting(userId int, update models.MeetingUpdate) error
	GetNextMeetings(params FilterParams) ([]models.Meeting, error)
	GetTopMeetings(params FilterParams) ([]models.Meeting, error)
	FilterLiked(params FilterParams) ([]models.Meeting, error)
	FilterRegistered(params FilterParams) ([]models.Meeting, error)
	FilterRecommended(params FilterParams) ([]models.Meeting, error)
	FilterTagged(params FilterParams, tagId int) ([]models.Meeting, error)
	FilterSimilar(params FilterParams, meetingId int) ([]models.Meeting, error)
}
