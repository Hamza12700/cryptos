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
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/google/uuid"
)

type route struct {
	Name    string
	handler http.HandlerFunc
}

type tmplData struct {
	RouteLen int
	Routes   []route
}

func router(routes []route) {
	for _, v := range routes {
		http.HandleFunc(v.Name, v.handler)
	}
}

func main() {

	apiRoutes := []route{
		{Name: "/md5", handler: md5Sum},
		{Name: "/sha1", handler: sha1Sum},
		{Name: "/sha256", handler: sha256Sum},
		{Name: "/sha224", handler: sha224Sum},
		{Name: "/uri-encode", handler: uriEncodor},
		{Name: "/uri-decode", handler: uriDecodor},
		{Name: "/uuid", handler: generateUUID},
		{Name: "/random-uuid", handler: randomUUID},
		{Name: "/base64", handler: encodeToBase64},
		{Name: "/decode-base64", handler: decodeToBase64},
		{Name: "/text-to-binary", handler: textToBinary},
		{Name: "/html-entities-escape", handler: escapeHtml},
		{Name: "/unescape-html-entities", handler: unescapeHtml},
	}
	router(apiRoutes)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./html-templates/static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./html-templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}

		data := tmplData{
			RouteLen: len(apiRoutes),
			Routes:   apiRoutes,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
	})

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
