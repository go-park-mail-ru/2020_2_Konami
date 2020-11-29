package repository

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"konami_backend/internal/pkg/meeting"
	"konami_backend/internal/pkg/models"
	profileRepo "konami_backend/internal/pkg/profile/repository"
	tagRepo "konami_backend/internal/pkg/tag/repository"
	"time"
)

type MeetingGormRepo struct {
	db *gorm.DB
}

func NewMeetingGormRepo(db *gorm.DB) meeting.Repository {
	return &MeetingGormRepo{db: db}
}

type Meeting struct {
	Id         int `gorm:"primaryKey;autoIncrement;"`
	AuthorId   int
	Title      string
	Text       string
	ImgSrc     string
	Tags       []tagRepo.Tag `gorm:"many2many:meeting_tags;"`
	City       string
	Address    string
	StartDate  time.Time
	EndDate    time.Time
	Seats      int
	SeatsLeft  int
	LikesCount int
	Regs       []Registration `gorm:"foreignKey:MeetingId"`
	Likes      []Like         `gorm:"foreignKey:MeetingId"`
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
		AuthorId:   data.AuthorId,
		Title:      data.Label.Title,
		Text:       data.Text,
		ImgSrc:     data.Label.Cover,
		City:       data.City,
		Address:    data.Address,
		Seats:      data.Seats,
		SeatsLeft:  data.SeatsLeft,
		LikesCount: data.LikesCount,
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
		Label:      &label,
		AuthorId:   obj.AuthorId,
		Text:       obj.Text,
		Address:    obj.Address,
		City:       obj.City,
		StartDate:  obj.StartDate.Format("2006-01-02T15:04:05.000Z0700"),
		EndDate:    obj.EndDate.Format("2006-01-02T15:04:05.000Z0700"),
		Seats:      obj.Seats,
		SeatsLeft:  obj.SeatsLeft,
		LikesCount: obj.LikesCount,
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

func (h *MeetingGormRepo) ToMeetingDetails(obj Meeting, userId int) (models.MeetingDetails, error) {
	card := ToMeetingCard(obj)
	m := models.MeetingDetails{Card: &card}
	if userId != -1 && h.LikeExists(obj.Id, userId) {
		m.Like = true
	}
	if userId != -1 && h.RegExists(obj.Id, userId) {
		m.Reg = true
	}
	var err error
	m.Registrations, err = h.GetRegistrations(obj)
	return m, err
}

func (h *MeetingGormRepo) ToMeetingList(meetings []Meeting, userId int) ([]models.Meeting, error) {
	result := make([]models.Meeting, len(meetings))
	for i, el := range meetings {
		result[i] = h.ToMeeting(el, userId)
	}
	return result, nil
}

func (h *MeetingGormRepo) GetRegistrations(m Meeting) ([]*models.ProfileLabel, error) {
	var p profileRepo.Profile
	res := make([]*models.ProfileLabel, len(m.Regs))
	for i, reg := range m.Regs {
		db := h.db.
			Where("id = ?", reg.UserId).
			First(&p)
		err := db.Error
		if err != nil {
			return nil, err
		}
		res[i] = &models.ProfileLabel{
			Id:     p.Id,
			Name:   p.Name,
			ImgSrc: p.ImgSrc,
		}
	}
	return res, nil
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
	newLike := false
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
		if db.Error == nil {
			newLike = true
		}
	}
	if !newLike {
		return db.Error
	}
	var m Meeting
	db = h.db.
		Where("id = ?", meetId).
		First(&m)
	if db.Error == nil {
		m.LikesCount += 1
		db = h.db.Save(m)
	}
	return db.Error
}

func (h *MeetingGormRepo) RemoveLike(meetId int, userId int) error {
	db := h.db.
		Where("meeting_id = ?", meetId).
		Where("user_id = ?", userId).
		Delete(Like{})
	if db.Error == nil && db.RowsAffected > 0 {
		var m Meeting
		db = h.db.
			Where("id = ?", meetId).
			First(&m)
		if db.Error == nil {
			m.LikesCount -= 1
			db = h.db.Save(m)
		}
	}
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
	var m Meeting
	db := h.db.
		Where("id = ?", meetId).
		First(&m)
	if db.Error != nil {
		return db.Error
	}
	if m.SeatsLeft == 0 {
		return meeting.ErrNoSeatsLeft
	}
	newReg := false
	db = h.db.
		Where("meeting_id = ?", meetId).
		Where("user_id = ?", userId).
		First(&l)
	if db.Error != nil {
		l = Registration{
			MeetingId: meetId,
			UserId:    userId,
		}
		db = h.db.Create(&l)
		newReg = db.Error == nil
	}
	if newReg {
		m.SeatsLeft -= 1
		db = h.db.Save(m)
	}
	return db.Error
}

func (h *MeetingGormRepo) RemoveReg(meetId int, userId int) error {
	db := h.db.
		Where("meeting_id = ?", meetId).
		Where("user_id = ?", userId).
		Delete(Registration{})
	if db.Error == nil && db.RowsAffected > 0 {
		var m Meeting
		db = h.db.
			Where("id = ?", meetId).
			First(&m)
		if db.Error == nil {
			m.SeatsLeft += 1
			db = h.db.Save(m)
		}
	}
	return db.Error
}

