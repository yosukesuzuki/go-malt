package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	s := r.PathPrefix("/admin").Subrouter()
	s.Handle("/", handlerFunc(adminIndex))
	s.Handle("/form/", handlerFunc(adminForm))
	s.HandleFunc("/rest/models", adminModels)
	s.HandleFunc("/rest/schema/{modelVar}", modelMetaData)
	s.HandleFunc("/rest/adminpage", handleAdminPage)
	s.HandleFunc("/rest/article", handleArticle)
	s.HandleFunc("/rest/{modelVar}/{keyName}", handleModelKeyName)
	s.HandleFunc("/image/upload/url", imageUploadURL)
	s.HandleFunc("/image/upload/handler", handleImageUpload)
	rest := r.PathPrefix("/rest").Subrouter()
	rest.HandleFunc("/search", handleSearch)
	r.Handle("/article/{keyName}", handlerFunc(articlePage))
	r.Handle("/search", handlerFunc(searchPage))
	r.Handle("/{keyName}", handlerFunc(generalPage))
	r.Handle("/", handlerFunc(index))
	http.Handle("/", r)
}
