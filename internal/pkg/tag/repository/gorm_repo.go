package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/tag"
	"strings"
)

type TagGormRepo struct {
	db *gorm.DB
}

func NewTagGormRepo(db *gorm.DB) tag.Repository {
	return &TagGormRepo{db: db}
}

type Tag struct {
	Id   int    `gorm:"primaryKey;autoIncrement;"`
	Name string `gorm:"unique;"`
}

func (t *Tag) TableName() string {
	return "tags"
}

func ToModel(t Tag) models.Tag {
	return models.Tag{
		TagId: t.Id,
		Name:  t.Name,
	}
}

func ToDbObject(tag models.Tag) Tag {
	return Tag{
		Id:   tag.TagId,
		Name: tag.Name,
	}
}

func (h *TagGormRepo) GetTagById(id int) (models.Tag, error) {
	var res Tag
	db := h.db.
		Table("tags").
		Where("id = ?", id).
		First(&res)

	err := db.Error
	if db.Error != nil {
		return models.Tag{}, err
	}
	return ToModel(res), nil
}

func (h *TagGormRepo) GetTagByName(name string) (models.Tag, error) {
	var res Tag
	db := h.db.
		Table("tags").
		Where("UPPER(name) = ?", strings.ToUpper(name)).
		First(&res)

	err := db.Error
	if db.Error != nil {
		return models.Tag{}, err
	}
	return ToModel(res), nil
}

func (h *TagGormRepo) CreateTag(name string) (models.Tag, error) {
	t := Tag{Name: name}
	db := h.db.Create(&t)
	err := db.Error
	if err != nil {
		return models.Tag{}, err
	}
	return ToModel(t), nil
}

func (h *TagGormRepo) GetOrCreateTag(name string) (models.Tag, error) {
	result, err := h.GetTagByName(name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		result, err = h.CreateTag(name)
	}
	if err != nil {
		return models.Tag{}, err
	}
	return result, nil
}

func (h *TagGormRepo) FilterTags(startsWith string) ([]models.Tag, error) {
	var tSlice []Tag
	db := h.db.
		Table("tags").
		Where("UPPER(name) LIKE ?", strings.ToUpper(startsWith)+"%").
		Find(&tSlice)
	err := db.Error
	if err != nil {
		return nil, err
	}
	res := make([]models.Tag, len(tSlice))
	for i, t := range tSlice {
		res[i] = ToModel(t)
	}
	return res, nil
}
