package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/zenazn/goji/web"
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

func generalPage(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "general page, %s!", c.URLParams["key_name"])
}

func articlePage(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "article, %s!", c.URLParams["key_name"])
}
