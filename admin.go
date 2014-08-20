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
	q := datastore.NewQuery(modelName).Order("-pageorder").Limit(perPage)
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
	q := datastore.NewQuery(modelName).Order("-pageorder").Limit(perPage)
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
			if typeOfT.Field(i).Name == "PageOrder" {
				maxPageOrder := getPageOrderMaxValue(c, modelVar)
				maxPageOrder = maxPageOrder + 1
				setPageOrderErr := reflections.SetField(modelStruct, "PageOrder", maxPageOrder)
				if setPageOrderErr != nil {
					c.Errorf("cannot set value of max page order") //use default value
				}
				pageOrderIncre(c, modelVar)
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
			dateyyyymmddhhmm, _ := time.Parse("2006-01-02 15:04", r.FormValue(typeOfT.Field(i).Tag.Get("json")))
			err := reflections.SetField(modelStruct, typeOfT.Field(i).Name, dateyyyymmddhhmm)
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
		setURLErr := reflections.SetField(modelStruct, "URL", keyName)
		if setURLErr != nil {
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
	r.ParseForm()
	var parameters map[string]bool
	parameters = make(map[string]bool)
	for parameterName := range r.Form {
		log.Println(parameterName)
		parameters[parameterName] = true
	}
	for i := 0; i < s.NumField(); i++ {
		// log.Println(typeOfT.Field(i).Name)
		// log.Println(typeOfT.Field(i).Tag.Get("datastore_type"))
		// specic process for verbose_name == "-"
		// - record update time
		// - update pageorder if pageorder value is sent
		if typeOfT.Field(i).Tag.Get("verbose_name") == "-" {
			switch typeOfT.Field(i).Name {
			case "Update":
				setDefaultErr := reflections.SetField(modelStruct, typeOfT.Field(i).Name, defaultValues[typeOfT.Field(i).Tag.Get("datastore_type")])
				if setDefaultErr != nil {
					http.Error(w, setDefaultErr.Error(), http.StatusInternalServerError)
					return
				}
				continue
			case "PageOrder":
				if r.FormValue(typeOfT.Field(i).Tag.Get("json")) != "" {
					tmpStringInt, _ := strconv.Atoi(r.FormValue(typeOfT.Field(i).Tag.Get("json")))
					err := reflections.SetField(modelStruct, typeOfT.Field(i).Name, tmpStringInt)
					if err != nil {
						setDefaultErr := reflections.SetField(modelStruct, typeOfT.Field(i).Name, 0)
						if setDefaultErr != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
					}
				}
				continue
			default:
				continue
			}
		}
		// skipped if field == URL
		if typeOfT.Field(i).Name == "URL" {
			continue
		}
		if val, ok := parameters[typeOfT.Field(i).Tag.Get("json")]; ok {
			log.Println(typeOfT.Field(i).Name)
			log.Println("value exist")
			log.Println(val)
		} else {
			log.Println(typeOfT.Field(i).Name)
			log.Println("no value=>")
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
			dateyyyymmddhhmm, _ := time.Parse("2006-01-02 15:04", r.FormValue(typeOfT.Field(i).Tag.Get("json")))
			err := reflections.SetField(modelStruct, typeOfT.Field(i).Name, dateyyyymmddhhmm)
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
	var images []map[string]string
	value, getErr := reflections.GetField(modelStruct, "Content")
	if getErr != nil {
		//log.Println("cannot get value of Content field")
		return
	}
	// log.Println(value)
	re, _ := regexp.Compile(`!\[(.*)\]\((.*)\)|!\[.*\]\[.*\]|\[.*\]: .*"".*""`)
	all := re.FindAllStringSubmatch(value.(string), -1)
	for _, v := range all {
		var inlineImage map[string]string
		inlineImage = make(map[string]string)
		inlineImage["filename"] = v[1]
		inlineImage["filepath"] = v[2]
		images = append(images, inlineImage)
		// log.Println(v)
	}
	// log.Println(images)
	jsonData, _ := json.Marshal(images)
	// log.Println(jsonData)
	jsonDataString := string(jsonData)
	if jsonDataString == "null" {
		setEmptyErr := reflections.SetField(modelStruct, "Images", "[]")
		if setEmptyErr != nil {
			log.Println("set images json error")
		}
	} else {
		setErr := reflections.SetField(modelStruct, "Images", jsonDataString)
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

func imageUploadURL(w http.ResponseWriter, r *http.Request) {
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
	imageURL, urlErr := image.ServingURL(c, file[0].BlobKey, &imageOptions)
	if urlErr != nil {
		http.Error(w, urlErr.Error(), http.StatusInternalServerError)
		return
	}
	modelName := "BlobStoreImage"
	key := datastore.NewKey(c, modelName, string(file[0].BlobKey), 0, nil)
	imageURLString := "//" + imageURL.Host + imageURL.Path
	modelStruct := BlobStoreImage{ImageURL: imageURLString, Update: time.Now(), Created: time.Now()}
	_, putErr := datastore.Put(c, key, &modelStruct)
	if putErr != nil {
		http.Error(w, putErr.Error(), http.StatusInternalServerError)
		return
	}
	executeJSON(w, 201, map[string]interface{}{"message": "created", "filename": imageURLString})
}

func getPageOrderMaxValue(c appengine.Context, modelVar string) int {
	var maxPageOrder MaxPageOrder
	modelName := modelNames[modelVar]
	key := datastore.NewKey(c, "MaxPageOrder", modelName, 0, nil)
	err := datastore.Get(c, key, &maxPageOrder)
	if err != nil {
		newMaxPageOrder := MaxPageOrder{ModelName: modelName, MaxOrder: 0}
		_, putErr := datastore.Put(c, key, &newMaxPageOrder)
		if putErr != nil {
			c.Errorf("cannot put default value of max page order")
		}
		return 0
	}
	return maxPageOrder.MaxOrder
}

func pageOrderIncre(c appengine.Context, modelVar string) {
	var maxPageOrder MaxPageOrder
	modelName := modelNames[modelVar]
	key := datastore.NewKey(c, "MaxPageOrder", modelName, 0, nil)
	err := datastore.Get(c, key, &maxPageOrder)
	if err != nil {
		newMaxPageOrder := MaxPageOrder{ModelName: modelName, MaxOrder: 0}
		_, putErr := datastore.Put(c, key, &newMaxPageOrder)
		if putErr != nil {
			c.Errorf("cannot put default value of max page order")
		}
	}
	maxPageOrder.MaxOrder = maxPageOrder.MaxOrder + 1
	_, putUpdateErr := datastore.Put(c, key, &maxPageOrder)
	if putUpdateErr != nil {
		c.Errorf("cannot put update value of max page order")
	}
}
