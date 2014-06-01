package main

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "this is index")
}

func generalPage(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "general page, %s!", c.URLParams["key_name"])
}

func articlePage(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "article, %s!", c.URLParams["key_name"])
}
