package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"konami_backend/internal/pkg/meeting"
	"konami_backend/internal/pkg/models"
	tagRepo "konami_backend/internal/pkg/tag/repository"
	"strings"
	"time"
)

type MeetingGormRepo struct {
	db *gorm.DB
}

func NewMeetingGormRepo(db *gorm.DB) meeting.Repository {
	return &MeetingGormRepo{db: db}
}

type Meeting struct {
	Id        int `gorm:"primaryKey;autoIncrement;"`
	AuthorId  int
	Title     string
	Text      string
	ImgSrc    string
	Tags      []tagRepo.Tag `gorm:"many2many:meeting_tags;"`
	City      string
	Address   string
	StartDate time.Time
	EndDate   time.Time
	Seats     int
	SeatsLeft int
	Regs      []Registration `gorm:"foreignKey:MeetingId"`
	Likes     []Like         `gorm:"foreignKey:MeetingId"`
}

type Registration struct {
	Id        int `gorm:"primaryKey;autoIncrement;"`
	MeetingId int
	UserId    int
}

type Like struct {
	Id        int `gorm:"primaryKey;autoIncrement;"`
	MeetingId int
	UserId    int
}

func (m *Meeting) TableName() string {
	return "meetings"
}

func (r *Registration) TableName() string {
	return "registrations"
}

func (l *Like) TableName() string {
	return "likes"
}

func ToDbObject(data models.MeetingCard) (Meeting, error) {
	m := Meeting{
		AuthorId:  data.AuthorId,
		Title:     data.Label.Title,
		Text:      data.Text,
		ImgSrc:    data.Label.Cover,
		City:      data.City,
		Address:   data.Address,
		Seats:     data.Seats,
		SeatsLeft: data.SeatsLeft,
	}
	m.Tags = make([]tagRepo.Tag, len(data.Tags))
	for i, val := range data.Tags {
		tag := tagRepo.ToDbObject(*val)
		m.Tags[i] = tag
	}
	layout := "2006-01-02T15:04:05.000Z0700"
	var errSt, errEnd error
	m.StartDate, errSt = time.Parse(layout, data.StartDate)
	m.EndDate, errEnd = time.Parse(layout, data.EndDate)
	if errSt != nil || errEnd != nil {
		return Meeting{}, errors.New("invalid datetime format")
	}
	return m, nil
}

func ToMeetingLabel(obj Meeting) models.MeetingLabel {
	m := models.MeetingLabel{
		Id:    obj.Id,
		Title: obj.Title,
		Cover: obj.ImgSrc,
	}
	return m
}

func ToMeetingCard(obj Meeting) models.MeetingCard {
	label := ToMeetingLabel(obj)
	m := models.MeetingCard{
		Label:     &label,
		AuthorId:  obj.AuthorId,
		Text:      obj.Text,
		Address:   obj.Address,
		City:      obj.City,
		StartDate: obj.StartDate.Format("2006-01-02T15:04:05.000Z0700"),
		EndDate:   obj.EndDate.Format("2006-01-02T15:04:05.000Z0700"),
		Seats:     obj.Seats,
		SeatsLeft: obj.SeatsLeft,
	}
	m.Tags = make([]*models.Tag, len(obj.Tags))
	for i, val := range obj.Tags {
		tag := tagRepo.ToModel(val)
		m.Tags[i] = &tag
	}
	return m
}

func (h *MeetingGormRepo) ToMeeting(obj Meeting, userId int) models.Meeting {
	card := ToMeetingCard(obj)
	m := models.Meeting{Card: &card}
	if userId != -1 && h.LikeExists(obj.Id, userId) {
		m.Like = true
	}
	if userId != -1 && h.RegExists(obj.Id, userId) {
		m.Reg = true
	}
	return m
}

func (h *MeetingGormRepo) LikeExists(meetId int, userId int) bool {
	var l Like
	db := h.db.
		Where("meeting_id = ?", meetId).
		Where("user_id = ?", userId).
		First(&l)
	if db.Error == nil {
		return true
	}
	return false
}

func (h *MeetingGormRepo) SetLike(meetId int, userId int) error {
	var l Like
	db := h.db.
		Where("meeting_id = ?", meetId).
		Where("user_id = ?", userId).
		First(&l)
	if db.Error != nil {
		l = Like{
			MeetingId: meetId,
			UserId:    userId,
		}
		db = h.db.Create(&l)
	}
	return db.Error
}

func (h *MeetingGormRepo) RemoveLike(meetId int, userId int) error {
	db := h.db.
		Where("meeting_id = ?", meetId).
		Where("user_id = ?", userId).
		Delete(Like{})
	return db.Error
}

func (h *MeetingGormRepo) RegExists(meetId int, userId int) bool {
	var l Registration
	db := h.db.
		Where("meeting_id = ?", meetId).
		Where("user_id = ?", userId).
		First(&l)
	if db.Error == nil {
		return true
	}
	return false
}

func (h *MeetingGormRepo) SetReg(meetId int, userId int) error {
	var l Registration
	db := h.db.
		Where("meeting_id = ?", meetId).
		Where("user_id = ?", userId).
		First(&l)
	if db.Error != nil {
		l = Registration{
			MeetingId: meetId,
			UserId:    userId,
		}
		db = h.db.Create(&l)
	}
	return db.Error
}

