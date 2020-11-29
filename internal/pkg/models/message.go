package models

type Message struct {
	Id        int    `json:"id"`
	AuthorId  int    `json:"authorId"`
	MeetingId string `json:"meetId"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}
