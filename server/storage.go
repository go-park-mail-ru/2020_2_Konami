package main

import (
	"fmt"
	"regexp"
	"strings"
)

type User struct {
	Id           int            `json:"id"`
	Name         string         `json:"name"`
	Gender       string         `json:"gender"`
	Birthday     string         `json:"birthday"`
	City         string         `json:"city"`
	Email        string         `json:"email"`
	Telegram     string         `json:"telegram"`
	Vk           string         `json:"vk"`
	MeetingTags  []string       `json:"meetingTags"`
	Education    string         `json:"education"`
	Job          string         `json:"job"`
	ImgSrc       string         `json:"imgSrc"`
	Aims         string         `json:"aims"`
	InterestTags []string       `json:"interestTags"`
	Interests    string         `json:"interests"`
	SkillTags    []string       `json:"skillTags"`
	Skills       string         `json:"skills"`
	Meetings     []*UserMeeting `json:"meetings"`
}

type Meeting struct {
	Id        int      `json:"id"`
	AuthorId  int      `json:"authorId"`
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	ImgSrc    string   `json:"imgSrc"`
	Tags      []string `json:"tags"`
	Place     string   `json:"place"`
	StartDate string   `json:"startDate"`
	EndDate   string   `json:"endDate"`
	Seats     int      `json:"seats"`
	SeatsLeft int      `json:"seatsLeft"`
	Like      bool     `json:"isLiked"`
	Reg       bool     `json:"isRegistered"`
}

type UserUpdate struct {
	Name        *string        `json:"name"`
	Gender      *string        `json:"gender"`
	City        *string        `json:"city"`
	Birthday    *string        `json:"birthday"`
	Email       *string        `json:"email"`
	Telegram    *string        `json:"telegram"`
	Vk          *string        `json:"vk"`
	MeetingTags []string       `json:"meetingTags"`
	Education   *string        `json:"education"`
	Job         *string        `json:"job"`
	Aims        *string        `json:"aims"`
	Interests   *string        `json:"interests"`
	Skills      *string        `json:"skills"`
	Meetings    []*UserMeeting `json:"meetings"`
}

type UserMeeting struct {
	Title  string `json:"text"`
	ImgSrc string `json:"imgSrc"`
	Link   string `json:"link"`
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	uId      int
}

type UserId struct {
	Uid int `json:"userId"`
}

type MeetingUpload struct {
	Address     string   `json:"address"`
	City        string   `json:"city"`
	Start       string   `json:"start"`
	End         string   `json:"end"`
	Description string   `json:"meet-description"`
	Tags        []string `json:"meetingTags"`
	Name        string   `json:"name"`
	Photo       string   `json:"photo"`
}

type ErrResponse struct {
	ResponseCode int    `json:"-"`
	ErrMsg       string `json:"error"`
}

var Sessions = map[string]int{}

type UsersByName []*User

func (u UsersByName) Len() int {
	return len(u)
}

func (u UsersByName) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (u UsersByName) Less(i, j int) bool {
	return u[i].Name < u[j].Name
}

type MeetingsByDate []*Meeting

func (m MeetingsByDate) Len() int {
	return len(m)
}
func (m MeetingsByDate) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m MeetingsByDate) Less(i, j int) bool {
	return m[i].StartDate < m[j].StartDate
}

