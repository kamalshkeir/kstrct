package kstrct

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// SetterCase represents a function that handles setting a field value for a specific reflect.Kind
type SetterCase func(fld reflect.Value, value reflect.Value, valueI any) error

var setterCases = make(map[reflect.Kind]SetterCase)

// NewSetterCase registers a new setter case for multiple reflect.Kinds
func NewSetterCase(handler SetterCase, kinds ...reflect.Kind) {
	for _, kind := range kinds {
		setterCases[kind] = handler
	}
}

// GetSetterCase retrieves the setter case for a specific reflect.Kind
func GetSetterCase(kind reflect.Kind) (SetterCase, bool) {
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
		if Debug {
			fmt.Printf("Skipping invalid value: kind=%v value=%v\n", fld.Kind(), value)
		}
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
		if Debug {
			fmt.Printf("Handling KV type: %+v\n", kv)
		}
		parts := strings.Split(kv.Key, ".")
		if len(parts) > 1 {
			// Try to parse the first part as an index if the field is a slice
			if fld.Kind() == reflect.Slice {
				index, err := strconv.Atoi(parts[0])
				if err == nil {
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
						Key:   strings.Join(parts[1:], "."),
						Value: kv.Value,
					}
					return SetRFValue(elem, nestedKV)
				}
			}

			// Get the field by name
			nestedField := fld.FieldByName(SnakeCaseToTitle(parts[0]))
			if !nestedField.IsValid() {
				return fmt.Errorf("field %s not found", parts[0])
			}

			// If it's a pointer, initialize it if nil and get the element
			if nestedField.Kind() == reflect.Ptr {
				if nestedField.IsNil() {
					nestedField.Set(reflect.New(nestedField.Type().Elem()))
				}
				nestedField = nestedField.Elem()
			}

			// Create a new KV with the remaining path
			nestedKV := KV{
				Key:   strings.Join(parts[1:], "."),
				Value: kv.Value,
			}

			// Handle different field types
			switch nestedField.Kind() {
			case reflect.Struct:
				return SetRFValue(nestedField, nestedKV)

			case reflect.Map:
				// For maps, use the first part as the key
				key := parts[0]
				if nestedField.IsNil() {
					nestedField.Set(reflect.MakeMap(nestedField.Type()))
				}
				keyValue := reflect.ValueOf(key)
				valueValue := reflect.ValueOf(kv.Value)
				if valueValue.Type().ConvertibleTo(nestedField.Type().Elem()) {
					nestedField.SetMapIndex(keyValue, valueValue.Convert(nestedField.Type().Elem()))
					return nil
				}
				return fmt.Errorf("setter_cases1: cannot convert value of type %T to map value type %v", kv.Value, nestedField.Type().Elem())

			case reflect.Slice:
				// For slices, try to parse the first part as an index
				index, err := strconv.Atoi(parts[0])
				if err != nil {
					// If it's not an index, handle it as a field name
					// For slices, handle string values as comma-separated lists
					if strValue, ok := kv.Value.(string); ok {
						// Split the string value by commas
						parts := strings.Split(strValue, ",")
						// Create a new slice with the correct length
						newSlice := reflect.MakeSlice(nestedField.Type(), len(parts), len(parts))
						// Convert and set each element
						for i, part := range parts {
							part = strings.TrimSpace(part)
							elem := newSlice.Index(i)
							if err := SetRFValue(elem, part); err != nil {
								return fmt.Errorf("error setting slice element %d: %v", i, err)
							}
						}
						nestedField.Set(newSlice)
						return nil
					}
					// For other value types, try direct assignment
					valueValue := reflect.ValueOf(kv.Value)
					if valueValue.Type().ConvertibleTo(nestedField.Type()) {
						nestedField.Set(valueValue.Convert(nestedField.Type()))
						return nil
					}
					return fmt.Errorf("setter_cases2: cannot convert value of type %T to slice type %v", kv.Value, nestedField.Type())
				}

				// Ensure slice has enough capacity
				if index >= nestedField.Len() {
					newSlice := reflect.MakeSlice(nestedField.Type(), index+1, index+1)
					if nestedField.Len() > 0 {
						reflect.Copy(newSlice, nestedField)
					}
					nestedField.Set(newSlice)
				}

				// Get the element at index
				elem := nestedField.Index(index)
				if elem.Kind() == reflect.Ptr {
					if elem.IsNil() {
						elem.Set(reflect.New(elem.Type().Elem()))
					}
					elem = elem.Elem()
				}

				// Create a new KV with the remaining path
				nestedKV := KV{
					Key:   strings.Join(parts[1:], "."),
					Value: kv.Value,
				}
				return SetRFValue(elem, nestedKV)

			default:
				return fmt.Errorf("cannot handle nested fields for type %v", nestedField.Kind())
			}
		}

		// If no dots in key, handle based on field type
		switch fld.Kind() {
		case reflect.Map:
			// For maps, use the key directly
			if fld.IsNil() {
				fld.Set(reflect.MakeMap(fld.Type()))
			}
			keyValue := reflect.ValueOf(kv.Key)
			valueValue := reflect.ValueOf(kv.Value)
			if valueValue.Type().ConvertibleTo(fld.Type().Elem()) {
				fld.SetMapIndex(keyValue, valueValue.Convert(fld.Type().Elem()))
				return nil
			}
			return fmt.Errorf("setter_cases3: cannot convert value of type %T to map value type %v", kv.Value, fld.Type().Elem())

		case reflect.Slice:
			// For slices, try to parse the key as an index
			index, err := strconv.Atoi(kv.Key)
			if err == nil {
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

				// Set the value
				return SetRFValue(elem, kv.Value)
			}

			// If not an index, handle string values as comma-separated lists
			if strValue, ok := kv.Value.(string); ok {
				// Split the string value by commas
				parts := strings.Split(strValue, ",")
				// Create a new slice with the correct length
				newSlice := reflect.MakeSlice(fld.Type(), len(parts), len(parts))
				// Convert and set each element
				for i, part := range parts {
					part = strings.TrimSpace(part)
					elem := newSlice.Index(i)
					if err := SetRFValue(elem, part); err != nil {
						return fmt.Errorf("error setting slice element %d: %v", i, err)
					}
				}
				fld.Set(newSlice)
				return nil
			}
			// For other value types, try direct assignment
			valueValue := reflect.ValueOf(kv.Value)
			if valueValue.Type().ConvertibleTo(fld.Type()) {
				fld.Set(valueValue.Convert(fld.Type()))
				return nil
			}
			return fmt.Errorf("setter_cases4: cannot convert value of type %T to slice type %v", kv.Value, fld.Type())

		case reflect.Struct:
			// For structs, use FieldByName
			field := fld.FieldByName(SnakeCaseToTitle(kv.Key))
			if !field.IsValid() {
				return fmt.Errorf("field %s not found", kv.Key)
			}
			return SetRFValue(field, kv.Value)

		default:
			return fmt.Errorf("cannot handle KV type for field type %v", fld.Kind())
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

func init() {
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
				return fmt.Errorf("cannot convert string %q to bool: %v", v, err)
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

	// Register Map handler
	NewSetterCase(func(fld reflect.Value, value reflect.Value, valueI any) error {
		// Initialize map if nil
		if fld.IsNil() {
			fld.Set(reflect.MakeMap(fld.Type()))
		}

		switch v := valueI.(type) {
		case map[string]any:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue, valValue.Convert(fld.Type().Elem()))
				} else {
					// Try to convert the value
					newVal := reflect.New(fld.Type().Elem()).Elem()
					if err := SetRFValue(newVal, val); err == nil {
						fld.SetMapIndex(keyValue, newVal)
					}
				}
			}
		case map[string]string:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue, valValue.Convert(fld.Type().Elem()))
				}
			}
		case map[string]int:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue, valValue.Convert(fld.Type().Elem()))
				}
			}
		case map[string]float64:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue, valValue.Convert(fld.Type().Elem()))
				}
			}
		case map[string]bool:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue, valValue.Convert(fld.Type().Elem()))
				}
			}
		case KV:
			parts := strings.Split(v.Key, ".")
			if len(parts) > 1 {
				key := parts[1]
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(v.Value)
				if valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue, valValue.Convert(fld.Type().Elem()))
				} else {
					// Try to convert the value
					newVal := reflect.New(fld.Type().Elem()).Elem()
					if err := SetRFValue(newVal, v.Value); err == nil {
						fld.SetMapIndex(keyValue, newVal)
					}
				}
			}
		default:
			return fmt.Errorf("cannot set map with value of type %T", valueI)
		}
		return nil
	}, reflect.Map)
}
