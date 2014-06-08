package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/article/{key_name}", articlePage)
	r.HandleFunc("/{key_name}", generalPage)
	r.HandleFunc("/", index)
	http.Handle("/", r)
}
