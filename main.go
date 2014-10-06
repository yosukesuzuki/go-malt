package main

import (
	"net/http"

	"appengine"
	"appengine/datastore"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
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
	c := appengine.NewContext(r)
	adminpage := &AdminPage{}
	key := datastore.NewKey(c, "AdminPage", keyName, 0, nil)
	err := datastore.Get(c, key, adminpage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	bodyHTML := string(blackfriday.MarkdownBasic([]byte(adminpage.Content)))
	data := map[string]interface{}{
		"title":       adminpage.Title,
		"description": "this is a starter app for GAE/Go",
		"update":      adminpage.Update,
		"body":        bodyHTML,
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
