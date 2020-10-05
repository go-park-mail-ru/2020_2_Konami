package main

type User struct {
	Id           int        `json:"id"`
	Name         string     `json:"name"`
	Gender       string     `json:"gender"`
	City         string     `json:"city"`
	Email        string     `json:"email"`
	Telegram     string     `json:"telegram"`
	Vk           string     `json:"vk"`
	Education    string     `json:"education"`
	Job          string     `json:"job"`
	ImgPath      string     `json:"imgPath"`
	Aims         string     `json:"aims"`
	InterestTags []string   `json:"interestTags"`
	Interests    string     `json:"interests"`
	SkillTags    []string   `json:"skillTags"`
	Skills       string     `json:"skills"`
	Meetings     []*Meeting `json:"meetings"`
}

type Meeting struct {
	Id     int      `json:"id"`
	Title  string   `json:"title"`
	Text   string   `json:"text"`
	ImgSrc string   `json:"imgSrc"`
	Tags   []string `json:"tags"`
	Place  string   `json:"place"`
	Date   string   `json:"date"`
}

type UserUpdate struct {
	Name         *string    `json:"name"`
	Gender       *string    `json:"gender"`
	City         *string    `json:"city"`
	Email        *string    `json:"email"`
	Telegram     *string    `json:"telegram"`
	Vk           *string    `json:"vk"`
	Education    *string    `json:"education"`
	Job          *string    `json:"job"`
	Aims         *string    `json:"aims"`
	InterestTags []string   `json:"interestTags"`
	Interests    *string    `json:"interests"`
	SkillTags    []string   `json:"skillTags"`
	Skills       *string    `json:"skills"`
	Meetings     []*Meeting `json:"meetings"`
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	uId      int
}

type UserId struct {
	Uid int `json:"userId"`
}

var Sessions = map[string]int{}

var UserStorage = map[int]*User{
	0: {
		Id:           0,
		Name:         "Александр",
		Gender:       "M",
		City:         "Нурсултан",
		Email:        "lucash@mail.ru",
		Telegram:     "",
		Vk:           "https://vk.com/id241926559",
		Education:    "МГТУ им. Н. Э. Баумана до 2010",
		Job:          "MAIL RU GROUP",
		ImgPath:      "assets/luckash.jpeg",
		Aims:         "Хочу от жизни всего",
		InterestTags: []string{"Картофель"},
		Interests:    "Люблю, когда встаешь утром, а на столе #Шыпшына и #Картофель",
		SkillTags:    []string{"Мелиорация"},
		Skills:       "#Мелиорация - это моя жизнь",
		Meetings:     []*Meeting{},
	},
}

var MeetingStorage = map[int]*Meeting{
	0: {
		Id:    0,
		Title: "Забив с++",
		Text: "Lorem ipsum dolor sit amet, " +
			"consectetur adipiscing elit, sed " +
			"do eiusmod tempor incididunt ut " +
			"labore et dolore magna aliqua. " +
			"Ut enim ad minim veniam, quis " +
			"nostrud exercitation ullamco labori",
		ImgSrc: "assets/paris.jpg",
		Tags:   []string{},
		Place:  "Дом Пушкина, улица Колотушкина",
		Date:   "12 сентября 2020",
	},
}

var CredStorage = map[string]*Credentials{
	"lukash@mail.ru": {
		Login:    "lukash@mail.ru",
		Password: "12345",
		uId:      0,
	},
}

func CommitUserUpdate(data *UserUpdate, userId int) bool {
	usr, exists := UserStorage[userId]
	if !exists {
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
	if data.Email != nil {
		usr.Email = *data.Email
	}
	if data.Telegram != nil {
		usr.Telegram = *data.Telegram
	}
	if data.Vk != nil {
		usr.Vk = *data.Vk
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
	if data.InterestTags != nil {
		usr.InterestTags = data.InterestTags
	}
	if data.Interests != nil {
		usr.Interests = *data.Interests
	}
	if data.SkillTags != nil {
		usr.SkillTags = data.SkillTags
	}
	if data.Skills != nil {
		usr.Skills = *data.Skills
	}
	if data.Meetings != nil {
		usr.Meetings = data.Meetings
	}
	return true
}
