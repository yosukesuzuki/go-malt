package main

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/oleiade/reflections"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
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
		listDataSet := getAdminPageList(w, r)
		executeJSON(w, 200, listDataSet)
	case "POST":
		setNewEntity(w, r, modelVar)
		executeJSON(w, 201, map[string]interface{}{"model_name": modelName, "message": "created"})
	}
}

func getAdminPageList(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	modelVar := "adminpage"
	modelName := modelNames[modelVar]
	c := appengine.NewContext(r)
	perPage := 10
	cursorKey := "adminpage_cursor"
	tmpOffsetValue := r.FormValue("offset")
	var offsetParam int
	if tmpOffsetValue != "" {
		offsetParam, _ = strconv.Atoi(tmpOffsetValue)
	} else {
		offsetParam = 0
	}
	cursorKeyCurrent := cursorKey + strconv.Itoa(offsetParam)
	q := datastore.NewQuery(modelName).Order("-update").Limit(perPage)
	item, err := memcache.Get(c, cursorKeyCurrent)
	if err == nil {
		cursor, err := datastore.DecodeCursor(string(item.Value))
		if err == nil {
			q = q.Start(cursor)
		}
	}
	var items []AdminPage
	var hasNext bool
	hasNext = true
	// Iterate over the results.
	t := q.Run(c)
	for {
		var ap AdminPage
		_, err := t.Next(&ap)
		if err == datastore.Done {
			break
		}
		if err != nil {
			c.Errorf("fetching next AdminPage: %v", err)
			hasNext = false
			break
		}
		items = append(items, ap)
	}
	nextOffset := offsetParam + perPage
	cursorKeyNext := cursorKey + strconv.Itoa(nextOffset)
	// Get updated cursor and store it for next time.
	if cursor, err := t.Cursor(); err == nil {
		memcache.Set(c, &memcache.Item{
			Key:   cursorKeyNext,
			Value: []byte(cursor.String()),
		})
	}
	listDataSet := map[string]interface{}{"items": items, "has_next": hasNext, "next_offset": nextOffset, "model_name": modelVar}
	return listDataSet
}

// handleArticlePage is REST handler of AdminPage struct.
func handleArticlePage(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//modelName := vars["modelName"]
	//var model = models[modelName]
	modelVar := "article"
	modelName := modelNames[modelVar]
	switch r.Method {
	case "GET":
		c := appengine.NewContext(r)
		q := datastore.NewQuery(modelName).Order("-update").Limit(20)
		var items []Article
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

// handleModelKeyName is REST handler of Model struct.
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
	var requiredProperties []string
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
			if strings.Contains(typeOfT.Field(i).Tag.Get("datastore"), "required") {
				requiredProperties = append(requiredProperties, typeOfT.Field(i).Tag.Get("json"))
			}
		}
	}
	schema["properties"] = properties
	schema["required"] = requiredProperties
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
	urlValue, err := reflections.GetField(modelStruct, "URL")
	if urlValue == "" {
		setUrlErr := reflections.SetField(modelStruct, "URL", keyName)
		if setUrlErr != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	_, putErr := datastore.Put(c, key, modelStruct)
	if putErr != nil {
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
