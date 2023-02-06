package kstrct

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

func ToSnakeCase(str string) string {
	var result strings.Builder
	var lastUpper bool
	for i, c := range str {
		if c >= 'A' && c <= 'Z' {
			if i > 0 && !lastUpper {
				result.WriteRune('_')
			}
			lastUpper = true
			result.WriteRune(unicode.ToLower(c))
		} else {
			lastUpper = false
			result.WriteRune(c)
		}
	}
	return result.String()
}

func SnakeCaseToTitle(inputUnderScoreStr string) (camelCase string) {
	// snake_case to camelCase
	var buffer strings.Builder
	var nextUpper bool
	for i, r := range inputUnderScoreStr {
		if r == '_' {
			nextUpper = true
			continue
		}
		if nextUpper || i == 0 {
			buffer.WriteRune(unicode.ToUpper(r))
			nextUpper = false
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

// GET INFO STRUCT

// fieldsPool is a sync.Pool that can be used to avoid allocating
// new slices on each call to the GetInfos function.
var fieldsPool = sync.Pool{
	New: func() interface{} {
		s := make([]string, 0, 32)
		return &s
	},
}

// fValuesPool is a sync.Pool that can be used to avoid allocating
// new maps on each call to the GetInfos function.
var fValuesPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]interface{})
	},
}

// fTagsPool is a sync.Pool that can be used to avoid allocating
// new maps on each call to the GetInfos function.
var fTagsPool = sync.Pool{
	New: func() interface{} {
		return make(map[string][]string)
	},
}

type Info struct {
	Fields []string
	Values map[string]interface{}
	Types  map[string]string
	Tags   map[string][]string
}

func GetInfos[T comparable](strct *T, tagsToCheck ...string) *Info {
	// Use the sync.Pool to get a slice and maps to use for the fields, values, and tags maps.
	fields := fieldsPool.Get().(*[]string)
	values := fValuesPool.Get().(map[string]interface{})
	tags := fTagsPool.Get().(map[string][]string)

	// Clear the slice and maps to reuse them.
	*fields = (*fields)[:0]
	for k := range values {
		delete(values, k)
	}
	for k := range tags {
		delete(tags, k)
	}

	s := reflect.ValueOf(strct).Elem()
	typeOfT := s.Type()
	numFields := s.NumField()

	// Pre-allocate the fields slice to avoid reallocations on each iteration.
	if cap((*fields)) < numFields {
		nn := []string{}
		fields = &nn
		// Create a new map to store the field types.
		types := make(map[string]string)

		for i := 0; i < numFields; i++ {
			f := s.Field(i)
			fname := typeOfT.Field(i).Name
			fname = ToSnakeCase(fname)
			fvalue := f.Interface()
			ftype := f.Type()

			*fields = append(*fields, fname)
			values[fname] = fvalue
			if ftype.Kind() == reflect.Ptr {
				types[fname] = ftype.Elem().String()
			} else {
				types[fname] = ftype.String()
			}
			for _, t := range tagsToCheck {
				if ftag, ok := typeOfT.Field(i).Tag.Lookup(t); ok {
					tagList := strings.Split(ftag, ";")
					tags[fname] = append(tags[fname], tagList...)
				}
			}
		}

		// Create a new Info struct and return it.
		info := &Info{
			Fields: *fields,
			Values: values,
			Types:  types,
			Tags:   tags,
		}

		// Return the fields slice and maps to the sync.Pool for reuse.
		fieldsPool.Put(fields)
		fValuesPool.Put(values)
		fTagsPool.Put(tags)

		return info
	}
	// If the capacity of the fields slice is not less than the number of fields,
	// we can use the existing slice and maps.
	types := make(map[string]string)

	for i := 0; i < numFields; i++ {
		f := s.Field(i)
		fname := typeOfT.Field(i).Name
		fname = ToSnakeCase(fname)
		fvalue := f.Interface()
		ftype := f.Type().Name()

		*fields = append(*fields, fname)
		values[fname] = fvalue
		types[fname] = ftype
		for _, t := range tagsToCheck {
			if ftag, ok := typeOfT.Field(i).Tag.Lookup(t); ok {
				tagList := strings.Split(ftag, ";")
				tags[fname] = append(tags[fname], tagList...)
			}
		}
	}

	// Create a new Info struct and return it.
	info := &Info{
		Fields: *fields,
		Values: values,
		Types:  types,
		Tags:   tags,
	}

	// Return the fields slice and maps to the sync.Pool for reuse.
	fieldsPool.Put(fields)
	fValuesPool.Put(values)
	fTagsPool.Put(tags)

	return info
}

