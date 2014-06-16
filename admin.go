package main

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/oleiade/reflections"
	"log"
	"net/http"
	"reflect"
	"time"
)

func executeJSON(w http.ResponseWriter, data map[string]interface{}) {
	jsonData, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(jsonData)
}

func handleAdminPage(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//modelName := vars["modelName"]
	//var model = models[modelName]
	modelVar := "adminpage"
	modelName := modelNames[modelVar]
	switch r.Method {
	case "GET":
		c := appengine.NewContext(r)
		q := datastore.NewQuery(modelName).Order("-update").Limit(20)
		var items []AdminPage
		if _, err := q.GetAll(c, &items); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		executeJSON(w, map[string]interface{}{"model_name": modelVar, "items": items})
	case "POST":
		setNewEntity(w, r, modelVar)
		executeJSON(w, map[string]interface{}{"model_name": modelName, "message": "created"})
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

func adminModels(w http.ResponseWriter, r *http.Request) {
	var itemList []string
	for k, _ := range models {
		itemList = append(itemList, k)
	}
	executeJSON(w, map[string]interface{}{"models": itemList})
}

func adminMetaData(w http.ResponseWriter, r *http.Request) {
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

func setNewEntity(w http.ResponseWriter, r *http.Request, modelVar string) {
	c := appengine.NewContext(r)
	modelName := modelNames[modelVar]
	modelStruct := models[modelVar]
	keyName := r.FormValue("url")
	if keyName == "" {
		keyName = time.Now().Format("20060102150405")
	}
	key := datastore.NewKey(c, modelName, keyName, 0, nil)
	s := reflect.ValueOf(modelStruct).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		log.Println(typeOfT.Field(i).Name)
		err := reflections.SetField(modelStruct, typeOfT.Field(i).Name, r.FormValue(typeOfT.Field(i).Tag.Get("json")))
		if err != nil {
			set_default_err := reflections.SetField(modelStruct, typeOfT.Field(i).Name, defaultValues[typeOfT.Field(i).Tag.Get("datastore_type")])
			if set_default_err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	_, err := datastore.Put(c, key, modelStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
