package main

import (
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
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

func migrateDBScheme() (db *gorm.DB) {
	// gorm postgres driver
	dsnDefinition := "host=localhost user=postgres password= dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsnDefinition), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	// Use pq package for array support in the db field
	// https://stackoverflow.com/questions/63256680/adding-an-array-of-integers-as-a-data-type-in-a-gorm-model
	type ChineseSentence struct {
		gorm.Model
		difficultyLevel int
		Chinese         string
		English         string
		Pinyin          string
		PinyinSlice     pq.StringArray `gorm:"type:text[]"`
	}
	db.AutoMigrate(&ChineseSentence{})
	return db
}

func openAndConnectToDB() (db *gorm.DB) {
	// gorm postgres driver
	dsnDefinition := "host=localhost user=postgres password= dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsnDefinition), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
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
		db := openAndConnectToDB()
		// Create
		db.Create(&ChineseSentence{difficultyLevel: 1, Chinese: "中文第二课", English: "First Chinese lesson",
			Pinyin: "zhong wen di yi ke", PinyinSlice: []string{"zhong", "wen", "di", "yi", "ke"}})
		w.Write([]byte("added sentence"))
	})

	http.HandleFunc("/list-sentence", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		db := openAndConnectToDB()
		type ChineseSentence struct {
			gorm.Model
			difficultyLevel int
			Chinese         string
			English         string
			Pinyin          string
			PinyinSlice     pq.StringArray `gorm:"type:text[]"`
		}
		chineseSentences := &[]ChineseSentence{}
		//db.First(&chineseSentence)
		db.Find(&chineseSentences)
		// Response with json
		// https://stackoverflow.com/questions/31622052/how-to-serve-up-a-json-response-using-go
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(&chineseSentences)
		if err != nil {
			panic("json failed to marshal data")
		}
		w.Write(data)
	})
	//http.HandleFunc("/remove-sentence", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("hello.world"))
	//})

	// Read this on heroku dynamic port number
	// https://stackoverflow.com/questions/56936448/deploying-a-golang-app-on-heroku-build-succeed-but-application-error
	//port := os.Getenv("PORT")
	//if err := http.ListenAndServe(":"+port, nil); err != nil {
	//	log.Fatal(err)
	//}
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
