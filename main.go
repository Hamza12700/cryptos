package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func hasher(w http.ResponseWriter, hashedText []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	hash := hex.EncodeToString(hashedText[:])
	jsonResq, err := json.Marshal(hash)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jsonResq)
}

func sha1Sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.URL.Query().Get("text")
	hashedText := sha1.Sum([]byte(textToHash))
	hasher(w, hashedText[:])
}

func sha256Sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.URL.Query().Get("text")
	hashedText := sha256.Sum256([]byte(textToHash))
	hasher(w, hashedText[:])
}

func sha224Sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.URL.Query().Get("text")	
	hashedText := sha256.Sum224([]byte(textToHash))
	hasher(w,hashedText[:])
}

func md5sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.URL.Query().Get("text")
	hashedText := md5.Sum([]byte(textToHash))
	hasher(w, hashedText[:])
}

func main() {
	http.HandleFunc("/md5", md5sum)
	http.HandleFunc("/sha1", sha1Sum)
	http.HandleFunc("/sha256", sha256Sum)
	http.HandleFunc("/sha224", sha224Sum)
	port := os.Getenv("PORT")
	if port == "" {
		port = "2323"
	}
	fmt.Println("Listening on", port)
	http.ListenAndServe(port, nil)
}
