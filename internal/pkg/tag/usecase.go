//go:generate mockgen -source=usecase.go -destination=./usecase_mock.go -package=tag
package tag

import (
	"konami_backend/internal/pkg/models"
)

type UseCase interface {
	GetTagById(id int) (models.Tag, error)
	GetTagByName(name string) (models.Tag, error)
	CreateTag(name string) (models.Tag, error)
	GetOrCreateTag(name string) (models.Tag, error)
	FilterTags(startsWith string) ([]models.Tag, error)
}
