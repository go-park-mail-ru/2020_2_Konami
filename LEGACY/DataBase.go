package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

func InitDB() (h *Handler) {
	db, err := gorm.Open("postgres", DBAuthData)
	if err != nil {
		log.Fatal("Can`t connect to DB")
	}
	db.LogMode(true)
	setup(db)
	h = &Handler{
		DB: db,
	}
	return h
}

/*
func main() {
	handler := InitDB()
	defer handler.DB.Close()
}
*/
