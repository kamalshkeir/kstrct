package kstrct

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func InitSetterStruct() {
	// Register Time handler
	NewSetterCase(func(fld reflect.Value, value reflect.Value, valueI any) error {
		// Only handle time.Time fields
		if fld.Type() != reflect.TypeOf(time.Time{}) {
			return fmt.Errorf("not a time.Time field")
		}
		if Debug {
			fmt.Printf("\n=== TIME HANDLER DEBUG ===\n")
			fmt.Printf("Field type: %v\n", fld.Type())
			fmt.Printf("Value type: %T\n", valueI)
			fmt.Printf("Value: %+v\n", valueI)
		}

		switch v := valueI.(type) {
		case time.Time:
			fld.Set(value)
			return nil
		case *time.Time:
			if v != nil {
				fld.Set(value)
			}
			return nil
		case string:
			if t, err := parseTimeString(v); err == nil {
				fld.Set(reflect.ValueOf(t))
				return nil
			}
		case []byte:
			if t, err := parseTimeString(string(v)); err == nil {
				fld.Set(reflect.ValueOf(t))
				return nil
			}
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
		case uint64:
			t := time.Unix(int64(v), 0)
			fld.Set(reflect.ValueOf(t))
			return nil
		case int32:
			t := time.Unix(int64(v), 0)
			fld.Set(reflect.ValueOf(t))
			return nil
		case uint32:
			t := time.Unix(int64(v), 0)
			fld.Set(reflect.ValueOf(t))
			return nil
		}
		return fmt.Errorf("cannot convert %T to time.Time", valueI)
	}, reflect.Struct)

	// Register Struct handler
	NewSetterCase(func(fld reflect.Value, value reflect.Value, valueI any) error {
		if Debug {
			fmt.Printf("\n=== STRUCT HANDLER DEBUG ===\n")
			fmt.Printf("Field type: %v\n", fld.Type())
			fmt.Printf("Value type: %T\n", valueI)
			fmt.Printf("Value: %+v\n", valueI)
		}

		// Handle SQL null types first
		if strings.HasPrefix(fld.Type().String(), "sql.Null") {
			if Debug {
				fmt.Printf("Handling SQL null type: %s\n", fld.Type().String())
			}
			switch fld.Type().String() {
			case "sql.NullString":
				if str, ok := valueI.(string); ok {
					fld.Set(reflect.ValueOf(sql.NullString{String: str, Valid: true}))
					return nil
				} else if valueI == nil {
					fld.Set(reflect.ValueOf(sql.NullString{Valid: false}))
					return nil
				}
			case "sql.NullInt64":
				switch v := valueI.(type) {
				case int64:
					fld.Set(reflect.ValueOf(sql.NullInt64{Int64: v, Valid: true}))
					return nil
				case int:
					fld.Set(reflect.ValueOf(sql.NullInt64{Int64: int64(v), Valid: true}))
					return nil
				case string:
					if i, err := strconv.ParseInt(v, 10, 64); err == nil {
						fld.Set(reflect.ValueOf(sql.NullInt64{Int64: i, Valid: true}))
						return nil
					}
				case nil:
					fld.Set(reflect.ValueOf(sql.NullInt64{Valid: false}))
					return nil
				}
			case "sql.NullFloat64":
				switch v := valueI.(type) {
				case float64:
					fld.Set(reflect.ValueOf(sql.NullFloat64{Float64: v, Valid: true}))
					return nil
				case float32:
					fld.Set(reflect.ValueOf(sql.NullFloat64{Float64: float64(v), Valid: true}))
					return nil
				case int:
					fld.Set(reflect.ValueOf(sql.NullFloat64{Float64: float64(v), Valid: true}))
					return nil
				case string:
					if f, err := strconv.ParseFloat(v, 64); err == nil {
						fld.Set(reflect.ValueOf(sql.NullFloat64{Float64: f, Valid: true}))
						return nil
					}
				case nil:
					fld.Set(reflect.ValueOf(sql.NullFloat64{Valid: false}))
					return nil
				}
			case "sql.NullBool":
				switch v := valueI.(type) {
				case bool:
					fld.Set(reflect.ValueOf(sql.NullBool{Bool: v, Valid: true}))
					return nil
				case string:
					if b, err := strconv.ParseBool(v); err == nil {
						fld.Set(reflect.ValueOf(sql.NullBool{Bool: b, Valid: true}))
						return nil
					}
				case nil:
					fld.Set(reflect.ValueOf(sql.NullBool{Valid: false}))
					return nil
				}
			case "sql.NullTime":
				switch v := valueI.(type) {
				case time.Time:
					fld.Set(reflect.ValueOf(sql.NullTime{Time: v, Valid: true}))
					return nil
				case string:
					if t, err := parseTimeString(v); err == nil {
						fld.Set(reflect.ValueOf(sql.NullTime{Time: t, Valid: true}))
						return nil
					}
				case nil:
					fld.Set(reflect.ValueOf(sql.NullTime{Valid: false}))
					return nil
				}
			}
			return fmt.Errorf("cannot convert %T to %s", valueI, fld.Type().String())
		}

		// Handle direct struct assignment
		if value.Type().Kind() == reflect.Struct {
			if Debug {
				fmt.Printf("Handling direct struct assignment\n")
			}
			// If types match exactly, do direct assignment
			if value.Type() == fld.Type() {
				if Debug {
					fmt.Printf("Types match exactly, doing direct assignment\n")
				}
				fld.Set(value)
				return nil
			}
		}

		// Special handling for time.Time
		if fld.Type() == reflect.TypeOf(time.Time{}) {
			if Debug {
				fmt.Printf("Handling time.Time field\n")
			}
			switch v := valueI.(type) {
			case time.Time:
				fld.Set(reflect.ValueOf(v))
				return nil
			case *time.Time:
				if v != nil {
					fld.Set(reflect.ValueOf((*v)))
				}
				return nil
			case string:
				if t, err := parseTimeString(v); err == nil {
					fld.Set(reflect.ValueOf(t))
					return nil
				} else {
					return fmt.Errorf("cannot parse time string %q: %v", v, err)
				}
			case int:
				fld.Set(reflect.ValueOf(time.Unix(int64(v), 0)))
				return nil
			case int64:
				fld.Set(reflect.ValueOf(time.Unix(v, 0)))
				return nil
			case int32:
				fld.Set(reflect.ValueOf(time.Unix(int64(v), 0)))
				return nil
			case uint:
				fld.Set(reflect.ValueOf(time.Unix(int64(v), 0)))
				return nil
			case uint64:
				fld.Set(reflect.ValueOf(time.Unix(int64(v), 0)))
				return nil
			case uint32:
				fld.Set(reflect.ValueOf(time.Unix(int64(v), 0)))
				return nil
			}
			return fmt.Errorf("cannot convert %T to time.Time", valueI)
		}

		// Handle map for struct
		if mapValue, ok := valueI.(map[string]any); ok {
			if Debug {
				fmt.Printf("Handling map for struct: %+v\n", mapValue)
			}
			for key, val := range mapValue {
				// Convert snake_case to TitleCase
				fieldName := SnakeCaseToTitle(key)
				if Debug {
					fmt.Printf("Converting field name from %s to %s\n", key, fieldName)
				}

				field := fld.FieldByName(fieldName)
				if !field.IsValid() {
					return fmt.Errorf("field %s not found", key)
				}

				// If the value is a map and the field is a struct, handle it recursively
				if nestedMap, ok := val.(map[string]any); ok && field.Kind() == reflect.Struct {
					if err := SetRFValue(field, nestedMap); err != nil {
						return fmt.Errorf("error setting field %s: %v", key, err)
					}
					continue
				}

				// Otherwise set the field directly
				if err := SetRFValue(field, val); err != nil {
					return fmt.Errorf("error setting field %s: %v", key, err)
				}
			}
			return nil
		}

		// Handle KV type for nested fields
		if kv, ok := valueI.(KV); ok {
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
					return fmt.Errorf("struct_cases1: cannot convert value of type %T to map value type %v", kv.Value, nestedField.Type().Elem())

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
						return fmt.Errorf("struct_cases2: cannot convert value of type %T to slice type %v", kv.Value, nestedField.Type())
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
				return fmt.Errorf("struct_cases3: cannot convert value of type %T to map value type %v", kv.Value, fld.Type().Elem())

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
				return fmt.Errorf("struct_cases4: cannot convert value of type %T to slice type %v", kv.Value, fld.Type())

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

		// Try to convert the value to a struct
		if reflect.TypeOf(valueI) != nil {
			if Debug {
				fmt.Printf("Converting value of type %T to struct\n", valueI)
			}
			value := reflect.ValueOf(valueI)
			if value.Type() == fld.Type() {
				fld.Set(value)
				return nil
			}
		}

		return fmt.Errorf("cannot set struct field with value of type %T", valueI)
	}, reflect.Struct)
}
