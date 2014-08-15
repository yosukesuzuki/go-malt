package main

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"appengine/image"
	"appengine/memcache"
	"appengine/search"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/oleiade/reflections"
	"log"
	"net/http"
	"reflect"
	"regexp"
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

// getAdminPageList returns entity list of AdminPage model
func getAdminPageList(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	modelVar := "adminpage"
	modelName := modelNames[modelVar]
	c := appengine.NewContext(r)
	perPage := 20
	cursorKey := modelVar + "_cursor"
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
	var items AdminPageList
	var hasNext bool
	hasNext = false
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
	if len(items) == perPage {
		hasNext = true
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
	listDataSet := map[string]interface{}{"items": items,
		"has_next": hasNext, "next_offset": nextOffset, "per_page": perPage, "model_name": modelVar}
	return listDataSet
}

// handleAdminPage is REST handler of AdminPage struct.
func handleArticle(w http.ResponseWriter, r *http.Request) {
	modelVar := "article"
	modelName := modelNames[modelVar]
	switch r.Method {
	case "GET":
		listDataSet := getArticleList(w, r)
		executeJSON(w, 200, listDataSet)
	case "POST":
		setNewEntity(w, r, modelVar)
		executeJSON(w, 201, map[string]interface{}{"model_name": modelName, "message": "created"})
	}
}

// getAdminPageList returns entity list of AdminPage model
func getArticleList(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	modelVar := "article"
	modelName := modelNames[modelVar]
	c := appengine.NewContext(r)
	perPage := 20
	cursorKey := modelVar + "_cursor"
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
	var items ArticleList
	var hasNext bool
	hasNext = false
	// Iterate over the results.
	t := q.Run(c)
	for {
		var ap Article
		_, err := t.Next(&ap)
		if err == datastore.Done {
			break
		}
		if err != nil {
			c.Errorf("fetching next Article: %v", err)
			hasNext = false
			break
		}
		items = append(items, ap)
	}
	if len(items) == perPage {
		hasNext = true
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
	listDataSet := map[string]interface{}{"items": items,
		"has_next": hasNext, "next_offset": nextOffset, "per_page": perPage, "model_name": modelVar}
	return listDataSet
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
		if typeOfT.Field(i).Tag.Get("verbose_name") == "-" {
			setDefaultErr := reflections.SetField(modelStruct, typeOfT.Field(i).Name, defaultValues[typeOfT.Field(i).Tag.Get("datastore_type")])
			if setDefaultErr != nil {
				continue
				//http.Error(w, setDefaultErr.Error(), http.StatusInternalServerError)
				//return
			}
			continue
		}
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
		case "DateTime":
			if r.FormValue(typeOfT.Field(i).Tag.Get("json")) == "" {
				err := reflections.SetField(modelStruct, typeOfT.Field(i).Name, defaultValues[typeOfT.Field(i).Tag.Get("datastore_type")])
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				continue
			}
			date_yyyymmddhhmm, _ := time.Parse("2006-01-02 15:04", r.FormValue(typeOfT.Field(i).Tag.Get("json")))
			err := reflections.SetField(modelStruct, typeOfT.Field(i).Name, date_yyyymmddhhmm)
			if err != nil {
				setDefaultErr := reflections.SetField(modelStruct, typeOfT.Field(i).Name, defaultValues[typeOfT.Field(i).Tag.Get("datastore_type")])
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
	beforePut(modelStruct)
	_, putErr := datastore.Put(c, key, modelStruct)
	if putErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	setDataSearchIndex(modelVar, keyName, modelStruct, c, w)
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
		// log.Println(typeOfT.Field(i).Name)
		// log.Println(typeOfT.Field(i).Tag.Get("datastore_type"))
		// log.Println(r.FormValue(typeOfT.Field(i).Tag.Get("json")))
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
		case "DateTime":
			if typeOfT.Field(i).Tag.Get("json") == "update" {
				setDefaultErr := reflections.SetField(modelStruct, typeOfT.Field(i).Name, defaultValues[typeOfT.Field(i).Tag.Get("datastore_type")])
				if setDefaultErr != nil {
					http.Error(w, setDefaultErr.Error(), http.StatusInternalServerError)
					return
				}
				continue
			}
			if typeOfT.Field(i).Tag.Get("json") == "created" {
				log.Println("skipped")
				continue
			}
			date_yyyymmddhhmm, _ := time.Parse("2006-01-02 15:04", r.FormValue(typeOfT.Field(i).Tag.Get("json")))
			err := reflections.SetField(modelStruct, typeOfT.Field(i).Name, date_yyyymmddhhmm)
			if err != nil {
				setDefaultErr := reflections.SetField(modelStruct, typeOfT.Field(i).Name, time.Now())
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
					//http.Error(w, err.Error(), http.StatusInternalServerError)
					//return
					continue
				}
			}
		}
	}
	beforePut(modelStruct)
	_, putErr := datastore.Put(c, key, modelStruct)
	if putErr != nil {
		http.Error(w, putErr.Error(), http.StatusInternalServerError)
		return
	}
	setDataSearchIndex(modelVar, keyName, modelStruct, c, w)
}

