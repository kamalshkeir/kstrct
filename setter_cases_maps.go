package kstrct

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func InitSetterMaps() {
	// Register Map handler
	NewSetterCase(func(fld reflect.Value, value reflect.Value, valueI any) error {
		if Debug {
			fmt.Printf("\n=== MAP HANDLER DEBUG ===\n")
			fmt.Printf("Field type: %v\n", fld.Type())
			fmt.Printf("Value type: %T\n", valueI)
			fmt.Printf("Value: %+v\n", valueI)
		}

		// Initialize map if nil
		if fld.IsNil() {
			if Debug {
				fmt.Printf("Initializing nil map of type: %v\n", fld.Type())
			}
			fld.Set(reflect.MakeMap(fld.Type()))
		}

		// Try direct assignment if types match
		if value.Type() == fld.Type() {
			if Debug {
				fmt.Printf("Direct type match, doing direct assignment\n")
			}
			fld.Set(value)
			return nil
		}

		// Try simple type conversion if possible
		if value.Type().ConvertibleTo(fld.Type()) {
			if Debug {
				fmt.Printf("Types convertible, doing conversion\n")
			}
			fld.Set(value.Convert(fld.Type()))
			return nil
		}

		if Debug {
			fmt.Printf("Handling value type: %T\n", valueI)
		}

		switch v := valueI.(type) {
		case string:
			// Get map's value type for type conversion
			mapType := fld.Type()
			keyType := mapType.Key()
			valueType := mapType.Elem()

			if Debug {
				fmt.Printf("field type %T %v\n", fld.Interface(), fld.Kind() == reflect.Pointer)
			}

			// Parse string into key-value pairs
			pairs := strings.Split(v, ",")
			for _, pair := range pairs {
				pair = strings.TrimSpace(pair)
				if pair == "" {
					continue
				}

				// Split into key and value
				kv := strings.SplitN(pair, ":", 2)
				if len(kv) != 2 {
					continue
				}

				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])

				// Remove surrounding quotes from value if present
				value = strings.TrimSpace(value) // First trim spaces
				if len(value) >= 2 {
					if (value[0] == '"' && value[len(value)-1] == '"') ||
						(value[0] == '\'' && value[len(value)-1] == '\'') {
						value = value[1 : len(value)-1]
					}
				}

				// Create key value based on map's key type
				var keyValue reflect.Value
				switch keyType.Kind() {
				case reflect.String:
					keyValue = reflect.ValueOf(key)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if i, err := strconv.ParseInt(key, 10, 64); err == nil {
						keyValue = reflect.ValueOf(i).Convert(keyType)
					} else {
						continue
					}
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if i, err := strconv.ParseUint(key, 10, 64); err == nil {
						keyValue = reflect.ValueOf(i).Convert(keyType)
					} else {
						continue
					}
				case reflect.Float32, reflect.Float64:
					if f, err := strconv.ParseFloat(key, 64); err == nil {
						keyValue = reflect.ValueOf(f).Convert(keyType)
					} else {
						continue
					}
				case reflect.Bool:
					if b, err := strconv.ParseBool(key); err == nil {
						keyValue = reflect.ValueOf(b)
					} else {
						continue
					}
				case reflect.Struct:
					// Special handling for time.Time as key
					if keyType == reflect.TypeOf(time.Time{}) {
						if t, err := parseTimeString(key); err == nil {
							keyValue = reflect.ValueOf(t)
						} else if i, err := strconv.ParseInt(key, 10, 64); err == nil {
							// Try parsing as Unix timestamp
							keyValue = reflect.ValueOf(time.Unix(i, 0))
						} else {
							continue
						}
					} else {
						// Try to convert the key string to the required type
						newKey := reflect.New(keyType).Elem()
						if err := SetRFValue(newKey, key); err == nil {
							keyValue = newKey
						} else {
							continue
						}
					}
				default:
					// Try to convert the key string to the required type
					newKey := reflect.New(keyType).Elem()
					if err := SetRFValue(newKey, key); err == nil {
						keyValue = newKey
					} else {
						continue
					}
				}

				if !keyValue.IsValid() {
					continue
				}

				// Handle different value types based on map's value type
				var finalValue reflect.Value
				switch valueType.Kind() {
				case reflect.String:
					finalValue = reflect.ValueOf(value)

				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
					if i, err := strconv.ParseUint(value, 10, 64); err == nil {
						// Convert to the specific uint type
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
					if i, err := strconv.ParseInt(value, 10, 64); err == nil {
						// Convert to the specific int type
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
					if f, err := strconv.ParseFloat(value, 64); err == nil {
						if valueType.Kind() == reflect.Float32 {
							finalValue = reflect.ValueOf(float32(f))
						} else {
							finalValue = reflect.ValueOf(f)
						}
					}

				case reflect.Bool:
					if b, err := strconv.ParseBool(value); err == nil {
						finalValue = reflect.ValueOf(b)
					}

				case reflect.Struct:
					// Special handling for time.Time
					if valueType == reflect.TypeOf(time.Time{}) {
						if t, err := parseTimeString(value); err == nil {
							finalValue = reflect.ValueOf(t)
						} else if i, err := strconv.ParseInt(value, 10, 64); err == nil {
							// Try parsing as Unix timestamp
							finalValue = reflect.ValueOf(time.Unix(i, 0))
						}
					}

				case reflect.Slice:
					// Handle slice types (e.g., []string)
					sliceValues := strings.Split(value, ",")
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
					// Try to convert to appropriate type
					if i, err := strconv.ParseInt(value, 10, 64); err == nil {
						finalValue = reflect.ValueOf(i)
					} else if f, err := strconv.ParseFloat(value, 64); err == nil {
						finalValue = reflect.ValueOf(f)
					} else if b, err := strconv.ParseBool(value); err == nil {
						finalValue = reflect.ValueOf(b)
					} else if t, err := parseTimeString(value); err == nil {
						finalValue = reflect.ValueOf(t)
					} else {
						finalValue = reflect.ValueOf(value)
					}

				default:
					// Create new value of the correct type
					newVal := reflect.New(valueType).Elem()
					if err := SetRFValue(newVal, value); err == nil {
						finalValue = newVal
					} else {
						continue
					}
				}

				// Set the map value if conversion was successful
				if finalValue.IsValid() && finalValue.Type().ConvertibleTo(valueType) {
					fld.SetMapIndex(keyValue, finalValue.Convert(valueType))
				}
			}

		case map[string]any:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				} else {
					// Try to convert the value
					newVal := reflect.New(fld.Type().Elem()).Elem()
					if err := SetRFValue(newVal, val); err == nil {
						fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), newVal)
					}
				}
			}
		case map[string]string:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				}
			}
		case map[string]int:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(fmt.Sprint(val)) // Convert int to string
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				}
			}
		case map[string]float64:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				} else {
					// Try to convert the value
					newVal := reflect.New(fld.Type().Elem()).Elem()
					if err := SetRFValue(newVal, val); err == nil {
						fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), newVal)
					}
				}
			}
		case map[string]bool:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				}
			}
		case KV:
			if Debug {
				fmt.Printf("\nHandling KV case\n")
				fmt.Printf("KV Key: %s\n", v.Key)
				fmt.Printf("KV Value: %+v\n", v.Value)
			}

			// Handle both dotted and non-dotted keys
			key := v.Key
			parts := strings.Split(v.Key, ".")
			if len(parts) > 1 {
				key = parts[1] // Use the part after the dot
			}

			keyValue := reflect.New(fld.Type().Key()).Elem()
			if err := SetRFValue(keyValue, key); err == nil {
				valValue := reflect.ValueOf(v.Value)
				if valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					if Debug {
						fmt.Printf("Setting map key '%v' to value '%v'\n", keyValue.Interface(), valValue.Interface())
					}
					fld.SetMapIndex(keyValue, valValue.Convert(fld.Type().Elem()))
				} else {
					// Try to convert the value
					newVal := reflect.New(fld.Type().Elem()).Elem()
					if err := SetRFValue(newVal, v.Value); err == nil {
						if Debug {
							fmt.Printf("Setting map key '%v' to converted value '%v'\n", keyValue.Interface(), newVal.Interface())
						}
						fld.SetMapIndex(keyValue, newVal)
					}
				}
			}
		case map[int]any:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				} else {
					// Try to convert the value
					newVal := reflect.New(fld.Type().Elem()).Elem()
					if err := SetRFValue(newVal, val); err == nil {
						fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), newVal)
					}
				}
			}
		case map[int]string:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				}
			}
		case map[int]int:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				}
			}
		case map[float64]any:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				} else {
					// Try to convert the value
					newVal := reflect.New(fld.Type().Elem()).Elem()
					if err := SetRFValue(newVal, val); err == nil {
						fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), newVal)
					}
				}
			}
		case map[time.Time]any:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				} else {
					// Try to convert the value
					newVal := reflect.New(fld.Type().Elem()).Elem()
					if err := SetRFValue(newVal, val); err == nil {
						fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), newVal)
					}
				}
			}
		case map[time.Time]string:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				}
			}
		case map[time.Time]int:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				}
			}
		case map[time.Time]time.Time:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				}
			}
		case map[string]time.Time:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				}
			}
		case map[int]time.Time:
			for key, val := range v {
				keyValue := reflect.ValueOf(key)
				valValue := reflect.ValueOf(val)
				if keyValue.Type().ConvertibleTo(fld.Type().Key()) && valValue.Type().ConvertibleTo(fld.Type().Elem()) {
					fld.SetMapIndex(keyValue.Convert(fld.Type().Key()), valValue.Convert(fld.Type().Elem()))
				}
			}
		default:
			// Try to handle any map type using reflection
			if value.Kind() == reflect.Map {
				iter := value.MapRange()
				for iter.Next() {
					key := iter.Key()
					val := iter.Value()

					if key.Type().ConvertibleTo(fld.Type().Key()) && val.Type().ConvertibleTo(fld.Type().Elem()) {
						fld.SetMapIndex(key.Convert(fld.Type().Key()), val.Convert(fld.Type().Elem()))
					} else {
						// Try to convert the value
						newVal := reflect.New(fld.Type().Elem()).Elem()
						if err := SetRFValue(newVal, val.Interface()); err == nil {
							if key.Type().ConvertibleTo(fld.Type().Key()) {
								fld.SetMapIndex(key.Convert(fld.Type().Key()), newVal)
							}
						}
					}
				}
				return nil
			}
			return fmt.Errorf("cannot set map with value of type %T", valueI)
		}
		return nil
	}, reflect.Map)
}
