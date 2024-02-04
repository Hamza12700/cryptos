package main

import (
	"crypto/md5"
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

func md5sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.URL.Query().Get("text")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	hashedText := md5.Sum([]byte(textToHash))
	md5Hashed := hex.EncodeToString(hashedText[:])
	resq := make(map[string]string)
	resq["md5sum"] = md5Hashed
	jsonResq, err := json.Marshal(resq)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jsonResq)
}

func main() {
	http.HandleFunc("/md5sum", corsMiddleWare(md5sum))
	fmt.Printf("Listening on %s", port)
	http.ListenAndServe(port, nil)
}