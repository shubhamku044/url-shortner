package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Url struct {
	Id        string `json:"id"`
	Url       string `json:"url"`
	Key       string `json:"key"`
	CreatedAt string `json:"created_at"`
}

var urls = make(map[string]Url)

func createUrl(url string) string {
	return "http://localhost:8080/re/" + url
}

func getUrl(key string) (Url, error) {
	url, ok := urls[key]
	if !ok {
		return Url{}, errors.New("Url not found")
	}
	return url, nil
}

func generateShortUrl(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)[:6]
	return hash
}

func handler(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Url string `json:"url"`
	}
	var myData data

	err := json.NewDecoder(r.Body).Decode(&myData)
	if err != nil || myData.Url == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	shortUrl := generateShortUrl(myData.Url)
	urls[shortUrl] = Url{
		Id:        shortUrl,
		Url:       myData.Url,
		Key:       shortUrl,
		CreatedAt: time.Now().String(),
	}
	url := createUrl(shortUrl)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"shortned_url": url})
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/re/"):]
	url, err := getUrl(key)
	if err != nil {
		http.Error(w, "Url not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.Url, http.StatusMovedPermanently)
}

func main() {
	fmt.Println("Starting server on port 3000...")
	http.HandleFunc("/", handler)
	http.HandleFunc("/re/", redirectHandler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
