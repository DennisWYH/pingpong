package db

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"pingpong/util"
	"strconv"
	"time"
)

type Article struct {
	ID            int64
	Title         string
	Content       string
	CreationDate  time.Time
	Tags          []string
	WordCount     int64
	Grade         string
	NumberOfRead  int64
	NumberOfFlash int64
}

func RunBBoltDB() {
	// open pingpong.db bbolt db
	db, err := bolt.Open("pingpong.db", 0666, nil)
	if err != nil {
		fmt.Println(err)
	}
	// update the db with article bucket
	db.Update(func(tx *bolt.Tx) error {
		// create articles bucket
		_, err := tx.CreateBucketIfNotExists([]byte("articles"))
		if err != nil {
			return err
		}
		return nil
	})
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("articles"))
		v := b.Get([]byte("2015-01-01"))
		fmt.Printf("%sn", v)
		tx.DeleteBucket([]byte("articles"))
		return nil
	})
	fmt.Println("closing db now")
	defer db.Close()
}

// CreateArticle to insert an article key value pair into articles bucket
func CreateArticle(title string, content string, grade string, tags []string) error {
	art := Article{}
	db, err := bolt.Open("pingpong.db", 0666, nil)
	if err != nil {
		fmt.Println(err)
	}
	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the articles bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte("articles"))
		fmt.Println("bucket article created.")
		// Generate ID for the article.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()
		fmt.Println("generate id, ", id)
		art.ID = int64(id)
		art.Title = title
		fmt.Println("title is, ", art.Title)
		art.Tags = tags
		art.Content = content
		fmt.Println("content is, ", art.Content)
		art.CreationDate = time.Now()
		art.WordCount = util.CountWords(content)
		art.Grade = grade
		art.NumberOfRead = 0
		art.NumberOfFlash = 0
		// Marshal article data into bytes.
		artBuf, err := json.Marshal(art)
		if err != nil {
			return err
		}
		// Persist bytes to articles bucket.
		fmt.Println("marchaled article struct", artBuf)
		err = b.Put([]byte(strconv.Itoa(int(art.ID))), artBuf)

		// retrieve the inserted value
		articleBucket := tx.Bucket([]byte("articles"))
		fmt.Println("articleBucket", articleBucket)
		value := articleBucket.Get([]byte(strconv.Itoa(int(art.ID))))
		fmt.Println("marchaled article struct", value)
		article := Article{}
		json.Unmarshal(value, &article)
		fmt.Println("unmarshalled article struct data", article)
		fmt.Println("article title", article.Title)
		return nil
	})
}

// DeleteArticle deletes an article key-value pair, given both bucket name and id key of the article.
func DeleteArticle(bucket string, key string) error {
	keyByte := []byte(key)
	db, err := bolt.Open("pingpong.db", 0666, nil)
	if err != nil {
		fmt.Println(err)
	}
	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the articles bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte(bucket))
		err := b.Delete(keyByte)
		if err != nil {
			return err
		}
		return nil
	})
}

func DeleteBucket(bucketName string) error {
	db, err := bolt.Open("pingpong.db", 0666, nil)
	if err != nil {
		fmt.Println(err)
	}
	return db.Update(func(tx *bolt.Tx) error {
		// Retrieve the articles bucket.
		// This should be created when the DB is first opened.
		err := tx.DeleteBucket([]byte(bucketName))
		if err != nil {
			return err
		}
		return nil
	})
}
