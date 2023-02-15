package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

type ChineseSentence struct {
	gorm.Model
	difficultyLevel int
	Chinese         string
	English         string
	Pinyin          string
	PinyinSlice     []string
}

func main() {
	// APIs for frontend actions
	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		fmt.Fprintf(w, "Hello, World!")
	})
	http.HandleFunc("/previous", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		fmt.Fprintf(w, "Hello, World!")
	})
	http.HandleFunc("/helloworld", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		w.Write([]byte("hello.world"))
	})

	// APIs for database CRUD management
	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		w.Write([]byte("this page displays a form where user can add sentences"))
	})
	http.HandleFunc("/add-sentence", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		// gorm postgres driver
		dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		// Migrate the schema
		type ChineseSentence struct {
			gorm.Model
			difficultyLevel int
			Chinese         string
			English         string
			Pinyin          string
			PinyinSlice     []string
		}

		db.AutoMigrate(&ChineseSentence{})
		// Create
		db.Create(&ChineseSentence{difficultyLevel: 1, Chinese: "中文第一课", English: "First Chinese lesson",
			Pinyin: "zhong wen di yi ke", PinyinSlice: []string{"zhong", "wen", "di", "yi", "ke"}})
		w.Write([]byte("added sentence"))
	})

	http.HandleFunc("/list-sentence", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		var chineseSentence ChineseSentence
		db.First(&chineseSentence, 1) // find chinese sentence with integer primary key

		w.Write([]byte(chineseSentence.Chinese))
	})
	//http.HandleFunc("/remove-sentence", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("hello.world"))
	//})

	// Read this on heroku dynamic port number
	// https://stackoverflow.com/questions/56936448/deploying-a-golang-app-on-heroku-build-succeed-but-application-error
	port := os.Getenv("PORT")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
