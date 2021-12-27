package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"pingpong/api"
)

// AddArticleTableEntry adds a db entry to Article table and returns the primary key: ID
func AddArticleTableEntry(title, content, grade string) (articleID int) {
	var newArticle api.Article
	newArticle.Title = title
	newArticle.Content = content
	newArticle.Grade = grade

	// Add the new article to the db table.
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})
	db.Create(&newArticle)
	return int(newArticle.ID)
}

// AddLookupTableEntry adds a db entry to Lookup table and returns the primary key: ID
func AddLookupTableEntry(hanzi, pinyin, enLookup string, articleID int) (lookupID int) {
	var newLookup api.Lookup
	newLookup.Hanzi = hanzi
	newLookup.Pinyin = pinyin
	newLookup.EnLookup = enLookup
	newLookup.ArticleID = articleID

	// Add the new article to the db table.
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})
	db.Create(&newLookup)
	return int(newLookup.ID)
}

// CreateDBTables creates and migrates Article and Lookup tables
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