func (h *MeetingGormRepo) RemoveReg(meetId int, userId int) error {
	db := h.db.
		Where("meeting_id = ?", meetId).
		Where("user_id = ?", userId).
		Delete(Registration{})
	return db.Error
}

func (h *MeetingGormRepo) CreateMeeting(data models.Meeting) (int, error) {
	m, err := ToDbObject(*data.Card)
	if err != nil {
		return 0, err
	}
	db := h.db.Create(&m)
	err = db.Error
	if err != nil {
		return 0, err
	}
	return m.Id, nil
}

func (h *MeetingGormRepo) GetMeeting(meetingId, userId int, authorized bool) (models.Meeting, error) {
	if !authorized {
		userId = -1
	}
	var m Meeting
	db := h.db.
		Where("id = ?", meetingId).
		First(&m)

	err := db.Error
	if err != nil {
		return models.Meeting{}, err
	}
	res := h.ToMeeting(m, userId)
	return res, nil
}

func (h *MeetingGormRepo) UpdateMeeting(userId int, update models.MeetingUpdate) error {
	_, err := h.GetMeeting(update.MeetId, -1, false)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return meeting.ErrMeetingNotFound
	}
	if err != nil {
		return err
	}
	if update.Fields.Like != nil && *update.Fields.Like == true {
		err = h.SetLike(update.MeetId, userId)
	} else if update.Fields.Like != nil && *update.Fields.Like == false {
		err = h.RemoveLike(update.MeetId, userId)
	}
	if err != nil {
		return err
	}
	if update.Fields.Reg != nil && *update.Fields.Reg == true {
		err = h.SetReg(update.MeetId, userId)
	} else if update.Fields.Reg != nil && *update.Fields.Reg == false {
		err = h.RemoveReg(update.MeetId, userId)
	}
	return err
}

func (h *MeetingGormRepo) GetAll(userId int) ([]models.Meeting, error) {
	var meetings []Meeting
	db := h.db.Find(&meetings)
	err := db.Error
	if err != nil {
		return []models.Meeting{}, err
	}
	result := make([]models.Meeting, len(meetings))
	for i, el := range meetings {
		result[i] = h.ToMeeting(el, userId)
	}
	return result, nil
}

func (h *MeetingGormRepo) FilterDate(dt time.Time, userId int) ([]models.Meeting, error) {
	dateSlice := strings.Split(dt.Format("2006-01-02"), "-")
	var meetings []Meeting
	db := h.db.Find(&meetings).
		Where("DATEPART(yy, StartDate) = ?", dateSlice[0]).
		Where("DATEPART(mm, StartDate) = ?", dateSlice[1]).
		Where("DATEPART(dd, StartDate) = ?", dateSlice[2]).
		Order("StartDate ASC")

	err := db.Error
	if err != nil {
		return []models.Meeting{}, err
	}
	result := make([]models.Meeting, len(meetings))
	for i, el := range meetings {
		result[i] = h.ToMeeting(el, userId)
	}
	return result, nil
}

func (h *MeetingGormRepo) FilterToday(userId int) ([]models.Meeting, error) {
	today := time.Now()
	return h.FilterDate(today, userId)
}

func (h *MeetingGormRepo) FilterTomorrow(userId int) ([]models.Meeting, error) {
	tomorrow := time.Now().Add(24 * time.Hour)
	return h.FilterDate(tomorrow, userId)
}

func (h *MeetingGormRepo) FilterFuture(userId int) ([]models.Meeting, error) {
	dateSlice := strings.Split(time.Now().Format("2006-01-02"), "-")
	var meetings []Meeting
	db := h.db.Find(&meetings).
		Where("DATEPART(yy, StartDate) >= ?", dateSlice[0]).
		Where("DATEPART(mm, StartDate) >= ?", dateSlice[1]).
		Where("DATEPART(dd, StartDate) >= ?", dateSlice[2]).
		Order("StartDate ASC")

	err := db.Error
	if err != nil {
		return []models.Meeting{}, err
	}
	result := make([]models.Meeting, len(meetings))
	for i, el := range meetings {
		result[i] = h.ToMeeting(el, userId)
	}
	return result, nil
}

func (h *MeetingGormRepo) FilterLiked(userId int) ([]models.Meeting, error) {
	var likes []Like
	db := h.db.Find(&likes).
		Where("UserId = ?", userId)

	if db.Error != nil {
		return nil, db.Error
	}
	result := make([]models.Meeting, len(likes))
	m := Meeting{}
	for i, el := range likes {
		db = h.db.
			Where("id = ?", el.MeetingId).
			First(&m)
		if db.Error != nil {
			return nil, db.Error
		}
		result[i] = h.ToMeeting(m, userId)
	}
	return result, nil
}

func (h *MeetingGormRepo) FilterRegistered(userId int) ([]models.Meeting, error) {
	var regs []Registration
	db := h.db.Find(&regs).
		Where("UserId = ?", userId)

	if db.Error != nil {
		return nil, db.Error
	}
	result := make([]models.Meeting, len(regs))
	m := Meeting{}
	for i, el := range regs {
		db = h.db.
			Where("id = ?", el.MeetingId).
			First(&m)

		if db.Error != nil {
			return nil, db.Error
		}
		result[i] = h.ToMeeting(m, userId)
	}
	return result, nil
}
