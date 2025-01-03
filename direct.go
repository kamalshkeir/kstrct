package kstrct

import (
	"reflect"
	"strconv"
	"time"
)

func (b *StructBuilder) handleDirectSet(fieldPath string, offset uintptr, value any, ftype reflect.Type) bool {
	switch ftype.Kind() {
	case reflect.String:
		switch vv := value.(type) {
		case string:
			SetStringField(b.ptr, offset, vv)
			return true
		case *string:
			if vv != nil {
				SetStringField(b.ptr, offset, *vv)
			}
			return true
		case []byte:
			SetStringField(b.ptr, offset, string(vv))
			return true
		case *[]byte:
			if vv != nil {
				SetStringField(b.ptr, offset, string(*vv))
			}
			return true
		}
	case reflect.Int:
		switch vv := value.(type) {
		case int:
			SetIntField(b.ptr, offset, vv)
			return true
		case *int:
			if vv != nil {
				SetIntField(b.ptr, offset, *vv)
			}
			return true
		case string:
			n, err := strconv.Atoi(vv)
			if err == nil {
				SetIntField(b.ptr, offset, n)
				return true
			}
		case *string:
			if vv != nil {
				n, err := strconv.Atoi(*vv)
				if err == nil {
					SetIntField(b.ptr, offset, n)
					return true
				}
			}
		case []byte:
			n, err := strconv.Atoi(string(vv))
			if err == nil {
				SetIntField(b.ptr, offset, n)
				return true
			}
		case *[]byte:
			if vv != nil {
				n, err := strconv.Atoi(string(*vv))
				if err == nil {
					SetIntField(b.ptr, offset, n)
					return true
				}
			}
		case int64:
			SetIntField(b.ptr, offset, int(vv))
			return true
		case *int64:
			if vv != nil {
				SetIntField(b.ptr, offset, int(*vv))
			}
			return true
		case uint:
			SetIntField(b.ptr, offset, int(vv))
			return true
		case *uint:
			if vv != nil {
				SetIntField(b.ptr, offset, int(*vv))
			}
			return true
		case uint8:
			SetIntField(b.ptr, offset, int(vv))
			return true
		case *uint8:
			if vv != nil {
				SetIntField(b.ptr, offset, int(*vv))
			}
			return true
		case int32:
			SetIntField(b.ptr, offset, int(vv))
			return true
		case *int32:
			if vv != nil {
				SetIntField(b.ptr, offset, int(*vv))
			}
			return true
		case float64:
			SetIntField(b.ptr, offset, int(vv))
			return true
		case *float64:
			if vv != nil {
				SetIntField(b.ptr, offset, int(*vv))
			}
			return true
		case uint64:
			SetIntField(b.ptr, offset, int(vv))
			return true
		case *uint64:
			if vv != nil {
				SetIntField(b.ptr, offset, int(*vv))
			}
			return true
		case float32:
			SetIntField(b.ptr, offset, int(vv))
			return true
		case *float32:
			if vv != nil {
				SetIntField(b.ptr, offset, int(*vv))
			}
			return true
		case int16:
			SetIntField(b.ptr, offset, int(vv))
			return true
		case *int16:
			if vv != nil {
				SetIntField(b.ptr, offset, int(*vv))
			}
			return true
		case uint32:
			SetIntField(b.ptr, offset, int(vv))
			return true
		case *uint32:
			if vv != nil {
				SetIntField(b.ptr, offset, int(*vv))
			}
			return true
		case uint16:
			SetIntField(b.ptr, offset, int(vv))
			return true
		case *uint16:
			if vv != nil {
				SetIntField(b.ptr, offset, int(*vv))
			}
			return true
		}
	case reflect.Uint:
		switch vv := value.(type) {
		case uint:
			SetUIntField(b.ptr, offset, vv)
			return true
		case *uint:
			if vv != nil {
				SetUIntField(b.ptr, offset, *vv)
			}
			return true
		case uint8:
			SetUIntField(b.ptr, offset, uint(vv))
			return true
		case *uint8:
			if vv != nil {
				SetUIntField(b.ptr, offset, uint(*vv))
			}
			return true
		case uint32:
			SetUIntField(b.ptr, offset, uint(vv))
			return true
		case *uint32:
			if vv != nil {
				SetUIntField(b.ptr, offset, uint(*vv))
			}
			return true
		case uint16:
			SetUIntField(b.ptr, offset, uint(vv))
			return true
		case *uint16:
			if vv != nil {
				SetUIntField(b.ptr, offset, uint(*vv))
			}
			return true

		case string:
			n, err := strconv.Atoi(vv)
			if err == nil {
				SetUIntField(b.ptr, offset, uint(n))
				return true
			}
		case *string:
			if vv != nil {
				n, err := strconv.Atoi(*vv)
				if err == nil {
					SetUIntField(b.ptr, offset, uint(n))
					return true
				}
			}
		case int:
			SetUIntField(b.ptr, offset, uint(vv))
			return true
		case *int:
			if vv != nil {
				SetUIntField(b.ptr, offset, uint(*vv))
			}
			return true
		case []byte:
			n, err := strconv.Atoi(string(vv))
			if err == nil {
				SetUIntField(b.ptr, offset, uint(n))
				return true
			}
		case *[]byte:
			if vv != nil {
				n, err := strconv.Atoi(string(*vv))
				if err == nil {
					SetUIntField(b.ptr, offset, uint(n))
					return true
				}
			}
		case int64:
			SetUIntField(b.ptr, offset, uint(vv))
			return true
		case *int64:
			if vv != nil {
				SetUIntField(b.ptr, offset, uint(*vv))
			}
			return true
		case int32:
			SetUIntField(b.ptr, offset, uint(vv))
			return true
		case *int32:
			if vv != nil {
				SetUIntField(b.ptr, offset, uint(*vv))
			}
			return true
		case float64:
			SetUIntField(b.ptr, offset, uint(vv))
			return true
		case *float64:
			if vv != nil {
				SetUIntField(b.ptr, offset, uint(*vv))
			}
			return true
		case uint64:
			SetUIntField(b.ptr, offset, uint(vv))
			return true
		case *uint64:
			if vv != nil {
				SetUIntField(b.ptr, offset, uint(*vv))
			}
			return true
		case float32:
			SetUIntField(b.ptr, offset, uint(vv))
			return true
		case *float32:
			if vv != nil {
				SetUIntField(b.ptr, offset, uint(*vv))
			}
			return true
		case int16:
			SetUIntField(b.ptr, offset, uint(vv))
			return true
		case *int16:
			if vv != nil {
				SetUIntField(b.ptr, offset, uint(*vv))
			}
			return true
		}
	case reflect.Bool:
		switch vv := value.(type) {
		case int:
			if vv == int(1) {
				SetBoolField(b.ptr, offset, true)
			} else if vv == int(0) {
				SetBoolField(b.ptr, offset, false)
			}
			return true
		case string:
			if vv == "1" {
				SetBoolField(b.ptr, offset, true)
				return true
			}
			if vv == "0" {
				SetBoolField(b.ptr, offset, false)
				return true
			}
			trfa, err := strconv.ParseBool(vv)
			if err == nil {
				SetBoolField(b.ptr, offset, trfa)
				return true
			}
		case []byte:
			if string(vv) == "1" {
				SetBoolField(b.ptr, offset, true)
				return true
			}
			if string(vv) == "0" {
				SetBoolField(b.ptr, offset, false)
				return true
			}
			trfa, err := strconv.ParseBool(string(vv))
			if err == nil {
				SetBoolField(b.ptr, offset, trfa)
				return true
			}
		case int64:
			if vv == int64(1) {
				SetBoolField(b.ptr, offset, true)
			} else if vv == int64(0) {
				SetBoolField(b.ptr, offset, false)
			}
			return true
		case uint:
			if vv == uint(1) {
				SetBoolField(b.ptr, offset, true)
			} else if vv == uint(0) {
				SetBoolField(b.ptr, offset, false)
			}
			return true
		case uint8:
			if vv == uint8(1) {
				SetBoolField(b.ptr, offset, true)
			} else if vv == uint8(0) {
				SetBoolField(b.ptr, offset, false)
			}
			return true
		case int32:
			if vv == int32(1) {
				SetBoolField(b.ptr, offset, true)
			} else if vv == int32(0) {
				SetBoolField(b.ptr, offset, false)
			}
			return true
		case float64:
			if vv == float64(1) {
				SetBoolField(b.ptr, offset, true)
			} else if vv == float64(0) {
				SetBoolField(b.ptr, offset, false)
			}
			return true
		case float32:
			if vv == float32(1) {
				SetBoolField(b.ptr, offset, true)
			} else if vv == float32(0) {
				SetBoolField(b.ptr, offset, false)
			}
			return true
		case int16:
			if vv == int16(1) {
				SetBoolField(b.ptr, offset, true)
			} else if vv == int16(0) {
				SetBoolField(b.ptr, offset, false)
			}
			return true
		case uint64:
			if vv == uint64(1) {
				SetBoolField(b.ptr, offset, true)
			} else if vv == uint64(0) {
				SetBoolField(b.ptr, offset, false)
			}
			return true
		case uint32:
			if vv == uint32(1) {
				SetBoolField(b.ptr, offset, true)
			} else if vv == uint32(0) {
				SetBoolField(b.ptr, offset, false)
			}
			return true
		case uint16:
			if vv == uint16(1) {
				SetBoolField(b.ptr, offset, true)
			} else if vv == uint16(0) {
				SetBoolField(b.ptr, offset, false)
			}
			return true
		}
	case reflect.Struct:
		if ftype == reflect.TypeOf(time.Time{}) {
			switch v := value.(type) {
			case time.Time:
				SetTimeField(b.ptr, offset, v)
				return true
			case *time.Time:
				if v != nil {
					SetTimeField(b.ptr, offset, *v)
				}
				return true
			case string:
				if t, err := parseTimeString(v); err == nil {
					SetTimeField(b.ptr, offset, t)
					return true
				}
			case []byte:
				if t, err := parseTimeString(string(v)); err == nil {
					SetTimeField(b.ptr, offset, t)
					return true
				}
			case int64:
				SetTimeField(b.ptr, offset, time.Unix(v, 0))
				return true
			case int:
				SetTimeField(b.ptr, offset, time.Unix(int64(v), 0))
				return true
			case uint:
				SetTimeField(b.ptr, offset, time.Unix(int64(v), 0))
				return true
			case uint64:
				SetTimeField(b.ptr, offset, time.Unix(int64(v), 0))
				return true
			case int32:
				SetTimeField(b.ptr, offset, time.Unix(int64(v), 0))
				return true
			case uint32:
				SetTimeField(b.ptr, offset, time.Unix(int64(v), 0))
				return true
			}
		}
	case reflect.Float64:
		switch vv := value.(type) {
		case float64:

			b.SetFloat64(fieldPath, vv)
		case int:
			b.SetFloat64(fieldPath, float64(vv))
			return true
		case string:
			trfa, err := strconv.Atoi(vv)
			if err == nil {
				b.SetFloat64(fieldPath, float64(trfa))
				return true
			}
		case []byte:
			trfa, err := strconv.Atoi(string(vv))
			if err == nil {
				b.SetFloat64(fieldPath, float64(trfa))
				return true
			}
		case int64:
			b.SetFloat64(fieldPath, float64(vv))
			return true
		case uint:
			b.SetFloat64(fieldPath, float64(vv))
			return true
		case uint8:
			b.SetFloat64(fieldPath, float64(vv))
			return true
		case int32:
			b.SetFloat64(fieldPath, float64(vv))
			return true
		case float32:
			b.SetFloat64(fieldPath, float64(vv))
			return true
		case int16:
			b.SetFloat64(fieldPath, float64(vv))
			return true
		case uint64:
			b.SetFloat64(fieldPath, float64(vv))
			return true
		case uint32:
			b.SetFloat64(fieldPath, float64(vv))
			return true
		case uint16:
			b.SetFloat64(fieldPath, float64(vv))
			return true
		}
	case reflect.Int64:
		switch vv := value.(type) {
		case int64:
			SetInt64Field(b.ptr, offset, int64(vv))
			return true
		case *int64:
			if vv != nil {
				SetInt64Field(b.ptr, offset, int64(*vv))
			}
			return true
		case int:
			SetInt64Field(b.ptr, offset, int64(vv))
			return true
		case *int:
			if vv != nil {
				SetInt64Field(b.ptr, offset, int64(*vv))
			}
			return true
		case int32:
			SetInt64Field(b.ptr, offset, int64(vv))
			return true
		case *int32:
			if vv != nil {
				SetInt64Field(b.ptr, offset, int64(*vv))
			}
			return true
		case uint8:
			SetInt64Field(b.ptr, offset, int64(vv))
			return true
		case *uint8:
			if vv != nil {
				SetInt64Field(b.ptr, offset, int64(*vv))
			}
			return true
		case uint:
			SetInt64Field(b.ptr, offset, int64(vv))
			return true
		case *uint:
			if vv != nil {
				SetInt64Field(b.ptr, offset, int64(*vv))
			}
			return true
		case uint32:
			SetInt64Field(b.ptr, offset, int64(vv))
			return true
		case *uint32:
			if vv != nil {
				SetInt64Field(b.ptr, offset, int64(*vv))
			}
			return true
		case uint16:
			SetInt64Field(b.ptr, offset, int64(vv))
			return true
		case *uint16:
			if vv != nil {
				SetInt64Field(b.ptr, offset, int64(*vv))
			}
			return true

		case string:
			n, err := strconv.Atoi(vv)
			if err == nil {
				SetInt64Field(b.ptr, offset, int64(n))
				return true
			}
		case *string:
			if vv != nil {
				n, err := strconv.Atoi(*vv)
				if err == nil {
					SetInt64Field(b.ptr, offset, int64(n))
					return true
				}
			}
		case []byte:
			n, err := strconv.Atoi(string(vv))
			if err == nil {
				SetInt64Field(b.ptr, offset, int64(n))
				return true
			}
		case *[]byte:
			if vv != nil {
				n, err := strconv.Atoi(string(*vv))
				if err == nil {
					SetInt64Field(b.ptr, offset, int64(n))
					return true
				}
			}
		case float64:
			SetInt64Field(b.ptr, offset, int64(vv))
			return true
		case *float64:
			if vv != nil {
				SetInt64Field(b.ptr, offset, int64(*vv))
			}
			return true
		case uint64:
			SetInt64Field(b.ptr, offset, int64(vv))
			return true
		case *uint64:
			if vv != nil {
				SetInt64Field(b.ptr, offset, int64(*vv))
			}
			return true
		case float32:
			SetInt64Field(b.ptr, offset, int64(vv))
			return true
		case *float32:
			if vv != nil {
				SetInt64Field(b.ptr, offset, int64(*vv))
			}
			return true
		case int16:
			SetInt64Field(b.ptr, offset, int64(vv))
			return true
		case *int16:
			if vv != nil {
				SetInt64Field(b.ptr, offset, int64(*vv))
			}
			return true
		}
	case reflect.Uint8:
		switch vv := value.(type) {
		case uint8:
			SetUInt8Field(b.ptr, offset, uint8(vv))
			return true
		case *uint8:
			if vv != nil {
				SetUInt8Field(b.ptr, offset, uint8(*vv))
			}
			return true
		case uint:
			SetUInt8Field(b.ptr, offset, uint8(vv))
			return true
		case *uint:
			if vv != nil {
				SetUInt8Field(b.ptr, offset, uint8(*vv))
			}
			return true
		case uint32:
			SetUInt8Field(b.ptr, offset, uint8(vv))
			return true
		case *uint32:
			if vv != nil {
				SetUInt8Field(b.ptr, offset, uint8(*vv))
			}
			return true
		case uint16:
			SetUInt8Field(b.ptr, offset, uint8(vv))
			return true
		case *uint16:
			if vv != nil {
				SetUInt8Field(b.ptr, offset, uint8(*vv))
			}
			return true

		case string:
			n, err := strconv.Atoi(vv)
			if err == nil {
				SetUInt8Field(b.ptr, offset, uint8(n))
				return true
			}
		case *string:
			if vv != nil {
				n, err := strconv.Atoi(*vv)
				if err == nil {
					SetUInt8Field(b.ptr, offset, uint8(n))
					return true
				}
			}
		case int:
			SetUInt8Field(b.ptr, offset, uint8(vv))
			return true
		case *int:
			if vv != nil {
				SetUInt8Field(b.ptr, offset, uint8(*vv))
			}
			return true
		case []byte:
			n, err := strconv.Atoi(string(vv))
			if err == nil {
				SetUInt8Field(b.ptr, offset, uint8(n))
				return true
			}
		case *[]byte:
			if vv != nil {
				n, err := strconv.Atoi(string(*vv))
				if err == nil {
					SetUInt8Field(b.ptr, offset, uint8(n))
					return true
				}
			}
		case int64:
			SetUInt8Field(b.ptr, offset, uint8(vv))
			return true
		case *int64:
			if vv != nil {
				SetUInt8Field(b.ptr, offset, uint8(*vv))
			}
			return true
		case int32:
			SetUInt8Field(b.ptr, offset, uint8(vv))
			return true
		case *int32:
			if vv != nil {
				SetUInt8Field(b.ptr, offset, uint8(*vv))
			}
			return true
		case float64:
			SetUInt8Field(b.ptr, offset, uint8(vv))
			return true
		case *float64:
			if vv != nil {
				SetUInt8Field(b.ptr, offset, uint8(*vv))
			}
			return true
		case uint64:
			SetUInt8Field(b.ptr, offset, uint8(vv))
			return true
		case *uint64:
			if vv != nil {
				SetUInt8Field(b.ptr, offset, uint8(*vv))
			}
			return true
		case float32:
			SetUInt8Field(b.ptr, offset, uint8(vv))
			return true
		case *float32:
			if vv != nil {
				SetUInt8Field(b.ptr, offset, uint8(*vv))
			}
			return true
		case int16:
			SetUInt8Field(b.ptr, offset, uint8(vv))
			return true
		case *int16:
			if vv != nil {
				SetUInt8Field(b.ptr, offset, uint8(*vv))
			}
			return true
		}
	case reflect.Uint64:
		switch vv := value.(type) {
		case uint8:
			SetUInt64Field(b.ptr, offset, uint64(vv))
			return true
		case *uint8:
			if vv != nil {
				SetUInt64Field(b.ptr, offset, uint64(*vv))
			}
			return true
		case uint:
			SetUInt64Field(b.ptr, offset, uint64(vv))
			return true
		case *uint:
			if vv != nil {
				SetUInt64Field(b.ptr, offset, uint64(*vv))
			}
			return true
		case uint32:
			SetUInt64Field(b.ptr, offset, uint64(vv))
			return true
		case *uint32:
			if vv != nil {
				SetUInt64Field(b.ptr, offset, uint64(*vv))
			}
			return true
		case uint16:
			SetUInt64Field(b.ptr, offset, uint64(vv))
			return true
		case *uint16:
			if vv != nil {
				SetUInt64Field(b.ptr, offset, uint64(*vv))
			}
			return true

		case string:
			n, err := strconv.Atoi(vv)
			if err == nil {
				SetUInt64Field(b.ptr, offset, uint64(n))
				return true
			}
		case *string:
			if vv != nil {
				n, err := strconv.Atoi(*vv)
				if err == nil {
					SetUInt64Field(b.ptr, offset, uint64(n))
					return true
				}
			}
		case int:
			SetUInt64Field(b.ptr, offset, uint64(vv))
			return true
		case *int:
			if vv != nil {
				SetUInt64Field(b.ptr, offset, uint64(*vv))
			}
			return true
		case []byte:
			n, err := strconv.Atoi(string(vv))
			if err == nil {
				SetUInt64Field(b.ptr, offset, uint64(n))
				return true
			}
		case *[]byte:
			if vv != nil {
				n, err := strconv.Atoi(string(*vv))
				if err == nil {
					SetUInt64Field(b.ptr, offset, uint64(n))
					return true
				}
			}
		case int64:
			SetUInt64Field(b.ptr, offset, uint64(vv))
			return true
		case *int64:
			if vv != nil {
				SetUInt64Field(b.ptr, offset, uint64(*vv))
			}
			return true
		case int32:
			SetUInt64Field(b.ptr, offset, uint64(vv))
			return true
		case *int32:
			if vv != nil {
				SetUInt64Field(b.ptr, offset, uint64(*vv))
			}
			return true
		case float64:
			SetUInt64Field(b.ptr, offset, uint64(vv))
			return true
		case *float64:
			if vv != nil {
				SetUInt64Field(b.ptr, offset, uint64(*vv))
			}
			return true
		case uint64:
			SetUInt64Field(b.ptr, offset, uint64(vv))
			return true
		case *uint64:
			if vv != nil {
				SetUInt64Field(b.ptr, offset, uint64(*vv))
			}
			return true
		case float32:
			SetUInt64Field(b.ptr, offset, uint64(vv))
			return true
		case *float32:
			if vv != nil {
				SetUInt64Field(b.ptr, offset, uint64(*vv))
			}
			return true
		case int16:
			SetUInt64Field(b.ptr, offset, uint64(vv))
			return true
		case *int16:
			if vv != nil {
				SetUInt64Field(b.ptr, offset, uint64(*vv))
			}
			return true
		}
	case reflect.Int32:
		switch vv := value.(type) {
		case int32:
			SetInt32Field(b.ptr, offset, int32(vv))
			return true
		case *int32:
			if vv != nil {
				SetInt32Field(b.ptr, offset, int32(*vv))
			}
			return true
		case int:
			SetInt32Field(b.ptr, offset, int32(vv))
			return true
		case *int:
			if vv != nil {
				SetInt32Field(b.ptr, offset, int32(*vv))
			}
			return true
		case int64:
			SetInt32Field(b.ptr, offset, int32(vv))
			return true
		case *int64:
			if vv != nil {
				SetInt32Field(b.ptr, offset, int32(*vv))
			}
			return true
		case uint8:
			SetInt32Field(b.ptr, offset, int32(vv))
			return true
		case *uint8:
			if vv != nil {
				SetInt32Field(b.ptr, offset, int32(*vv))
			}
			return true
		case uint:
			SetInt32Field(b.ptr, offset, int32(vv))
			return true
		case *uint:
			if vv != nil {
				SetInt32Field(b.ptr, offset, int32(*vv))
			}
			return true
		case uint32:
			SetInt32Field(b.ptr, offset, int32(vv))
			return true
		case *uint32:
			if vv != nil {
				SetInt32Field(b.ptr, offset, int32(*vv))
			}
			return true
		case uint16:
			SetInt32Field(b.ptr, offset, int32(vv))
			return true
		case *uint16:
			if vv != nil {
				SetInt32Field(b.ptr, offset, int32(*vv))
			}
			return true

		case string:
			n, err := strconv.Atoi(vv)
			if err == nil {
				SetInt32Field(b.ptr, offset, int32(n))
				return true
			}
		case *string:
			if vv != nil {
				n, err := strconv.Atoi(*vv)
				if err == nil {
					SetInt32Field(b.ptr, offset, int32(n))
					return true
				}
			}

		case []byte:
			n, err := strconv.Atoi(string(vv))
			if err == nil {
				SetInt32Field(b.ptr, offset, int32(n))
				return true
			}
		case *[]byte:
			if vv != nil {
				n, err := strconv.Atoi(string(*vv))
				if err == nil {
					SetInt32Field(b.ptr, offset, int32(n))
					return true
				}
			}
		case float64:
			SetInt32Field(b.ptr, offset, int32(vv))
			return true
		case *float64:
			if vv != nil {
				SetInt32Field(b.ptr, offset, int32(*vv))
			}
			return true
		case uint64:
			SetInt32Field(b.ptr, offset, int32(vv))
			return true
		case *uint64:
			if vv != nil {
				SetInt32Field(b.ptr, offset, int32(*vv))
			}
			return true
		case float32:
			SetInt32Field(b.ptr, offset, int32(vv))
			return true
		case *float32:
			if vv != nil {
				SetInt32Field(b.ptr, offset, int32(*vv))
			}
			return true
		case int16:
			SetInt32Field(b.ptr, offset, int32(vv))
			return true
		case *int16:
			if vv != nil {
				SetInt32Field(b.ptr, offset, int32(*vv))
			}
			return true
		}
	case reflect.Uint32:
		switch vv := value.(type) {
		case uint32:
			SetUInt32Field(b.ptr, offset, uint32(vv))
			return true
		case *uint32:
			if vv != nil {
				SetUInt32Field(b.ptr, offset, uint32(*vv))
			}
			return true
		case uint:
			SetUInt32Field(b.ptr, offset, uint32(vv))
			return true
		case *uint:
			if vv != nil {
				SetUInt32Field(b.ptr, offset, uint32(*vv))
			}
			return true
		case uint8:
			SetUInt32Field(b.ptr, offset, uint32(vv))
			return true
		case *uint8:
			if vv != nil {
				SetUInt32Field(b.ptr, offset, uint32(*vv))
			}
			return true
		case uint16:
			SetUInt32Field(b.ptr, offset, uint32(vv))
			return true
		case *uint16:
			if vv != nil {
				SetUInt32Field(b.ptr, offset, uint32(*vv))
			}
			return true
		case int32:
			SetUInt32Field(b.ptr, offset, uint32(vv))
			return true
		case *int32:
			if vv != nil {
				SetUInt32Field(b.ptr, offset, uint32(*vv))
			}
			return true
		case int:
			SetUInt32Field(b.ptr, offset, uint32(vv))
			return true
		case *int:
			if vv != nil {
				SetUInt32Field(b.ptr, offset, uint32(*vv))
			}
			return true
		case int64:
			SetUInt32Field(b.ptr, offset, uint32(vv))
			return true
		case *int64:
			if vv != nil {
				SetUInt32Field(b.ptr, offset, uint32(*vv))
			}
			return true
		case string:
			n, err := strconv.Atoi(vv)
			if err == nil {
				SetUInt32Field(b.ptr, offset, uint32(n))
				return true
			}
		case *string:
			if vv != nil {
				n, err := strconv.Atoi(*vv)
				if err == nil {
					SetUInt32Field(b.ptr, offset, uint32(n))
					return true
				}
			}

		case []byte:
			n, err := strconv.Atoi(string(vv))
			if err == nil {
				SetUInt32Field(b.ptr, offset, uint32(n))
				return true
			}
		case *[]byte:
			if vv != nil {
				n, err := strconv.Atoi(string(*vv))
				if err == nil {
					SetUInt32Field(b.ptr, offset, uint32(n))
					return true
				}
			}
		case float64:
			SetUInt32Field(b.ptr, offset, uint32(vv))
			return true
		case *float64:
			if vv != nil {
				SetUInt32Field(b.ptr, offset, uint32(*vv))
			}
			return true
		case uint64:
			SetUInt32Field(b.ptr, offset, uint32(vv))
			return true
		case *uint64:
			if vv != nil {
				SetUInt32Field(b.ptr, offset, uint32(*vv))
			}
			return true
		case float32:
			SetUInt32Field(b.ptr, offset, uint32(vv))
			return true
		case *float32:
			if vv != nil {
				SetUInt32Field(b.ptr, offset, uint32(*vv))
			}
			return true
		case int16:
			SetUInt32Field(b.ptr, offset, uint32(vv))
			return true
		case *int16:
			if vv != nil {
				SetUInt32Field(b.ptr, offset, uint32(*vv))
			}
			return true
		}
	case reflect.Float32:
		switch vv := value.(type) {
		case float32:
			b.SetFloat32(fieldPath, vv)
			return true
		case int:
			b.SetFloat32(fieldPath, float32(vv))
			return true
		case string:
			trfa, err := strconv.Atoi(vv)
			if err == nil {
				b.SetFloat32(fieldPath, float32(trfa))
				return true
			}
		case []byte:
			trfa, err := strconv.Atoi(string(vv))
			if err == nil {
				b.SetFloat32(fieldPath, float32(trfa))
				return true
			}
		case int64:
			b.SetFloat32(fieldPath, float32(vv))
			return true
		case uint:
			b.SetFloat32(fieldPath, float32(vv))
			return true
		case uint8:
			b.SetFloat32(fieldPath, float32(vv))
			return true
		case int32:
			b.SetFloat32(fieldPath, float32(vv))
			return true
		case int16:
			b.SetFloat32(fieldPath, float32(vv))
			return true
		case uint64:
			b.SetFloat32(fieldPath, float32(vv))
			return true
		case uint32:
			b.SetFloat32(fieldPath, float32(vv))
			return true
		case uint16:
			b.SetFloat32(fieldPath, float32(vv))
			return true
		}
	}
	return false
}
