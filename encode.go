package ledger

import (
	"bytes"
	"encoding/gob"
	"errors"
	"reflect"
)

// Encodes a generic object to gob encoding.
// Stores encrypted value in a byte slice.
// Object must have fixed size in memory.
func EncryptData[T any](data T) ([]byte, error) {
	// Verify that data != nil
	if v := reflect.ValueOf(data); (v.Kind() == reflect.Ptr ||
		v.Kind() == reflect.Interface ||
		v.Kind() == reflect.Slice ||
		v.Kind() == reflect.Map ||
		v.Kind() == reflect.Chan ||
		v.Kind() == reflect.Func) && v.IsNil() {
		return nil, errors.New("unable to encrypt nil values")
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(data)

	return buf.Bytes(), nil
}

// Decodes a generic object from gob encoding.
// Type T must match the type used to encrypt the object.
func DecryptData[T any](data []byte) (T, error) {
	buf := bytes.NewReader(data)
	decoder := gob.NewDecoder(buf)
	var value T
	decoder.Decode(&value)

	return value, nil
}
