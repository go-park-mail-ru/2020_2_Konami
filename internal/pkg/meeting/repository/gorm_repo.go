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
	Regs      []Registration
	Likes     []Like
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
	layout := "2006-01-02 15:04:05"
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
		StartDate: obj.StartDate.Format("2006-01-02 15:04:05"),
		EndDate:   obj.EndDate.Format("2006-01-02 15:04:05"),
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

func likeIndex(s []Like, target int) int {
	for i, el := range s {
		if el.UserId == target {
			return i
		}
	}
	return -1
}

func regIndex(s []Registration, target int) int {
	for i, el := range s {
		if el.UserId == target {
			return i
		}
	}
	return -1
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
	var m Meeting
	db := h.db.First(&m).
		Where("id = ?", meetingId)

	err := db.Error
	if err != nil {
		return models.Meeting{}, err
	}
	card := ToMeetingCard(m)
	res := models.Meeting{Card: &card}
	if authorized && likeIndex(m.Likes, userId) != -1 {
		res.Like = true
	}
	if authorized && regIndex(m.Regs, userId) != -1 {
		res.Reg = true
	}
	return res, nil
}

func (h *MeetingGormRepo) UpdateMeeting(userId int, update models.MeetingUpdate) error {
	m := Meeting{}
	db := h.db.First(&m).
		Where("id = ?", update.MeetId)

	err := db.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return meeting.ErrMeetingNotFound
	}
	if err != nil {
		return err
	}
	if update.Fields.Like != nil && *update.Fields.Like == true {
		like := Like{MeetingId: m.Id, UserId: userId}
		if likeIndex(m.Likes, userId) == -1 {
			m.Likes = append(m.Likes, like)
		}
	} else if update.Fields.Like != nil && *update.Fields.Like == false {
		li := likeIndex(m.Likes, userId)
		if li != -1 {
			m.Likes = append(m.Likes[:li], m.Likes[li+1:]...)
		}
	}
	if update.Fields.Reg != nil && *update.Fields.Reg == true {
		reg := Registration{MeetingId: m.Id, UserId: userId}
		if regIndex(m.Regs, userId) == -1 {
			m.Regs = append(m.Regs, reg)
		}
	} else if update.Fields.Reg != nil && *update.Fields.Reg == false {
		reg := regIndex(m.Regs, userId)
		if reg != -1 {
			m.Regs = append(m.Regs[:reg], m.Regs[reg+1:]...)
		}
	}
	db = h.db.Save(&m)
	return db.Error
}

func (h *MeetingGormRepo) GetAll() ([]models.MeetingCard, error) {
	var meetings []Meeting
	db := h.db.Find(&meetings)
	err := db.Error
	if err != nil {
		return []models.MeetingCard{}, err
	}
	result := make([]models.MeetingCard, len(meetings))
	for i, el := range meetings {
		result[i] = ToMeetingCard(el)
	}
	return result, nil
}

func (h *MeetingGormRepo) FilterDate(dt time.Time) ([]models.MeetingCard, error) {
	dateSlice := strings.Split(dt.Format("2006-01-02"), "-")
	var meetings []Meeting
	db := h.db.Find(&meetings).
		Where("DATEPART(yy, StartDate) = ?", dateSlice[0]).
		Where("DATEPART(mm, StartDate) = ?", dateSlice[1]).
		Where("DATEPART(dd, StartDate) = ?", dateSlice[2]).
		Order("StartDate ASC")

	err := db.Error
	if err != nil {
		return []models.MeetingCard{}, err
	}
	result := make([]models.MeetingCard, len(meetings))
	for i, el := range meetings {
		result[i] = ToMeetingCard(el)
	}
	return result, nil
}

func (h *MeetingGormRepo) FilterToday() ([]models.MeetingCard, error) {
	today := time.Now()
	return h.FilterDate(today)
}

func (h *MeetingGormRepo) FilterTomorrow() ([]models.MeetingCard, error) {
	tomorrow := time.Now().Add(24 * time.Hour)
	return h.FilterDate(tomorrow)
}

func (h *MeetingGormRepo) FilterFuture() ([]models.MeetingCard, error) {
	dateSlice := strings.Split(time.Now().Format("2006-01-02"), "-")
	var meetings []Meeting
	db := h.db.Find(&meetings).
		Where("DATEPART(yy, StartDate) >= ?", dateSlice[0]).
		Where("DATEPART(mm, StartDate) >= ?", dateSlice[1]).
		Where("DATEPART(dd, StartDate) >= ?", dateSlice[2]).
		Order("StartDate ASC")

	err := db.Error
	if err != nil {
		return []models.MeetingCard{}, err
	}
	result := make([]models.MeetingCard, len(meetings))
	for i, el := range meetings {
		result[i] = ToMeetingCard(el)
	}
	return result, nil
}

func (h *MeetingGormRepo) FilterLiked(userId int) ([]models.MeetingCard, error) {
	var likes []Like
	db := h.db.Find(&likes).
		Where("UserId = ?", userId)

	if db.Error != nil {
		return nil, db.Error
	}
	result := make([]models.MeetingCard, len(likes))
	for i, el := range likes {
		db = h.db.
			First(&result[i]).
			Where("id = ?", el.MeetingId)

		if db.Error != nil {
			return nil, db.Error
		}
	}
	return result, nil
}

func (h *MeetingGormRepo) FilterRegistered(userId int) ([]models.MeetingCard, error) {
	var regs []Registration
	db := h.db.Find(&regs).
		Where("UserId = ?", userId)

	if db.Error != nil {
		return nil, db.Error
	}
	result := make([]models.MeetingCard, len(regs))
	for i, el := range regs {
		db = h.db.
			First(&result[i]).
			Where("id = ?", el.MeetingId)

		if db.Error != nil {
			return nil, db.Error
		}
	}
	return result, nil
}
