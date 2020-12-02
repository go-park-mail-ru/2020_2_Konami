//go:generate easyjson meeting_update.go
package models

type MeetUpdateFields struct {
	Reg  *bool `json:"isRegistered"`
	Like *bool `json:"isLiked"`
}

//easyjson:json
type MeetingUpdate struct {
	MeetId int               `json:"meetId"`
	Fields *MeetUpdateFields `json:"fields"`
}

//easyjson:json
type MeetingData struct {
	Address   *string  `json:"address"`
	City      *string  `json:"city"`
	Start     *string  `json:"start"`
	End       *string  `json:"end"`
	Text      *string  `json:"meet-description"`
	Tags      []string `json:"meetingTags"`
	Title     *string  `json:"name"`
	Photo     *string  `json:"photo"`
	Seats     *int     `json:"seats"`
	SeatsLeft *int     `json:"seatsLeft"`
}
