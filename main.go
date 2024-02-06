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

func middleware(w http.ResponseWriter, sendResquest string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	jsonResq, err := json.Marshal(sendResquest)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jsonResq)
}
func hexToString(w http.ResponseWriter, hashedText []byte) {
	hash := hex.EncodeToString(hashedText[:])
	middleware(w, hash)
}

func sha1Sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.URL.Query().Get("text")
	hashedText := sha1.Sum([]byte(textToHash))
	hexToString(w, hashedText[:])
}

func sha256Sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.URL.Query().Get("text")
	hashedText := sha256.Sum256([]byte(textToHash))
	hexToString(w, hashedText[:])
}

func sha224Sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.URL.Query().Get("text")
	hashedText := sha256.Sum224([]byte(textToHash))
	hexToString(w, hashedText[:])
}

func md5sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.URL.Query().Get("text")
	hashedText := md5.Sum([]byte(textToHash))
	hexToString(w, hashedText[:])
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
	http.ListenAndServe(":"+port, nil)
}
