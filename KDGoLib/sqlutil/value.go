package sqlutil

import (
	"bytes"
	"database/sql/driver"
	"reflect"
	"strings"

	"github.com/viethqc/gogstash/KDGoLib/errutil"
	"github.com/viethqc/gogstash/KDGoLib/jsonex"
)

// errors
var (
	ErrorUnsupportedScanType1 = errutil.NewFactory("unsupported scan type: %T")
	ErrorInvalidObjType1      = errutil.NewFactory("invalid obj type: %T")
	ErrorNoValueFound1        = errutil.NewFactory("no value found for key %v")
)

// SQLScanJSON set obj to value's JSON representation
func SQLScanJSON(obj interface{}, value interface{}) (err error) {
	switch val := value.(type) {
	case []byte:
		return jsonex.Unmarshal(val, obj)
	case nil:
		return nil
	default:
		return ErrorUnsupportedScanType1.New(nil, value)
	}
}

// SQLScanStrictJSON set obj to value's JSON representation, all field should exist in obj
func SQLScanStrictJSON(obj interface{}, value interface{}) (err error) {
	switch val := value.(type) {
	case []byte:
		decoder := jsonex.NewDecoder(bytes.NewBuffer(val))
		decoder.DisallowUnknownFields()
		return decoder.Decode(obj)
	case nil:
		return nil
	default:
		return ErrorUnsupportedScanType1.New(nil, value)
	}
}

// SQLValueJSON return obj's JSON representation which implements driver.Value
func SQLValueJSON(obj interface{}) (value driver.Value, err error) {
	if obj == nil {
		return
	}
	jsondata, err := jsonex.Marshal(obj)
	if err != nil {
		return
	}
	if bytes.Equal([]byte("{}"), jsondata) {
		return
	}
	return jsondata, nil
}

func ensureScanValue(obj interface{}) (refval reflect.Value, err error) {
	// Get the value of obj and make sure it's either a pointer or nil
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return rv, ErrorInvalidObjType1.New(nil, obj)
	}
	// So we can get its actual value
	pv := reflect.Indirect(rv)

	return pv, nil
}

// SQLScanString set obj to value's string representation
func SQLScanString(obj interface{}, value interface{}) (err error) {
	if value == nil {
		return
	}

	pv, err := ensureScanValue(obj)
	if err != nil {
		return
	}

	switch value.(type) {
	case []byte:
		pv.SetString(string(value.([]byte)))
		return
	case string:
		pv.SetString(value.(string))
		return
	default:
		return ErrorUnsupportedScanType1.New(nil, value)
	}
}

// SQLScanEnumString set obj to value's enum in stringMapEnum representation
func SQLScanEnumString(obj interface{}, value interface{}, stringMapEnum map[string]interface{}) (err error) {
	if value == nil {
		return
	}

	pv, err := ensureScanValue(obj)
	if err != nil {
		return
	}

	var enumstr string
	switch value.(type) {
	case []byte:
		enumstr = string(value.([]byte))
	case string:
		enumstr = value.(string)
	default:
		return ErrorUnsupportedScanType1.New(nil, value)
	}

	enumval, ok := stringMapEnum[enumstr]
	if !ok {
		return ErrorNoValueFound1.New(nil, value)
	}
	ev := reflect.ValueOf(enumval)
	pv.Set(ev)
	return
}

// SQLScanStringSlice set obj to value's stringslice representation
func SQLScanStringSlice(obj interface{}, value interface{}) (err error) {
	if value == nil {
		return
	}

	pv, err := ensureScanValue(obj)
	if err != nil {
		return
	}

	switch value.(type) {
	case []byte:
		stringslice := parseArray(string(value.([]byte)))
		pv.Set(reflect.ValueOf(stringslice))
		return
	case string:
		stringslice := parseArray(value.(string))
		pv.Set(reflect.ValueOf(stringslice))
		return
	default:
		return ErrorUnsupportedScanType1.New(nil, value)
	}
}

// SQLValueStringSlice return s postgresql representation which implements driver.Value
func SQLValueStringSlice(obj interface{}) (value driver.Value, err error) {
	if obj == nil {
		return
	}

	rv := reflect.ValueOf(obj)

	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			return
		}
		return SQLValueStringSlice(reflect.Indirect(rv).Interface())
	case reflect.Slice:
		switch rv.Type().Elem().Kind() {
		case reflect.String:
			strs := []string{}
			for i, rvlen := 0, rv.Len(); i < rvlen; i++ {
				rvelem := rv.Index(i)
				strelem := `"` + strings.Replace(strings.Replace(rvelem.String(), `\`, `\\\`, -1), `"`, `\"`, -1) + `"`
				strs = append(strs, strelem)
			}
			return "{" + strings.Join(strs, ",") + "}", nil
		}
	}

	return "", ErrorInvalidObjType1.New(nil, obj)
}
