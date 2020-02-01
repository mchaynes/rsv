package rsv

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// UnmarshalRow unmarshalls a row into the input interface
// the input interface must be a ptr to a struct, and will panic if not a ptr.
//
// UnmarshalRow supports all builtin types, will parse input strings into appropriate type
// To specify index of the row to use for this data, use the "idx" tag
//     type A struct {
//         B string `idx:"0"`
//         C string `idx:"2"`
//         I int    `idx:"3"`
//     }
// You can also specify `omitempty` on pointers to keep the value nil if the input string is empty
//     type A struct {
//         B *string `idx:"0,omitempty"`
//     }
// If input row is `[]string{""}`, A.B will be `nil`
//
// UnmarshalRow returns ErrFailedToParse if any of the values are not parsable into their respective types
func UnmarshalRow(row []string, v interface{}) error {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		panic("v should be a pointer to a struct")
	}
	// deref the pointer
	t = t.Elem()
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < t.NumField(); i++ {
		fv := val.Field(i)
		idx, omitempty, err := getTagInfo(t.Field(i).Tag)
		if err != nil && !errors.Is(err, errIdxTagNotSet{}) {
			return err
		}
		if errors.Is(err, errIdxTagNotSet{}) && fv.Kind() != reflect.Struct {
			continue
		}
		err = set(fv, row, idx, omitempty)
		if err != nil {
			return err
		}
	}
	return nil
}

func set(v reflect.Value, row []string, idx int, omitempty bool) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(row[idx])
	case reflect.Int, reflect.Int64, reflect.Int32:
		i, err := strconv.ParseInt(row[idx], 10, 64)
		if err != nil {
			return ErrFailedToParse{Value: row[idx], Err: err}
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(row[idx], 10, 64)
		if err != nil {
			return ErrFailedToParse{Value: row[idx], Err: err}
		}
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		var (
			f   float64
			err error
		)
		if len(row[idx]) > 0 {
			f, err = strconv.ParseFloat(row[idx], 64)
			if err != nil {
				return ErrFailedToParse{Value: row[idx], Err: err}
			}
		}
		v.SetFloat(f)
	case reflect.Ptr:
		// ignore empty values for ptrs if omitempty
		if omitempty && len(row[idx]) == 0 {
			return nil
		}
		// create a new ptr to the type of this pointer. i.e. create string if *string
		v.Set(reflect.New(v.Type().Elem()))
		// call set again, but with deref'd value. note that this will recurse down for ptr to ptr's (and so on)
		return set(v.Elem(), row, idx, omitempty)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			idx, omitempty, err := getTagInfo(v.Type().Field(i).Tag)
			if err != nil && !errors.Is(err, errIdxTagNotSet{}) {
				return err
			}
			if errors.Is(err, errIdxTagNotSet{}) && fv.Kind() != reflect.Struct {
				continue
			}
			err = set(fv, row, idx, omitempty)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getTagInfo(tag reflect.StructTag) (int, bool, error) {
	t, ok := tag.Lookup("idx")
	var omitempty bool
	if strings.Contains(t, ",omitempty") {
		omitempty = true
		t = strings.Split(t, ",")[0]
	}
	if !ok {
		return -1, omitempty, errIdxTagNotSet{}
	}
	idx, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return -1, omitempty, ErrFailedToParse{Value: t, Err: err}
	}
	return int(idx), omitempty, nil
}
