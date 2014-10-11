package main

import (
	"net/http"

	"appengine"
	"appengine/search"
)

func getSearchResult(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	c := appengine.NewContext(r)
	var items ArticleSearchList
	// Iterate over the results.
	index, openErr := search.Open("global")
	if openErr != nil {
		http.Error(w, openErr.Error(), http.StatusInternalServerError)
	}
	keyword := r.FormValue("keyword")
	for t := index.Search(c, keyword, nil); ; {
		var as ArticleSearch
		_, err := t.Next(&as)
		if err == search.Done {
			break
		}
		if err != nil {
			c.Errorf("fetching next ArticleSearch: %v", err)
			break
		}
		items = append(items, as)
	}
	listDataSet := map[string]interface{}{"items": items}
	return listDataSet
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	listDataSet := getSearchResult(w, r)
	executeJSON(w, 200, listDataSet)
}
