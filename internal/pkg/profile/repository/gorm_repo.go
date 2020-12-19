package repository

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	meetingRepo "konami_backend/internal/pkg/meeting/repository"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/profile"
	"konami_backend/internal/pkg/tag"
	tagRepo "konami_backend/internal/pkg/tag/repository"
	"time"
)

type ProfileGormRepo struct {
	db      *gorm.DB
	TagRepo tag.Repository
}

func NewProfileGormRepo(db *gorm.DB) profile.Repository {
	return &ProfileGormRepo{db: db}
}

type Profile struct {
	Id           int `gorm:"primaryKey;autoIncrement;"`
	Name         string
	ImgSrc       string
	Job          string
	MeetingTags  []tagRepo.Tag `gorm:"many2many:profile_meetingTags;"`
	InterestTags []InterestTag `gorm:"many2many:profile_interestTags;"`
	SkillTags    []SkillTag    `gorm:"many2many:profile_skillTags;"`
	Gender       string
	Birthday     time.Time
	City         string
	Login        string `gorm:"unique;"`
	PwdHash      string
	Telegram     string
	Vk           string
	Education    string
	Aims         string
	Interests    string
	Skills       string
	Meetings     []meetingRepo.Meeting `gorm:"foreignKey:AuthorId;"`
}

type InterestTag struct {
	Id   int    `gorm:"primaryKey;autoIncrement;"`
	Name string `gorm:"unique;"`
}

func (t *InterestTag) TableName() string {
	return "InterestTags"
}

type SkillTag struct {
	Id   int    `gorm:"primaryKey;autoIncrement;"`
	Name string `gorm:"unique;"`
}

func (t *SkillTag) TableName() string {
	return "SkillTags"
}

type Subscription struct {
	Id       int `gorm:"primaryKey;autoIncrement;"`
	AuthorId int
	TargetId int
}

func (t *Subscription) TableName() string {
	return "Subscriptions"
}

func (h *ProfileGormRepo) GetUserSubscriptions(userId int) ([]models.ProfileCard, error) {
	var subs []Subscription
	db := h.db.
		Where("AuthorId = ?", userId).
		Find(&subs)
	err := db.Error
	if err != nil {
		return nil, err
	}
	result := make([]models.ProfileCard, len(subs))
	for i, sub := range subs {
		var p Profile
		db := h.db.
			Where("id = ?", sub.TargetId).
			Preload("MeetingTags").
			Preload("InterestTags").
			Preload("SkillTags").
			Preload("Meetings").
			First(&p)
		err := db.Error
		if err != nil {
			return nil, err
		}
		result[i] = ToProfileCard(p)
	}
	return result, nil
}

func (h *ProfileGormRepo) CreateSubscription(authorId int, targetId int) (int, error) {
	var p Profile
	db := h.db.
		Where("id = ?", targetId).
		First(&p)
	err := db.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return -1, profile.ErrUserNonExistent
	}
	if err != nil {
		return -1, db.Error
	}
	s := Subscription{
		AuthorId: authorId,
		TargetId: targetId,
	}
	db = h.db.Create(&s)
	err = db.Error
	if err != nil {
		return -1, err
	}
	return s.Id, nil
}

func (h *ProfileGormRepo) RemoveSubscription(authorId int, targetId int) error {
	db := h.db.
		Where("AuthorId = ?", authorId).
		Where("TargetId = ?", targetId).
		Delete(&Subscription{})
	return db.Error
}

func (h *ProfileGormRepo) GetSkillByName(name string) (SkillTag, error) {
	var res SkillTag
	db := h.db.
		Where("name = ?", name).
		First(&res)

	err := db.Error
	if err != nil {
		return SkillTag{}, err
	}
	return SkillTag{Id: res.Id, Name: res.Name}, nil
}

func (h *ProfileGormRepo) CreateSkill(name string) (SkillTag, error) {
	t := SkillTag{Name: name}
	db := h.db.Create(&t)
	err := db.Error
	if err != nil {
		return SkillTag{}, err
	}
	return SkillTag{Id: t.Id, Name: t.Name}, nil
}

