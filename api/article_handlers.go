package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-pinyin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"pingpong/database"
	"pingpong/util"
	"strconv"
	"strings"
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

// GetArticlesHandler returns all the articles in the Article table
// API: curl localhost:3456/articles
func GetArticlesHandler(c *gin.Context) {

	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})

	var articles []Article
	db.Find(&articles)

	//c.IndentedJSON(http.StatusOK, &articles)
	c.HTML(http.StatusOK, "viewArticles.html", gin.H{
		"articles": &articles,
	})
}

// GetLookups returns all the lookups data in the Lookup table
// API: curl localhost:3456/lookups
func GetLookups() {

	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})

	var lookups []Lookup
	db.Find(&lookups)

	for _, lookup := range lookups {
		fmt.Println("article id = ", lookup.ArticleID)
		fmt.Println("article Hanzi = ", lookup.Hanzi)
	}
}

// DeleteArticleByIDHandler deletes article given the article ID
// API: curl -X DELETE localhost:3456/article/id/:articleID
func DeleteArticleByIDHandler(c *gin.Context) {
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})

	var article *Article
	id := c.Param("articleID")
	intID, _ := strconv.Atoi(id)
	db.Delete(&article, intID)

	var articles *[]Article
	db.Find(&articles)
	c.IndentedJSON(http.StatusOK, &articles)
	c.HTML(http.StatusOK, "viewArticles.html", gin.H{
		"articles": &articles,
	})
}

// GetFocusedArticlesHandler handles the request and renders viewFocusedRead html
// API: curl localhost:3456/focusedRead/id/:articleID
func GetFocusedArticlesHandler(c *gin.Context) {
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})
	id := c.Param("articleID")
	intID, _ := strconv.Atoi(id)

	var article *Article
	db.First(&article, "ID=?", intID)

	var lookups []Lookup
	db.Where("article_id = ?", article.ID).Find(&lookups)

	// for each article content, we first tokenize it
	tokens, err := util.Tokenizer(article.Content)
	pinyins := Tokens_to_pinyins(tokens)

	if err != nil {
		fmt.Print("There is an error in tokenizing the article content", err)
	}

	c.HTML(http.StatusOK, "viewFocusedRead.html", gin.H{
		"tokens":  tokens,
		"pinyins": pinyins,
		"article": &article,
		"lookups": &lookups,
	})
}

// DeleteAllArticleHandler deletes all articles for testing purpose
// API: curl -X DELETE localhost:3456/articles
func DeleteAllArticleHandler(c *gin.Context) {
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})

	// query all the articles
	var articles *[]Article
	db.Find(&articles)

	// delete
	db.Delete(&articles)

	// get all articles to see if there is articles left in the db.
	GetArticlesHandler(c)
}

// UpdateArticleByIDHandler updates an article given its ID
// API: curl -X PUT -d "content=?" localhost:3456/update/article/id/:id
func UpdateArticleByIDHandler(c *gin.Context) {
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})

	var article *Article
	id := c.Param("id")
	intID, _ := strconv.Atoi(id)

	db.First(&article, "ID=?", intID)

	article.Content = c.PostForm("content")
	db.Save(&article)

	c.IndentedJSON(http.StatusOK, &article)
	c.HTML(http.StatusOK, "viewArticles.html", gin.H{
		"articles": &article,
	})
}

// GetArticleByIDHanlder returns article given its ID
// API: localhost:3456/article/id/:id
func GetArticleByIDHandler(c *gin.Context) {
	id := c.Param("articleID")
	intID, _ := strconv.Atoi(id)
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})

	var article *Article
	db.First(&article, "ID=?", intID)

	articleStruct := *article
	content := articleStruct.Content
	a := pinyin.NewArgs()
	a.Style = pinyin.Tone
	contentPinyins := pinyin.Pinyin(content, a)
	fmt.Println("pin yin is ", contentPinyins)

	slicedContent := strings.Split(content, "")

	hanziPinyins := make(map[string][]string)
	for i := 0; i < len(slicedContent); i++ {
		key := slicedContent[i]
		value := contentPinyins[i]
		hanziPinyins[key] = value
	}

	c.HTML(http.StatusOK, "viewArticleById.html", gin.H{
		"hanzi":        content,
		"hanziPinyins": hanziPinyins,
	})
}

// GetArticleByGradeHandler returns the articles given by the grade
// API: localhost:3456/article/grade/:grade
func GetArticleByGradeHandler(c *gin.Context) {
	grade := c.Param("grade")
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})

	var articles *[]Article
	db.Find(&articles, "Grade=?", grade)

	var hanzis []string
	var pinyins []string
	var tokenizedContents [][]string
	var words []string
	var wordsEns [][]string

	for _, article := range *articles {
		articleStruct := article
		content := articleStruct.Content

		hanzis = append(hanzis, content)
		pinyins = append(pinyins, util.HanziToPinyins(content))
		tokenizedContent, err := util.Tokenizer(content)
		if err != nil {
			c.Error(err)
		}
		tokenizedContents = append(tokenizedContents, tokenizedContent)
		words = util.ExtractWords(tokenizedContent)

	}
	for _, word := range words {
		wordsEn, err := util.Cn_en_lookup(word)
		if err != nil {
			fmt.Println("")
		}
		wordsEns = append(wordsEns, wordsEn)
	}

	c.HTML(http.StatusOK, "viewArticleByGrade.html", gin.H{
		"hanzis":            hanzis,
		"pinyins":           pinyins,
		"tokenizedContents": tokenizedContents,
		"words":             words,
		"wordsEns":          wordsEns,
	})
}

func Tokens_to_pinyins(tokens []string) []string {
	var pinyins []string
	for _, val := range tokens {
		pinyin := util.HanziToPinyins(val)
		pinyins = append(pinyins, pinyin)
	}
	return pinyins
}

// AddArticleHandler addes entry to the article table as well as lookup table
// API: curl -X POST -H "Content-Type: application/x-www-form-urlencoded"
//  -d "title=new&content=entry&grade=white" localhost:3456/addSimpleArticle
// gin context documentation: https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
func AddArticleHandler(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	grade := c.PostForm("grade")

	// for each article content, we first tokenize it
	tokens, err := util.Tokenizer(content)
	if err != nil {
		fmt.Print("There is an error in tokenizing the article content", err)
	}

	articleID := database.AddArticleTableEntry(title, content, grade)

	// for the tokens []string slice, get rid of the entries if they are symbols.
	tokensWithoutSymbols := []string{}
	for _, token := range tokens {
		if !util.CheckIfSymbols(token) {
			tokensWithoutSymbols = append(tokensWithoutSymbols, token)
		}
	}

	// for each of the token (hanzi), we find out its pinyin and its enLookup
	// then we save the lookup entry into lookup table
	for _, hanzi := range tokensWithoutSymbols {
		pinyin := util.HanziToPinyins(hanzi)
		enLookup, err := util.Cn_en_lookup(hanzi)
		if err == nil {
			firstEnLookup := enLookup[0]
			database.AddLookupTableEntry(hanzi, pinyin, firstEnLookup, articleID)
		}
	}
	// display the article and the lookup in viewFocusedRead.templ
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})
	var article Article
	db.Where("ID", articleID).Find(&article)

	var lookups []Lookup
	db.Where("article_id", articleID).Find(&lookups)

	pinyins := Tokens_to_pinyins(tokens)

	c.HTML(http.StatusCreated, "viewFocusedRead.html", gin.H{
		"articles": &article,
		"lookups":  &lookups,
		"tokens":   tokens,
		"pinyins":  pinyins,
	})
}
