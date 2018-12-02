package model

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

const nullSQL = "null"

// Initializer is for model when the have need extra initialize on save/update
type Initializer interface {
	// Initialize is the method to call att save/update
	Initialize()
}

// JSONB is an aggregate interface for json fields
type JSONB interface {
	json.Marshaler
	json.Unmarshaler
	driver.Valuer
	sql.Scanner
}

// Int64Slice is simple slice to handle it for json field
type Int64Slice []int64

// Int64Array is used to handle real array in database
type Int64Array []int64

// NullTime is null-time for json in null
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// NullInt64 is null int64 for json in null
type NullInt64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

// NullString is json friendly version of sql.NullString
type NullString struct {
	Valid  bool
	String string
}

// GenericJSONField is used to handle generic json data in postgres
type GenericJSONField map[string]interface{}

// MarshalJSON try to save this into json
func (gjf *GenericJSONField) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(*gjf))
}

// UnmarshalJSON try to unmarshal this from a json string
func (gjf *GenericJSONField) UnmarshalJSON(d []byte) error {
	tmp := make(map[string]interface{})
	err := json.Unmarshal(d, &tmp)
	if err != nil {
		return err
	}
	*gjf = tmp
	return nil
}

// StringJSONArray is use to handle string to string map in postgres
type StringJSONArray map[string]string

// Scan convert the json array ino string slice
func (is *Int64Slice) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}

	return json.Unmarshal(b, is)
}

// Value try to get the string slice representation in database
func (is Int64Array) Value() (driver.Value, error) {
	b, err := json.Marshal(is)
	if err != nil {
		return nil, err
	}
	// Its time to change [] to {}
	b = bytes.Replace(b, []byte("["), []byte("{"), 1)
	b = bytes.Replace(b, []byte("]"), []byte("}"), 1)

	return b, nil
}

// Scan convert the json array ino string slice
func (is *Int64Array) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	b = bytes.Replace(b, []byte("{"), []byte("["), 1)
	b = bytes.Replace(b, []byte("}"), []byte("]"), 1)

	return json.Unmarshal(b, is)
}

// Value try to get the string slice representation in database
func (is Int64Slice) Value() (driver.Value, error) {
	return json.Marshal(is)
}

// MarshalJSON try to marshaling to json
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return nt.Time.MarshalJSON()
	}

	return []byte(nullSQL), nil
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	tmp := &pq.NullTime{}
	err := tmp.Scan(value)
	if err != nil {
		return err
	}
	nt.Time, nt.Valid = tmp.Time, tmp.Valid
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// UnmarshalJSON try to unmarshal dae from input
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	text := strings.ToLower(string(b))
	if text == nullSQL {
		nt.Valid = false
		nt.Time = time.Time{}
		return nil
	}

	err := json.Unmarshal(b, &nt.Time)
	if err != nil {
		return err
	}

	nt.Valid = true
	return nil
}

// MarshalJSON try to marshaling to json
func (nt NullInt64) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return []byte(fmt.Sprintf(`%d`, nt.Int64)), nil
	}

	return []byte(nullSQL), nil
}

// UnmarshalJSON try to unmarshal dae from input
func (nt *NullInt64) UnmarshalJSON(b []byte) error {
	text := strings.ToLower(string(b))
	if text == nullSQL {
		nt.Valid = false

		return nil
	}

	err := json.Unmarshal(b, &nt.Int64)
	if err != nil {
		return err
	}

	nt.Valid = true
	return nil
}

// Scan implements the Scanner interface.
func (nt *NullInt64) Scan(value interface{}) error {
	inn := &sql.NullInt64{}
	err := inn.Scan(value)
	if err != nil {
		return err
	}
	nt.Int64 = inn.Int64
	nt.Valid = inn.Valid
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullInt64) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Int64, nil
}

// Scan convert the json array ino string slice
func (gjf *GenericJSONField) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}

	return json.Unmarshal(b, gjf)
}

// Value try to get the string slice representation in database
func (gjf GenericJSONField) Value() (driver.Value, error) {
	return json.Marshal(gjf)
}

// Scan convert the json array ino string slice
func (ss *StringJSONArray) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}

	return json.Unmarshal(b, ss)
}

// Value try to get the string slice representation in database
func (ss StringJSONArray) Value() (driver.Value, error) {
	return json.Marshal(ss)
}

// Scan implements the Scanner interface.
func (ns *NullString) Scan(value interface{}) error {
	tmp := &sql.NullString{}
	err := tmp.Scan(value)
	if err != nil {
		return err
	}
	ns.Valid = tmp.Valid
	ns.String = tmp.String
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// MarshalJSON try to marshaling to json
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}

	return []byte(nullSQL), nil
}

// UnmarshalJSON try to unmarshal dae from input
func (ns NullString) UnmarshalJSON(b []byte) error {
	text := strings.ToLower(string(b))
	if text == nullSQL {
		ns.Valid = false
		ns.String = ""
		return nil
	}

	err := json.Unmarshal(b, &ns.String)
	if err != nil {
		return err
	}

	ns.Valid = true
	return nil
}

// JSONBWrapper is a simple type to handle postgres jsonb
type JSONBWrapper struct {
	inner interface{}
}

// MarshalJSON try to save it to json
func (w *JSONBWrapper) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.inner)
}

// UnmarshalJSON try to load it from json
func (w *JSONBWrapper) UnmarshalJSON(b []byte) error {
	if w.inner == nil {
		tmp := make(GenericJSONField)
		w.inner = &tmp
	}
	return json.Unmarshal(b, w.inner)
}

// Value try to return database friendly value
func (w *JSONBWrapper) Value() (driver.Value, error) {
	return json.Marshal(w.inner)
}

// Scan try to scan value from database
func (w *JSONBWrapper) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return fmt.Errorf("unsupported type %T", src)
	}
	if w.inner == nil {
		tmp := make(GenericJSONField)
		w.inner = &tmp
	}

	return json.Unmarshal(b, w.inner)
}

// NewJSONB return the new JSONBWrapper for store and get from db
func NewJSONB(inner interface{}) *JSONBWrapper {
	return &JSONBWrapper{
		inner: inner,
	}
}
