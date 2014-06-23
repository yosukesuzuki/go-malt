package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	s := r.PathPrefix("/admin").Subrouter()
	s.Handle("/", handlerFunc(adminIndex))
	s.HandleFunc("/rest/metadata", adminModels)
	s.HandleFunc("/rest/metadata/{modelName}", adminMetaData)
	s.HandleFunc("/rest/adminpage", handleAdminPage)
	s.HandleFunc("/rest/adminpage/{keyName}", handleAdminPageKeyName)
	r.Handle("/article/{keyName}", handlerFunc(articlePage))
	r.Handle("/{keyName}", handlerFunc(generalPage))
	r.Handle("/", handlerFunc(index))
	http.Handle("/", r)
}
