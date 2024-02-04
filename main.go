package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

const port string = ":2323"

func corsMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

		next(w, r)
	}
}

func sha1Sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.URL.Query().Get("text")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	hashedText := sha1.Sum([]byte(textToHash))
	sha1Hash := hex.EncodeToString(hashedText[:])
	jsonResq, err := json.Marshal(sha1Hash)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jsonResq)
}

func sha256Sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.URL.Query().Get("text")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	hashedText := sha256.Sum256([]byte(textToHash))
	sha1Hash := hex.EncodeToString(hashedText[:])
	jsonResq, err := json.Marshal(sha1Hash)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jsonResq)
}

func md5sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.URL.Query().Get("text")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	hashedText := md5.Sum([]byte(textToHash))
	md5Hashed := hex.EncodeToString(hashedText[:])
	jsonResq, err := json.Marshal(md5Hashed)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jsonResq)
}

func main() {
	http.HandleFunc("/md5sum", corsMiddleWare(md5sum))
	http.HandleFunc("/sha1sum", corsMiddleWare(sha1Sum))
	http.HandleFunc("/sha256sum", corsMiddleWare(sha256Sum))
	fmt.Println("Listening on", port)
	http.ListenAndServe(port, nil)
}
