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
	handler func(http.ResponseWriter, *http.Request)
}

type tmplData struct {
	RouteLen int
	Routes   *[]route
}

func main() {

	router := http.NewServeMux()

	apiRoutes := []route{
		{Name: "/md5/{text}", handler: Md5Sum},
		{Name: "/sha1/{text}", handler: Sha1Sum},
		{Name: "/sha256/{text}", handler: Sha256Sum},
		{Name: "/sha224/{text}", handler: Sha224Sum},
		{Name: "/uri/encode/{text}", handler: UriEncodor},
		{Name: "/uri/decode/{text}", handler: UriDecodor},
		{Name: "/uuid/{text}", handler: GenerateUUID},
		{Name: "/uuid/random", handler: RandomUUID},
		{Name: "/base64/encode/{text}", handler: EncodeToBase64},
		{Name: "/base64/decode/{text}", handler: DecodeToBase64},
		{Name: "/text-to-binary/{text}", handler: TextToBinary},
		{Name: "/html-entities/escape/{text}", handler: EscapeHtml},
		{Name: "/html-entities/unescape/{text}", handler: UnescapeHtml},
	}
	for _, v := range apiRoutes {
		router.HandleFunc(v.Name, v.handler)
	}

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./html-templates/static"))))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	server.ListenAndServe()
}
