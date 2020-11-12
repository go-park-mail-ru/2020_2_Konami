package models

type Meeting struct {
	Card *MeetingCard `json:"card"`
	Like bool         `json:"isLiked"`
	Reg  bool         `json:"isRegistered"`
}
