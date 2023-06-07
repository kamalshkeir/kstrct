package kstrct

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var Debug = false

func SetReflectFieldValue(fld reflect.Value, value any, isTime ...bool) error {
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

	var errReturn error
	switch fld.Kind() {
	case reflect.Pointer:
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
		return errReturn
	case reflect.Bool:
		switch v := value.(type) {
		case bool:
			fld.SetBool(v)
		case string:
			if v == "1" || v == "true" {
				fld.SetBool(true)
			} else if v == "0" || v == "false" {
				fld.SetBool(false)
			} else {
				errReturn = fmt.Errorf("invalid bool string value: %v", v)
			}
		case int, int64, int32, uint, uint64, float32, float64, uint32:
			// Convert numeric values to boolean values
			if Debug {
				fmt.Printf("value: %v, typeValue: %T %v \n", v, v, v == 0)
			}
			if vToSet.Int() == 0 {
				fld.SetBool(false)
			} else {
				fld.SetBool(true)
			}
			if float32(vToSet.Float()) != float32(0) {
				fld.SetBool(true)
			} else {
				fld.SetBool(false)
			}
		case KV:
			return SetReflectFieldValue(fld, v.Value)
		default:
			if vToSet.IsValid() {
				fld.Set(vToSet)
			} else {
				errReturn = fmt.Errorf("zero value setted: cannot assign value of type %s to field of type %s", vToSet.Type(), fld.Type())
			}
		}
		return errReturn
	case reflect.String:
		switch v := value.(type) {
		case string:
			fld.SetString(v)
		case time.Time:
			fld.SetString(v.String())
		case float64, float32, int64, int32, uint, int, uint64, uint32:
			fld.SetString(fmt.Sprintf("%v", v))
		case KV:
			return SetReflectFieldValue(fld, v.Value)
		default:
			if vToSet.IsValid() {
				fld.Set(vToSet)
			} else {
				errReturn = fmt.Errorf("sprintf value setted: cannot assign value of type %s to field of type %s", vToSet.Type(), fld.Type())
				fld.SetString(fmt.Sprintf("%v", v))
			}
		}
		return errReturn
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
		case float64:
			fld.SetUint(uint64(v))
		case float32:
			fld.SetUint(uint64(v))
		case KV:
			return SetReflectFieldValue(fld, v.Value)
		default:
			if vToSet.IsValid() {
				fld.Set(vToSet)
			} else {
				errReturn = fmt.Errorf("zero value setted: cannot assign value of type %s to field of type %s", vToSet.Type(), fld.Type())
			}
		}
		return errReturn
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
		case float64:
			fld.SetUint(uint64(v))
		case float32:
			fld.SetUint(uint64(v))
		case KV:
			return SetReflectFieldValue(fld, v.Value)
		default:
			if vToSet.IsValid() {
				fld.Set(vToSet)
			} else {
				errReturn = fmt.Errorf("zero value setted: cannot assign value of type %s to field of type %s", vToSet.Type(), fld.Type())
			}
		}
		return errReturn
	case reflect.Struct:
		switch v := value.(type) {
		case map[string]any:
			for i := 0; i < fld.NumField(); i++ {
				field := fld.Field(i)
				tfield := fld.Type().Field(i)
				fname := tfield.Name
				if neww, ok := v[ToSnakeCase(fname)]; ok {
					err := SetReflectFieldValue(field, neww)
					if err != nil {
						return err
					}
				} else if kstrctTag, ok := tfield.Tag.Lookup("kname"); ok {
					if kstrctTag == "-" {
						continue
					}
					if val, ok := v[kstrctTag]; ok {
						err := SetReflectFieldValue(field, val)
						if err != nil {
							return err
						}
					}
				} else {
					sl := []KV{}
					for k, vv := range v {
						sl = append(sl, KV{Key: k, Value: vv})
					}
					err := FillFromKV(fld.Addr().Interface(), sl, true)
					if err != nil {
						return err
					}
				}
			}
		case map[int]any:
			for i, value := range v {
				field := fld.Field(i)
				err := SetReflectFieldValue(field, value)
				if err != nil {
					return err
				}
			}
		case []KV:
			err := FillFromKV(fld.Addr().Interface(), v, true)
			if err != nil {
				return err
			}
		case KV:
			return SetReflectFieldValue(fld, v.Value)
		case time.Time:
			fld.Set(reflect.ValueOf(v))
			return nil
		case int64:
			t := time.Unix(v, 0)
			fld.Set(reflect.ValueOf(t))
			return nil
		case int:
			t := time.Unix(int64(v), 0)
			fld.Set(reflect.ValueOf(t))
			return nil
		case uint:
			t := time.Unix(int64(v), 0)
			fld.Set(reflect.ValueOf(t))
			return nil
		case string:
			// Use a regular expression to match the desired date format
			if strings.Contains(v, ":") {
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
						fmt.Println("error set reflect time long:", err)
						return err
					}
					fld.Set(reflect.ValueOf(t))
				} else {
					t, err := time.Parse("2006-01-02 15:04", v)
					if err != nil {
						fmt.Println("error set reflect time short:", err)
						return err
					}
					fld.Set(reflect.ValueOf(t))
				}
			} else {
				fld.Set(vToSet)
			}
			return errReturn

		case []any:
			// Walk the fields
			for i := 0; i < fld.NumField(); i++ {
				if err := SetReflectFieldValue(fld.Field(i), v[i]); err != nil {
					return err
				}
			}
			return errReturn
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
		return errReturn
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
		return errReturn
	case reflect.Interface:
		unwrapped := fld.Elem()
		return SetReflectFieldValue(unwrapped, value)
	case reflect.Slice:
		targetType := fld.Type()
		typeName := targetType.String()
		if strings.HasPrefix(typeName, "[") || strings.HasPrefix(typeName, "*[") {
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
					strctElement := reflect.New(fld.Type().Elem()).Elem()
					err := SetReflectFieldValue(strctElement, vToSet.Interface())
					if err != nil {
						fmt.Println("err set nested:", err)
						return err
					}
					array = reflect.Append(fld, strctElement)
				}
			}
			fld.Set(array)
		}
		return errReturn
	}
	if errPanic != nil {
		return errPanic
	}
	if errReturn != nil {
		return errReturn
	}
	return nil
}