var MeetingStorage = map[int]*Meeting{
	0: {
		Id:       0,
		AuthorId: 0,
		Title:    "Забив с++",
		Text: "Lorem ipsum dolor sit amet, " +
			"consectetur adipiscing elit, sed " +
			"do eiusmod tempor incididunt ut " +
			"labore et dolore magna aliqua. " +
			"Ut enim ad minim veniam, quis " +
			"nostrud exercitation ullamco labori",
		ImgSrc:    "assets/paris.jpg",
		Tags:      []string{"C++"},
		Place:     "Москва, улица Колотушкина, дом Пушкина",
		StartDate: "2020-11-09 19:00:00",
		EndDate:   "2020-11-09 21:00:00",
		Seats:     100,
		SeatsLeft: 2,
	},
	1: {
		Id:       1,
		AuthorId: 0,
		Title:    "Python for Web",
		Text: "Lorem ipsum dolor sit amet, " +
			"consectetur adipiscing elit, sed " +
			"do eiusmod tempor incididunt ut " +
			"labore et dolore magna aliqua. " +
			"Ut enim ad minim veniam, quis " +
			"nostrud exercitation ullamco labori",
		ImgSrc:    "assets/paris.jpg",
		Tags:      []string{"Python", "Web"},
		Place:     "СПБ, улица Вязов, д.1",
		StartDate: "2020-11-10 19:00:00",
		EndDate:   "2020-11-10 21:00:00",
		Seats:     10,
		SeatsLeft: 0,
	},
	2: {
		Id:       2,
		AuthorId: 1,
		Title:    "GoLang for Web",
		Text: "Lorem ipsum dolor sit amet, " +
			"consectetur adipiscing elit, sed " +
			"do eiusmod tempor incididunt ut " +
			"labore et dolore magna aliqua. " +
			"Ut enim ad minim veniam, quis " +
			"nostrud exercitation ullamco labori",
		ImgSrc:    "assets/paris.jpg",
		Tags:      []string{"Go", "Web"},
		Place:     "СПБ, улица Вязов, д.1",
		StartDate: "2020-11-01 19:00:00",
		EndDate:   "2020-11-01 21:00:00",
		Seats:     10,
		SeatsLeft: 10,
	},
}

var UserStorage map[int]*User = map[int]*User{
	0: {
		Id:           0,
		Name:         "Александр",
		Gender:       "M",
		City:         "Нурсултан",
		Birthday:     "1990-09-12",
		Email:        "lucash@mail.ru",
		Telegram:     "",
		Vk:           "https://vk.com/id241926559",
		MeetingTags:  []string{"RandomTag1", "RandomTag5"},
		Education:    "МГТУ им. Н. Э. Баумана до 2010",
		Job:          "MAIL RU GROUP",
		ImgSrc:       "assets/luckash.jpeg",
		Aims:         "Хочу от жизни всего",
		InterestTags: []string{"Шыпшына", "Бульба"},
		Interests:    "Люблю, когда встаешь утром, а на столе #Шыпшына и #Бульба",
		SkillTags:    []string{"Мелиорация"},
		Skills:       "#Мелиорация - это моя жизнь",
		Meetings: []*UserMeeting{
			&UserMeeting{
				Title:  MeetingStorage[0].Title,
				ImgSrc: MeetingStorage[0].ImgSrc,
				Link:   fmt.Sprintf("/meet?meetId=%d", MeetingStorage[0].Id),
			},
			&UserMeeting{
				Title:  MeetingStorage[1].Title,
				ImgSrc: MeetingStorage[1].ImgSrc,
				Link:   fmt.Sprintf("/meet?meetId=%d", MeetingStorage[1].Id),
			},
		},
	},
	1: {
		Id:           1,
		Name:         "Роман",
		Gender:       "M",
		City:         "Москва",
		Birthday:     "2000-09-10",
		Email:        "lucash2@mail.ru",
		Telegram:     "",
		Vk:           "https://vk.com/id420",
		MeetingTags:  []string{"RandomTag1", "RandomTag5"},
		Job:          "HH.ru",
		ImgSrc:       "assets/luckash.jpg",
		InterestTags: []string{"ДВП", "ДСП"},
		Interests:    "Люблю клеить #ДВП и #ДСП",
		SkillTags:    []string{"Деревообработка"},
		Skills:       "Моя жизнь - это #Деревообработка",
		Meetings: []*UserMeeting{
			&UserMeeting{
				Title:  MeetingStorage[2].Title,
				ImgSrc: MeetingStorage[2].ImgSrc,
				Link:   fmt.Sprintf("/meet?meetId=%d", MeetingStorage[2].Id),
			},
		},
	},
}

var CredStorage = map[string]*Credentials{
	"lukash@mail.ru": {
		Login:    "lukash@mail.ru",
		Password: "$2a$04$7aVIDD36QgWr2L6iFgHGtesm0elmggbTryERfPruKS1e9R8CHadHi",
		uId:      0,
	},
	"lukash2@mail.ru": {
		Login:    "lukash2@mail.ru",
		Password: "$2a$04$7aVIDD36QgWr2L6iFgHGtesm0elmggbTryERfPruKS1e9R8CHadHi",
		uId:      1,
	},
}

