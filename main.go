package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/helloworld", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	fmt.Printf("Server running (port=8080), route: http://localhost:8080/helloworld\n")
	port := os.Getenv("PORT")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
