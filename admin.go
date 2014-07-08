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
	"strconv"
	"time"
)

// return JSON with Content type header
func executeJSON(w http.ResponseWriter, statusCode int, data map[string]interface{}) {
	jsonData, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}

// handleAdminPage is REST handler of AdminPage struct.
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
		executeJSON(w, 200, map[string]interface{}{"model_name": modelVar, "items": items})
	case "POST":
		setNewEntity(w, r, modelVar)
		executeJSON(w, 201, map[string]interface{}{"model_name": modelName, "message": "created"})
	}
}

// handleAdminPage is REST handler of AdminPage struct.
func handleModelKeyName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	keyName := vars["keyName"]
	modelVar := vars["modelVar"]
	modelName := modelNames[modelVar]
	switch r.Method {
	case "GET":
		c := appengine.NewContext(r)
		modelStruct := models[modelVar]
		key := datastore.NewKey(c, modelName, keyName, 0, nil)
		err := datastore.Get(c, key, modelStruct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		executeJSON(w, 200, map[string]interface{}{"model_name": modelVar, "item": modelStruct})
	case "PUT":
		updateEntity(w, r, modelVar)
		executeJSON(w, 200, map[string]interface{}{"model_name": modelVar, "message": "updated"})
	case "DELETE":
		c := appengine.NewContext(r)
		key := datastore.NewKey(c, modelName, keyName, 0, nil)
		err := datastore.Delete(c, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		executeJSON(w, 200, map[string]interface{}{"model_name": modelVar, "message": "deleted"})
	}
}

// adminIndex renders index page for admin
func adminIndex(w http.ResponseWriter, r *http.Request) error {
	data := map[string]interface{}{
		"title":       "admin top",
		"description": "this is a starter app for GAE/Go",
		"body":        "admin page",
	}
	return executeTemplate(w, "adminIndex", 200, data)
}

// adminIndex renders index page for admin
func adminForm(w http.ResponseWriter, r *http.Request) error {
	data := map[string]interface{}{
		"title":       "admin form",
		"description": "this is a starter app for GAE/Go",
		"body":        "admin form",
	}
	return executeTemplate(w, "form", 200, data)
}

// adminModels returns list of models
func adminModels(w http.ResponseWriter, r *http.Request) {
	var itemList []map[string]string
	for k := range models {
		itemList = append(itemList, map[string]string{"name": modelNames[k], "description": modelDescriptions[k], "id": k})
	}
	executeJSON(w, 200, map[string]interface{}{"models": itemList})
}

// adminMetaData returns Fields data of a struct
func modelMetaData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	modelVar := vars["modelVar"]
	var model = models[modelVar]
	modelName := modelNames[modelVar]
	schema := map[string]interface{}{"type": "object", "title": modelName}
	var properties map[string]interface{}
	properties = make(map[string]interface{})
	s := reflect.ValueOf(model).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		if typeOfT.Field(i).Tag.Get("verbose_name") != "-" {
			jsonSchemaType := jsonSchemaTypes[typeOfT.Field(i).Tag.Get("datastore_type")]
			property := map[string]interface{}{
				"type":       jsonSchemaType,
				"title":      typeOfT.Field(i).Tag.Get("verbose_name"),
				"fieldOrder": i,
			}
			switch typeOfT.Field(i).Tag.Get("datastore_type") {
			case "String":
				property["maxLength"] = 500
			case "DateTime":
				property["format"] = "date-time"
			}
			properties[typeOfT.Field(i).Tag.Get("json")] = property
		}
	}
	schema["properties"] = properties
	executeJSON(w, 200, map[string]interface{}{"schema": schema})
}

