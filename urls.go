package main

import (
	"net/http"

	"github.com/zenazn/goji"
)

func init() {
	http.Handle("/", goji.DefaultMux)
	goji.Get("/article/:key_name", articlePage)
	goji.Get("/:key_name", generalPage)
	goji.Get("/", index)
}
