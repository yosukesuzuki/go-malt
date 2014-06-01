package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "this is index")
}

func generalPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "this is ", r.URL.Path[6:])
}

func articlePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "this is article ", r.URL.Path[9:])
}
