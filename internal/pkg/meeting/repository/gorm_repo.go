package repository

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"konami_backend/internal/pkg/meeting"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/profile"
	tagRepo "konami_backend/internal/pkg/tag/repository"
	"regexp"
	"time"
)

type MeetingGormRepo struct {
	db       *gorm.DB
	profRepo profile.Repository
}

func NewMeetingGormRepo(db *gorm.DB, profileRepo profile.Repository) meeting.Repository {
	return &MeetingGormRepo{db: db, profRepo: profileRepo}
}

func NewMeetingGormRepoLite(db *gorm.DB) meeting.Repository {
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
	m.Registrations = make([]*models.ProfileLabel, len(obj.Regs))
	for i, reg := range obj.Regs {
		label, err := h.profRepo.GetLabel(reg.UserId)
		if err != nil {
			return models.MeetingDetails{}, err
		}
		m.Registrations[i] = &label
	}
	return m, err
}

func (h *MeetingGormRepo) ToMeetingList(meetings []Meeting, userId int) ([]models.Meeting, error) {
	result := make([]models.Meeting, len(meetings))
	for i, el := range meetings {
		result[i] = h.ToMeeting(el, userId)
	}
	return result, nil
}

func (h *MeetingGormRepo) LikeExists(meetId int, userId int) bool {
	var l Like
	db := h.db.
		Where("meeting_id = ?", meetId).
		Where("user_id = ?", userId).
		First(&l)
	return db.Error == nil
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
	return db.Error == nil
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

func (h *MeetingGormRepo) FilterQuery(params meeting.FilterParams) *gorm.DB {
	return h.db.
		Where("Id > ?", params.PrevId).
		Where("Start_Date >= ?::date ", params.StartDate.Format("2006-01-02")).
		Where("End_Date <= ?::date", params.EndDate.Format("2006-01-02")).
		Preload("Tags").
		Preload("Regs").
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
	l := Registration{
		MeetingId: m.Id,
		UserId:    m.AuthorId,
	}
	db = h.db.Create(&l)
	return m.Id, db.Error
}

func (h *MeetingGormRepo) GetMeeting(meetingId, userId int, authorized bool) (models.MeetingDetails, error) {
	if !authorized {
		userId = -1
	}
	var m Meeting
	db := h.db.
		Where("id = ?", meetingId).
		Preload("Tags").
		Preload("Regs").
		First(&m)
	err := db.Error
	if err != nil {
		return models.MeetingDetails{}, err
	}
	return h.ToMeetingDetails(m, userId)
}

func (h *MeetingGormRepo) UpdateMeeting(update models.MeetingCard) error {
	obj, err := ToDbObject(update)
	if err != nil {
		return err
	}
	obj.Id = update.Label.Id
	db := h.db.Omit(clause.Associations).Save(&obj)
	err = db.Error
	if err == nil {
		err = h.db.Model(&obj).Association("Tags").Replace(obj.Tags)
	}
	return err
}

func (h *MeetingGormRepo) GetNextMeetings(params meeting.FilterParams) ([]models.Meeting, error) {
	var meetings []Meeting
	db := h.FilterQuery(params).Order("Start_Date ASC").Find(&meetings)
	err := db.Error
	if err != nil {
		return []models.Meeting{}, err
	}
	return h.ToMeetingList(meetings, params.UserId)
}

func (h *MeetingGormRepo) GetTopMeetings(params meeting.FilterParams) ([]models.Meeting, error) {
	var meetings []Meeting
	db := h.FilterQuery(params).Order("Likes_Count DESC").Find(&meetings)
	err := db.Error
	if err != nil {
		return []models.Meeting{}, err
	}
	return h.ToMeetingList(meetings, params.UserId)
}

func (h *MeetingGormRepo) FilterSubsLiked(params meeting.FilterParams) ([]models.Meeting, error) {
	subs, err := h.profRepo.GetUserSubscriptionIds(params.UserId)
	if err != nil {
		return nil, err
	}
	meetMap := make(map[int]*models.Meeting)
	for _, sub := range subs {
		params.UserId = sub
		subLiked, err := h.FilterLiked(params)
		if err != nil {
			return nil, err
		}
		params.CountLimit -= len(subLiked)
		for _, meet := range subLiked {
			meetMap[meet.Card.Label.Id] = &meet
		}
	}
	result := []models.Meeting{}
	for _, meetPtr := range meetMap {
		result = append(result, *meetPtr)
	}
	return result, nil
}

func (h *MeetingGormRepo) FilterLiked(params meeting.FilterParams) ([]models.Meeting, error) {
	var likes []Like
	db := h.db.
		Where("User_Id = ?", params.UserId).
		Where("Meeting_Id > ?", params.PrevId).
		Order("Meeting_Id ASC").
		Find(&likes)

	if db.Error != nil {
		return nil, db.Error
	}
	result := []models.Meeting{}
	total := 0
	for _, el := range likes {
		if total >= params.CountLimit {
			break
		}
		var m Meeting
		db = h.FilterQuery(params).
			Where("id = ?", el.MeetingId).
			First(&m)
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			continue
		}
		if db.Error != nil {
			return nil, db.Error
		}
		result = append(result, h.ToMeeting(m, params.UserId))
		total += 1
	}
	return result, nil
}

func (h *MeetingGormRepo) FilterLikedTags(userId int) (map[int]bool, error) {
	var likes []Like
	db := h.db.
		Where("User_Id = ?", userId).
		Order("Meeting_Id ASC").
		Find(&likes)

	if db.Error != nil {
		return nil, db.Error
	}
	result := map[int]bool{}
	for _, el := range likes {
		var m Meeting
		db = h.db.
			Where("id = ?", el.MeetingId).
			Preload("Tags").
			First(&m)
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			continue
		}
		if db.Error != nil {
			return nil, db.Error
		}
		for _, t := range m.Tags {
			result[t.Id] = true
		}
	}
	return result, nil
}

