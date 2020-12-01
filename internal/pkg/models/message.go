//go:generate easyjson message.go
package models

//easyjson:json
type Message struct {
	Id        int    `json:"id"`
	AuthorId  int    `json:"authorId"`
	MeetingId int    `json:"meetId"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}
