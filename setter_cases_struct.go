package kstrct

import (
	"database/sql"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func InitSetterStruct() {
	// Register Time handler
	NewSetterCase(func(fld reflect.Value, value reflect.Value, valueI any) error {
		if Debug {
			fmt.Printf("\n=== TIME HANDLER DEBUG ===\n")
			fmt.Printf("Field type: %v\n", fld.Type())
			fmt.Printf("Value type: %T\n", valueI)
			fmt.Printf("Value: %+v\n", valueI)
		}
		// Only handle time.Time fields
		if fld.Type() != reflect.TypeOf(time.Time{}) {
			return fmt.Errorf("not a time.Time field")
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
				fmt.Printf("Found SQL type\n")
			}
			if err := handleSqlNull(fld, valueI); err != nil {
				if Debug {
					fmt.Printf("SQL handler returned error: %v\n", err)
				}
				return err
			}
			return nil
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
				if Debug {
					fmt.Printf("Setting time value: %v\n", v)
				}
				fld.Set(reflect.ValueOf(v))
				return nil
			case *time.Time:
				if v != nil {
					if Debug {
						fmt.Printf("Setting pointer time value: %v\n", *v)
					}
					fld.Set(reflect.ValueOf((*v)))
				}
				return nil
			case string:
				if t, err := parseTimeString(v); err == nil {
					if Debug {
						fmt.Printf("Setting parsed time value: %v\n", t)
					}
					fld.Set(reflect.ValueOf(t))
					return nil
				} else {
					return fmt.Errorf("cannot parse time string %q: %v", v, err)
				}
			case int:
				if Debug {
					fmt.Printf("Setting int value as time: %d\n", v)
				}
				fld.Set(reflect.ValueOf(time.Unix(int64(v), 0)))
				return nil
			case int64:
				if Debug {
					fmt.Printf("Setting int64 value as time: %d\n", v)
				}
				fld.Set(reflect.ValueOf(time.Unix(v, 0)))
				return nil
			case int32:
				if Debug {
					fmt.Printf("Setting int32 value as time: %d\n", v)
				}
				fld.Set(reflect.ValueOf(time.Unix(int64(v), 0)))
				return nil
			case uint:
				if Debug {
					fmt.Printf("Setting uint value as time: %d\n", v)
				}
				fld.Set(reflect.ValueOf(time.Unix(int64(v), 0)))
				return nil
			case uint64:
				if Debug {
					fmt.Printf("Setting uint64 value as time: %d\n", v)
				}
				fld.Set(reflect.ValueOf(time.Unix(int64(v), 0)))
				return nil
			case uint32:
				if Debug {
					fmt.Printf("Setting uint32 value as time: %d\n", v)
				}
				fld.Set(reflect.ValueOf(time.Unix(int64(v), 0)))
				return nil
			}
			return fmt.Errorf("cannot convert %T to time.Time", valueI)
		}

		switch vTyped := valueI.(type) {
		case map[string]any:
			if Debug {
				fmt.Printf("Handling map[string]any: %+v\n", vTyped)
			}
			for key, val := range vTyped {
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

		case []KV:
			if Debug {
				fmt.Printf("Handling []KV: %+v\n", vTyped)
			}
			// Check if we're setting a slice field with an index
			if fld.Kind() == reflect.Slice {
				// Try to parse the key as an index
				if index, err := strconv.Atoi(vTyped[0].Key); err == nil {
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
					// Set the values
					for _, val := range vTyped {
						if err := SetRFValue(elem, val); err != nil {
							return fmt.Errorf("error setting slice element at index %d: %v", index, err)
						}
					}
					return nil
				}
			}

			// If not a slice or no index, handle normally
			for _, val := range vTyped {
				// Convert snake_case to TitleCase
				fieldName := SnakeCaseToTitle(val.Key)
				if Debug {
					fmt.Printf("Converting field name from %s to %s\n", val.Key, fieldName)
				}

				field := fld.FieldByName(fieldName)
				if !field.IsValid() {
					return fmt.Errorf("field %s not found", val.Key)
				}

				// If the value is a map and the field is a struct, handle it recursively
				if nestedMap, ok := val.Value.(map[string]any); ok && field.Kind() == reflect.Struct {
					if err := SetRFValue(field, nestedMap); err != nil {
						return fmt.Errorf("error setting field %s %v: %v", val.Key, val.Value, err)
					}
					continue
				}

				// Otherwise set the field directly
				if err := SetRFValue(field, val.Value); err != nil {
					return fmt.Errorf("error setting field %s: %v", val.Key, err)
				}
			}
			return nil

		case KV:
			if Debug {
				fmt.Printf("Handling KV: %+v\n", vTyped)
			}
			parts := strings.Split(vTyped.Key, ".")
			if len(parts) > 1 {
				if Debug {
					fmt.Printf("Found nested key: %v\n", parts)
				}
				// Try to parse the first part as an index if the field is a slice
				if fld.Kind() == reflect.Slice {
					index, err := strconv.Atoi(parts[0])
					if err == nil {
						if Debug {
							fmt.Printf("Found slice index: %d\n", index)
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
							Key:   strings.Join(parts[1:], "."),
							Value: vTyped.Value,
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

				// If it's a slice and we're not accessing by index, we need to create a new element
				if nestedField.Kind() == reflect.Slice {
					if Debug {
						fmt.Printf("Found slice field: %v\n", nestedField.Type())
					}
					// Create a new element
					elemType := nestedField.Type().Elem()
					var elem reflect.Value

					// If the slice is empty, create a new element
					if nestedField.Len() == 0 {
						elem = reflect.New(elemType).Elem()
						nestedField.Set(reflect.Append(nestedField, elem))
						elem = nestedField.Index(0)
					} else {
						// Use the last element
						elem = nestedField.Index(nestedField.Len() - 1)
					}

					// Create a new KV with just the field name
					nestedKV := KV{
						Key:   parts[1],
						Value: vTyped.Value,
					}

					// Set the field on the element
					if err := SetRFValue(elem, nestedKV); err != nil {
						return fmt.Errorf("error setting field %s: %v", parts[1], err)
					}
					return nil
				}

				// Create a new KV with the remaining path
				nestedKV := KV{
					Key:   strings.Join(parts[1:], "."),
					Value: vTyped.Value,
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
					valueValue := reflect.ValueOf(vTyped.Value)
					if valueValue.Type().ConvertibleTo(nestedField.Type().Elem()) {
						nestedField.SetMapIndex(keyValue, valueValue.Convert(nestedField.Type().Elem()))
						return nil
					}
					return fmt.Errorf("struct_cases1: cannot convert value of type %T to map value type %v", vTyped.Value, nestedField.Type().Elem())
				default:
					return fmt.Errorf("cannot handle nested fields for type %v", nestedField.Kind())
				}
			}

			// If no dots in key, handle based on field type
			switch fld.Kind() {
			case reflect.Map:
				if Debug {
					fmt.Printf("Handling map field\n")
				}
				// For maps, use the key directly
				if fld.IsNil() {
					fld.Set(reflect.MakeMap(fld.Type()))
				}

				// Get map's value type for type conversion
				mapType := fld.Type()
				keyType := mapType.Key()
				valueType := mapType.Elem()

				// Handle string value
				if strValue, ok := vTyped.Value.(string); ok {
					if Debug {
						fmt.Printf("Handling string value: %s\n", strValue)
					}
					// Create key value based on map's key type
					var keyValue reflect.Value
					key := vTyped.Key

					switch keyType.Kind() {
					case reflect.String:
						keyValue = reflect.ValueOf(key)
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						if i, err := strconv.ParseInt(key, 10, 64); err == nil {
							keyValue = reflect.ValueOf(i).Convert(keyType)
						}
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						if i, err := strconv.ParseUint(key, 10, 64); err == nil {
							keyValue = reflect.ValueOf(i).Convert(keyType)
						}
					case reflect.Float32, reflect.Float64:
						if f, err := strconv.ParseFloat(key, 64); err == nil {
							keyValue = reflect.ValueOf(f).Convert(keyType)
						}
					case reflect.Bool:
						if b, err := strconv.ParseBool(key); err == nil {
							keyValue = reflect.ValueOf(b)
						}
					case reflect.Struct:
						// Special handling for time.Time as key
						if keyType == reflect.TypeOf(time.Time{}) {
							if t, err := parseTimeString(key); err == nil {
								keyValue = reflect.ValueOf(t)
							} else if i, err := strconv.ParseInt(key, 10, 64); err == nil {
								keyValue = reflect.ValueOf(time.Unix(i, 0))
							}
						} else {
							newKey := reflect.New(keyType).Elem()
							if err := SetRFValue(newKey, key); err == nil {
								keyValue = newKey
							}
						}
					default:
						newKey := reflect.New(keyType).Elem()
						if err := SetRFValue(newKey, key); err == nil {
							keyValue = newKey
						}
					}

					if !keyValue.IsValid() {
						return fmt.Errorf("invalid key conversion for type %v", keyType)
					}

					// Handle different value types based on map's value type
					var finalValue reflect.Value
					switch valueType.Kind() {
					case reflect.String:
						finalValue = reflect.ValueOf(strValue)

					case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
						if i, err := strconv.ParseUint(strValue, 10, 64); err == nil {
							switch valueType.Kind() {
							case reflect.Uint8:
								if i <= math.MaxUint8 {
									finalValue = reflect.ValueOf(uint8(i))
								}
							case reflect.Uint16:
								if i <= math.MaxUint16 {
									finalValue = reflect.ValueOf(uint16(i))
								}
							case reflect.Uint32:
								if i <= math.MaxUint32 {
									finalValue = reflect.ValueOf(uint32(i))
								}
							default: // Uint64 or Uint
								finalValue = reflect.ValueOf(i)
							}
						}

					case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
						if i, err := strconv.ParseInt(strValue, 10, 64); err == nil {
							switch valueType.Kind() {
							case reflect.Int8:
								if i >= math.MinInt8 && i <= math.MaxInt8 {
									finalValue = reflect.ValueOf(int8(i))
								}
							case reflect.Int16:
								if i >= math.MinInt16 && i <= math.MaxInt16 {
									finalValue = reflect.ValueOf(int16(i))
								}
							case reflect.Int32:
								if i >= math.MinInt32 && i <= math.MaxInt32 {
									finalValue = reflect.ValueOf(int32(i))
								}
							default: // Int64 or Int
								finalValue = reflect.ValueOf(i)
							}
						}

					case reflect.Float32, reflect.Float64:
						if f, err := strconv.ParseFloat(strValue, 64); err == nil {
							if valueType.Kind() == reflect.Float32 {
								finalValue = reflect.ValueOf(float32(f))
							} else {
								finalValue = reflect.ValueOf(f)
							}
						}

					case reflect.Bool:
						if b, err := strconv.ParseBool(strValue); err == nil {
							finalValue = reflect.ValueOf(b)
						}

					case reflect.Struct:
						// Special handling for time.Time
						if valueType == reflect.TypeOf(time.Time{}) {
							if t, err := parseTimeString(strValue); err == nil {
								finalValue = reflect.ValueOf(t)
							} else if i, err := strconv.ParseInt(strValue, 10, 64); err == nil {
								finalValue = reflect.ValueOf(time.Unix(i, 0))
							}
						}

					case reflect.Slice:
						sliceValues := strings.Split(strValue, ",")
						sliceType := valueType
						newSlice := reflect.MakeSlice(sliceType, 0, len(sliceValues))

						for _, sv := range sliceValues {
							sv = strings.TrimSpace(sv)
							elemValue := reflect.New(sliceType.Elem()).Elem()
							if err := SetRFValue(elemValue, sv); err == nil {
								newSlice = reflect.Append(newSlice, elemValue)
							}
						}
						finalValue = newSlice

					case reflect.Interface:
						if i, err := strconv.ParseInt(strValue, 10, 64); err == nil {
							finalValue = reflect.ValueOf(i)
						} else if f, err := strconv.ParseFloat(strValue, 64); err == nil {
							finalValue = reflect.ValueOf(f)
						} else if b, err := strconv.ParseBool(strValue); err == nil {
							finalValue = reflect.ValueOf(b)
						} else if t, err := parseTimeString(strValue); err == nil {
							finalValue = reflect.ValueOf(t)
						} else {
							finalValue = reflect.ValueOf(strValue)
						}

					default:
						newVal := reflect.New(valueType).Elem()
						if err := SetRFValue(newVal, strValue); err == nil {
							finalValue = newVal
						}
					}

					if finalValue.IsValid() && finalValue.Type().ConvertibleTo(valueType) {
						fld.SetMapIndex(keyValue, finalValue.Convert(valueType))
						return nil
					}
					return fmt.Errorf("struct_cases3: cannot convert string value %q to map value type %v", strValue, valueType)
				}

				// Original handling for non-string values
				keyValue := reflect.ValueOf(vTyped.Key)
				valueValue := reflect.ValueOf(vTyped.Value)
				if valueValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue, valueValue.Convert(fld.Type().Elem()))
					return nil
				}
				return fmt.Errorf("struct_cases3: cannot convert value of type %T to map value type %v", vTyped.Value, fld.Type().Elem())

			case reflect.Slice:
				if Debug {
					fmt.Printf("Handling slice field\n")
				}
				// For slices, try to parse the key as an index
				index, err := strconv.Atoi(vTyped.Key)
				if err == nil {
					if Debug {
						fmt.Printf("Found slice index: %d\n", index)
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

					// Set the value
					return SetRFValue(elem, vTyped.Value)
				}

				// If not an index, handle string values as comma-separated lists
				if strValue, ok := vTyped.Value.(string); ok {
					if Debug {
						fmt.Printf("Handling comma-separated list: %s\n", strValue)
					}
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
				valueValue := reflect.ValueOf(vTyped.Value)
				if valueValue.Type().ConvertibleTo(fld.Type()) {
					fld.Set(valueValue.Convert(fld.Type()))
					return nil
				}
				return fmt.Errorf("struct_cases4: cannot convert value of type %T to slice type %v", vTyped.Value, fld.Type())

			case reflect.Struct:
				if Debug {
					fmt.Printf("Handling struct field\n")
				}
				// For structs, use FieldByName
				field := fld.FieldByName(SnakeCaseToTitle(vTyped.Key))
				if !field.IsValid() {
					return fmt.Errorf("field %s not found", vTyped.Key)
				}
				return SetRFValue(field, vTyped.Value)

			default:
				return fmt.Errorf("cannot handle KV type for field type %v", fld.Kind())
			}
		}

		// Try to convert the value to a struct
		if reflect.TypeOf(valueI) != nil {
			if Debug {
				fmt.Printf("Converting value to struct\n")
			}
			value := reflect.ValueOf(valueI)

			// If it's a KV type, try to set the field
			if kv, ok := valueI.(KV); ok {
				field := fld.FieldByName(SnakeCaseToTitle(kv.Key))
				if !field.IsValid() {
					return fmt.Errorf("field %s not found", kv.Key)
				}
				return SetRFValue(field, kv.Value)
			}

			// If types match exactly, do direct assignment
			if value.Type() == fld.Type() {
				if Debug {
					fmt.Printf("Types match exactly, doing direct assignment\n")
				}
				fld.Set(value)
				return nil
			}

			// For other types, try to set the "Id" field if it exists
			if fld.Kind() == reflect.Struct {
				field := fld.FieldByName("Id")
				if field.IsValid() && field.CanSet() {
					return SetRFValue(field, valueI)
				}
			}
		}

		return fmt.Errorf("cannot set struct field with value of type %T", valueI)
	}, reflect.Struct)
}

func handleSqlNull(fld reflect.Value, valueI any) error {
	// Handle SQL null types first
	if strings.HasPrefix(fld.Type().String(), "sql.Null") {
		if Debug {
			fmt.Printf("\n=== STRUCT:handleSqlNull HANDLER DEBUG ===\n")
			fmt.Printf("SQL type: %s\n", fld.Type().String())
			fmt.Printf("Value type: %T\n", valueI)
			fmt.Printf("Value: %+v\n", valueI)
		}

		// If it's a KV type, use its Value field
		if kv, ok := valueI.(KV); ok {
			if Debug {
				fmt.Printf("Found KV type, Value: %+v\n", kv.Value)
			}
			valueI = kv.Value
		}

		switch fld.Type().String() {
		case "sql.NullString":
			if Debug {
				fmt.Printf("Handling NullString\n")
			}
			if str, ok := valueI.(string); ok {
				if Debug {
					fmt.Printf("Setting string value: %s\n", str)
				}
				fld.Set(reflect.ValueOf(sql.NullString{String: str, Valid: true}))
				return nil
			} else if valueI == nil {
				fld.Set(reflect.ValueOf(sql.NullString{Valid: false}))
				return nil
			}
		case "sql.NullInt64":
			if Debug {
				fmt.Printf("Handling NullInt64\n")
			}
			switch v := valueI.(type) {
			case int64:
				if Debug {
					fmt.Printf("Setting int64 value: %d\n", v)
				}
				fld.Set(reflect.ValueOf(sql.NullInt64{Int64: v, Valid: true}))
				return nil
			case int:
				if Debug {
					fmt.Printf("Setting int value: %d\n", v)
				}
				fld.Set(reflect.ValueOf(sql.NullInt64{Int64: int64(v), Valid: true}))
				return nil
			case float64:
				if Debug {
					fmt.Printf("Setting float64 value as int64: %f\n", v)
				}
				fld.Set(reflect.ValueOf(sql.NullInt64{Int64: int64(v), Valid: true}))
				return nil
			case string:
				if i, err := strconv.ParseInt(v, 10, 64); err == nil {
					if Debug {
						fmt.Printf("Setting parsed int64 value: %d\n", i)
					}
					fld.Set(reflect.ValueOf(sql.NullInt64{Int64: i, Valid: true}))
					return nil
				}
			case nil:
				fld.Set(reflect.ValueOf(sql.NullInt64{Valid: false}))
				return nil
			}
		case "sql.NullFloat64":
			if Debug {
				fmt.Printf("Handling NullFloat64\n")
			}
			switch v := valueI.(type) {
			case float64:
				if Debug {
					fmt.Printf("Setting float64 value: %f\n", v)
				}
				fld.Set(reflect.ValueOf(sql.NullFloat64{Float64: v, Valid: true}))
				return nil
			case float32:
				if Debug {
					fmt.Printf("Setting float32 value: %f\n", v)
				}
				fld.Set(reflect.ValueOf(sql.NullFloat64{Float64: float64(v), Valid: true}))
				return nil
			case int:
				if Debug {
					fmt.Printf("Setting int value as float64: %d\n", v)
				}
				fld.Set(reflect.ValueOf(sql.NullFloat64{Float64: float64(v), Valid: true}))
				return nil
			case int64:
				if Debug {
					fmt.Printf("Setting int64 value as float64: %d\n", v)
				}
				fld.Set(reflect.ValueOf(sql.NullFloat64{Float64: float64(v), Valid: true}))
				return nil
			case string:
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					if Debug {
						fmt.Printf("Setting parsed float64 value: %f\n", f)
					}
					fld.Set(reflect.ValueOf(sql.NullFloat64{Float64: f, Valid: true}))
					return nil
				}
			case nil:
				fld.Set(reflect.ValueOf(sql.NullFloat64{Valid: false}))
				return nil
			}
		case "sql.NullBool":
			if Debug {
				fmt.Printf("Handling NullBool\n")
			}
			switch v := valueI.(type) {
			case bool:
				if Debug {
					fmt.Printf("Setting bool value: %v\n", v)
				}
				fld.Set(reflect.ValueOf(sql.NullBool{Bool: v, Valid: true}))
				return nil
			case string:
				if b, err := strconv.ParseBool(v); err == nil {
					if Debug {
						fmt.Printf("Setting parsed bool value: %v\n", b)
					}
					fld.Set(reflect.ValueOf(sql.NullBool{Bool: b, Valid: true}))
					return nil
				}
			case int:
				if Debug {
					fmt.Printf("Setting int value as bool: %d\n", v)
				}
				fld.Set(reflect.ValueOf(sql.NullBool{Bool: v != 0, Valid: true}))
				return nil
			case nil:
				fld.Set(reflect.ValueOf(sql.NullBool{Valid: false}))
				return nil
			}
		case "sql.NullTime":
			if Debug {
				fmt.Printf("Handling NullTime\n")
			}
			switch v := valueI.(type) {
			case time.Time:
				if Debug {
					fmt.Printf("Setting time value: %v\n", v)
				}
				fld.Set(reflect.ValueOf(sql.NullTime{Time: v, Valid: true}))
				return nil
			case string:
				if t, err := parseTimeString(v); err == nil {
					if Debug {
						fmt.Printf("Setting parsed time value: %v\n", t)
					}
					fld.Set(reflect.ValueOf(sql.NullTime{Time: t, Valid: true}))
					return nil
				}
			case int64:
				if Debug {
					fmt.Printf("Setting int64 value as time: %d\n", v)
				}
				fld.Set(reflect.ValueOf(sql.NullTime{Time: time.Unix(v, 0), Valid: true}))
				return nil
			case nil:
				fld.Set(reflect.ValueOf(sql.NullTime{Valid: false}))
				return nil
			}
		}
		if Debug {
			fmt.Printf("Failed to handle SQL type\n")
		}
	}
	return nil
}
