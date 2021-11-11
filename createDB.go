package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateDBTables() {
	db, err := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	article := Article{
		Title: "title",
		Content: "content",
	}
	db.AutoMigrate(&article)

	db.Create(&article)
}
