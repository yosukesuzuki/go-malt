package main

import (
	"net/http"
)

func init() {
	http.HandleFunc("/article/", articlePage)
	http.HandleFunc("/page/", generalPage)
	http.HandleFunc("/", index)
}
