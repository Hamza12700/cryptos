package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"os"

	"github.com/google/uuid"
)

type route struct {
	name    string
	handler http.HandlerFunc
}

func router(routes []route) {
	for i := 0; i < len(routes); i++ {
		http.HandleFunc(routes[i].name, routes[i].handler)
	}
}

func main() {

	apiRoutes := []route{
		{name: "/md5", handler: md5Sum},
		{name: "/sha1", handler: sha1Sum},
		{name: "/sha256", handler: sha256Sum},
		{name: "/sha224", handler: sha224Sum},
		{name: "/uri-encode", handler: uriEncodor},
		{name: "/uri-decode", handler: uriDecodor},
		{name: "/uuid", handler: generateUUID},
		{name: "/random-uuid", handler: randomUUID},
		{name: "/base64", handler: encodeToBase64},
		{name: "/decode-base64", handler: decodeToBase64},
		{name: "/text-to-binary", handler: textToBinary},
		{name: "/html-entities-escape", handler: escapeHtml},
		{name: "/unescape-html-entities", handler: unescapeHtml},
	}
	router(apiRoutes)

	rootDir := http.FileServer(http.Dir("./html-templates"))
	http.Handle("/", rootDir)

	port := os.Getenv("PORT")
	if port == "" {
		port = "2323"
	}

	fmt.Println("Listening on", port)
	http.ListenAndServe(":"+port, nil)
}

func middleware(w http.ResponseWriter, sendResquest string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
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

func md5Sum(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	middleware(w, randomlyGenUUID.String())
}

func encodeToBase64(w http.ResponseWriter, r *http.Request) {
	textToConvert := r.URL.Query().Get("text")
	encodedText := base64.StdEncoding.EncodeToString([]byte(textToConvert))
	middleware(w, encodedText)
}

func decodeToBase64(w http.ResponseWriter, r *http.Request) {
	textToDecode := r.URL.Query().Get("text")
	decodedText, err := base64.StdEncoding.DecodeString(textToDecode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	middleware(w, string(decodedText))
}

func textToBinary(w http.ResponseWriter, r *http.Request) {
	textToConvert := r.URL.Query().Get("text")
	var binString string
	for _, i := range textToConvert {
		binString = fmt.Sprintf("%s%08b", binString, i)
	}
	middleware(w, binString)
}

func escapeHtml(w http.ResponseWriter, r *http.Request) {
	htmlToEscape := r.URL.Query().Get("text")
	escapedHtmlEntities := html.EscapeString(htmlToEscape)
	middleware(w, escapedHtmlEntities)
}

func unescapeHtml(w http.ResponseWriter, r *http.Request) {
	escapeHtml := r.URL.Query().Get("text")
	isBase64 := r.URL.Query().Get("base64")
	convertToBase64, err := base64.StdEncoding.DecodeString(escapeHtml)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	unscapeHtml := html.UnescapeString(string(convertToBase64))
	if isBase64 == "true" {
		unscapeHtmlBase64 := base64.StdEncoding.EncodeToString([]byte(unscapeHtml))
		middleware(w, unscapeHtmlBase64)
		return
	}
	middleware(w, unscapeHtml)
}
