package main

import (
	"github.com/oleiade/reflections"
	"testing"
)

func TestImagesFromText(t *testing.T) {
	modelStruct := models["adminpage"]
	contentString := `
	![file](//localhost:8080/_ah/img/l6VnJ25q7BizPxLc4OQelQ==)

	![file](//localhost:8080/_ah/img/RIXogd9LBBsMmoqZXZ_eEQ==)
	`
	setErr := reflections.SetField(modelStruct, "Content", contentString)
	if setErr != nil {
		t.Errorf("Connot set data for test")
	}
	imagesFromText(modelStruct)
	jsonValue, _ := reflections.GetField(modelStruct, "Images")
	assertValue := `[{"filename":"file","filepath":"//localhost:8080/_ah/img/l6VnJ25q7BizPxLc4OQelQ=="},{"filename":"file","filepath":"//localhost:8080/_ah/img/RIXogd9LBBsMmoqZXZ_eEQ=="}]`
	if jsonValue.(string) != assertValue {
		t.Errorf("imagesFromText is not working properly")
	}
}
