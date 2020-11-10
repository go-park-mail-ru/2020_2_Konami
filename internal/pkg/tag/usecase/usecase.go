package usecase

import (
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/tag"
)

type TagUseCase struct {
	repo tag.Repository
}

func NewTagUseCase(tagRepo tag.Repository) tag.UseCase {
	return TagUseCase{repo: tagRepo}
}

func (u TagUseCase) GetTagById(id int) (models.Tag, error) {
	return u.repo.GetTagById(id)
}
func (u TagUseCase) GetTagByName(name string) (models.Tag, error) {
	return u.repo.GetTagByName(name)
}
func (u TagUseCase) CreateTag(name string) (models.Tag, error) {
	return u.repo.CreateTag(name)
}
func (u TagUseCase) GetOrCreateTag(name string) (models.Tag, error) {
	return u.repo.GetOrCreateTag(name)
}

func (u TagUseCase) FilterTags(startsWith string) ([]models.Tag, error) {
	return u.repo.FilterTags(startsWith)
}