func (h *ProfileGormRepo) GetOrCreateSkill(name string) (SkillTag, error) {
	result, err := h.GetSkillByName(name)
	if err == gorm.ErrRecordNotFound {
		result, err = h.CreateSkill(name)
	}
	if err != nil {
		return SkillTag{}, err
	}
	return result, nil
}

func (h *ProfileGormRepo) GetInterestByName(name string) (InterestTag, error) {
	var res InterestTag
	db := h.db.
		Where("name = ?", name).
		First(&res)

	err := db.Error
	if db.Error != nil {
		return InterestTag{}, err
	}
	return InterestTag{Id: res.Id, Name: res.Name}, nil
}

func (h *ProfileGormRepo) CreateInterest(name string) (InterestTag, error) {
	t := InterestTag{Name: name}
	db := h.db.Create(&t)
	err := db.Error
	if err != nil {
		return InterestTag{}, err
	}
	return InterestTag{Id: t.Id, Name: t.Name}, nil
}

func (h *ProfileGormRepo) GetOrCreateInterest(name string) (InterestTag, error) {
	result, err := h.GetInterestByName(name)
	if err == gorm.ErrRecordNotFound {
		result, err = h.CreateInterest(name)
	}
	if err != nil {
		return InterestTag{}, err
	}
	return result, nil
}

func ToProfileCard(obj Profile) models.ProfileCard {
	p := models.ProfileCard{
		Label: &models.ProfileLabel{
			Id:     obj.Id,
			Name:   obj.Name,
			ImgSrc: obj.ImgSrc,
		},
		Job: obj.Job,
	}
	p.InterestTags = make([]string, len(obj.InterestTags))
	for i, val := range obj.InterestTags {
		p.InterestTags[i] = val.Name
	}
	p.SkillTags = make([]string, len(obj.SkillTags))
	for i, val := range obj.SkillTags {
		p.SkillTags[i] = val.Name
	}
	return p
}

func ToProfile(obj Profile) models.Profile {
	card := ToProfileCard(obj)
	p := models.Profile{
		Card:      &card,
		Gender:    obj.Gender,
		City:      obj.City,
		Login:     obj.Login,
		PwdHash:   obj.PwdHash,
		Telegram:  obj.Telegram,
		Vk:        obj.Vk,
		Education: obj.Education,
		Aims:      obj.Aims,
		Interests: obj.Interests,
		Skills:    obj.Skills,
	}
	if obj.Birthday.Unix() != (time.Time{}).Unix() {
		p.Birthday = obj.Birthday.Format("2006-01-02")
	}
	p.MeetingTags = make([]*models.Tag, len(obj.MeetingTags))
	for i, val := range obj.MeetingTags {
		t := tagRepo.ToModel(val)
		p.MeetingTags[i] = &t
	}
	p.Meetings = make([]*models.MeetingLabel, len(obj.Meetings))
	for i, val := range obj.Meetings {
		m := meetingRepo.ToMeetingLabel(val)
		p.Meetings[i] = &m
	}
	return p
}

func ToDbObject(p models.Profile) (Profile, error) {
	obj := Profile{
		Id:        p.Card.Label.Id,
		Name:      p.Card.Label.Name,
		ImgSrc:    p.Card.Label.ImgSrc,
		Job:       p.Card.Job,
		Gender:    p.Gender,
		City:      p.City,
		Login:     p.Login,
		PwdHash:   p.PwdHash,
		Telegram:  p.Telegram,
		Vk:        p.Vk,
		Education: p.Education,
		Aims:      p.Aims,
		Interests: p.Interests,
		Skills:    p.Skills,
	}
	obj.MeetingTags = make([]tagRepo.Tag, len(p.MeetingTags))
	for i, el := range p.MeetingTags {
		obj.MeetingTags[i] = tagRepo.Tag{Id: el.TagId, Name: el.Name}
	}
	obj.InterestTags = make([]InterestTag, len(p.Card.InterestTags))
	for i, el := range p.Card.InterestTags {
		obj.InterestTags[i] = InterestTag{Name: el}
	}
	obj.SkillTags = make([]SkillTag, len(p.Card.SkillTags))
	for i, el := range p.Card.SkillTags {
		obj.SkillTags[i] = SkillTag{Name: el}
	}
	if p.Birthday != "" {
		var err error
		layout := "2006-01-02"
		obj.Birthday, err = time.Parse(layout, p.Birthday)
		if err != nil {
			return Profile{}, err
		}
	}
	return obj, nil
}