func (h *MeetingGormRepo) GetQuery(dest *[]Meeting, params meeting.FilterParams) *gorm.DB {
	return h.db.Find(dest).
		Where("Id > ?", params.PrevId).
		Where("StartDate >= ? ", params.StartDate).
		Where("EndDate <= ?", params.EndDate).
		Limit(params.CountLimit)
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

func (h *MeetingGormRepo) GetMeeting(meetingId, userId int, authorized bool) (models.MeetingDetails, error) {
	if !authorized {
		userId = -1
	}
	var m Meeting
	db := h.db.
		Preload("Tags").
		Preload("Regs").
		Where("id = ?", meetingId).
		First(&m)
	err := db.Error
	if err != nil {
		return models.MeetingDetails{}, err
	}
	return h.ToMeetingDetails(m, userId)
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

func (h *MeetingGormRepo) GetNextMeetings(params meeting.FilterParams) ([]models.Meeting, error) {
	var meetings []Meeting
	db := h.GetQuery(&meetings, params).Order("StartDate ASC")
	err := db.Error
	if err != nil {
		return []models.Meeting{}, err
	}
	return h.ToMeetingList(meetings, params.UserId)
}

func (h *MeetingGormRepo) GetTopMeetings(params meeting.FilterParams) ([]models.Meeting, error) {
	var meetings []Meeting
	db := h.GetQuery(&meetings, params).Order("LikesCount DESC")
	err := db.Error
	if err != nil {
		return []models.Meeting{}, err
	}
	return h.ToMeetingList(meetings, params.UserId)
}

func (h *MeetingGormRepo) FilterLiked(params meeting.FilterParams) ([]models.Meeting, error) {
	var likes []Like
	db := h.db.Find(&likes).
		Where("UserId = ?", params.UserId).
		Where("MeetingId > ?", params.PrevId).
		Order("MeetingId ASC").
		Limit(params.CountLimit)

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
		result[i] = h.ToMeeting(m, params.UserId)
	}
	return result, nil
}

func (h *MeetingGormRepo) FilterRegistered(params meeting.FilterParams) ([]models.Meeting, error) {
	var regs []Registration
	db := h.db.Find(&regs).
		Where("UserId = ?", params.UserId).
		Where("MeetingId > ?", params.PrevId).
		Order("MeetingId ASC").
		Limit(params.CountLimit)

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
		result[i] = h.ToMeeting(m, params.UserId)
	}
	return result, nil
}

func (h *MeetingGormRepo) ExtractMeetingsFromRows(params meeting.FilterParams, rows *sql.Rows) ([]Meeting, error) {
	meetings := []Meeting{}
	total := 0
	var meetBuf Meeting
	var meetId int
	var err error
	for err == nil && rows.Next() && total < params.CountLimit {
		err = rows.Scan(&meetId)
		db := h.db.
			Where("id = ?", meetId).
			First(&meetBuf)
		if db.Error == nil && (meetBuf.StartDate.Before(params.StartDate) || meetBuf.EndDate.After(params.EndDate)) {
			continue
		}
		meetings = append(meetings, meetBuf)
		total += 1
		err = db.Error
	}
	return meetings, err
}

func (h *MeetingGormRepo) FilterRecommended(params meeting.FilterParams) ([]models.Meeting, error) {
	var userProfile profileRepo.Profile
	db := h.db.
		Where("id = ?", params.UserId).
		Preload("MeetingTags").
		First(&userProfile)
	err := db.Error
	if err != nil {
		return nil, err
	}
	tagIds := []int{} // Tags to which user is subscribed
	for _, tag := range userProfile.MeetingTags {
		tagIds = append(tagIds, tag.Id)
	}
	// Meetings with whose tags
	rows, err := h.db.Table("meeting_tags").
		Where("tag_id IN ?", tagIds).
		Where("meeting_id > ?", params.PrevId).
		Order("meeting_id ASC").
		Select("meeting_id").Rows()
	defer rows.Close()
	meetings, err := h.ExtractMeetingsFromRows(params, rows)
	if err == nil {
		return h.ToMeetingList(meetings, params.UserId)
	}
	return nil, err
}

func (h *MeetingGormRepo) FilterTagged(params meeting.FilterParams, tagId int) ([]models.Meeting, error) {
	rows, err := h.db.Table("meeting_tags").
		Where("tag_id = ?", tagId).
		Where("meeting_id > ?", params.PrevId).
		Order("meeting_id ASC").
		Select("meeting_id").Rows()
	defer rows.Close()
	meetings, err := h.ExtractMeetingsFromRows(params, rows)
	if err == nil {
		return h.ToMeetingList(meetings, params.UserId)
	}
	return nil, err
}

func (h *MeetingGormRepo) FilterSimilar(params meeting.FilterParams, meetingId int) ([]models.Meeting, error) {
	var m Meeting
	db := h.db.Preload("Tags").
		Where("id = ?", meetingId).
		First(&m)
	err := db.Error
	if err != nil {
		return nil, err
	}
	tagIds := []int{}
	for _, tag := range m.Tags {
		tagIds = append(tagIds, tag.Id)
	}
	rows, err := h.db.Table("meeting_tags").
		Where("tag_id IN ?", tagIds).
		Where("meeting_id > ?", params.PrevId).
		Order("meeting_id ASC").
		Select("meeting_id").Rows()
	defer rows.Close()
	meetings, err := h.ExtractMeetingsFromRows(params, rows)
	if err == nil {
		return h.ToMeetingList(meetings, params.UserId)
	}
	return nil, err
}

func (h *MeetingGormRepo) SearchMeetings(params meeting.FilterParams, meetingName string) ([]models.Meeting, error) {
	var res []Meeting
	err := h.db.Table("meetings").
		Where("title @@ to_tsquery(?)", meetingName).
		Or("text @@ to_tsquery(?)", meetingName).
		Find(&res).Error

	if err != nil {
		return nil, err
	}

	return h.ToMeetingList(res, params.UserId)
}
