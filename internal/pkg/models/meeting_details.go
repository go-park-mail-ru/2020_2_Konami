package models

type MeetingDetails struct {
	Card          *MeetingCard    `json:"card"`
	Like          bool            `json:"isLiked"`
	Reg           bool            `json:"isRegistered"`
	Registrations []*ProfileLabel `json:"registrations"`
}