func (h *MeetingGormRepo) FilterSubsRegistered(params meeting.FilterParams) ([]models.Meeting, error) {
	subs, err := h.profRepo.GetUserSubscriptionIds(params.UserId)
	if err != nil {
		return nil, err
	}
	meetMap := make(map[int]*models.Meeting)
	for _, sub := range subs {
		params.UserId = sub
		subLiked, err := h.FilterRegistered(params)
		if err != nil {
			return nil, err
		}
		params.CountLimit -= len(subLiked)
		for _, meet := range subLiked {
			meetMap[meet.Card.Label.Id] = &meet
		}
	}
	result := []models.Meeting{}
	for _, meetPtr := range meetMap {
		result = append(result, *meetPtr)
	}
	return result, nil
}

func (h *MeetingGormRepo) FilterRegistered(params meeting.FilterParams) ([]models.Meeting, error) {
	var regs []Registration
	db := h.db.
		Where("User_Id = ?", params.UserId).
		Where("Meeting_Id > ?", params.PrevId).
		Order("Meeting_Id ASC").
		Limit(params.CountLimit).Find(&regs)

	if db.Error != nil {
		return nil, db.Error
	}
	result := []models.Meeting{}
	total := 0
	for _, el := range regs {
		if total >= params.CountLimit {
			break
		}
		var m Meeting
		db = h.FilterQuery(params).
			Where("id = ?", el.MeetingId).
			First(&m)
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			continue
		}
		if db.Error != nil {
			return nil, db.Error
		}
		result = append(result, h.ToMeeting(m, params.UserId))
	}
	return result, nil
}

func (h *MeetingGormRepo) ExtractMeetingsFromRows(params meeting.FilterParams, rows *sql.Rows) ([]Meeting, error) {
	meetings := []Meeting{}
	total := 0
	var meetId int
	var err error
	for err == nil && rows.Next() && total < params.CountLimit {
		var meetBuf Meeting
		err = rows.Scan(&meetId)
		db := h.db.
			Where("id = ?", meetId).
			Preload("Tags").
			Preload("Regs").
			First(&meetBuf)
		// Meetings now running are also displayed (hence EndDate.Before(params.StartDate))
		if db.Error == nil && (meetBuf.EndDate.Before(params.StartDate) || meetBuf.EndDate.After(params.EndDate)) {
			continue
		}
		meetings = append(meetings, meetBuf)
		total += 1
		err = db.Error
	}
	return meetings, err
}

func (h *MeetingGormRepo) FilterRecommended(params meeting.FilterParams) ([]models.Meeting, error) {
	subscriptions, err := h.profRepo.GetTagSubscriptions(params.UserId)
	if err != nil {
		return nil, err
	}
	likedTags, err := h.FilterLikedTags(params.UserId)
	if err != nil {
		return nil, err
	}
	for _, el := range subscriptions {
		likedTags[el] = true
	}
	likedTagsList := []int{}
	for k := range likedTags {
		likedTagsList = append(likedTagsList, k)
	}
	// Meetings with whose tags
	rows, err := h.db.Table("meeting_tags").
		Where("tag_id IN ?", likedTagsList).
		Where("meeting_id > ?", params.PrevId).
		Order("meeting_id ASC").
		Distinct("meeting_id").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	meetings, err := h.ExtractMeetingsFromRows(params, rows)
	if err != nil {
		return nil, err
	}
	return h.ToMeetingList(meetings, params.UserId)
}

func (h *MeetingGormRepo) FilterTagged(params meeting.FilterParams, tagId int) ([]models.Meeting, error) {
	rows, err := h.db.Table("meeting_tags").
		Where("tag_id = ?", tagId).
		Where("meeting_id > ?", params.PrevId).
		Order("meeting_id ASC").
		Distinct("meeting_id").Rows()
	if err != nil {
		return nil, err
	}
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
		Where("meeting_id <> ?", meetingId).
		Order("meeting_id ASC").
		Distinct("meeting_id").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	meetings, err := h.ExtractMeetingsFromRows(params, rows)
	if err == nil {
		return h.ToMeetingList(meetings, params.UserId)
	}
	return nil, err
}

func (h *MeetingGormRepo) SearchMeetings(params meeting.FilterParams,
	searchQuery string, limit int) ([]models.Meeting, error) {
	var res []Meeting
	nonWord := regexp.MustCompile(`([!&$()*+.:<=>?[\\\]^{|}-])`)
	searchQuery = nonWord.ReplaceAllString(searchQuery, "\\$1")
	space := regexp.MustCompile(`\s+`)
	searchQuery = space.ReplaceAllString(searchQuery, ":* & ") + ":*"
	db := h.db.Table("meetings").Where(`
(setweight(to_tsvector('russian', title), 'A') || setweight(to_tsvector('english', title), 'A') ||
setweight(to_tsvector('russian', text), 'B') || setweight(to_tsvector('english', text), 'B') || 
setweight(to_tsvector('russian', city), 'C') || setweight(to_tsvector('english', city), 'C') ||
setweight(to_tsvector('russian', address), 'D') || setweight(to_tsvector('english', address), 'D')
) @@ 
(to_tsquery('russian', ?) || to_tsquery('english', ?))`, searchQuery, searchQuery)
	if limit > 0 {
		db = db.Limit(limit)
	}
	db = db.Find(&res)

	if db.Error != nil {
		return nil, db.Error
	}

	return h.ToMeetingList(res, params.UserId)
}
