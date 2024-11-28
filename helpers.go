package kstrct

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var Debug = false

type Istring interface {
	String() string
}
type Ibyte interface {
	Byte() []byte
}
type Iuint interface {
	Uint() uint
}
type Iuint64 interface {
	Uint64() uint64
}
type Iuint32 interface {
	Uint32() uint32
}
type Iuint16 interface {
	Uint16() uint16
}
type Iuint8 interface {
	Uint8() uint8
}
type Ifloat64 interface {
	Float64() float64
}
type Ifloat32 interface {
	Float32() float32
}
type Ibool interface {
	Bool() bool
}
type Iint interface {
	Int() int
}
type Iint64 interface {
	Int64() int64
}
type Iint32 interface {
	Int32() int32
}
type Iint16 interface {
	Int16() int16
}
type Iint8 interface {
	Int8() int8
}
type Itime interface {
	Time() time.Time
}

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
	if vToSet.Kind() == fld.Kind() && fld.Kind() != reflect.Slice {
		fld.Set(vToSet)
		return nil
	}
	var errReturn error
	switch fld.Kind() {
	case reflect.Pointer:
		unwrapped := fld.Elem()
		if unwrapped.Kind() == reflect.Slice {
			return fmt.Errorf("field of type pointer to slice not handled")
		}
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
		case Ibool:
			fld.SetBool(v.Bool())
		case Istring:
			if v.String() == "1" || v.String() == "true" {
				fld.SetBool(true)
			} else if v.String() == "0" || v.String() == "false" {
				fld.SetBool(false)
			} else {
				errReturn = fmt.Errorf("invalid bool string value: %v", v)
			}
		case bool:
			fld.SetBool(v)
		case *bool:
			fld.SetBool(*v)
		case string:
			if v == "1" || v == "true" {
				fld.SetBool(true)
			} else if v == "0" || v == "false" {
				fld.SetBool(false)
			} else {
				errReturn = fmt.Errorf("invalid bool string value: %v", v)
			}
		case *string:
			if *v == "1" || *v == "true" {
				fld.SetBool(true)
			} else if *v == "0" || *v == "false" {
				fld.SetBool(false)
			} else {
				errReturn = fmt.Errorf("invalid bool string value: %v", v)
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, Iint, Iint16, Iint32, Iint64, Iint8, Iuint, Iuint64, Iuint32, Iuint16, Iuint8, Ifloat64, Ifloat32:
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
		case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64, *float32, *float64:
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
				return errReturn
			}
		}
		return errReturn
	case reflect.String:
		switch v := value.(type) {
		case Istring:
			fld.SetString(v.String())
		case Ibyte:
			fld.SetString(string(v.Byte()))
		case Itime:
			fld.SetString(v.Time().String())
		case string:
			fld.SetString(v)
		case []byte:
			fld.SetString(string(v))
		case float64, float32, int64, int32, uint, int, uint64, uint32, Iint, Iint16, Iint32, Iint64, Iint8, Iuint, Iuint64, Iuint32, Iuint16, Iuint8, Ifloat64, Ifloat32:
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
		case Iuint:
			fld.SetUint(uint64(v.Uint()))
		case uint:
			fld.SetUint(uint64(v))
		case Iuint64:
			fld.SetUint(v.Uint64())
		case uint64:
			fld.SetUint(v)
		case Iint:
			fld.SetUint(uint64(v.Int()))
		case int:
			fld.SetUint(uint64(v))
		case Iint64:
			fld.SetUint(uint64(v.Int64()))
		case int64:
			fld.SetUint(uint64(v))
		case Iuint32:
			fld.SetUint(uint64(v.Uint32()))
		case uint32:
			fld.SetUint(uint64(v))
		case Iint32:
			fld.SetUint(uint64(v.Int32()))
		case int32:
			fld.SetUint(uint64(v))
		case Istring:
			if v, err := strconv.Atoi(v.String()); err == nil {
				fld.SetUint(uint64(v))
			}
		case string:
			if v, err := strconv.Atoi(v); err == nil {
				fld.SetUint(uint64(v))
			}
		case Ibyte:
			if v, err := strconv.Atoi(string(v.Byte())); err == nil {
				fld.SetUint(uint64(v))
			}
		case []byte:
			if v, err := strconv.Atoi(string(v)); err == nil {
				fld.SetUint(uint64(v))
			}
		case Ifloat64:
			fld.SetUint(uint64(v.Float64()))
		case Ifloat32:
			fld.SetUint(uint64(v.Float32()))
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
		case Iint:
			fld.SetInt(int64(v.Int()))
		case int:
			fld.SetInt(int64(v))
		case Iint64:
			fld.SetInt(v.Int64())
		case int64:
			fld.SetInt(v)
		case Iuint:
			fld.SetInt(int64(v.Uint()))
		case uint:
			fld.SetInt(int64(v))
		case Iuint64:
			fld.SetInt(int64(v.Uint64()))
		case uint64:
			fld.SetInt(int64(v))
		case Iint32:
			fld.SetInt(int64(v.Int32()))
		case int32:
			fld.SetInt(int64(v))
		case Iint16:
			fld.SetInt(int64(v.Int16()))
		case int16:
			fld.SetInt(int64(v))
		case Iint8:
			fld.SetInt(int64(v.Int8()))
		case int8:
			fld.SetInt(int64(v))
		case Istring:
			if v, err := strconv.Atoi(v.String()); err == nil {
				fld.SetInt(int64(v))
			}
		case string:
			if v, err := strconv.Atoi(v); err == nil {
				fld.SetInt(int64(v))
			}
		case Ibyte:
			if v, err := strconv.Atoi(string(v.Byte())); err == nil {
				fld.SetInt(int64(v))
			}
		case []byte:
			if v, err := strconv.Atoi(string(v)); err == nil {
				fld.SetInt(int64(v))
			}
		case Ifloat64:
			fld.SetUint(uint64(v.Float64()))
		case float64:
			fld.SetUint(uint64(v))
		case Ifloat32:
			fld.SetUint(uint64(v.Float32()))
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
		if v, ok := fld.Addr().Interface().(sql.Scanner); ok {
			// nulltype
			if !vToSet.IsZero() {
				err := v.Scan(value)
				if err != nil && strings.Contains(err.Error(), "time.Time") {
					if vv, ok := value.(int64); ok {
						t := time.Unix(vv, 0)
						return v.Scan(t)
					}
				}
				return err
			}
			return nil
		}
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
		case Itime:
			fld.Set(reflect.ValueOf(v.Time()))
			return nil
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
		switch v := value.(type) {
		case Ifloat64:
			fld.SetFloat(v.Float64())
		case float64:
			fld.SetFloat(v)
		case Ifloat32:
			fld.SetFloat(float64(v.Float32()))
		case float32:
			fld.SetFloat(float64(v))
		case Istring:
			f64, err := strconv.ParseFloat(v.String(), 64)
			if err == nil {
				fld.SetFloat(f64)
			}
		case string:
			f64, err := strconv.ParseFloat(v, 64)
			if err == nil {
				fld.SetFloat(f64)
			}
		case Ibyte:
			f64, err := strconv.ParseFloat(string(v.Byte()), 64)
			if err == nil {
				fld.SetFloat(f64)
			}
		case []byte:
			f64, err := strconv.ParseFloat(string(v), 64)
			if err == nil {
				fld.SetFloat(f64)
			}
		default:
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
		if typeName[0] == '[' {
			array := reflect.New(targetType).Elem()

			switch vToSet.Kind() {
			case reflect.String:
				// valueToSet string comma separated
				item := array.Type().Elem()
				switch item.Kind() {
				case reflect.String:
					array.Set(reflect.ValueOf(strings.Split(value.(string), ",")))
				case reflect.Int:
					for _, v := range strings.Split(value.(string), ",") {
						if vv, err := strconv.Atoi(v); err != nil {
							return err
						} else {
							array.Set(reflect.Append(array, reflect.ValueOf(vv)))
						}
					}
				case reflect.Int64:
					for _, v := range strings.Split(value.(string), ",") {
						if vv, err := strconv.Atoi(v); err != nil {
							return err
						} else {
							array.Set(reflect.Append(array, reflect.ValueOf(int64(vv))))
						}
					}
				case reflect.Int32:
					for _, v := range strings.Split(value.(string), ",") {
						if vv, err := strconv.Atoi(v); err != nil {
							return err
						} else {
							array.Set(reflect.Append(array, reflect.ValueOf(int32(vv))))
						}
					}
				case reflect.Uint:
					for _, v := range strings.Split(value.(string), ",") {
						if vv, err := strconv.Atoi(v); err != nil {
							return err
						} else {
							array.Set(reflect.Append(array, reflect.ValueOf(uint(vv))))
						}
					}
				case reflect.Uint64:
					for _, v := range strings.Split(value.(string), ",") {
						if vv, err := strconv.Atoi(v); err != nil {
							return err
						} else {
							array.Set(reflect.Append(array, reflect.ValueOf(uint64(vv))))
						}
					}
				case reflect.Uint32:
					for _, v := range strings.Split(value.(string), ",") {
						if vv, err := strconv.Atoi(v); err != nil {
							return err
						} else {
							array.Set(reflect.Append(array, reflect.ValueOf(uint32(vv))))
						}
					}
				case reflect.Int16:
					for _, v := range strings.Split(value.(string), ",") {
						if vv, err := strconv.Atoi(v); err != nil {
							return err
						} else {
							array.Set(reflect.Append(array, reflect.ValueOf(int16(vv))))
						}
					}
				case reflect.Uint16:
					for _, v := range strings.Split(value.(string), ",") {
						if vv, err := strconv.Atoi(v); err != nil {
							return err
						} else {
							array.Set(reflect.Append(array, reflect.ValueOf(uint16(vv))))
						}
					}
				case reflect.Int8:
					for _, v := range strings.Split(value.(string), ",") {
						if vv, err := strconv.Atoi(v); err != nil {
							return err
						} else {
							array.Set(reflect.Append(array, reflect.ValueOf(int8(vv))))
						}
					}
				case reflect.Uint8:
					for _, v := range strings.Split(value.(string), ",") {
						if vv, err := strconv.Atoi(v); err != nil {
							return err
						} else {
							array.Set(reflect.Append(array, reflect.ValueOf(uint8(vv))))
						}
					}
				case reflect.Float64:
					for _, v := range strings.Split(value.(string), ",") {
						if vv, err := strconv.ParseFloat(v, 64); err != nil {
							return err
						} else {
							array.Set(reflect.Append(array, reflect.ValueOf(vv)))
						}
					}
				case reflect.Float32:
					for _, v := range strings.Split(value.(string), ",") {
						if vv, err := strconv.ParseFloat(v, 32); err != nil {
							return err
						} else {
							array.Set(reflect.Append(array, reflect.ValueOf(float32(vv))))
						}
					}
				case reflect.Bool:
					for _, v := range strings.Split(value.(string), ",") {
						if vv, err := strconv.ParseBool(v); err != nil {
							return err
						} else {
							array.Set(reflect.Append(array, reflect.ValueOf(vv)))
						}
					}

				default:
					return fmt.Errorf("unsupported slice type (comma separated): %v", item.Kind())
				}

			case reflect.Slice:
				for i := 0; i < vToSet.Len(); i++ {
					elem := reflect.New(fld.Type().Elem()).Elem() // Create a new instance of strctElement for each iteration
					err := SetReflectFieldValue(elem, vToSet.Index(i).Interface())
					if err != nil {
						fmt.Println("err set nested:", err)
						return err
					}
					array = reflect.Append(array, elem)
				}
			default:
				return fmt.Errorf("value to set is neither string or slice, got: %v", vToSet.Kind())
			}

			fld.Set(array)
		}
		return errReturn
	}
	if errPanic != nil {
		return errPanic
	}
	return nil
}