// userID -> Liked meetings
var Likes = map[int][]int{
	0: []int{0},
}

// userID -> Liked meetings
var Registrations = map[int][]int{
	0: []int{0},
}

func contains(s []int, target int) bool {
	for _, el := range s {
		if el == target {
			return true
		}
	}
	return false
}

func UserLikes(userId, meetId int) bool {
	userLikes, lOk := Likes[userId]
	if lOk && contains(userLikes, meetId) {
		return true
	}
	return false
}

func SetEl(userId, meetId int, storage map[int][]int) bool {
	userElements, lOk := storage[userId]
	if !lOk {
		storage[userId] = []int{meetId}
		return true
	}
	if !contains(userElements, meetId) {
		storage[userId] = append(userElements, meetId)
		return true
	}
	return false
}

func RemoveEl(userId, meetId int, storage map[int][]int) bool {
	userElements, lOk := storage[userId]
	if lOk {
		for i, el := range userElements {
			if el == meetId {
				storage[userId] = append(userElements[:i], userElements[i+1:]...)
				return true
			}
		}
	}
	return false
}

func UserRegistered(userId, meetId int) bool {
	userRegs, rOk := Registrations[userId]
	if rOk && contains(userRegs, meetId) {
		return true
	}
	return false
}

type MeetUpdateFields struct {
	Reg  bool `json:"isRegistered"`
	Like bool `json:"isLiked"`
}

type MeetingUpdate struct {
	MeetId int               `json:"meetId"`
	Fields *MeetUpdateFields `json:"fields"`
}

func CommitUserUpdate(data *UserUpdate, usr *User) bool {
	ISOdt := regexp.MustCompile(`^(?:[1-9]\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\d|2[0-8])` +
		`|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\d(?:0[48]|[2468][048]` +
		`|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)$`)
	reEmail := regexp.MustCompile(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@` +
		`((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	if data.Birthday != nil {
		bDay := *data.Birthday
		tInd := strings.Index(bDay, "T")
		if tInd != -1 {
			bDay = bDay[:tInd]
		}
		if !ISOdt.MatchString(bDay) {
			return false
		}
		usr.Birthday = bDay
	}
	if data.Gender != nil && *data.Gender != "M" && *data.Gender != "F" && *data.Gender != "" {
		return false
	}
	if data.Email != nil && *data.Email != "" && !reEmail.MatchString(*data.Email) {
		return false
	}
	if data.Name != nil {
		usr.Name = *data.Name
	}
	if data.Gender != nil {
		usr.Gender = *data.Gender
	}
	if data.City != nil {
		usr.City = *data.City
	}
	if data.Birthday != nil {
		usr.Birthday = *data.Birthday
	}
	if data.Email != nil {
		usr.Email = *data.Email
	}
	if data.Telegram != nil {
		usr.Telegram = *data.Telegram
	}
	if data.Vk != nil {
		usr.Vk = *data.Vk
	}
	if data.MeetingTags != nil {
		usr.MeetingTags = data.MeetingTags
	}
	if data.Education != nil {
		usr.Education = *data.Education
	}
	if data.Job != nil {
		usr.Job = *data.Job
	}
	if data.Aims != nil {
		usr.Aims = *data.Aims
	}
	reMatch := regexp.MustCompile(`\#(?:([a-zA-Z0-9_а-яА-Яё\+\-*]{3,20})|(?:\(([a-zA-Z0-9_а-яА-Яё\ ]{3,20})\)))`)
	reSub := regexp.MustCompile(`[#()]`)
	if data.Interests != nil {
		usr.Interests = *data.Interests
		res := reMatch.FindAllString(usr.Interests, -1)
		usr.InterestTags = make([]string, len(res))
		for i, str := range res {
			usr.InterestTags[i] = reSub.ReplaceAllString(str, "")
		}
	}
	if data.Skills != nil {
		usr.Skills = *data.Skills
		res := reMatch.FindAllString(usr.Skills, -1)
		usr.SkillTags = make([]string, len(res))
		for i, str := range res {
			usr.SkillTags[i] = reSub.ReplaceAllString(str, "")
		}
	}
	if data.Meetings != nil {
		usr.Meetings = data.Meetings
	}
	return true
}
