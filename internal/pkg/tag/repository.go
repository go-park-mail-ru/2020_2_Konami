package tag

import (
	"konami_backend/internal/pkg/models"
)

type Repository interface {
	GetTagById(id int) (models.Tag, error)
	GetTagByName(name string) (models.Tag, error)
	CreateTag(name string) (models.Tag, error)
	GetOrCreateTag(name string) (models.Tag, error)
	FilterTags(startsWith string) ([]models.Tag, error)
}
