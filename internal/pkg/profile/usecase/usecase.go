package usecase

import (
	"io"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/profile"
	"konami_backend/internal/pkg/tag"
	"konami_backend/internal/pkg/utils/uploads_handler"
)

type ProfileUseCase struct {
	ProfileRepo    profile.Repository
	UploadsHandler uploads_handler.UploadsHandler
	TagRepo        tag.Repository
	ProfilePicsDir string
	defaultImgSrc  string
}

func NewProfileUseCase(ProfileRepo profile.Repository,
	UploadsHandler uploads_handler.UploadsHandler,
	TagRepo tag.Repository,
	ProfilePicsDir string,
	defaultImgSrc string) profile.UseCase {

	return &ProfileUseCase{
		ProfileRepo:    ProfileRepo,
		UploadsHandler: UploadsHandler,
		TagRepo:        TagRepo,
		ProfilePicsDir: ProfilePicsDir,
		defaultImgSrc:  defaultImgSrc,
	}
}

func (p ProfileUseCase) GetAll() ([]models.ProfileCard, error) {
	panic("implement me")
}

func (p ProfileUseCase) GetProfile(userId int) (models.Profile, error) {
	panic("implement me")
}

func (p ProfileUseCase) EditProfile(userId int, update models.ProfileUpdate) error {
	panic("implement me")
}

func (p ProfileUseCase) UploadProfilePic(userId int, filename string, img io.Reader) error {
	panic("implement me")
}

func (p ProfileUseCase) SignUp(cred models.Credentials) (int, error) {
	panic("implement me")
}

func (p ProfileUseCase) Validate(cred models.Credentials) (int, error) {
	panic("implement me")
}
