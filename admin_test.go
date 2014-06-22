package main

import (
	"appengine/aetest"
	"appengine/datastore"
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSetNewEntity(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	data := url.Values{}
	data.Set("displaypage", "true")
	data.Add("title", "title")
	data.Add("url", "url")
	data.Add("pageorder", "1")
	data.Add("content", "content")
	r, _ := http.NewRequest("POST", "/admin/rest/adminpage", bytes.NewBufferString(data.Encode()))
	w := httptest.NewRecorder()
	modelVar := "adminpage"
	modelName := modelNames[modelVar]
	key := datastore.NewKey(c, modelName, "url", 0, nil)
	setNewEntity(w, r, modelVar)
	modelStruct := AdminPage{}
	if err := datastore.Get(c, key, &modelStruct); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, modelStruct.URL, "url", "url value check")

}
