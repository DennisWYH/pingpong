package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-pinyin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"pingpong/database"
	"pingpong/util"
	"strconv"
	"strings"
)

type Article struct {
	Title   string
	Content string
	Grade   string
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
	c.HTML(http.StatusOK, "viewArticles.tmpl", gin.H{
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
	c.HTML(http.StatusOK, "viewArticles.tmpl", gin.H{
		"articles": &articles,
	})
}

// GetFocusedArticlesHandler handles the request and renders viewFocusedRead tmpl
// API: curl localhost:3456/focusedRead/id/:articleID
func GetFocusedArticlesHandler(c *gin.Context) {
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})
	id := c.Param("articleID")
	intID, _ := strconv.Atoi(id)

	var article *Article
	db.First(&article, "ID=?", intID)

	var lookups []Lookup
	var lookup *Lookup
	fmt.Println("article ID is ", article.ID)
	fmt.Println("the model of lookup table is, ", db.Model(&lookup))
	db.Where("article_id = ?", article.ID).Find(&lookups)

	fmt.Println("article lookup is", lookups)

	c.HTML(http.StatusOK, "viewFocusedRead.tmpl", gin.H{
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
	c.HTML(http.StatusOK, "viewArticles.tmpl", gin.H{
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

	c.HTML(http.StatusOK, "viewArticleById.tmpl", gin.H{
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

	c.HTML(http.StatusOK, "viewArticleByGrade.tmpl", gin.H{
		"hanzis":            hanzis,
		"pinyins":           pinyins,
		"tokenizedContents": tokenizedContents,
		"words":             words,
		"wordsEns":          wordsEns,
	})
}

// AddArticleHandler addes entry to the article table as well as lookup table
// API: curl -X POST -H "Content-Type: application/x-www-form-urlencoded"
//  -d "title=new&content=entry&grade=white" localhost:3456/addSimpleArticle
// gin context documentation: https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
func AddArticleHandler(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	grade := c.PostForm("grade")
	articleID := database.AddArticleTableEntry(title, content, grade)

	// for each article content, we first tokenize it
	tokens, err := util.Tokenizer(content)
	if err != nil {
		c.Error(err)
	}

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
		if err != nil {
			fmt.Printf("error in cn_en lookup, there is no result in lookup %s", err)
		} else {
			firstEnLookup := enLookup[0]
			database.AddLookupTableEntry(hanzi, pinyin, firstEnLookup, articleID)
		}
	}

	// display the article and the lookup in viewFocusedRead.templ
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})
	var articles []Article
	db.Where("ID", articleID).Find(&articles)

	var lookups []Lookup
	db.Where("article_id", articleID).Find(&lookups)

	c.HTML(http.StatusCreated, "viewFocusedRead.tmpl", gin.H{
		"articles": &articles,
		"lookups":  &lookups,
	})
}

func addTestArticle(title, content, grade string) {
	requestURL := "http://localhost:3456/addSimpleArticle"
	requestForm := url.Values{}
	requestForm.Add("title", title)
	requestForm.Add("content", content)
	requestForm.Add("grade", grade)
	req, err := http.NewRequest(http.MethodPost, requestURL, strings.NewReader(requestForm.Encode()))
	if err != nil {
		fmt.Println("BatchAddTestArticleData Error in request:", err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

// BatchAddTestArticleDataHandler adds some test articles for testing
// API: curl -X POST localhost:3456/batchAddArticles
func BatchAddTestArticleDataHandler(c *gin.Context) {
	addTestArticle("第一篇文章", "瑞士政府当地时间17日宣布新的防疫措施，以应对目前严峻的新冠肺炎疫情形势。从本月20日起，"+
		"未接种疫苗者将不能进入餐馆、酒吧以及文化、体育、休闲等室内公共活动场所；恢复所有人在家工作的要求，一些必须到工作场所进行的工作除外；"+
		"室内聚会人数不能超过30人，如果聚会中有未接种疫苗者，则不能超过10人。据悉，该措施将持续到明年1月24日。17日，"+
		"瑞士新增新冠肺炎确诊病例9941例，目前该国累计有294例新冠肺炎患者在医院接受重症监护。瑞士政府担心，随着奥密克戎毒株的传播，"+
		"医院重症监护病房可能会出现超负荷运转。", "blue")
	addTestArticle("第二篇文章", "我和小丽是好朋友。", "white")
	addTestArticle("第三篇文章", "太阳很晒。", "black")
	GetArticlesHandler(c)
}

func addTestLookup(hanzi string, pinyin string, enLookup string, articleID int) {
	var newLookup Lookup
	newLookup.Hanzi = hanzi
	newLookup.Pinyin = pinyin
	newLookup.EnLookup = enLookup
	newLookup.ArticleID = articleID

	// Add the newLookup to the db Lookup table.
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})
	db.Create(&newLookup)

	// show the article table after adding an entry
	var lookups []Lookup
	db.Find(&lookups)

	fmt.Println(&lookups)
}

// BatchAddTestLookupData adds some test lookup data for testing
// API: curl -X POST localhost:3456/batchAddLookup
func BatchAddTestLookupData() {
	addTestLookup("文章", "wen zhang", "article", 8)
	addTestLookup("我", "wo", "me", 8)
	addTestLookup("内容", "nei rong", "content", 8)
	GetLookups()
}
