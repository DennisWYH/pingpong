package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"pingpong/api"
)

func CreateDBTables() {
	db, err := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	var article *api.Article
	db.AutoMigrate(&article)
	db.Create(&article)

	// Some look up entries belong to an article
	var lookup *api.Lookup
	db.AutoMigrate(&lookup)
	db.Create(&lookup)
}
