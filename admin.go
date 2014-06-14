package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

func executeJSON(w http.ResponseWriter, data map[string]interface{}) {
	jsonData, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(jsonData)
}

func handleAdminPage(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	switch r.Method {
	case "GET":
		data := map[string]interface{}{
			"title":       "admin top",
			"description": "this is a starter app for GAE/Go",
			"body":        "admin page",
		}
		executeJSON(w, data)
	}
}

// adminIndex renders index page for admin
func adminIndex(w http.ResponseWriter, r *http.Request) error {
	data := map[string]interface{}{
		"title":       "admin top",
		"description": "this is a starter app for GAE/Go",
		"body":        "admin page",
	}
	return executeTemplate(w, "index", 200, data)
}

// ModelField is struct to return metadata of a Model
type ModelField struct {
	FieldName   string `json:"field_name"`
	FieldType   string `json:"field_type"`
	VerboseName string `json:"verbose_name"`
}

//Map for Models which can be used in restful API
var models = map[string]interface{}{"adminpage": &AdminPage{}, "article": &Article{}}

func modelMetaData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	modelName := vars["modelName"]
	var model = models[modelName]
	var itemList []ModelField
	s := reflect.ValueOf(model).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		modelField := ModelField{typeOfT.Field(i).Tag.Get("json"), typeOfT.Field(i).Tag.Get("datastore_type"), typeOfT.Field(i).Tag.Get("verbose_name")}
		itemList = append(itemList, modelField)
	}
	executeJSON(w, map[string]interface{}{"model_name": modelName, "fields": itemList})
}
