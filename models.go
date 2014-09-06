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
	"String":     "",
	"Text":       "",
	"Integer":    0,
	"DateTime":   time.Now(),
	"StringList": [...]string{""},
}

//Map for Models which can be used in restful API
var models = map[string]interface{}{"adminpage": &AdminPage{}, "article": &Article{}}
var searchModels = map[string]interface{}{"adminpage": &AdminPageSearch{}, "article": &ArticleSearch{}}
var modelNames = map[string]string{"adminpage": "AdminPage", "article": "Article"}
var modelDescriptions = map[string]string{
	"adminpage": "model for storing general information",
	"article":   "model for storing article",
}

// AdminPage stores content for general pages.
// URL is reserved for Key Name.
// do not use "draft" for property name, it also reserved
type AdminPage struct {
	URL         string    `datastore:"url" json:"url" datastore_type:"String" verbose_name:"URL=Key Name"`
	DisplayPage bool      `datastore:"displaypage" json:"displaypage" datastore_type:"Boolean" verbose_name:"Display this page"`
	Title       string    `datastore:"title,required" json:"title" datastore_type:"String" verbose_name:"Title"`
	PageOrder   int       `datastore:"pageorder" json:"pageorder" datastore_type:"Integer" verbose_name:"-"`
	Content     string    `datastore:"content,noindex" json:"content" datastore_type:"Text" verbose_name:"Content"`
	Images      string    `datastore:"images,noindex" json:"images" datastore_type:"Text" verbose_name:"-"`
	ExternalURL string    `datastore:"externalurl" json:"externalurl" datastore_type:"String" verbose_name:"Link to ..."`
	Archive     bool      `datastore:"archive" json:"archive" datastore_type:"Boolean" verbose_name:"-"`
	Revision    int       `datastore:"revision" json:"revision" datastore_type:"Integer" verbose_name:"-"`
	Update      time.Time `datastore:"update" json:"update" datastore_type:"DateTime" verbose_name:"-"`
	Created     time.Time `datastore:"created" json:"created" datastore_type:"DateTime" verbose_name:"-"`
}

// AdminPageList is the array type of AdminPage
type AdminPageList []AdminPage

// AdminPageSearch is the type for indexing the content of AdminPage
type AdminPageSearch struct {
	URL         string    `datastore:"url" json:"url" datastore_type:"String" verbose_name:"URL=Key Name"`
	DisplayPage string    `datastore:"displaypage" json:"displaypage" datastore_type:"Boolean" verbose_name:"Display this page"`
	Title       string    `datastore:"title,required" json:"title" datastore_type:"String" verbose_name:"Title"`
	PageOrder   float64   `datastore:"pageorder" json:"pageorder" datastore_type:"Integer" verbose_name:"-"`
	Content     string    `datastore:"content,noindex" json:"content" datastore_type:"Text" verbose_name:"Content"`
	Images      string    `datastore:"images,noindex" json:"images" datastore_type:"Text" verbose_name:"-"`
	ExternalURL string    `datastore:"externalurl" json:"externalurl" datastore_type:"String" verbose_name:"Link to ..."`
	Update      time.Time `datastore:"update" json:"update" datastore_type:"DateTime" verbose_name:"-"`
	Created     time.Time `datastore:"created" json:"created" datastore_type:"DateTime" verbose_name:"-"`
}

