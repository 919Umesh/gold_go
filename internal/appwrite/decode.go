package appwrite

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"
)

func Decode(model interface{}, target interface{}) error {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Access unexported 'data' field
	dataField := v.FieldByName("data")
	if !dataField.IsValid() {
		return fmt.Errorf("field 'data' not found in model %T", model)
	}

	ptr := reflect.NewAt(dataField.Type(), unsafe.Pointer(dataField.UnsafeAddr())).Elem()
	dataBytes, ok := ptr.Interface().([]byte)
	if !ok {
		return fmt.Errorf("field 'data' is not []byte in model %T", model)
	}

	if len(dataBytes) == 0 {
		return fmt.Errorf("model %T has empty data (SDK limitation in lists)", model)
	}

	return json.Unmarshal(dataBytes, target)
}

func DecodeListItem(list interface{}, index int, target interface{}) error {
	var rawDocs struct {
		Documents []interface{} `json:"documents"`
	}

	if err := Decode(list, &rawDocs); err != nil {
		return err
	}

	if index < 0 || index >= len(rawDocs.Documents) {
		return fmt.Errorf("document at index %d not found in list", index)
	}

	docBytes, err := json.Marshal(rawDocs.Documents[index])
	if err != nil {
		return err
	}

	return json.Unmarshal(docBytes, target)
}