func (h ProfileGormRepo) GetAll() ([]models.ProfileCard, error) {
	var profiles []Profile
	db := h.db.
		Preload("MeetingTags").
		Preload("InterestTags").
		Preload("SkillTags").
		Preload("Meetings").
		Find(&profiles)
	err := db.Error
	if err != nil {
		return nil, err
	}
	result := make([]models.ProfileCard, len(profiles))
	for i, el := range profiles {
		result[i] = ToProfileCard(el)
	}
	return result, nil
}

func (h ProfileGormRepo) GetProfile(userId int) (models.Profile, error) {
	var p Profile
	db := h.db.
		Where("id = ?", userId).
		Preload("MeetingTags").
		Preload("InterestTags").
		Preload("SkillTags").
		Preload("Meetings").
		First(&p)
	err := db.Error
	if err != nil {
		return models.Profile{}, err
	}
	return ToProfile(p), nil
}

func (h ProfileGormRepo) EditProfile(update models.Profile) error {
	updatedObj, err := ToDbObject(update)
	if err != nil {
		return err
	}
	target := Profile{Id: update.Card.Label.Id}
	db := h.db.Omit(clause.Associations).Save(&updatedObj)

	if db.Error == nil {
		err = h.db.Model(&target).Association("MeetingTags").Replace(updatedObj.MeetingTags)
	}
	if err == nil {
		buf := []InterestTag{}
		for _, val := range update.Card.InterestTags {
			tg, err := h.GetOrCreateInterest(val)
			if err == nil {
				buf = append(buf, tg)
			} else {
				return err
			}
		}
		err = h.db.Model(&target).Association("InterestTags").Replace(buf)
	}
	if err == nil {
		buf := []SkillTag{}
		for _, val := range update.Card.SkillTags {
			tg, err := h.GetOrCreateSkill(val)
			if err == nil {
				buf = append(buf, tg)
			} else {
				return err
			}
		}
		err = h.db.Model(&target).Association("SkillTags").Replace(buf)
	}
	if err == nil {
		err = h.db.Model(&target).Association("MeetingTags").Replace(updatedObj.MeetingTags)
	}
	return err
}

func (h ProfileGormRepo) EditProfilePic(userId int, imgSrc string) error {
	var obj Profile
	db := h.db.
		Where("id = ?", userId).
		First(&obj)
	err := db.Error
	if err != nil {
		return err
	}
	obj.ImgSrc = imgSrc
	db = h.db.Save(obj)
	return db.Error
}

func (h ProfileGormRepo) Create(p models.Profile) (int, error) {
	obj, err := ToDbObject(p)
	if err != nil {
		return 0, err
	}
	db := h.db.Create(&obj)
	return obj.Id, db.Error
}

func (h ProfileGormRepo) GetCredentials(login string) (int, string, error) {
	var obj Profile
	db := h.db.
		Where("Login = ?", login).
		First(&obj)
	err := db.Error
	if err == gorm.ErrRecordNotFound {
		return 0, "", profile.ErrUserNonExistent
	}
	if err != nil {
		return 0, "", err
	}
	return obj.Id, obj.PwdHash, nil
}

func (h *ProfileGormRepo) GetLabel(userId int) (models.ProfileLabel, error) {
	var p Profile
	db := h.db.
		Where("id = ?", userId).
		First(&p)
	err := db.Error
	if err != nil {
		return models.ProfileLabel{}, err
	}
	return models.ProfileLabel{
		Id:     p.Id,
		Name:   p.Name,
		ImgSrc: p.ImgSrc,
	}, nil
}

func (h *ProfileGormRepo) GetTagSubscriptions(userId int) (tagIds []int, err error) {
	var userProfile Profile
	db := h.db.
		Where("id = ?", userId).
		Preload("MeetingTags").
		First(&userProfile)
	err = db.Error
	if err != nil {
		return nil, err
	}
	// Tags to which user is subscribed
	tagIds = make([]int, len(userProfile.MeetingTags))
	for i, t := range userProfile.MeetingTags {
		tagIds[i] = t.Id
	}
	return tagIds, nil
}
