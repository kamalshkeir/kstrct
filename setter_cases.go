package kstrct

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// SetterCase represents a function that handles setting a field value for a specific reflect.Kind
type SetterCase func(fld reflect.Value, value reflect.Value, valueI any) error

var (
	setterCases = make(map[reflect.Kind]SetterCase)
	setterMu    sync.RWMutex
)

// NewSetterCase registers a new setter case for multiple reflect.Kinds
func NewSetterCase(handler SetterCase, kinds ...reflect.Kind) {
	setterMu.Lock()
	defer setterMu.Unlock()
	for _, kind := range kinds {
		setterCases[kind] = handler
	}
}

// GetSetterCase retrieves the setter case for a specific reflect.Kind
func GetSetterCase(kind reflect.Kind) (SetterCase, bool) {
	setterMu.RLock()
	defer setterMu.RUnlock()
	handler, ok := setterCases[kind]
	return handler, ok
}

// SetRFValue sets a value to a reflect.Value, handling type conversions
func SetRFValue(fld reflect.Value, value any) error {
	if !fld.IsValid() || !fld.CanSet() {
		return fmt.Errorf("field is not valid or cannot be set")
	}

	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return nil
	}
	// Handle pointer types
	if fld.Kind() == reflect.Pointer {
		if value == nil {
			fld.Set(reflect.Zero(fld.Type()))
			return nil
		}
		if fld.IsNil() {
			fld.Set(reflect.New(fld.Type().Elem()))
		}
		// If the value is already a pointer, dereference it
		if v.Kind() == reflect.Ptr {
			return SetRFValue(fld.Elem(), v.Elem().Interface())
		}
		return SetRFValue(fld.Elem(), value)
	}

	// Handle KV type for nested fields
	if kv, ok := value.(KV); ok {
		// If no dots in key, handle based on field type
		switch fld.Kind() {
		case reflect.Struct:
			// If the key contains dots, handle nested fields
			if strings.Contains(kv.Key, ".") {
				parts := strings.SplitN(kv.Key, ".", 2)
				field := fld.FieldByName(SnakeCaseToTitle(parts[0]))
				if !field.IsValid() {
					return fmt.Errorf("field %s not found", parts[0])
				}

				// If the field is a slice, validate the index before proceeding
				if field.Kind() == reflect.Slice {
					nextParts := strings.SplitN(parts[1], ".", 2)
					if _, err := strconv.Atoi(nextParts[0]); err != nil {
						return fmt.Errorf("invalid slice index: %v", nextParts[0])
					}
				}

				// Create a new KV with the remaining path
				nestedKV := KV{
					Key:   parts[1],
					Value: kv.Value,
				}

				// Call SetRFValue and return any error
				if err := SetRFValue(field, nestedKV); err != nil {
					return err
				}
				return nil
			}

			// Handle non-nested field
			field := fld.FieldByName(SnakeCaseToTitle(kv.Key))
			if !field.IsValid() {
				return fmt.Errorf("field %s not found", kv.Key)
			}

			// If the field is a map, verify the value type before setting
			if field.Kind() == reflect.Map && field.Type().Key().Kind() != reflect.String {
				if m, ok := kv.Value.(map[string]string); ok {
					// Try to convert all keys to verify they're valid
					for k := range m {
						if _, err := convertToType(k, field.Type().Key()); err != nil {
							return fmt.Errorf("invalid map key type: cannot convert '%v' to %v", k, field.Type().Key())
						}
					}
				}
			}

			return SetRFValue(field, kv.Value)

		case reflect.Slice:
			// For slices, try to parse the key as an index
			if strings.Contains(kv.Key, ".") {
				parts := strings.SplitN(kv.Key, ".", 2)
				// Handle []KV value first
				if kvs, ok := kv.Value.([]KV); ok {
					for _, val := range kvs {
						if err := SetRFValue(fld, val); err != nil {
							return fmt.Errorf("error setting slice element: %v", err)
						}
					}
					return nil
				}
				// Try to parse the first part as an index
				index, err := strconv.Atoi(parts[0])
				if err != nil {
					return fmt.Errorf("invalid slice index: %v", parts[0])
				}

				if index < 0 {
					return fmt.Errorf("negative slice index: %d", index)
				}

				// Get the element type
				elemType := fld.Type().Elem()
				if elemType.Kind() == reflect.Ptr {
					elemType = elemType.Elem()
				}

				// Validate that we can handle nested access on the element type
				switch elemType.Kind() {
				case reflect.Map, reflect.Struct:
				default:
					return fmt.Errorf("cannot access nested field on slice element of type %v", elemType)
				}

				// Ensure slice has enough capacity
				if index >= fld.Len() {
					newSlice := reflect.MakeSlice(fld.Type(), index+1, index+1)
					if fld.Len() > 0 {
						reflect.Copy(newSlice, fld)
					}
					fld.Set(newSlice)
				}

				// Get the element at index
				elem := fld.Index(index)
				if elem.Kind() == reflect.Ptr {
					if elem.IsNil() {
						elem.Set(reflect.New(elem.Type().Elem()))
					}
					elem = elem.Elem()
				}

				// Create a new KV with the remaining path
				nestedKV := KV{
					Key:   parts[1],
					Value: kv.Value,
				}

				// Call SetRFValue and return any error
				if err := SetRFValue(elem, nestedKV); err != nil {
					return err
				}
				return nil
			} else {
				// Handle non-nested index
				if index, err := strconv.Atoi(kv.Key); err == nil {
					if index < 0 {
						return fmt.Errorf("negative slice index: %d", index)
					}
					// Ensure slice has enough capacity
					if index >= fld.Len() {
						newSlice := reflect.MakeSlice(fld.Type(), index+1, index+1)
						if fld.Len() > 0 {
							reflect.Copy(newSlice, fld)
						}
						fld.Set(newSlice)
					}
					// Get element at index
					elem := fld.Index(index)
					if elem.Kind() == reflect.Ptr {
						if elem.IsNil() {
							elem.Set(reflect.New(elem.Type().Elem()))
						}
						elem = elem.Elem()
					}
					return SetRFValue(elem, kv.Value)
				}
			}
		}
	}

	// Try to use registered handler first
	if handler, ok := GetSetterCase(fld.Kind()); ok {
		return handler(fld, v, value)
	}

	// Try direct assignment if types match
	if v.Type().AssignableTo(fld.Type()) {
		fld.Set(v)
		return nil
	}

	// Try type conversion
	if v.Type().ConvertibleTo(fld.Type()) {
		fld.Set(v.Convert(fld.Type()))
		return nil
	}

	return fmt.Errorf("cannot set field of type %s with value of type %T", fld.Type(), value)
}

