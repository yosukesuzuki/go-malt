package main

import (
	//"appengine"
	//"doc"
	"encoding/json"
	"net/http"
)

func executeJSON(w http.ResponseWriter, data map[string]interface{}) {
	jsonData, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(jsonData)
}

func adminPageEntityList(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title":       "admin top",
		"description": "this is a starter app for GAE/Go",
		"body":        "admin page",
	}
	executeJSON(w, data)
}

// adminIndex renders index page for admin
func adminIndex(w http.ResponseWriter, r *http.Request) error {
	data := map[string]interface{}{
		"title":       "admin top",
		"description": "this is a starter app for GAE/Go",
		"body":        "admin page",
	}
	return executeTemplate(w, "index", 200, data)
}
