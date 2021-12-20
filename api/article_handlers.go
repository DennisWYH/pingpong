package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-pinyin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"pingpong/util"
	"strconv"
	"strings"
)

type Article struct {
	Title   string
	Content string
	Grade   string
	//Tags          []string
	//WordCount     int64
	//NumberOfRead  int64
	//NumberOfFlash int64
	gorm.Model
}

type Lookup struct {
	Hanzi    string
	Pinyin   string
	EnLookup string
	CnLookup string
}

// API: curl localhost:3456/articles
func GetArticles(c *gin.Context) {

	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})

	var articles []Article
	db.Find(&articles)

	//c.IndentedJSON(http.StatusOK, &articles)
	c.HTML(http.StatusOK, "viewArticles.tmpl", gin.H{
		"articles": &articles,
	})
}

// API: curl -X DELETE localhost:3456/article/id/:id
func DeleteArticleByID(c *gin.Context) {
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})

	var article *Article
	id := c.Param("id")
	intID, _ := strconv.Atoi(id)
	db.Delete(&article, intID)

	var articles *[]Article
	db.Find(&articles)
	c.IndentedJSON(http.StatusOK, &articles)
	c.HTML(http.StatusOK, "viewArticles.tmpl", gin.H{
		"articles": &articles,
	})
}

// API: curl localhost:3456/focusedRead
func GetFocusedArticles(c *gin.Context) {
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})

	var article *Article
	db.First(&article)

	lookups := []Lookup{}
	lookup1 := Lookup{"中国", "zhong1 guo2", "china; middle-country", "国名，中国"}
	lookup2 := Lookup{"水杯", "shui3 bei1", "container for water", "装水用的容器"}
	lookups = append(lookups, lookup1)
	lookups = append(lookups, lookup2)
	c.HTML(http.StatusOK, "viewFocusedRead.tmpl", gin.H{
		"article": &article,
		"lookups": &lookups,
	})
}

// API: curl -X DELETE localhost:3456/articles
func DeleteAllArticle(c *gin.Context) {
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})

	// query all the articles
	var articles *[]Article
	db.Find(&articles)

	// delete
	db.Delete(&articles)

	// get all articles to see if there is articles left in the db.
	GetArticles(c)
}

// API: curl -X PUT -d "content=?" localhost:3456/update/article/id/:id
func UpdateArticleByID(c *gin.Context) {
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

// API: localhost:3456/article/id/:id
func GetArticleByID(c *gin.Context) {
	id := c.Param("id")
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

// API: localhost:3456/article/grade/:grade
func GetArticleByGrade(c *gin.Context) {
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
		tokenizedContent := util.Tokenizer(content)
		tokenizedContents = append(tokenizedContents, tokenizedContent)
		words = util.ExtractWords(tokenizedContent)

	}
	for _, word := range words {
		wordsEn := util.Cn_en_lookup(word)
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

// API: curl -X POST -H "Content-Type: application/x-www-form-urlencoded"
//  -d "title=new&content=entry" localhost:3456/addSimpleArticle
// gin context documentation: https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
func AddArticle(c *gin.Context) {
	//Todo: the content from user input has to be chinese,
	// for later pinyin convert.

	//Todo: what if there are English words in the paragraph...
	// solution: create a map, when English recognized, put blank in there or display english itself.
	var newArticle Article

	newArticle.Title = c.PostForm("title")
	newArticle.Content = c.PostForm("content")
	newArticle.Grade = c.PostForm("grade")

	// Add the new article to the db table.
	db, _ := gorm.Open(sqlite.Open("pingpong.db"), &gorm.Config{})
	db.Create(&newArticle)

	// show the article table after adding an entry
	var articles []Article
	db.Find(&articles)

	c.HTML(http.StatusCreated, "viewArticles.tmpl", gin.H{
		"articles": &articles,
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

// API: curl -X POST localhost:3456/batchAddArticles
func BatchAddTestArticleData(c *gin.Context) {
	addTestArticle("第一篇文章", "瑞士政府当地时间17日宣布新的防疫措施，以应对目前严峻的新冠肺炎疫情形势。从本月20日起，未接种疫苗者将不能进入餐馆、酒吧以及文化、体育、休闲等室内公共活动场所；恢复所有人在家工作的要求，一些必须到工作场所进行的工作除外；室内聚会人数不能超过30人，如果聚会中有未接种疫苗者，则不能超过10人。据悉，该措施将持续到明年1月24日。17日，瑞士新增新冠肺炎确诊病例9941例，目前该国累计有294例新冠肺炎患者在医院接受重症监护。瑞士政府担心，随着奥密克戎毒株的传播，医院重症监护病房可能会出现超负荷运转。", "blue")
	addTestArticle("第二篇文章", "我和小丽是好朋友。", "white")
	addTestArticle("第三篇文章", "太阳很晒。", "black")
	GetArticles(c)
}
