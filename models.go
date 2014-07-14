package main

import (
	//    "log"
	//    "fmt"
	"time"
	//    "strings"
	//    "strconv"
	//    "appengine"
	//    "appengine/urlfetch"
	//    "appengine/memcache"
	//    "html/template"
	//    "net/http"
	//    "github.com/go-martini/martini"
	//    "github.com/martini-contrib/render"
	//    "github.com/PuerkitoBio/goquery"
)

var jsonSchemaTypes = map[string]string{"Boolean": "boolean", "String": "string", "Text": "string", "Integer": "integer", "DateTime": "string"}
var jsonSchemaFields = map[string]string{"Text": "textarea"}
var defaultValues = map[string]interface{}{"Boolean": false,
	"String":   "",
	"Text":     "",
	"Integer":  0,
	"DateTime": time.Now(),
}

//Map for Models which can be used in restful API
var models = map[string]interface{}{"adminpage": &AdminPage{}, "article": &Article{}}
var modelNames = map[string]string{"adminpage": "AdminPage", "article": "Article"}
var modelDescriptions = map[string]string{
	"adminpage": "model for storing general information",
	"article":   "model for storing article",
}

// AdminPage stores content for general pages
type AdminPage struct {
	DisplayPage bool      `datastore:"displaypage" json:"displaypage" datastore_type:"Boolean" verbose_name:"Display this page"`
	Title       string    `datastore:"title,required" json:"title" datastore_type:"String" verbose_name:"Title"`
	URL         string    `datastore:"url" json:"url" datastore_type:"String" verbose_name:"URL"`
	PageOrder   int       `datastore:"pageorder" json:"pageorder" datastore_type:"Integer" verbose_name:"Page Order"`
	Content     string    `datastore:"content,noindex" json:"content" datastore_type:"Text" verbose_name:"Content"`
	Images      string    `datastore:"images,noindex" json:"images" datastore_type:"Text" verbose_name:"-"`
	ExternalURL string    `datastore:"externalurl" json:"externalurl" datastore_type:"String" verbose_name:"Link to ..."`
	Update      time.Time `datastore:"update" json:"update" datastore_type:"DateTime" verbose_name:"-"`
	Create      time.Time `datastore:"created" json:"created" datastore_type:"DateTime" verbose_name:"-"`
}

type AdminPageList struct {
	Items []AdminPage
}

// Article stores daily update contents
type Article struct {
	DisplayPage bool      `datastore:"displaypage" json:"displaypage" datastore_type:"Boolean" verbose_name:"Display this page"`
	Title       string    `datastore:"title,required" json:"title" datastore_type:"String" verbose_name:"Title"`
	URL         string    `datastore:"url" json:"url" datastore_type:"String" verbose_name:"URL"`
	PageOrder   int       `datastore:"pageorder" json:"pageorder" datastore_type:"Integer" verbose_name:"Page Order"`
	Content     string    `datastore:"content,noindex" json:"content" datastore_type:"Text" verbose_name:"Body Content"`
	Images      string    `datastore:"images,noindex" json:"images" datastore_type:"Text" verbose_name:"-"`
	TagString   string    `datastore:"tagstring,noindex" json:"tagstring" datastore_type:"Text" verbose_name:"TagString"`
	Tags        []string  `datastore:"tags" json:"tags" datastore_type:"StringList" verbose_name:"-"`
	ExternalURL string    `datastore:"externalurl" json:"externalurl" datastore_type:"String" verbose_name:"Link to ..."`
	Update      time.Time `datastore:"update" json:"update" datastore_type:"DateTime" verbose_name:"-"`
	Create      time.Time `datastore:"created" json:"created" datastore_type:"DateTime" verbose_name:"-"`
}

type ArticleList struct {
	Items []Article
}
