package usecase

import (
	"errors"
	"github.com/google/uuid"
	"konami_backend/internal/pkg/meeting"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/tag"
	"konami_backend/internal/pkg/utils/uploads_handler"
)

type MeetingUseCase struct {
	MeetRepo         meeting.Repository
	UploadsHandler   uploads_handler.UploadsHandler
	TagRepo          tag.Repository
	MeetingCoversDir string
	defaultImgSrc    string
}

func NewMeetingUseCase(MeetRepo meeting.Repository,
	UploadsHandler uploads_handler.UploadsHandler,
	TagRepo tag.Repository,
	MeetingCoversDir string,
	defaultImgSrc string) meeting.UseCase {

	return &MeetingUseCase{
		MeetRepo:         MeetRepo,
		UploadsHandler:   UploadsHandler,
		TagRepo:          TagRepo,
		MeetingCoversDir: MeetingCoversDir,
		defaultImgSrc:    defaultImgSrc,
	}
}

func (uc *MeetingUseCase) CreateMeeting(authorId int, data models.MeetingData) (int, error) {
	if data.Title == nil || data.Text == nil || data.Address == nil || data.City == nil ||
		data.Start == nil || data.End == nil || (data.Seats != nil && *data.Seats < 0) ||
		*data.End < *data.Start || *data.Title == "" {
		return 0, errors.New("invalid meeting data")
	}
	imgSrc := uc.defaultImgSrc
	var err error
	if data.Photo != nil {
		imgSrc = uc.MeetingCoversDir + "/" + uuid.New().String()
		imgSrc, err = uc.UploadsHandler.UploadBase64Image(imgSrc, data.Photo)
		if err != nil {
			return 0, err
		}
	}
	m := models.Meeting{
		Card: &models.MeetingCard{
			Label: &models.MeetingLabel{
				Title: *data.Title,
				Cover: imgSrc,
			},
			AuthorId:  authorId,
			Text:      *data.Text,
			Tags:      []*models.Tag{},
			City:      *data.City,
			Address:   *data.Address,
			StartDate: *data.Start,
			EndDate:   *data.End,
			Seats:     1000 * 1000 * 1000,
		},
	}
	if data.Seats != nil {
		m.Card.Seats = *data.Seats
	}
	m.Card.SeatsLeft = m.Card.Seats
	if data.Tags != nil {
		for _, tagName := range data.Tags {
			t, err := uc.TagRepo.GetOrCreateTag(tagName)
			if err != nil {
				return 0, err
			}
			m.Card.Tags = append(m.Card.Tags, &t)
		}
	}
	return uc.MeetRepo.CreateMeeting(m)
}

func (uc *MeetingUseCase) GetMeeting(meetingId, userId int, authorized bool) (models.Meeting, error) {
	return uc.MeetRepo.GetMeeting(meetingId, userId, authorized)
}

func (uc *MeetingUseCase) UpdateMeeting(userId int, update models.MeetingUpdate) error {
	return uc.MeetRepo.UpdateMeeting(userId, update)
}

func (uc *MeetingUseCase) GetAll() ([]models.MeetingCard, error) {
	return uc.MeetRepo.GetAll()
}

func (uc *MeetingUseCase) FilterToday() ([]models.MeetingCard, error) {
	return uc.MeetRepo.FilterToday()
}

func (uc *MeetingUseCase) FilterTomorrow() ([]models.MeetingCard, error) {
	return uc.MeetRepo.FilterTomorrow()
}

func (uc *MeetingUseCase) FilterFuture() ([]models.MeetingCard, error) {
	return uc.MeetRepo.FilterFuture()
}

func (uc *MeetingUseCase) FilterLiked(userId int) ([]models.MeetingCard, error) {
	return uc.MeetRepo.FilterLiked(userId)
}

func (uc *MeetingUseCase) FilterRegistered(userId int) ([]models.MeetingCard, error) {
	return uc.MeetRepo.FilterRegistered(userId)
}
