package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	. "github.com/hamza12700/cryptos/endpoints"
)

type route struct {
	Name    string
	handler http.HandlerFunc
}

type tmplData struct {
	RouteLen int
	Routes   *[]route
}

func main() {

	apiRoutes := []route{
		{Name: "/md5", handler: Md5Sum},
		{Name: "/sha1", handler: Sha1Sum},
		{Name: "/sha256", handler: Sha256Sum},
		{Name: "/sha224", handler: Sha224Sum},
		{Name: "/uri-encode", handler: UriEncodor},
		{Name: "/uri-decode", handler: UriDecodor},
		{Name: "/uuid", handler: GenerateUUID},
		{Name: "/random-uuid", handler: RandomUUID},
		{Name: "/base64", handler: EncodeToBase64},
		{Name: "/decode-base64", handler: DecodeToBase64},
		{Name: "/text-to-binary", handler: TextToBinary},
		{Name: "/html-entities-escape", handler: EscapeHtml},
		{Name: "/unescape-html-entities", handler: UnescapeHtml},
	}
	for _, v := range apiRoutes {
		http.HandleFunc(v.Name, v.handler)
	}

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
			Routes:   &apiRoutes,
		}

		err = tmpl.Execute(w, &data)
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
