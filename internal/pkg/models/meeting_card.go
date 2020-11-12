package models

type MeetingCard struct {
	Label     *MeetingLabel `json:"label"`
	AuthorId  int           `json:"authorId"`
	Text      string        `json:"text"`
	Tags      []*Tag        `json:"tags"`
	Address   string        `json:"address"`
	City      string        `json:"city"`
	StartDate string        `json:"startDate"`
	EndDate   string        `json:"endDate"`
	Seats     int           `json:"seats"`
	SeatsLeft int           `json:"seatsLeft"`
}
