// Copyright 2011 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// +build appengine

package main

import (
	"appengine"
	//"doc"
	"errors"
	"fmt"
	"github.com/gorilla/site/doc"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"time"
)

type handlerFunc func(http.ResponseWriter, *http.Request) error

func (f handlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := f(w, r)
	if err != nil {
		appengine.NewContext(r).Errorf("Error %s", err.Error())
		if e, ok := err.(doc.GetError); ok {
			http.Error(w, "Error getting files from "+e.Host+".", http.StatusInternalServerError)
		} else if appengine.IsCapabilityDisabled(err) || appengine.IsOverQuota(err) {
			http.Error(w, "Internal error: "+err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
	}
}

func mapFmt(kvs ...interface{}) (map[string]interface{}, error) {
	if len(kvs)%2 != 0 {
		return nil, errors.New("map requires even number of arguments")
	}
	m := make(map[string]interface{})
	for i := 0; i < len(kvs); i += 2 {
		s, ok := kvs[i].(string)
		if !ok {
			return nil, errors.New("even args to map must be strings")
		}
		m[s] = kvs[i+1]
	}
	return m, nil
}

// relativePathFmt formats an import path as HTML.
func relativePathFmt(importPath string, parentPath interface{}) string {
	if p, ok := parentPath.(string); ok && p != "" && strings.HasPrefix(importPath, p) {
		importPath = importPath[len(p)+1:]
	}
	return urlFmt(importPath)
}

// importPathFmt formats an import with zero width space characters to allow for breeaks.
func importPathFmt(importPath string) string {
	importPath = urlFmt(importPath)
	if len(importPath) > 45 {
		// Allow long import paths to break following "/"
		importPath = strings.Replace(importPath, "/", "/&#8203;", -1)
	}
	return importPath
}

// relativeTime formats the time t in nanoseconds as a human readable relative
// time.
func relativeTime(t time.Time) string {
	const day = 24 * time.Hour
	d := time.Now().Sub(t)
	switch {
	case d < time.Second:
		return "just now"
	case d < 2*time.Second:
		return "one second ago"
	case d < time.Minute:
		return fmt.Sprintf("%d seconds ago", d/time.Second)
	case d < 2*time.Minute:
		return "one minute ago"
	case d < time.Hour:
		return fmt.Sprintf("%d minutes ago", d/time.Minute)
	case d < 2*time.Hour:
		return "one hour ago"
	case d < day:
		return fmt.Sprintf("%d hours ago", d/time.Hour)
	case d < 2*day:
		return "one day ago"
	}
	return fmt.Sprintf("%d days ago", d/day)
}

var (
	h3Open     = []byte("<h3 ")
	h4Open     = []byte("<h2 ")
	h3Close    = []byte("</h3>")
	h4Close    = []byte("</h2>")
	rfcRE      = regexp.MustCompile(`RFC\s+(\d{3,4})`)
	rfcReplace = []byte(`<a href="http://tools.ietf.org/html/rfc$1">$0</a>`)
)

func urlFmt(path string) string {
	u := url.URL{Path: path}
	return u.String()
}

func executeTemplate(w http.ResponseWriter, name string, status int, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	return tpls[name].ExecuteTemplate(w, "base", data)
}

// eq reports whether the first argument is equal to
// any of the remaining arguments.
// https://groups.google.com/group/golang-nuts/msg/a720bf35f454288b
func eq(args ...interface{}) bool {
	if len(args) == 0 {
		return false
	}
	x := args[0]
	switch x := x.(type) {
	case string, int, int64, byte, float32, float64:
		for _, y := range args[1:] {
			if x == y {
				return true
			}
		}
		return false
	}
	for _, y := range args[1:] {
		if reflect.DeepEqual(x, y) {
			return true
		}
	}
	return false
}

// GetBaseURL return base url for google analytics
func GetBaseURL() string {
	//Change This to Your Google Analytics ID
	const url = "goappstarter.appspot.com"
	return url
}

// GetGaID return id for google analytics
func GetGaID() string {
	//Change This to Your Google Analytics ID
	const gaID = "UA-51746203-1"
	return gaID
}

// GetFbAppID returns id for Facebook
func GetFbAppID() string {
	//Change This to Your FB App ID
	const fbAppID = "1491547807729755"
	return fbAppID
}

var tpls = map[string]*template.Template{
	"404":          newTemplate("templates/base.html", "templates/404.html"),
	"index":        newTemplate("templates/base.html", "templates/index.html"),
	"page":         newTemplate("templates/base.html", "templates/page.html"),
	"articleIndex": newTemplate("templates/base.html", "templates/article_index.html"),
}

var funcs = template.FuncMap{
	"eq":           eq,
	"map":          mapFmt,
	"relativePath": relativePathFmt,
	"relativeTime": relativeTime,
	"importPath":   importPathFmt,
	"url":          urlFmt,
	"baseUrl":      GetBaseURL,
	"gaID":         GetGaID,
	"fbAppID":      GetFbAppID,
}

func newTemplate(files ...string) *template.Template {
	return template.Must(template.New("*").Funcs(funcs).ParseFiles(files...))
}
