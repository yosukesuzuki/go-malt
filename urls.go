package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	s := r.PathPrefix("/admin").Subrouter()
	s.Handle("/", handlerFunc(adminIndex))
	s.HandleFunc("/rest/adminpage", adminPageEntityList).Methods("GET")
	//s.Handle("/rest/adminpage/{keyName}", handlerFunc(adminPageGetEntity)).Methods("GET")
	//s.Handle("/rest/adminpage/{keyName}", handlerFunc(adminPageUpdateEntity)).Methods("POST")
	//s.Handle("/rest/adminpage", handlerFunc(adminPageNewEntity)).Methods("PUT")
	//s.Handle("/rest/adminpage/{keyName}", handlerFunc(adminPageDeleteEntity)).Methods("DELETE")
	r.Handle("/article/{keyName}", handlerFunc(articlePage))
	r.Handle("/{keyName}", handlerFunc(generalPage))
	r.Handle("/", handlerFunc(index))
	http.Handle("/", r)
}
