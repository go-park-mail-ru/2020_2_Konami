package usecase

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io"
	"konami_backend/internal/pkg/models"
	"konami_backend/internal/pkg/profile"
	"konami_backend/internal/pkg/tag"
	"konami_backend/internal/pkg/utils/uploads_handler"
	"regexp"
	"strconv"
	"strings"
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

func (h ProfileUseCase) UpdateProfileCard(card *models.ProfileCard, data models.ProfileUpdate) {
	if data.Name != nil {
		card.Label.Name = *data.Name
	}
	if data.Job != nil {
		card.Job = *data.Job
	}
	reMatch := regexp.MustCompile(`#(?:([a-zA-Z0-9_а-яА-Яё+\-*]{3,20})|(?:\(([a-zA-Z0-9_а-яА-Яё ]{3,20})\)))`)
	reSub := regexp.MustCompile(`[#()]`)
	if data.Interests != nil {
		res := reMatch.FindAllString(*data.Interests, -1)
		card.InterestTags = make([]string, len(res))
		for i, str := range res {
			card.InterestTags[i] = reSub.ReplaceAllString(str, "")
		}
	}
	if data.Skills != nil {
		res := reMatch.FindAllString(*data.Skills, -1)
		card.SkillTags = make([]string, len(res))
		for i, str := range res {
			card.SkillTags[i] = reSub.ReplaceAllString(str, "")
		}
	}
}

func (h ProfileUseCase) GetAll(params profile.FilterParams) ([]models.ProfileCard, error) {
	return h.ProfileRepo.GetAll(params)
}

func (h ProfileUseCase) GetUserSubscriptions(params profile.FilterParams) ([]models.ProfileCard, error) {
	return h.ProfileRepo.GetUserSubscriptions(params)
}

func (h ProfileUseCase) CreateSubscription(authorId int, targetId int) (int, error) {
	return h.ProfileRepo.CreateSubscription(authorId, targetId)
}

func (h ProfileUseCase) RemoveSubscription(authorId int, targetId int) error {
	return h.ProfileRepo.RemoveSubscription(authorId, targetId)
}

func (h ProfileUseCase) GetProfile(reqAuthorId, userId int) (models.Profile, error) {
	return h.ProfileRepo.GetProfile(reqAuthorId, userId)
}

func (h ProfileUseCase) EditProfile(userId int, data models.ProfileUpdate) error {
	p, err := h.ProfileRepo.GetProfile(-1, userId)
	if err != nil {
		return err
	}
	if data.Birthday != nil {
		p.Birthday = *data.Birthday
	}
	if data.Gender != nil {
		if *data.Gender != "M" && *data.Gender != "F" && *data.Gender != "" {
			return errors.New("non-binary gender not allowed :-)")
		}
		p.Gender = *data.Gender
	}
	if data.City != nil {
		p.City = *data.City
	}
	if data.Telegram != nil {
		p.Telegram = *data.Telegram
	}
	if data.Vk != nil {
		p.Vk = *data.Vk
	}
	if data.MeetingTags != nil {
		mTags := make([]*models.Tag, len(data.MeetingTags))
		for i, el := range data.MeetingTags {
			t, err := h.TagRepo.GetOrCreateTag(el)
			if err != nil {
				return errors.New("unable to parse tags")
			}
			mTags[i] = &t
		}
		p.MeetingTags = mTags
	}
	if data.Education != nil {
		p.Education = *data.Education
	}
	if data.Aims != nil {
		p.Aims = *data.Aims
	}
	if data.Interests != nil {
		p.Interests = *data.Interests
	}
	if data.Skills != nil {
		p.Skills = *data.Skills
	}
	h.UpdateProfileCard(p.Card, data)
	return h.ProfileRepo.EditProfile(p)
}

func (h ProfileUseCase) UploadProfilePic(userId int, filename string, img io.Reader) error {
	fnameSlice := strings.Split(filename, ".")
	ext := fnameSlice[len(fnameSlice)-1]
	imgPath := h.ProfilePicsDir + strconv.Itoa(userId) + "." + ext
	imgPath, err := h.UploadsHandler.UploadImage(imgPath, img)
	if err != nil {
		return err
	}
	return h.ProfileRepo.EditProfilePic(userId, imgPath)
}

func (h ProfileUseCase) SignUp(cred models.Credentials) (int, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(cred.Password), bcrypt.MinCost)
	if err != nil {
		return 0, err
	}
	p := models.Profile{
		Card: &models.ProfileCard{
			Label: &models.ProfileLabel{
				Name:   "Пользователь",
				ImgSrc: h.defaultImgSrc,
			},
			InterestTags: []string{},
			SkillTags:    []string{},
		},
		Login:       cred.Login,
		PwdHash:     string(hashed),
		MeetingTags: []*models.Tag{},
		Meetings:    []*models.MeetingLabel{},
	}
	return h.ProfileRepo.Create(p)
}

func (h ProfileUseCase) Validate(cred models.Credentials) (int, error) {
	userId, pwdHash, err := h.ProfileRepo.GetCredentials(cred.Login)
	if err != nil {
		return 0, err
	}
	cmpRes := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(cred.Password))
	if cmpRes != nil {
		return 0, profile.ErrInvalidCredentials
	}
	return userId, nil
}