func beforePut(modelStruct interface{}) {
	imagesFromText(modelStruct)
}

func imagesFromText(modelStruct interface{}) {
	//var images []map[string]string
	s := reflect.ValueOf(modelStruct).Elem()
	typeOfT := s.Type()
	var images []map[string]string
	for i := 0; i < s.NumField(); i++ {
		log.Println(typeOfT.Field(i).Name)
		value, getErr := reflections.GetField(modelStruct, typeOfT.Field(i).Name)
		if getErr != nil {
			continue
		}
		log.Println(value)
		switch typeOfT.Field(i).Tag.Get("datastore_type") {
		case "Text":
			re, _ := regexp.Compile(`!\[(.*)\]\((.*)\)|!\[.*\]\[.*\]|\[.*\]: .*"".*""`)
			all := re.FindAllStringSubmatch(value.(string), -1)
			for _, v := range all {
				var inlineImage map[string]string
				inlineImage = make(map[string]string)
				inlineImage["filename"] = v[1]
				inlineImage["filepath"] = v[2]
				images = append(images, inlineImage)
				log.Println(v)
			}

		}
	}
	log.Println(images)
	jsonData, _ := json.Marshal(images)
	log.Println(jsonData)
	if string(jsonData) == "null" {
		setEmptyErr := reflections.SetField(modelStruct, "Images", "[]")
		if setEmptyErr != nil {
			log.Println("set images json error")
		}
	} else {
		setErr := reflections.SetField(modelStruct, "Images", string(jsonData))
		if setErr != nil {
			log.Println("set images json error")
		}
	}
}

func setDataSearchIndex(modelVar string, keyName string, modelStruct interface{}, c appengine.Context, w http.ResponseWriter) {
	docID := modelVar + "_" + keyName
	searchStruct := searchModels[modelVar]
	s := reflect.ValueOf(searchStruct).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		//log.Println(typeOfT.Field(i).Name)
		value, getErr := reflections.GetField(modelStruct, typeOfT.Field(i).Name)
		if getErr != nil {
			continue
		}
		setErr := reflections.SetField(searchStruct, typeOfT.Field(i).Name, value)
		if setErr != nil {
			continue
		}
	}
	index, err := search.Open("global")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = index.Put(c, docID, searchStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func imageUploadUrl(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	uploadURL, err := blobstore.UploadURL(c, "/admin/image/upload/handler", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	executeJSON(w, 200, map[string]interface{}{"uploadurl": uploadURL.Path})
}

func handleImageUpload(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	blobs, _, err := blobstore.ParseUpload(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file := blobs["file"]
	if len(file) == 0 {
		c.Errorf("no file uploaded")
		executeJSON(w, 200, map[string]interface{}{"message": "no file uploaded"})
		return
	}
	var imageOptions image.ServingURLOptions
	imageUrl, urlErr := image.ServingURL(c, file[0].BlobKey, &imageOptions)
	if urlErr != nil {
		http.Error(w, urlErr.Error(), http.StatusInternalServerError)
		return
	}
	modelName := "BlobStoreImage"
	key := datastore.NewKey(c, modelName, string(file[0].BlobKey), 0, nil)
	imageUrlString := "//" + imageUrl.Host + imageUrl.Path
	modelStruct := BlobStoreImage{ImageUrl: imageUrlString, Update: time.Now(), Created: time.Now()}
	_, putErr := datastore.Put(c, key, &modelStruct)
	if putErr != nil {
		http.Error(w, putErr.Error(), http.StatusInternalServerError)
		return
	}
	executeJSON(w, 201, map[string]interface{}{"message": "created", "filename": imageUrlString})
}
