package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	s := r.PathPrefix("/admin").Subrouter()
	s.Handle("/", handlerFunc(adminIndex))
	s.Handle("/form", handlerFunc(adminForm))
	s.HandleFunc("/rest/models", adminModels)
	s.HandleFunc("/rest/schema/{modelVar}", modelMetaData)
	s.HandleFunc("/rest/adminpage", handleAdminPage)
	s.HandleFunc("/rest/{modelVar}/{keyName}", handleModelKeyName)
	r.Handle("/article/{keyName}", handlerFunc(articlePage))
	r.Handle("/{keyName}", handlerFunc(generalPage))
	r.Handle("/", handlerFunc(index))
	http.Handle("/", r)
}
