package main

type UserProfile struct {
	ImgSrc    string 	`json:"img_src"`
	Name      string 	`json:"name"`
	City      string 	`json:"city"`
	Telegram  string 	`json:"telegram"`
	Vk        string 	`json:"vk"`
	Meetings  []meeting `json:"meetings"`
	Interest  string	`json:"interest"`
	Skills    string	`json:"skills"`
	Education string	`json:"education"`
	Job       string	`json:"job"`
	Aims      string	`json:"aims"`
}

type meeting struct {
	ImgSrc string	`json:"img_src"`
	Text   string 	`json:"text"`
}

type userCard struct {
	CardId    int      `json:"card_id"`
	ImgSrc    string   `json:"img_src"`
	Name      string   `json:"name"`
	Job       string   `json:"job"`
	Interests []string `json:"interests"`
	Skills    []string `json:"skills"`
}

type meetCard struct {
	ImgSrc		string		`json:"img_src"`
	CardId		int			`json:"card_id"`
	Text		string		`json:"text"`
	Labels		[]string	`json:"labels"`
	Title 		string		`json:"title"`
	Place 		string		`json:"place"`
	Date 		string		`json:"date"`
}

type userChanges struct {
	Field 		string	`json:"field"`
	Text 		string	`json:"text"`
}

type LogPwd struct {
	Login 	string `json:"login"`
	Pwd 	string `json:"pwd"`
}

