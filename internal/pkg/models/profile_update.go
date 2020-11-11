package models

type ProfileUpdate struct {
	Name        *string  `json:"name"`
	Gender      *string  `json:"gender"`
	City        *string  `json:"city"`
	Birthday    *string  `json:"birthday"`
	Telegram    *string  `json:"telegram"`
	Vk          *string  `json:"vk"`
	MeetingTags []string `json:"meetingTags"`
	Education   *string  `json:"education"`
	Job         *string  `json:"job"`
	Aims        *string  `json:"aims"`
	Interests   *string  `json:"interests"`
	Skills      *string  `json:"skills"`
}
