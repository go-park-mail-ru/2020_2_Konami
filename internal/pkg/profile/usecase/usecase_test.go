package usecase

import (
	"errors"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/profile"
	"konami_backend/internal/pkg/tag"
	uploadsHandlerPkg "konami_backend/internal/pkg/utils/uploads_handler"
	"strings"
	"testing"
)

func TestTag(t *testing.T) {
	t.Run("TestValidateProfile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		proRepo := profile.NewMockRepository(ctrl)
		uploadsHandler := uploadsHandlerPkg.NewUploadsHandler("uploadsDir")

		tagRepo := tag.NewMockRepository(ctrl)

		p := NewProfileUseCase(proRepo, uploadsHandler, tagRepo, "", "")

		hashed, _ := bcrypt.GenerateFromPassword([]byte("qwerty"), bcrypt.MinCost)

		proRepo.EXPECT().
			GetCredentials("qwerty").
			Return(1, string(hashed), nil)

		_, _ = p.Validate(models.Credentials{
			Login:    "qwerty",
			Password: "qwerty",
		})
	})

	t.Run("TestValidateProfileErrorLogin", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		proRepo := profile.NewMockRepository(ctrl)
		uploadsHandler := uploadsHandlerPkg.NewUploadsHandler("uploadsDir")

		tagRepo := tag.NewMockRepository(ctrl)

		p := NewProfileUseCase(proRepo, uploadsHandler, tagRepo, "", "")

		hashed, _ := bcrypt.GenerateFromPassword([]byte("qwerty"), bcrypt.MinCost)

		proRepo.EXPECT().
			GetCredentials("qwerty").
			Return(1, string(hashed), errors.New("ERROR"))

		_, _ = p.Validate(models.Credentials{
			Login:    "qwerty",
			Password: "qwerty",
		})
	})


	t.Run("TestValidateProfileErrorWrongHash", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		proRepo := profile.NewMockRepository(ctrl)
		uploadsHandler := uploadsHandlerPkg.NewUploadsHandler("uploadsDir")

		tagRepo := tag.NewMockRepository(ctrl)

		p := NewProfileUseCase(proRepo, uploadsHandler, tagRepo, "", "")

		proRepo.EXPECT().
			GetCredentials("qwerty").
			Return(1, "qwerty", nil)

		_, _ = p.Validate(models.Credentials{
			Login:    "qwerty",
			Password: "qwerty",
		})
	})

	t.Run("TestUploadErr", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		proRepo := profile.NewMockRepository(ctrl)
		uploadsHandler := uploadsHandlerPkg.NewUploadsHandler("uploadsDir")

		tagRepo := tag.NewMockRepository(ctrl)

		p := NewProfileUseCase(proRepo, uploadsHandler, tagRepo, "", "")

		r := strings.NewReader("abcde")

		_ = p.UploadProfilePic(1, "file", r)
	})

	t.Run("TestUpdateNothing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		proRepo := profile.NewMockRepository(ctrl)
		uploadsHandler := uploadsHandlerPkg.NewUploadsHandler("uploadsDir")

		tagRepo := tag.NewMockRepository(ctrl)

		p := NewProfileUseCase(proRepo, uploadsHandler, tagRepo, "", "")

		testProfile := models.Profile{
			Card:        nil,
			Gender:      "",
			Birthday:    "",
			City:        "",
			Login:       "",
			PwdHash:     "",
			Telegram:    "",
			Vk:          "",
			Education:   "",
			MeetingTags: nil,
			Aims:        "",
			Interests:   "",
			Skills:      "",
			Meetings:    nil,
		}

		proRepo.EXPECT().GetProfile(1).Return(testProfile, nil)
		proRepo.EXPECT().EditProfile(testProfile).Return(nil)

		_ = p.EditProfile(1, models.ProfileUpdate{
			Name:        nil,
			Gender:      nil,
			City:        nil,
			Birthday:    nil,
			Telegram:    nil,
			Vk:          nil,
			MeetingTags: nil,
			Education:   nil,
			Job:         nil,
			Aims:        nil,
			Interests:   nil,
			Skills:      nil,
		})
	})

	t.Run("TestUpdateAll", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		proRepo := profile.NewMockRepository(ctrl)
		uploadsHandler := uploadsHandlerPkg.NewUploadsHandler("uploadsDir")

		tagRepo := tag.NewMockRepository(ctrl)

		p := NewProfileUseCase(proRepo, uploadsHandler, tagRepo, "", "")

		testProfile := models.Profile{
			Card:        nil,
			Gender:      "",
			Birthday:    "",
			City:        "",
			Login:       "",
			PwdHash:     "",
			Telegram:    "",
			Vk:          "",
			Education:   "",
			MeetingTags: nil,
			Aims:        "",
			Interests:   "",
			Skills:      "",
			Meetings:    nil,
		}

		proRepo.EXPECT().GetProfile(1).Return(testProfile, nil)
		testProfile.Birthday = "M"
		testProfile.Gender = "M"
		testProfile.City = "M"
		testProfile.Telegram = "M"
		testProfile.Vk = "M"
		testProfile.Education = "M"
		testProfile.Aims = "M"

		proRepo.EXPECT().EditProfile(testProfile).Return(nil)

		upd := "M"

		_ = p.EditProfile(1, models.ProfileUpdate{
			Name:        nil,
			Gender:      &upd,
			City:        &upd,
			Birthday:    &upd,
			Telegram:    &upd,
			Vk:          &upd,
			MeetingTags: nil,
			Education:   &upd,
			Job:         nil,
			Aims:        &upd,
			Interests:   nil,
			Skills:      nil,
		})
	})

	t.Run("TestOtherProfileUtils", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		proRepo := profile.NewMockRepository(ctrl)
		uploadsHandler := uploadsHandlerPkg.NewUploadsHandler("uploadsDir")

		tagRepo := tag.NewMockRepository(ctrl)

		p := NewProfileUseCase(proRepo, uploadsHandler, tagRepo, "", "")

		proRepo.EXPECT().
			GetAll().
			Return([]models.ProfileCard{}, nil)

		proRepo.EXPECT().
			GetProfile(1).
			Return(models.Profile{}, nil)

		_, _ = p.GetAll()
		_, _ = p.GetProfile(1)
	})
}