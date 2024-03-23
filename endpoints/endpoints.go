package endpoints

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

func middleware(w http.ResponseWriter, sendResquest string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.WriteHeader(http.StatusCreated)
	jsonResq, err := json.Marshal(sendResquest)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(jsonResq)
}

func hexToString(hashedText []byte) string {
	hash := hex.EncodeToString(hashedText[:])
	return hash
}

func Sha1Sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.PathValue("text")
	hashedText := sha1.Sum([]byte(textToHash))
	hash := hexToString(hashedText[:])
	middleware(w, hash)
}

func Sha256Sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.PathValue("text")
	hashedText := sha256.Sum256([]byte(textToHash))
	hash := hexToString(hashedText[:])
	middleware(w, hash)
}

func Sha224Sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.PathValue("text")
	hashedText := sha256.Sum224([]byte(textToHash))
	hash := hexToString(hashedText[:])
	middleware(w, hash)
}

func Md5Sum(w http.ResponseWriter, r *http.Request) {
	textToHash := r.PathValue("text")
	hashedText := md5.Sum([]byte(textToHash))
	hash := hexToString(hashedText[:])
	middleware(w, hash)
}

func UriEncodor(w http.ResponseWriter, r *http.Request) {
	textToEncode := r.PathValue("text")
	encodedURI := url.PathEscape(textToEncode)
	middleware(w, encodedURI)
}

func UriDecodor(w http.ResponseWriter, r *http.Request) {
	textToDecode := r.PathValue("text")
	decodedURI, err := url.PathUnescape(textToDecode)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	middleware(w, decodedURI)
}

func GenerateUUID(w http.ResponseWriter, r *http.Request) {
	textToConvert := r.PathValue("text")
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
	uuidStr := uuidBytes.String()

	middleware(w, uuidStr)
}

func RandomUUID(w http.ResponseWriter, r *http.Request) {
	randomlyGenUUID, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	randomUuidStr := randomlyGenUUID.String()
	middleware(w, randomUuidStr)
}

func EncodeToBase64(w http.ResponseWriter, r *http.Request) {
	textToConvert := r.PathValue("text")
	encodedText := base64.StdEncoding.EncodeToString([]byte(textToConvert))
	middleware(w, encodedText)
}

func DecodeToBase64(w http.ResponseWriter, r *http.Request) {
	textToDecode := r.URL.Query().Get("text")
	decodedText, err := base64.StdEncoding.DecodeString(textToDecode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	decodedTextStr := string(decodedText)
	middleware(w, decodedTextStr)
}

func TextToBinary(w http.ResponseWriter, r *http.Request) {
	textToConvert := r.URL.Query().Get("text")
	var binString string
	for _, i := range textToConvert {
		binString = fmt.Sprintf("%s%08b", binString, i)
	}
	middleware(w, binString)
}

func EscapeHtml(w http.ResponseWriter, r *http.Request) {
	htmlToEscape := r.PathValue("text")
	escapedHtmlEntities := html.EscapeString(htmlToEscape)
	middleware(w, escapedHtmlEntities)
}

func UnescapeHtml(w http.ResponseWriter, r *http.Request) {
	escapeHtml := r.PathValue("text")
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
