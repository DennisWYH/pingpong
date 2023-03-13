package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

// Global Db connection that can be shared among different handlers.
var dbConnection *gorm.DB

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

type ChineseSentence struct {
	Id uint64 `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	gorm.Model
	DifficultyLevel    int
	Chinese            string
	EnglishTranslation string
	Pinyin             string
	//PinyinSlice        pq.StringArray `gorm:"type:text[]"`
}

func migrateDBScheme() (db *gorm.DB) {
	// gorm postgres driver
	dsnDefinition := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsnDefinition), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	// Use pq package for array support in the db field
	// https://stackoverflow.com/questions/63256680/adding-an-array-of-integers-as-a-data-type-in-a-gorm-model
	type ChineseSentence struct {
		Id uint64 `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
		gorm.Model
		DifficultyLevel    int
		Chinese            string
		EnglishTranslation string
		Pinyin             string
		//pinyinSlice        pq.StringArray `gorm:"type:text[]"`
	}
	db.AutoMigrate(&ChineseSentence{})
	return db
}

func openAndConnectToDB() (db *gorm.DB) {
	// gorm postgres driver
	// https://medium.com/@awshawka12/go-postgres-heroku-tutorial-801b3d69e116
	dsnDefinition := os.Getenv("DATABASE_URL")
	//dsnDefinition := "host=localhost user=postgres password= dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsnDefinition), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func validateSentenceData(data ChineseSentence) {
	fmt.Println("Validating: The chinese is:", data.Chinese)
	fmt.Println("Validating: The english-translation is:", data.EnglishTranslation)
	fmt.Println("Validating: The difficulty-level is:", data.DifficultyLevel)
	fmt.Println("Validating: The pinyins are:", data.Pinyin)
}

func main() {
	// APIs for webinterface
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		chineseSentence := &ChineseSentence{}
		// Get one record, no specified order
		dbConnection.Take(&chineseSentence)
		// SELECT * FROM users LIMIT 1;		// Response with json
		// https://stackoverflow.com/questions/31622052/how-to-serve-up-a-json-response-using-go
		w.Header().Set("Content-Type", "application/json")
		marshaledData, err := json.Marshal(&chineseSentence)
		if err != nil {
			panic("json failed to marshal data")
		}
		w.Write(marshaledData)
	})

	http.HandleFunc("/getById", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		type ChineseSentence struct {
			ID uint64 `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
			gorm.Model
			DifficultyLevel    int
			Chinese            string
			EnglishTranslation string
			Pinyin             string
			//PinyinSlice        pq.StringArray `gorm:"type:text[]"`
		}
		chineseSentence := &ChineseSentence{}
		id := r.URL.Query().Get("id")
		// Get one record
		err := dbConnection.First(&chineseSentence, id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If it's already the last ID, we query from the 1st one again.
			dbConnection.First(&chineseSentence, 1)
		}
		// SELECT * FROM users LIMIT 1;		// Response with json
		// https://stackoverflow.com/questions/31622052/how-to-serve-up-a-json-response-using-go
		w.Header().Set("Content-Type", "application/json")
		marshaledData, err := json.Marshal(&chineseSentence)
		if err != nil {
			panic("json failed to marshal data")
		}
		w.Write(marshaledData)
	})

	http.HandleFunc("/add-sentence", func(w http.ResponseWriter, r *http.Request) {
		// Test curl
		// curl -v -X POST http://localhost:8080/add-sentence -d '{"chinese":"中文第二课", "pinyin": "testpinyin",
		// "englishTranslation":"test", "difficultyLevel":"9"}'

		// Allowing cross-domain request
		enableCors(&w)

		var data ChineseSentence
		if r.Body != nil {
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&data)
			if err != nil {
				fmt.Println("an error has occured while decoding request body: ", err)
			}
			fmt.Println("The data decoded from http request body is:", data)
		}
		validateSentenceData(data)
		//Create an entry in the db
		dbConnection.Create(&ChineseSentence{DifficultyLevel: data.DifficultyLevel,
			Chinese: data.Chinese, EnglishTranslation: data.EnglishTranslation,
			Pinyin: data.Pinyin},
		)
		w.WriteHeader(http.StatusCreated)
		marshaledData, err := json.Marshal(&data)
		if err != nil {
			fmt.Println("an error has occured while marshalling data: ", err)
		}
		w.Write(marshaledData)
	})

	http.HandleFunc("/list-sentence", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		type ChineseSentence struct {
			gorm.Model
			DifficultyLevel    int
			Chinese            string
			EnglishTranslation string
			Pinyin             string
			//PinyinSlice        pq.StringArray `gorm:"type:text[]"`
		}
		chineseSentences := &[]ChineseSentence{}
		//db.First(&chineseSentence)
		dbConnection.Find(&chineseSentences)
		// Response with json
		// https://stackoverflow.com/questions/31622052/how-to-serve-up-a-json-response-using-go
		w.Header().Set("Content-Type", "application/json")
		marshaledData, err := json.Marshal(&chineseSentences)
		if err != nil {
			panic("json failed to marshal data")
		}
		w.Write(marshaledData)
	})

	// Database schema migration
	//migrateDBScheme()

	// Read this on heroku dynamic port number
	// https://stackoverflow.com/questions/56936448/deploying-a-golang-app-on-heroku-build-succeed-but-application-error
	port := os.Getenv("PORT")
	dbConnection = openAndConnectToDB()
	// When testing locally port number is 8080
	//if err := http.ListenAndServe(":8080", nil); err != nil {
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

	// Documentations:
	// How to import pg.dump file into heroku's postgres db adds-on
	// heroku pg:backups:restore "s3StorageAddress" --app pingpong-fun
}
