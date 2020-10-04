package main

type UserCard struct {
	CardId       uint      `json:"cardId"`
	Name         string   `json:"name"`
	ImgSrc       string   `json:"imgSrc"`
	Job          string   `json:"job"`
	Interestings []string `json:"interestings"`
	Skills       []string `json:"skills"`
}

type MeetCard struct {
	CardId uint      `json:"cardId"`
	Title  string   `json:"title"`
	Text   string   `json:"text"`
	ImgSrc string   `json:"imgSrc"`
	Labels []string `json:"labels"`
	Place  string   `json:"place"`
	Date   string   `json:"date"`
}

type Meeting struct {
	ImgSrc string `json:"imgSrc"`
	Text   string `json:"text"`
}

type UserProfile struct {
	ImgSrc       string    `json:"imgSrc"`
	Name         string    `json:"name"`
	City         string    `json:"city"`
	Telegram     string    `json:"telegram"`
	Vk           string    `json:"vk"`
	Meetings     []Meeting `json:"meetings"`
	Interestings string    `json:"interestings"`
	Skills       string    `json:"skills"`
	Education    string    `json:"education"`
	Job          string    `json:"job"`
	Aims         string    `json:"aims"`
}

type Credentials struct {
	Id       uint    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserId struct {
	Uid uint `json:"userId"`
}

type LogPassPair struct {
	Login 	    string `json:"login"`
	Password 	string `json:"password"`
}

var Sessions = map[string]uint{}

var UserCards = map[uint]UserCard{
	0: {
		CardId:       0,
		Name:         "Александр",
		ImgSrc:       "assets/luckash.jpeg",
		Job:          "Главный чекист КГБ",
		Interestings: []string{"Картофель", "Хоккей"},
		Skills:       []string{"Разгон митингов", "Сбор урожая"},
	},
}

var MeetCards = map[uint]MeetCard{
	0: {
		CardId: 0,
		Title:  "Забив с++",
		Text: "Lorem ipsum dolor sit amet, " +
			"consectetur adipiscing elit, sed " +
			"do eiusmod tempor incididunt ut " +
			"labore et dolore magna aliqua. " +
			"Ut enim ad minim veniam, quis " +
			"nostrud exercitation ullamco labori",
		ImgSrc: "assets/paris.jpg",
		Labels: []string{"Rust", "Забив", "В падике"},
		Place:  "Дом Пушкина, улица Колотушкина",
		Date:   "12 сентября 2020",
	},
}

var UserProfiles = map[uint]UserProfile{
	0: {
		ImgSrc:   "assets/luckash.jpeg",
		Name:     "Александр Лукашенко",
		City:     "Петрозаводск",
		Telegram: "",
		Vk:       "https://vk.com/id241926559",
		Meetings: []Meeting{
			{
				ImgSrc: "assets/vk.png",
				Text:   "Александр Лукашенко",
			},
			{
				ImgSrc: "assets/vk.png",
				Text:   "Александр Лукашенко",
			},
		},
		Interestings: `
                Lorem ipsum dolor sit amet, 
                consectetur adipiscing elit, sed 
                do eiusmod tempor incididunt ut 
                labore et dolore magna aliqua. 
                Ut enim ad minim veniam, quis 
                nostrud exercitation ullamco 
                laboris nisi ut aliquip ex ea 
                commodo consequat. Duis aute 
                irure dolor in reprehenderit 
                in voluptate velit esse cillum 
        `,
		Skills: `Lorem ipsum dolor sit amet, 
                consectetur adipiscing elit, sed 
                do eiusmod tempor incididunt ut 
                labore et dolore magna aliqua. 
                Ut enim ad minim veniam, quis 
                nostrud exercitation ullamco`,
		Education: "МГТУ им. Н. Э. Баумана до 2010",
		Job:       "MAIL GROUP до 2008",
		Aims:      "Хочу от жизни всего",
	},
}

var CredStorage = map[string]Credentials{
	"lukash@mail.ru": {
		Login:    "lukash@mail.ru",
		Password: "12345",
		Id:       0,
	},
}

