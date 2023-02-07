package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// APIs for frontend actions
	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.HandleFunc("/previous", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	http.HandleFunc("/helloworld", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello.world"))
	})

	// APIs for database CRUD management
	http.HandleFunc("/add-sentence", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this page displays a form where user can add sentences"))
	})

	//http.HandleFunc("/remove-sentence", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("hello.world"))
	//})
	//http.HandleFunc("/list-sentence", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("hello.world"))
	//})

	// Read this on heroku dynamic port number
	// https://stackoverflow.com/questions/56936448/deploying-a-golang-app-on-heroku-build-succeed-but-application-error
	port := os.Getenv("PORT")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