func SetReflectFieldValue(fld reflect.Value, value interface{}) error {
	var errPanic error
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				errPanic = errors.New(x)
			case error:
				errPanic = x
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				errPanic = fmt.Errorf("%v", r)
			}
		}
	}()
	vToSet := reflect.ValueOf(value)
	if vToSet.Kind() == fld.Kind() {
		fld.Set(vToSet)
		return nil
	}

	switch fld.Kind() {
	case reflect.Bool:
		switch v := value.(type) {
		case bool:
			fld.SetBool(v)
		case int, int64, int32, uint, uint64, float32, float64, uint32:
			fld.SetBool(v != 0)
		case string:
			if v == "1" || v == "true" {
				fld.SetBool(true)
			} else if v == "0" || v == "false" {
				fld.SetBool(false)
			} else {
				return fmt.Errorf("invalid bool string value: %v", v)
			}
		default:
			if vToSet.IsValid() {
				fld.Set(vToSet)
			} else {
				fld.Set(reflect.Zero(fld.Type()))
			}
		}
		return nil
	case reflect.String:
		switch v := value.(type) {
		case string:
			fld.SetString(v)
		case time.Time:
			fld.SetString(v.String())
		default:
			if vToSet.IsValid() {
				fld.Set(vToSet)
			} else {
				fld.Set(reflect.Zero(fld.Type()))
			}
		}
		return nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		switch v := value.(type) {
		case uint:
			fld.SetUint(uint64(v))
		case uint64:
			fld.SetUint(v)
		case int:
			fld.SetUint(uint64(v))
		case int64:
			fld.SetUint(uint64(v))
		case uint32:
			fld.SetUint(uint64(v))
		case int32:
			fld.SetUint(uint64(v))
		case string:
			if v, err := strconv.Atoi(v); err == nil {
				fld.SetUint(uint64(v))
			}
		default:
			if vToSet.IsValid() {
				fld.Set(vToSet)
			} else {
				fld.Set(reflect.Zero(fld.Type()))
			}
		}
		return nil
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		switch v := value.(type) {
		case int:
			fld.SetInt(int64(v))
		case int64:
			fld.SetInt(v)
		case uint:
			fld.SetInt(int64(v))
		case uint64:
			fld.SetInt(int64(v))
		case int32:
			fld.SetInt(int64(v))
		case int16:
			fld.SetInt(int64(v))
		case int8:
			fld.SetInt(int64(v))
		case string:
			if v, err := strconv.Atoi(v); err == nil {
				fld.SetInt(int64(v))
			}
		default:
			if vToSet.IsValid() {
				fld.Set(vToSet)
			} else {
				fld.Set(reflect.Zero(fld.Type()))
			}
		}
		return nil
	case reflect.Struct:
		switch v := value.(type) {
		case time.Time:
			if len(v.String()) >= len("2006-01-02 15:04:05") {
				v.Format("2006-01-02 15:04:05")
			} else {
				v.Format("2006-01-02 15:04")
			}
			fld.Set(reflect.ValueOf(v))
			return nil
		case int:
			t := time.Unix(int64(v), 0)
			fld.Set(reflect.ValueOf(t))
		case uint:
			t := time.Unix(int64(v), 0)
			fld.Set(reflect.ValueOf(t))
		case int64:
			t := time.Unix(v, 0)
			fld.Set(reflect.ValueOf(t))
		case string:
			// Use a regular expression to match the desired date format
			v = strings.ReplaceAll(v, "T", " ")
			long := false
			if len(v) >= len("2006-01-02 15:04:05") {
				long = true
				v = v[:len("2006-01-02 15:04:05")]
			} else {
				v = v[:len("2006-01-02 15:04")]
			}
			if long {
				t, err := time.Parse("2006-01-02 15:04:05", v)
				if err != nil {
					return err
				}
				fld.Set(reflect.ValueOf(t))
			} else {
				t, err := time.Parse("2006-01-02 15:04", v)
				if err != nil {
					return err
				}
				fld.Set(reflect.ValueOf(t))
			}
		case *string:
			// Use a regular expression to match the desired date format
			*v = strings.ReplaceAll(*v, "T", " ")
			long := false
			if len(*v) >= len("2006-01-02 15:04:05") {
				long = true
				*v = (*v)[:len("2006-01-02 15:04:05")]
			} else {
				*v = (*v)[:len("2006-01-02 15:04")]
			}
			if long {
				t, err := time.Parse("2006-01-02 15:04:05", *v)
				if err != nil {
					return err
				}
				fld.Set(reflect.ValueOf(t))
			} else {
				t, err := time.Parse("2006-01-02 15:04", *v)
				if err != nil {
					return err
				}
				fld.Set(reflect.ValueOf(t))
			}
		case []interface{}:
			// Walk the fields
			for i := 0; i < fld.NumField(); i++ {
				if err := SetReflectFieldValue(fld.Field(i), v[i]); err != nil {
					return err
				}
			}
			return nil
		default:
			if vToSet.Type().AssignableTo(fld.Type()) {
				fld.Set(vToSet)
			} else if vToSet.Kind() == reflect.Slice {
				// Convert the value slice to a slice of the correct element type
				sliceType := reflect.SliceOf(fld.Type().Elem())
				convertedSlice := reflect.MakeSlice(sliceType, vToSet.Len(), vToSet.Cap())
				reflect.Copy(convertedSlice, vToSet)
				fld.Set(convertedSlice)
			} else {
				return fmt.Errorf("cannot assign value of type %s to field of type %s", vToSet.Type(), fld.Type())
			}
		}
		return nil
	case reflect.Float64, reflect.Float32:
		if v, ok := value.(float64); ok {
			fld.SetFloat(v)
		} else if v, ok := value.(float32); ok {
			fld.SetFloat(float64(v))
		} else if v, ok := value.(string); ok {
			f64, err := strconv.ParseFloat(v, 64)
			if err == nil {
				fld.SetFloat(f64)
			}
		} else if v, ok := value.([]byte); ok {
			f64, err := strconv.ParseFloat(string(v), 64)
			if err == nil {
				fld.SetFloat(f64)
			}
		} else {
			if vToSet.IsValid() {
				fld.Set(vToSet)
			} else {
				fld.Set(reflect.Zero(fld.Type()))
			}
		}
		return nil
	case reflect.Interface:
		unwrapped := fld.Elem()
		return SetReflectFieldValue(unwrapped, value)
	case reflect.Slice:
		targetType := fld.Type()
		typeName := targetType.String()
		if strings.HasPrefix(typeName, "[]") {
			array := reflect.New(targetType).Elem()
			for _, v := range strings.Split(fmt.Sprintf("%v", value), ",") {
				switch typeName[2:] {
				case "string":
					array = reflect.Append(array, reflect.ValueOf(v))
				case "int":
					if vv, err := strconv.Atoi(v); err == nil {
						array = reflect.Append(array, reflect.ValueOf(vv))
					}
				case "uint":
					if vv, err := strconv.ParseUint(v, 10, 64); err == nil {
						array = reflect.Append(array, reflect.ValueOf(uint(vv)))
					}
				case "float64":
					if vv, err := strconv.ParseFloat(v, 64); err == nil {
						array = reflect.Append(array, reflect.ValueOf(vv))
					}
				default:
					return fmt.Errorf("filling slice received:%s", typeName)
				}
			}
			fld.Set(array)
		}
		return nil
	case reflect.Ptr:
		unwrapped := fld.Elem()
		if !unwrapped.IsValid() {
			newUnwrapped := reflect.New(fld.Type().Elem())
			if err := SetReflectFieldValue(newUnwrapped, value); err != nil {
				return err
			}
			fld.Set(newUnwrapped)
		} else {
			if err := SetReflectFieldValue(unwrapped, value); err != nil {
				return err
			}
		}
		return nil
	}
	if errPanic != nil {
		return errPanic
	}
	return nil
}
