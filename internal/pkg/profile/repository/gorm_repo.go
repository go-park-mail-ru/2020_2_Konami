package repository

import (
	"github.com/jinzhu/gorm"
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
		Birthday:  obj.Birthday.Format("2006-01-02"),
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
		obj.MeetingTags[i] = tagRepo.ToDbObject(*el)
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
		layout := "2006-01-02 15:04:05"
		obj.Birthday, err = time.Parse(layout, p.Birthday)
		if err != nil {
			return Profile{}, err
		}
	}
	return obj, nil
}

func (h ProfileGormRepo) GetAll() ([]models.ProfileCard, error) {
	var profiles []Profile
	db := h.db.Find(&profiles)
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
		First(&p)
	err := db.Error
	if err != nil {
		return models.Profile{}, err
	}
	return ToProfile(p), nil
}

func (h ProfileGormRepo) EditProfile(update models.Profile) error {
	var old Profile
	db := h.db.
		Where("id = ?", update.Card.Label.Id).
		First(&old)
	err := db.Error
	if err != nil {
		return err
	}
	updatedObj, err := ToDbObject(update)
	if err != nil {
		return err
	}
	updatedObj.Id = old.Id
	updatedObj.Meetings = old.Meetings
	db = h.db.Save(updatedObj)
	return db.Error
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
	if gorm.IsRecordNotFoundError(err) {
		return 0, "", profile.ErrUserNonExistent
	}
	if err != nil {
		return 0, "", err
	}
	return obj.Id, obj.PwdHash, nil
}
