package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) error {
	data := map[string]interface{}{
		"title":       "top",
		"description": "this is a starter app for GAE/Go",
		"body":        "hello world",
	}
	return executeTemplate(w, "index", 200, data)
}

func generalPage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	keyName := vars["keyName"]
	data := map[string]interface{}{
		"title":       keyName,
		"description": "this is a starter app for GAE/Go",
		"body":        keyName,
	}
	return executeTemplate(w, "page", 200, data)
}

func articlePage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	keyName := vars["keyName"]
	data := map[string]interface{}{
		"title":       keyName,
		"description": "this is a starter app for GAE/Go",
		"body":        keyName,
	}
	return executeTemplate(w, "page", 200, data)
}
