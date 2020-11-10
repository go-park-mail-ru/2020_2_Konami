package models

type Profile struct {
	Card        *ProfileCard    `json:"card"`
	Gender      string          `json:"gender"`
	Birthday    string          `json:"birthday"`
	City        string          `json:"city"`
	Login       string          `json:"login"`
	PwdHash     string          `json:"-"`
	Telegram    string          `json:"telegram"`
	Vk          string          `json:"vk"`
	Education   string          `json:"education"`
	MeetingTags []*Tag          `json:"meetingTags"`
	Aims        string          `json:"aims"`
	Interests   string          `json:"interests"`
	Skills      string          `json:"skills"`
	Meetings    []*MeetingLabel `json:"meetings"`
}
