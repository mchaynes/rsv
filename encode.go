package rsv

import (
	"errors"
	"reflect"
)

// MarshalRow marshals an interface into a slice of interface{}
// Will create a row that is as long as the maximum idx tag, with
// nil values for unspecified indices.
// for example, the following struct
//     type A struct {
//         Min int `idx:"0"`
//         Max int `idx:"3"`
//     }
// if A.Min = 1, and A.Max = 2, MarshalRow will produce the following row:
//    []interface{1, nil, nil, 2}
func MarshalRow(i interface{}) (row []interface{}, err error) {
	v := reflect.ValueOf(i)
	t := reflect.TypeOf(i)
	kind := v.Kind()
	if kind == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
		kind = v.Kind()
	}
	if kind != reflect.Struct {
		panic("i should be a struct or ptr to a struct")
	}
	m := make(map[int]interface{})
	for j := 0; j < v.NumField(); j++ {
		fv := v.Field(j)
		idx, _, err := getTagInfo(t.Field(j).Tag)
		if err != nil && !errors.Is(err, errIdxTagNotSet{}) {
			return row, err
		}
		if errors.Is(err, errIdxTagNotSet{}) && fv.Kind() != reflect.Struct {
			continue
		}
		m, err = get(fv, idx, m)
		if err != nil {
			return nil, err
		}
	}
	return buildSlice(m), nil
}

func get(v reflect.Value, idx int, m map[int]interface{}) (map[int]interface{}, error) {
	switch v.Kind() {
	case reflect.String:
		m[idx] = v.String()
	case reflect.Int, reflect.Int64, reflect.Int32:
		m[idx] = v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		m[idx] = v.Uint()
	case reflect.Float32, reflect.Float64:
		m[idx] = v.Float()
	case reflect.Ptr:
		if v.CanAddr() {
			return get(v.Elem(), idx, m)
		}
		m[idx] = nil
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			idx, _, err := getTagInfo(v.Type().Field(i).Tag)
			if err != nil && !errors.Is(err, errIdxTagNotSet{}) {
				return nil, err
			}
			if errors.Is(err, errIdxTagNotSet{}) && fv.Kind() != reflect.Struct {
				continue
			}
			m, err = get(fv, idx, m)
			if err != nil {
				return m, err
			}
		}
	}
	return m, nil
}

func buildSlice(m map[int]interface{}) []interface{} {
	size := max(m) + 1
	s := make([]interface{}, size)
	for k, v := range m {
		s[k] = v
	}
	return s
}

func max(m map[int]interface{}) int {
	max := 0
	for i := range m {
		if max < i {
			max = i
		}
	}
	return max
}
