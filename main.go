package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/google/uuid"
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

func uriEncodor(w http.ResponseWriter, r *http.Request) {
	textToEncode := r.URL.Query().Get("text")
	encodedURI := url.PathEscape(textToEncode)
	middleware(w, encodedURI)
}

func uriDecodor(w http.ResponseWriter, r *http.Request) {
	textToDecode := r.URL.Query().Get("text")
	decodedURI, err := url.PathUnescape(textToDecode)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	middleware(w, decodedURI)
}

func generateUUID(w http.ResponseWriter, r *http.Request) {
	textToConvert := r.URL.Query().Get("text")
	hasher := md5.New()
	_, err := hasher.Write([]byte(textToConvert))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	md5Hash := hex.EncodeToString(hasher.Sum(nil))
	uuidBytes, err := uuid.FromBytes([]byte(md5Hash[:16]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	middleware(w, uuidBytes.String())
}

func randomUUID(w http.ResponseWriter, r *http.Request) {
	randomlyGenUUID, err := uuid.NewRandom()
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
	middleware(w,randomlyGenUUID.String())
}

func main() {
	http.HandleFunc("/md5", md5sum)
	http.HandleFunc("/sha1", sha1Sum)
	http.HandleFunc("/sha256", sha256Sum)
	http.HandleFunc("/sha224", sha224Sum)
	http.HandleFunc("/uri-encode", uriEncodor)
	http.HandleFunc("/uri-decode", uriDecodor)
	http.HandleFunc("/uuid", generateUUID)
	http.HandleFunc("/random-uuid", randomUUID)
	port := os.Getenv("PORT")
	if port == "" {
		port = "2323"
	}
	fmt.Println("Listening on", port)
	http.ListenAndServe(":"+port, nil)
}
