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

	var article *Article
	db.AutoMigrate(&article)
	db.Create(&article)
}