// setNewEntity put a new entity to datastore which is sent in FormValue
func setNewEntity(w http.ResponseWriter, r *http.Request, modelVar string) {
	c := appengine.NewContext(r)
	modelName := modelNames[modelVar]
	modelStruct := models[modelVar]
	s := reflect.ValueOf(modelStruct).Elem()
	typeOfT := s.Type()
	keyName := r.FormValue("url")
	if keyName == "" {
		keyName = time.Now().Format("20060102150405")
	}
	key := datastore.NewKey(c, modelName, keyName, 0, nil)
	for i := 0; i < s.NumField(); i++ {
		//log.Println(typeOfT.Field(i).Name)
		//log.Println(typeOfT.Field(i).Tag.Get("datastore_type"))
		//log.Println(r.FormValue(typeOfT.Field(i).Tag.Get("json")))
		switch typeOfT.Field(i).Tag.Get("datastore_type") {
		case "Boolean":
			err := reflections.SetField(modelStruct, typeOfT.Field(i).Name, r.FormValue(typeOfT.Field(i).Tag.Get("json")) == "on")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case "Integer":
			tmpStringInt, _ := strconv.Atoi(r.FormValue(typeOfT.Field(i).Tag.Get("json")))
			err := reflections.SetField(modelStruct, typeOfT.Field(i).Name, tmpStringInt)
			if err != nil {
				setDefaultErr := reflections.SetField(modelStruct, typeOfT.Field(i).Name, 0)
				if setDefaultErr != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		default:
			err := reflections.SetField(modelStruct, typeOfT.Field(i).Name, r.FormValue(typeOfT.Field(i).Tag.Get("json")))
			if err != nil {
				setDefaultErr := reflections.SetField(modelStruct, typeOfT.Field(i).Name, defaultValues[typeOfT.Field(i).Tag.Get("datastore_type")])
				if setDefaultErr != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}
	_, err := datastore.Put(c, key, modelStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateEntity(w http.ResponseWriter, r *http.Request, modelVar string) {
	c := appengine.NewContext(r)
	modelName := modelNames[modelVar]
	modelStruct := models[modelVar]
	vars := mux.Vars(r)
	keyName := vars["keyName"]
	s := reflect.ValueOf(modelStruct).Elem()
	typeOfT := s.Type()
	key := datastore.NewKey(c, modelName, keyName, 0, nil)
	getErr := datastore.Get(c, key, modelStruct)
	if getErr != nil {
		http.Error(w, getErr.Error(), http.StatusInternalServerError)
		return
	}
	for i := 0; i < s.NumField(); i++ {
		//log.Println(typeOfT.Field(i).Name)
		//log.Println(typeOfT.Field(i).Tag.Get("datastore_type"))
		//log.Println(r.FormValue(typeOfT.Field(i).Tag.Get("json")))
		switch typeOfT.Field(i).Tag.Get("datastore_type") {
		case "Boolean":
			err := reflections.SetField(modelStruct, typeOfT.Field(i).Name, r.FormValue(typeOfT.Field(i).Tag.Get("json")) == "on")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case "Integer":
			tmpStringInt, _ := strconv.Atoi(r.FormValue(typeOfT.Field(i).Tag.Get("json")))
			err := reflections.SetField(modelStruct, typeOfT.Field(i).Name, tmpStringInt)
			if err != nil {
				setDefaultErr := reflections.SetField(modelStruct, typeOfT.Field(i).Name, 0)
				if setDefaultErr != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		default:
			err := reflections.SetField(modelStruct, typeOfT.Field(i).Name, r.FormValue(typeOfT.Field(i).Tag.Get("json")))
			if err != nil {
				if typeOfT.Field(i).Tag.Get("json") == "update" {
					setDefaultErr := reflections.SetField(modelStruct, typeOfT.Field(i).Name, defaultValues[typeOfT.Field(i).Tag.Get("datastore_type")])
					if setDefaultErr != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				} else if typeOfT.Field(i).Tag.Get("json") == "created" {
					log.Println("skipped")
					continue
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
	}
	_, putErr := datastore.Put(c, key, modelStruct)
	if putErr != nil {
		http.Error(w, putErr.Error(), http.StatusInternalServerError)
		return
	}
}