// Helper function to convert string to various types
func convertToType(value string, targetType reflect.Type) (reflect.Value, error) {
	switch targetType.Kind() {
	case reflect.String:
		return reflect.ValueOf(value), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v).Convert(targetType), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v).Convert(targetType), nil
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v).Convert(targetType), nil
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(v), nil
	default:
		return reflect.Value{}, fmt.Errorf("unsupported type conversion: %v", targetType)
	}
}

func init() {
	InitSetterSQL()
	// Register String handler with pointer support
	NewSetterCase(func(fld reflect.Value, value reflect.Value, valueI any) error {
		if ptr, ok := valueI.(*string); ok {
			if ptr != nil {
				fld.SetString(*ptr)
				return nil
			}
			return nil
		}
		fld.SetString(fmt.Sprint(valueI))
		return nil
	}, reflect.String)

	// Register Bool handler with pointer support
	NewSetterCase(func(fld reflect.Value, value reflect.Value, valueI any) error {
		var boolVal bool
		switch v := valueI.(type) {
		case *bool:
			if v != nil {
				boolVal = *v
			}
		case bool:
			boolVal = v
		case string:
			if b, err := strconv.ParseBool(v); err == nil {
				boolVal = b
			} else {
				if v == "1" {
					boolVal = true
				} else if v == "0" {
					boolVal = false
				} else {
					return fmt.Errorf("cannot convert string %q to bool: %v", v, err)
				}
			}
		case int, int64, uint, uint64, uint32, uint16, uint8, int32, int16, int8:
			if v == 1 {
				boolVal = true
			} else if v == 0 {
				boolVal = false
			} else {
				return fmt.Errorf("cannot convert number %d to bool", v)
			}
		case *int, *int64, *uint, *uint64, *uint32, *uint16, *uint8, *int32, *int16, *int8:
			if v == 1 {
				boolVal = true
			} else if v == 0 {
				boolVal = false
			} else {
				return fmt.Errorf("cannot convert *number %v to bool", v)
			}
		case Iint:
			if v.Int() == 1 {
				boolVal = true
			} else if v.Int() == 0 {
				boolVal = false
			} else {
				return fmt.Errorf("cannot convert Iint %d to bool", v.Int())
			}

		case Ibool:
			boolVal = v.Bool()
		case Istring:
			if b, err := strconv.ParseBool(v.String()); err == nil {
				boolVal = b
			} else {
				if v.String() == "1" {
					boolVal = true
				} else if v.String() == "0" {
					boolVal = false
				} else {
					return fmt.Errorf("cannot convert Istring %q to bool: %v", v.String(), err)
				}
			}
		default:
			if value.Type().ConvertibleTo(fld.Type()) {
				fld.Set(value.Convert(fld.Type()))
				return nil
			}
			return fmt.Errorf("cannot convert %T to bool (value: %v)", valueI, valueI)
		}
		fld.SetBool(boolVal)
		return nil
	}, reflect.Bool)

	InitSetterNums()

	InitSetterStruct()

	InitSetterSlice()

	InitSetterMaps()
}
