package kstrct

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func SetReflectFieldValue(fld reflect.Value, value any) error {
	var errRecover error
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				errRecover = errors.New(x)
			case error:
				errRecover = x
			default:
				errRecover = fmt.Errorf("unknown panic: %v", r)
			}
		}
	}()
	valueToSet := reflect.ValueOf(value)
	switch fld.Kind() {
	case valueToSet.Kind():
		fld.Set(valueToSet)
	case reflect.Ptr:
		unwrapped := fld.Elem()
		if !unwrapped.IsValid() {
			newUnwrapped := reflect.New(fld.Type().Elem())
			errRecover = SetReflectFieldValue(newUnwrapped, value)
			fld.Set(newUnwrapped)
		} else {
			errRecover = SetReflectFieldValue(unwrapped, value)
		}
	case reflect.Interface:
		unwrapped := fld.Elem()
		errRecover = SetReflectFieldValue(unwrapped, value)
	case reflect.Struct:
		switch v := value.(type) {
		case string:
			if strings.Contains(v, ":") || strings.Contains(v, "-") {
				l := len("2006-01-02T15:04")
				if strings.Contains(v[:l], "T") {
					if len(v) >= l {
						t, err := time.Parse("2006-01-02T15:04", v[:l])
						if err != nil {
							fld.Set(reflect.ValueOf(t))
						}
					}
				} else if len(v) >= len("2006-01-02 15:04:05") {
					t, err := time.Parse("2006-01-02 15:04:05", v[:len("2006-01-02 15:04:05")])
					if err == nil {
						fld.Set(reflect.ValueOf(t))
					}
				} else {
					errRecover = fmt.Errorf("setFieldValue Struct: doesn't match any case: %v", v)
				}
			}
		case time.Time:
			fld.Set(valueToSet)

		case []any:
			// walk the fields
			for i := 0; i < fld.NumField(); i++ {
				errRecover = SetReflectFieldValue(fld.Field(i), v[i])
			}
		default:
			errRecover = fmt.Errorf("case struct SetReflectFieldValue got value: %v is not valid for, fieldName : %s", valueToSet.Interface(), fld.Type().Name())
		}
	case reflect.String:
		switch valueToSet.Kind() {
		case reflect.String:
			fld.SetString(valueToSet.String())
		case reflect.Struct:
			fld.SetString(valueToSet.String())
		default:
			if valueToSet.IsValid() {
				fld.Set(valueToSet)
			} else {
				errRecover = fmt.Errorf("case string SetReflectFieldValue got value: %v is not valid for, fieldName : %s", valueToSet.Interface(), fld.Type().Name())
			}
		}
	case reflect.Int:
		switch v := value.(type) {
		case int64:
			fld.SetInt(v)
		case string:
			if v, err := strconv.Atoi(v); err == nil {
				fld.SetInt(int64(v))
			}
		case int:
			fld.SetInt(int64(v))
		}
	case reflect.Int64:
		switch v := value.(type) {
		case int64:
			fld.SetInt(v)
		case string:
			if v, err := strconv.Atoi(v); err == nil {
				fld.SetInt(int64(v))
			}
		case []byte:
			if v, err := strconv.Atoi(string(v)); err != nil {
				fld.SetInt(int64(v))
			}
		case int:
			fld.SetInt(int64(v))
		}
	case reflect.Bool:
		switch valueToSet.Kind() {
		case reflect.Int:
			if value == 1 {
				fld.SetBool(true)
			}
		case reflect.Int64:
			if value == int64(1) {
				fld.SetBool(true)
			}
		case reflect.Uint64:
			if value == uint64(1) {
				fld.SetBool(true)
			}
		case reflect.String:
			if value == "1" {
				fld.SetBool(true)
			} else if value == "true" {
				fld.SetBool(true)
			}
		}
	case reflect.Uint:
		switch v := value.(type) {
		case uint:
			fld.SetUint(uint64(v))
		case uint64:
			fld.SetUint(v)
		case int64:
			fld.SetUint(uint64(v))
		case int:
			fld.SetUint(uint64(v))
		}
	case reflect.Uint64:
		switch v := value.(type) {
		case uint:
			fld.SetUint(uint64(v))
		case uint64:
			fld.SetUint(v)
		case int64:
			fld.SetUint(uint64(v))
		case int:
			fld.SetUint(uint64(v))
		}
	case reflect.Float64:
		if v, ok := value.(float64); ok {
			fld.SetFloat(v)
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
			errRecover = fmt.Errorf("cannot set float64 from setReflectFieldValue :%T", value)
		}
	case reflect.Slice:
		targetType := fld.Type()
		typeName := targetType.String()
		if strings.HasPrefix(typeName, "[]") {
			array := reflect.New(targetType).Elem()
			for _, v := range strings.Split(valueToSet.String(), ",") {
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
					errRecover = fmt.Errorf("filling slice received:%s", typeName)
				}
			}
			fld.Set(array)
		}
	default:
		switch v := value.(type) {
		case []byte:
			fld.SetString(string(v))
		default:
			errRecover = fmt.Errorf("setFieldValue: case not handled , unable to fill struct,field kind: %v, value to fill : %v", fld.Kind(), value)
		}
	}
	return errRecover
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func SnakeCaseToTitle(inputUnderScoreStr string) (camelCase string) {
	//snake_case to camelCase
	if strings.Contains(inputUnderScoreStr, "_") {
		sp := strings.Split(inputUnderScoreStr, "_")
		for i := range sp {
			sp[i] = strings.ToUpper(string(sp[i][0])) + sp[i][1:]
		}
		inputUnderScoreStr = strings.Join(sp, "")
	} else {
		fl := strings.ToUpper(inputUnderScoreStr)[0]
		inputUnderScoreStr = string(fl) + inputUnderScoreStr[1:]
	}
	return inputUnderScoreStr
}

func GetInfos[T comparable](strct *T, tags ...string) (fields []string, fValues map[string]any, fTypes map[string]string, fTags map[string][]string) {
	fields = []string{}
	fValues = map[string]any{}
	fTypes = map[string]string{}
	fTags = map[string][]string{}

	s := reflect.ValueOf(strct).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fname := typeOfT.Field(i).Name
		fname = ToSnakeCase(fname)
		fvalue := f.Interface()
		ftype := f.Type().Name()

		fields = append(fields, fname)
		fTypes[fname] = ftype
		fValues[fname] = fvalue
		for _, t := range tags {
			if ftag, ok := typeOfT.Field(i).Tag.Lookup(t); ok {
				tags := strings.Split(ftag, ";")
				fTags[fname] = append(fTags[fname], tags...)
			}
		}
	}
	return fields, fValues, fTypes, fTags
}
