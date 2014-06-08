package main

import (
	"appengine"
	"doc"
	"net/http"

	"github.com/gorilla/mux"
)

type handlerFunc func(http.ResponseWriter, *http.Request) error

func (f handlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := f(w, r)
	if err != nil {
		appengine.NewContext(r).Errorf("Error %s", err.Error())
		if e, ok := err.(doc.GetError); ok {
			http.Error(w, "Error getting files from "+e.Host+".", http.StatusInternalServerError)
		} else if appengine.IsCapabilityDisabled(err) || appengine.IsOverQuota(err) {
			http.Error(w, "Internal error: "+err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
	}
}

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
