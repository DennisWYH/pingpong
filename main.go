package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"pingpong/db"
)

// to be implemented later
//func countGrade() string {
//}

//func deleteKey(key []byte) error {
//
//}
func runServer() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	router.LoadHTMLGlob("templates/*")
	router.GET("/home", func(c *gin.Context) {
		// Call the HTML method of the Context to render a template
		c.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the home.html template
			"home.html",
			// Pass the data that the page uses (in this case, 'title')
			gin.H{
				"title": "Home Page",
				"page":  "Home Page",
			},
		)
	})
	router.Run(":8083")
}

// lookupC2C given a chinese word
// look it up in word.json and return the explanation
func lookupC2C(c string) {
	dictionaryPath := "./chinese-xinhua/data/word.json"
	file, _ := ioutil.ReadFile(dictionaryPath)
	var data = make([]Word, 100)
	err := json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(data); i++ {
		if data[i].Word == c {
			fmt.Println(data[i].Explanation)
		}
	}

}

type Word struct {
	Word        string `json:"word"`
	Oldword     string `json:"oldword"`
	Strokes     string `json:"strokes"`
	Pinyin      string `json:"pinyin"`
	Radicals    string `json:"radicals"`
	Explanation string `json:"explanation"'`
	More        string `json:"more"`
}

func Init() {
	// create pingpong db, update an articles bucket, add some data and delete the data.
	db.RunBBoltDB()
	// the server shall keep running persistently.
	//runServer()
	// update an articles bucket with some data, and print it out.
	db.CreateArticle("美好一天", "今天我非常高兴。", "whileA", []string{"日记", "生活"})
	//delete an article key-value pair
	//db.DeleteArticle("articles", "15")
	db.DeleteBucket("articles")

}

func main() {
	Init()
}
