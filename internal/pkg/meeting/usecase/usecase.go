package usecase

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"konami_backend/internal/pkg/meeting"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/tag"
	"konami_backend/internal/pkg/utils/uploads_handler"
	"strings"
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

func (uc *MeetingUseCase) GetMeeting(meetingId, userId int, authorized bool) (models.MeetingDetails, error) {
	return uc.MeetRepo.GetMeeting(meetingId, userId, authorized)
}

func (uc *MeetingUseCase) UpdateMeeting(userId int, update models.MeetingUpdate) error {
	if update.Fields == nil {
		return errors.New("invalid update data")
	}
	m, err := uc.MeetRepo.GetMeeting(update.MeetId, -1, false)
	if err != nil {
		return errors.New("invalid meeting id")
	}
	if update.Fields.Card != nil && update.Fields.Card.Photo != nil {
		imgSrc := uc.MeetingCoversDir + "/" + uuid.New().String()
		if !strings.HasSuffix(m.Card.Label.Cover, uc.defaultImgSrc) {
			imgSrc = strings.TrimPrefix(m.Card.Label.Cover, uc.UploadsHandler.UploadsDir+"/")
		}
		m.Card.Label.Cover, err = uc.UploadsHandler.UploadBase64Image(imgSrc, update.Fields.Card.Photo)
		if err != nil {
			return err
		}
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return meeting.ErrMeetingNotFound
	}
	if err != nil {
		return err
	}
	if update.Fields.Like != nil && *update.Fields.Like {
		err = uc.MeetRepo.SetLike(update.MeetId, userId)
	} else if update.Fields.Like != nil && !*update.Fields.Like {
		err = uc.MeetRepo.RemoveLike(update.MeetId, userId)
	}
	if err != nil {
		return err
	}
	if update.Fields.Reg != nil && *update.Fields.Reg {
		err = uc.MeetRepo.SetReg(update.MeetId, userId)
	} else if update.Fields.Reg != nil && !*update.Fields.Reg {
		err = uc.MeetRepo.RemoveReg(update.MeetId, userId)
	}
	if update.Fields.Card == nil {
		return err
	}
	if update.Fields.Card.Address != nil {
		m.Card.Address = *update.Fields.Card.Address
	}
	if update.Fields.Card.City != nil {
		m.Card.City = *update.Fields.Card.City
	}
	if update.Fields.Card.Start != nil {
		m.Card.StartDate = *update.Fields.Card.Start
	}
	if update.Fields.Card.End != nil {
		m.Card.EndDate = *update.Fields.Card.End
	}
	if update.Fields.Card.Seats != nil {
		occupied := m.Card.Seats - m.Card.SeatsLeft
		m.Card.Seats = *update.Fields.Card.Seats
		m.Card.SeatsLeft = m.Card.Seats - occupied
	}
	if update.Fields.Card.Text != nil {
		m.Card.Text = *update.Fields.Card.Text
	}
	if update.Fields.Card.Title != nil {
		m.Card.Label.Title = *update.Fields.Card.Title
	}
	if update.Fields.Card.Tags != nil {
		m.Card.Tags = []*models.Tag{}
		for _, tagName := range update.Fields.Card.Tags {
			t, err := uc.TagRepo.GetOrCreateTag(tagName)
			if err != nil {
				return err
			}
			m.Card.Tags = append(m.Card.Tags, &t)
		}
	}
	return uc.MeetRepo.UpdateMeeting(*m.Card)
}

func (uc *MeetingUseCase) GetNextMeetings(params meeting.FilterParams) ([]models.Meeting, error) {
	return uc.MeetRepo.GetNextMeetings(params)
}

func (uc *MeetingUseCase) GetTopMeetings(params meeting.FilterParams) ([]models.Meeting, error) {
	return uc.MeetRepo.GetTopMeetings(params)
}

func (uc *MeetingUseCase) FilterLiked(params meeting.FilterParams) ([]models.Meeting, error) {
	return uc.MeetRepo.FilterLiked(params)
}

func (uc *MeetingUseCase) FilterRegistered(params meeting.FilterParams) ([]models.Meeting, error) {
	return uc.MeetRepo.FilterRegistered(params)
}

func (uc *MeetingUseCase) FilterSubsLiked(params meeting.FilterParams) ([]models.Meeting, error) {
	return uc.MeetRepo.FilterSubsLiked(params)
}

func (uc *MeetingUseCase) FilterSubsRegistered(params meeting.FilterParams) ([]models.Meeting, error) {
	return uc.MeetRepo.FilterSubsRegistered(params)
}

func (uc *MeetingUseCase) FilterRecommended(params meeting.FilterParams) ([]models.Meeting, error) {
	return uc.MeetRepo.FilterRecommended(params)
}

func (uc *MeetingUseCase) FilterTagged(params meeting.FilterParams, tagId int) ([]models.Meeting, error) {
	return uc.MeetRepo.FilterTagged(params, tagId)
}

func (uc *MeetingUseCase) FilterSimilar(params meeting.FilterParams, meetingId int) ([]models.Meeting, error) {
	return uc.MeetRepo.FilterSimilar(params, meetingId)
}

func (uc *MeetingUseCase) SearchMeetings(params meeting.FilterParams,
	meetingName string, limit int) ([]models.Meeting, error) {
	return uc.MeetRepo.SearchMeetings(params, meetingName, limit)
}
