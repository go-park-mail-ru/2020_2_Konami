package models

type Message struct {
	Id        int    `json:"id"`
	AuthorId  int    `json:"authorId"`
	MeetingId int    `json:"meetId"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}
