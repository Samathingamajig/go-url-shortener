package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

type Slug = string

type RegisteredRedirect struct {
	slug      Slug
	targetUrl string
	createdAt time.Time
	uses      uint64
}

var redirects = make(map[Slug]RegisteredRedirect)

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
