package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title":       "top",
		"description": "this is a starter app for GAE/Go",
		"body":        "hello world",
	}
	t, _ := template.ParseFiles("templates/base.html", "templates/index.html")
	t.ExecuteTemplate(w, "base", data)
}

func generalPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key_name := vars["key_name"]
	fmt.Fprintf(w, "general page, %s!", key_name)
}

func articlePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key_name := vars["key_name"]
	fmt.Fprintf(w, "article, %s!", key_name)
}
