package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Slug = string

type RegisteredRedirect struct {
	Slug      Slug      `json:"slug"`
	TargetUrl string    `json:"targetUrl"`
	CreatedAt time.Time `json:"createdAt"`
	Uses      uint64    `json:"uses"`
}

type CreationRequest struct {
	Slug      Slug   `json:"slug"`
	TargetUrl string `json:"targetUrl"`
}

var redirects = make(map[Slug]RegisteredRedirect)

type RedirectAlreadyExistsError struct {
	slug Slug
}

func (r *RedirectAlreadyExistsError) Error() string {
	return fmt.Sprintf("Redirect already exists for slug '%s'", r.slug)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	io.WriteString(w, "Hello, World!")
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	var creationRequest CreationRequest

	err := json.NewDecoder(r.Body).Decode(&creationRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	log.Printf("Received: %+v\n", creationRequest)

	err = createAndInsertRedirect(creationRequest.Slug, creationRequest.TargetUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func createAndInsertRedirect(slug Slug, targetUrl string) error {
	if _, ok := redirects[slug]; ok {
		return &RedirectAlreadyExistsError{
			slug: slug,
		}
	}

	redirects[slug] = RegisteredRedirect{
		Slug:      slug,
		TargetUrl: targetUrl,
		CreatedAt: time.Now(),
		Uses:      0,
	}

	return nil
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(redirects)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func slugRedirectHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	redirect, ok := redirects[slug]
	if !ok {
		http.Error(w, "Redirect not found.", http.StatusNotFound)
		return
	}

	redirect.Uses++
	redirects[slug] = redirect

	http.Redirect(w, r, redirect.TargetUrl, http.StatusMovedPermanently)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/{slug}", slugRedirectHandler)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
