package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Article struct {
	Title   string
	Content string
	Grade   string
	Tokens  string
	Pinyins string
	gorm.Model
	//Tags          []string
	//WordCount     int64
	//NumberOfRead  int64
	//NumberOfFlash int64
}

type Lookup struct {
	Hanzi     string
	Pinyin    string
	EnLookup  string
	ArticleID int
	Article   Article
	gorm.Model
}

// AddArticleTableEntry adds a db entry to Article table and returns the primary key: ID
func AddArticleTableEntry(title, content, grade string) (articleID int) {
	var newArticle Article
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
	var newLookup Lookup
	newLookup.Hanzi = hanzi
	newLookup.Pinyin = pinyin
	newLookup.EnLookup = enLookup
	newLookup.ArticleID = articleID

	// Add the new Lookup struct to the db table.
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

	var article *Article
	err = db.AutoMigrate(&article)
	if err != nil {
		fmt.Println(err)
	}
	db.Create(&article)

	// Some look up entries belong to an article
	var lookup *Lookup
	err = db.AutoMigrate(&lookup)
	if err != nil {
		fmt.Println(err)
	}
	db.Create(&lookup)
}