// Article stores daily update contents
// URL is reserved for Key Name.
// do not use "draft" for property name, it also reserved
type Article struct {
	URL         string    `datastore:"url" json:"url" datastore_type:"String" verbose_name:"URL=Key Name"`
	DisplayPage bool      `datastore:"displaypage" json:"displaypage" datastore_type:"Boolean" verbose_name:"Display this page"`
	Title       string    `datastore:"title,required" json:"title" datastore_type:"String" verbose_name:"Title"`
	DisplayTime time.Time `datastore:"displaytime" json:"displaytime" datastore_type:"DateTime" verbose_name:"Display Time(UTC)"`
	PageOrder   int       `datastore:"pageorder" json:"pageorder" datastore_type:"Integer" verbose_name:"-"`
	Content     string    `datastore:"content,noindex" json:"content" datastore_type:"Text" verbose_name:"Body Content"`
	Images      string    `datastore:"images,noindex" json:"images" datastore_type:"Text" verbose_name:"-"`
	TagString   string    `datastore:"tagstring,noindex" json:"tagstring" datastore_type:"String" verbose_name:"TagString(comma separated)"`
	Tags        []string  `datastore:"tags" json:"tags" datastore_type:"StringList" verbose_name:"-"`
	ExternalURL string    `datastore:"externalurl" json:"externalurl" datastore_type:"String" verbose_name:"Link to ..."`
	Archive     bool      `datastore:"archive" json:"archive" datastore_type:"Boolean" verbose_name:"-"`
	Revision    int       `datastore:"revision" json:"revision" datastore_type:"Integer" verbose_name:"-"`
	Update      time.Time `datastore:"update" json:"update" datastore_type:"DateTime" verbose_name:"-"`
	Created     time.Time `datastore:"created" json:"created" datastore_type:"DateTime" verbose_name:"-"`
}

// ArticleList is the array type of Article
type ArticleList []Article

// ArticleSearch is the type for indexing the content of Article
type ArticleSearch struct {
	URL         string    `datastore:"url" json:"url" datastore_type:"String" verbose_name:"URL=Key Name"`
	DisplayPage string    `datastore:"displaypage" json:"displaypage" datastore_type:"Boolean" verbose_name:"Display this page"`
	Title       string    `datastore:"title,required" json:"title" datastore_type:"String" verbose_name:"Title"`
	DisplayTime time.Time `datastore:"displaytime" json:"displaytime" datastore_type:"DateTime" verbose_name:"Display Time"`
	PageOrder   float64   `datastore:"pageorder" json:"pageorder" datastore_type:"Integer" verbose_name:"-"`
	Content     string    `datastore:"content,noindex" json:"content" datastore_type:"Text" verbose_name:"Body Content"`
	Images      string    `datastore:"images,noindex" json:"images" datastore_type:"Text" verbose_name:"-"`
	TagString   string    `datastore:"tagstring,noindex" json:"tagstring" datastore_type:"String" verbose_name:"TagString"`
	ExternalURL string    `datastore:"externalurl" json:"externalurl" datastore_type:"String" verbose_name:"Link to ..."`
	Update      time.Time `datastore:"update" json:"update" datastore_type:"DateTime" verbose_name:"-"`
	Created     time.Time `datastore:"created" json:"created" datastore_type:"DateTime" verbose_name:"-"`
}

// BlobStoreImage is the type for storing meta data of blobstore contents
type BlobStoreImage struct {
	Title    string    `datastore:"title" json:"title" datastore_type:"String" verbose_name:"Title"`
	Note     string    `datastore:"note,noindex" json:"content" datastore_type:"Text" verbose_name:"Description"`
	BlobKey  string    `datastore:"blob_key" json:"tagstring" datastore_type:"String" verbose_name:"Blobkey"`
	ImageURL string    `datastore:"image_url" json:"image_url" datastore_type:"String" verbose_name:"Image Url"`
	Update   time.Time `datastore:"update" json:"update" datastore_type:"DateTime" verbose_name:"-"`
	Created  time.Time `datastore:"created" json:"created" datastore_type:"DateTime" verbose_name:"-"`
}

// MaxPageOrder is the type for storing max number of page order of a struct
type MaxPageOrder struct {
	ModelName string `datastore:"modelname" json:"modelname" datastore_type:"String" verbose_name:"Model Name"`
	MaxOrder  int    `datastore:"maxorder" json:"maxorder" datastore_type:"Integer" verbose_name:"Max Order"`
}
