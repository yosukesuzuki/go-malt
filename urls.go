package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	s := r.PathPrefix("/admin").Subrouter()
	s.Handle("/", handlerFunc(adminIndex))
	s.HandleFunc("/rest/adminpage", handleAdminPage)
	r.Handle("/article/{keyName}", handlerFunc(articlePage))
	r.Handle("/{keyName}", handlerFunc(generalPage))
	r.Handle("/", handlerFunc(index))
	http.Handle("/", r)
}
