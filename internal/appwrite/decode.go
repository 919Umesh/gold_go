package appwrite

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"
)

// Decode extracts the raw JSON data from an Appwrite model (Document, DocumentList, etc.)
// and unmarshals it into the target. This bypasses the SDK's broken Decode() method.
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

// DecodeListItem is a helper to extract a specific document from a ListDocuments response.
// Indices are 0-based. This is necessary because individual Documents inside a List
// response do NOT have their own 'data' fields populated by the SDK.
func DecodeListItem(list interface{}, index int, target interface{}) error {
	var rawDocs struct {
		Documents []interface{} `json:"documents"`
	}

	// Decode the list into our temporary struct to get the array
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
