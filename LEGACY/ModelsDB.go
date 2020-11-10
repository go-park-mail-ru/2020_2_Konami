package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
	"time"
)

type Handler struct {
	DB *gorm.DB
}

type User struct {
	gorm.Model
	Name      string
	Gender    string
	City      string
	Email     string
	Telegram  string
	Vk        string
	Education string
	Job       string
	ImgSrc    string
	Aims      string
	Interest  pq.StringArray `gorm:"type:varchar[]"`
	Skills    pq.StringArray `gorm:"type:varchar[]"`
	Meetings  pq.Int64Array  `gorm:"type:integer[]"`
}

type Meeting struct {
	gorm.Model
	Title  string
	Text   string
	ImgSrc string
	Place  string
	Tags   pq.StringArray `gorm:"type:varchar[]"`
	Date   time.Time
}

type Session struct {
	gorm.Model
	UserID int
	Token  string
}

type UserOnMeeting struct {
	gorm.Model
	UserID    int
	MeetingID int
}

type UserVote struct {
	gorm.Model
	UserID    int
	MeetingID int
}

func setup(db *gorm.DB) {
	db.AutoMigrate(&Meeting{},
		&User{},
		&UserOnMeeting{},
		&Session{},
		&UserVote{})
}
