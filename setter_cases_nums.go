package kstrct

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func InitSetterNums() {
	// Register Int handlers with pointer support
	NewSetterCase(func(fld reflect.Value, value reflect.Value, valueI any) error {
		var intVal int64
		if Debug {
			fmt.Printf("DEBUG: Setting int field. Field type: %v, Value type: %T, Value: %v\n", fld.Type(), valueI, valueI)
		}

		switch v := valueI.(type) {
		case *int:
			if v != nil {
				intVal = int64(*v)
			}
		case *int64:
			if v != nil {
				intVal = *v
			}
		case *int32:
			if v != nil {
				intVal = int64(*v)
			}
		case *int16:
			if v != nil {
				intVal = int64(*v)
			}
		case *int8:
			if v != nil {
				intVal = int64(*v)
			}
		case int:
			intVal = int64(v)
		case int64:
			intVal = v
		case int32:
			intVal = int64(v)
		case int16:
			intVal = int64(v)
		case int8:
			intVal = int64(v)
		case float64:
			intVal = int64(v)
		case string:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				intVal = i
			} else {
				return fmt.Errorf("cannot convert string %q to int: %v", v, err)
			}
		default:
			if value.Type().ConvertibleTo(fld.Type()) {
				fld.Set(value.Convert(fld.Type()))
				return nil
			}
			return fmt.Errorf("cannot convert %T to int (value: %v)", valueI, valueI)
		}
		if Debug {
			fmt.Printf("DEBUG: Converting to int64: %v\n", intVal)
		}
		fld.SetInt(intVal)
		return nil
	}, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64)

	// Register Uint handlers
	NewSetterCase(func(fld reflect.Value, value reflect.Value, valueI any) error {
		var uintVal uint64
		switch v := valueI.(type) {
		case int:
			if v < 0 {
				return fmt.Errorf("cannot convert negative value %d to uint", v)
			}
			uintVal = uint64(v)
		case int64:
			if v < 0 {
				return fmt.Errorf("cannot convert negative value %d to uint", v)
			}
			uintVal = uint64(v)
		case int32:
			if v < 0 {
				return fmt.Errorf("cannot convert negative value %d to uint", v)
			}
			uintVal = uint64(v)
		case int16:
			if v < 0 {
				return fmt.Errorf("cannot convert negative value %d to uint", v)
			}
			uintVal = uint64(v)
		case int8:
			if v < 0 {
				return fmt.Errorf("cannot convert negative value %d to uint", v)
			}
			uintVal = uint64(v)
		case uint:
			uintVal = uint64(v)
		case uint64:
			uintVal = v
		case uint32:
			uintVal = uint64(v)
		case uint16:
			uintVal = uint64(v)
		case uint8:
			uintVal = uint64(v)
		case float32:
			if v < 0 {
				return fmt.Errorf("cannot convert negative value %f to uint", v)
			}
			uintVal = uint64(v)
		case float64:
			if v < 0 {
				return fmt.Errorf("cannot convert negative value %f to uint", v)
			}
			uintVal = uint64(v)
		case string:
			if i, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64); err == nil {
				uintVal = i
			} else {
				return fmt.Errorf("cannot convert string %q to uint", v)
			}
		default:
			return fmt.Errorf("cannot convert %T to uint", valueI)
		}
		fld.SetUint(uintVal)
		return nil
	}, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64)

	// Register Float handlers with pointer support
	NewSetterCase(func(fld reflect.Value, value reflect.Value, valueI any) error {
		var floatVal float64
		if Debug {
			fmt.Printf("DEBUG: Setting float field. Field type: %v, Value type: %T, Value: %v\n", fld.Type(), valueI, valueI)
		}

		switch v := valueI.(type) {
		case *float64:
			if v != nil {
				floatVal = *v
			}
		case *float32:
			if v != nil {
				floatVal = float64(*v)
			}
		case float64:
			floatVal = v
		case float32:
			floatVal = float64(v)
		case int:
			floatVal = float64(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				floatVal = f
			} else {
				return fmt.Errorf("cannot convert string %q to float: %v", v, err)
			}
		default:
			if value.Type().ConvertibleTo(fld.Type()) {
				fld.Set(value.Convert(fld.Type()))
				return nil
			}
			return fmt.Errorf("cannot convert %T to float (value: %v)", valueI, valueI)
		}
		if Debug {
			fmt.Printf("DEBUG: Converting to float64: %v\n", floatVal)
		}
		fld.SetFloat(floatVal)
		return nil
	}, reflect.Float32, reflect.Float64)
}
