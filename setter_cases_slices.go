package kstrct

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Cache for string builders and slices
var (
	slicePool = sync.Pool{
		New: func() any {
			s := make([]int, 0, 50)
			return &s
		},
	}
	// Pool for small slices to reduce allocations
	smallSlicePool = sync.Pool{
		New: func() any {
			// Pre-allocate a slice that can hold common cases (e.g. "a,b,c")
			s := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf("")), 0, 8)
			return &s
		},
	}
)

// getSmallSlice gets a pre-allocated slice from the pool if the type matches
func getSmallSlice(typ reflect.Type, size int) (reflect.Value, bool) {
	if size <= 8 { // Only use pool for small slices
		if p := smallSlicePool.Get().(*reflect.Value); p != nil {
			if (*p).Type().Elem() == typ {
				newSlice := *p
				newSlice = newSlice.Slice(0, 0)
				if newSlice.Cap() >= size {
					return newSlice, true
				}
			}
			smallSlicePool.Put(p)
		}
	}
	return reflect.Value{}, false
}

// putSmallSlice returns a slice to the pool if it meets the criteria
func putSmallSlice(v reflect.Value) {
	if v.Cap() <= 8 && v.Len() <= 8 {
		smallSlicePool.Put(&v)
	}
}

func InitSetterSlice() {
	// Register Slice handler
	NewSetterCase(func(fld reflect.Value, value reflect.Value, valueI any) error {
		// Handle direct slice value
		if value.Kind() == reflect.Slice || reflect.TypeOf(valueI).Kind() == reflect.Slice {
			// If types match exactly, do direct assignment
			if value.Type() == fld.Type() {
				fld.Set(value)
				return nil
			}

			// Try to get a pre-allocated slice from the pool
			newSlice, ok := getSmallSlice(fld.Type().Elem(), value.Len())
			if !ok {
				newSlice = reflect.MakeSlice(fld.Type(), 0, value.Len())
			}

			elemType := fld.Type().Elem()
			for i := 0; i < value.Len(); i++ {
				sourceElem := value.Index(i).Interface()
				// Grow slice as needed
				newSlice = reflect.Append(newSlice, reflect.Zero(elemType))
				elem := newSlice.Index(i)

				// Use direct memory access for basic types
				switch elemType.Kind() {
				case reflect.String:
					elem.SetString(fmt.Sprint(sourceElem))
				case reflect.Int:
					if v, ok := sourceElem.(int); ok {
						elem.SetInt(int64(v))
					} else if v, ok := sourceElem.(string); ok {
						if i, err := strconv.Atoi(v); err == nil {
							elem.SetInt(int64(i))
						}
					}
				case reflect.Bool:
					if v, ok := sourceElem.(bool); ok {
						elem.SetBool(v)
					} else if v, ok := sourceElem.(string); ok {
						if b, err := strconv.ParseBool(v); err == nil {
							elem.SetBool(b)
						}
					}
				case reflect.Float64:
					if v, ok := sourceElem.(float64); ok {
						elem.SetFloat(v)
					} else if v, ok := sourceElem.(string); ok {
						if f, err := strconv.ParseFloat(v, 64); err == nil {
							elem.SetFloat(f)
						}
					}
				default:
					if err := SetRFValue(elem, sourceElem); err != nil {
						putSmallSlice(newSlice)
						return fmt.Errorf("error converting slice element %d: %v", i, err)
					}
				}
			}
			fld.Set(newSlice)
			return nil
		}

		// Handle comma-separated string values for slices
		if str, ok := valueI.(string); ok {
			elemType := fld.Type().Elem()

			// Special handling for slice of maps or structs
			switch elemType.Kind() {
			case reflect.Map:
				if err := handleMapSliceFromString(fld, str); err != nil {
					return err
				}
				return nil
			case reflect.Struct:
				if err := handleStructSliceFromString(fld, str); err != nil {
					return err
				}
				return nil
			}

			// Fast path for empty string
			if str == "" {
				fld.Set(reflect.MakeSlice(fld.Type(), 0, 0))
				return nil
			}

			// Count commas to pre-allocate slice with exact size
			count := 1
			for i := 0; i < len(str); i++ {
				if str[i] == ',' {
					count++
				}
			}

			// Pre-allocate result slice with exact size
			newSlice := reflect.MakeSlice(fld.Type(), count, count)

			// Parse elements without allocating substrings
			idx := 0
			start := 0
			for i := 0; i <= len(str); i++ {
				if i == len(str) || str[i] == ',' {
					// Trim spaces without allocation
					end := i
					for start < end && str[start] == ' ' {
						start++
					}
					for end > start && str[end-1] == ' ' {
						end--
					}

					if start < end {
						elem := newSlice.Index(idx)
						part := str[start:end]

						// Use direct memory access for basic types
						switch elemType.Kind() {
						case reflect.String:
							elem.SetString(part)
						case reflect.Int:
							if v, err := strconv.Atoi(part); err == nil {
								elem.SetInt(int64(v))
							}
						case reflect.Bool:
							if v, err := strconv.ParseBool(part); err == nil {
								elem.SetBool(v)
							}
						case reflect.Float64:
							if v, err := strconv.ParseFloat(part, 64); err == nil {
								elem.SetFloat(v)
							}
						default:
							if elemType == reflect.TypeOf(time.Time{}) {
								if t, err := parseTimeString(part); err == nil {
									elem.Set(reflect.ValueOf(t))
								}
							} else {
								SetRFValue(elem, part)
							}
						}
						idx++
					}
					start = i + 1
				}
			}

			// If we found fewer non-empty elements than commas, resize the slice
			if idx < count {
				newSlice = newSlice.Slice(0, idx)
			}

			fld.Set(newSlice)
			return nil
		}

		// handle MAPS,KVs Fill struct
		switch vval := valueI.(type) {
		case []map[string]any:
			newSlice := reflect.MakeSlice(fld.Type(), len(vval), len(vval))
			for i, m := range vval {
				elem := newSlice.Index(i)
				if err := SetRFValue(elem, m); err != nil {
					return fmt.Errorf("error setting slice element from map: %v", err)
				}
			}
			fld.Set(newSlice)
			return nil

		case KV:
			// Find dot position without allocating
			dotPos := -1
			for i := 0; i < len(vval.Key); i++ {
				if vval.Key[i] == '.' {
					dotPos = i
					break
				}
			}

			// Check if first part is numeric for indexed access
			if dotPos >= 0 {
				// Try to parse index without allocating
				index := 0
				for i := 0; i < dotPos; i++ {
					if vval.Key[i] >= '0' && vval.Key[i] <= '9' {
						index = index*10 + int(vval.Key[i]-'0')
					} else {
						index = -1
						break
					}
				}

				if index >= 0 {
					// Ensure slice has enough capacity
					if index >= fld.Len() {
						newSlice := reflect.MakeSlice(fld.Type(), index+1, index+1)
						reflect.Copy(newSlice, fld)
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

					// Handle nested field
					nestedKV := KV{
						Key:   vval.Key[dotPos+1:],
						Value: vval.Value,
					}
					return SetRFValue(elem, nestedKV)
				}
			}

			// If no index found or not numeric, try to append
			newSlice := reflect.MakeSlice(fld.Type(), fld.Len()+1, fld.Len()+1)
			reflect.Copy(newSlice, fld)
			elem := newSlice.Index(fld.Len())

			// Use direct memory access for basic types
			elemType := elem.Type()
			switch elemType.Kind() {
			case reflect.String:
				if str, ok := vval.Value.(string); ok {
					elem.SetString(str)
				} else {
					elem.SetString(fmt.Sprint(vval.Value))
				}
			case reflect.Int:
				if v, ok := vval.Value.(int); ok {
					elem.SetInt(int64(v))
				} else if str, ok := vval.Value.(string); ok {
					if v, err := strconv.Atoi(str); err == nil {
						elem.SetInt(int64(v))
					}
				}
			case reflect.Bool:
				if v, ok := vval.Value.(bool); ok {
					elem.SetBool(v)
				} else if str, ok := vval.Value.(string); ok {
					if v, err := strconv.ParseBool(str); err == nil {
						elem.SetBool(v)
					}
				}
			case reflect.Float64:
				if v, ok := vval.Value.(float64); ok {
					elem.SetFloat(v)
				} else if str, ok := vval.Value.(string); ok {
					if v, err := strconv.ParseFloat(str, 64); err == nil {
						elem.SetFloat(v)
					}
				}
			default:
				if err := SetRFValue(elem, vval.Value); err != nil {
					return fmt.Errorf("error setting slice element: %v", err)
				}
			}
			fld.Set(newSlice)
			return nil

		case map[string]any, map[string]string, map[string]int, map[string]uint,
			map[string]bool, map[string]time.Time, map[string]int64, map[string]uint8:
			// Pre-allocate new slice with space for one more element
			newSlice := reflect.MakeSlice(fld.Type(), fld.Len()+1, fld.Len()+1)
			reflect.Copy(newSlice, fld)
			elem := newSlice.Index(fld.Len())
			if err := SetRFValue(elem, vval); err != nil {
				return fmt.Errorf("error setting slice element from map: %v", err)
			}
			fld.Set(newSlice)
			return nil
		}

		return fmt.Errorf("cannot assign value of type %T to slice of type %s", valueI, fld.Type())
	}, reflect.Slice)
}

// Helper functions at the end of InitSetterSlice
func handleMapSliceFromString(fld reflect.Value, str string) error {
	// Get positions slice from pool
	posPtr := slicePool.Get().(*[]int)
	positions := *posPtr
	positions = positions[:0]
	defer func() {
		*posPtr = positions[:0]
		slicePool.Put(posPtr)
	}()

	// First, split the input into individual maps
	var maps []string
	current := ""
	inQuote := false
	quoteChar := byte(0)
	for i := 0; i < len(str); i++ {
		if (str[i] == '"' || str[i] == '\'') && (i == 0 || str[i-1] != '\\') {
			if !inQuote {
				inQuote = true
				quoteChar = str[i]
			} else if str[i] == quoteChar {
				inQuote = false
			}
		}
		if !inQuote && str[i] == ';' {
			if current != "" {
				maps = append(maps, current)
				current = ""
			}
		} else {
			current += string(str[i])
		}
	}
	if current != "" {
		maps = append(maps, current)
	}

	// If no explicit map separation, treat the whole string as one map
	if len(maps) == 0 {
		if strings.TrimSpace(str) == "" {
			return nil // Handle empty string case
		}
		maps = []string{str}
	}

	// Create slice with the right size
	newSlice := reflect.MakeSlice(fld.Type(), len(maps), len(maps))
	elemType := fld.Type().Elem()

	// Process each map
	for i, mapStr := range maps {
		elem := newSlice.Index(i)
		mapElem := reflect.MakeMap(elemType)

		// Split into key-value pairs
		pairs := strings.Split(mapStr, ",")
		hasValidPair := false
		for _, pair := range pairs {
			pair = strings.TrimSpace(pair)
			if pair == "" {
				continue
			}

			// Count colons to detect invalid format
			colonCount := 0
			for _, ch := range pair {
				if ch == ':' {
					colonCount++
				}
			}
			if colonCount != 1 {
				return fmt.Errorf("invalid key-value pair format (should contain exactly one colon): %s", pair)
			}

			// Split into key and value
			kv := strings.SplitN(pair, ":", 2)
			if len(kv) != 2 {
				return fmt.Errorf("invalid key-value pair format: %s", pair)
			}

			key := strings.TrimSpace(kv[0])
			if key == "" {
				return fmt.Errorf("empty key in pair: %s", pair)
			}

			value := strings.TrimSpace(kv[1])
			if value == "" {
				return fmt.Errorf("empty value in pair: %s", pair)
			}
			hasValidPair = true

			// Remove surrounding quotes if present
			if len(value) >= 2 {
				if (value[0] == '"' && value[len(value)-1] == '"') ||
					(value[0] == '\'' && value[len(value)-1] == '\'') {
					value = value[1 : len(value)-1]
				}
			}

			// Create key and value reflect.Values
			keyValue := reflect.ValueOf(key)
			if elemType.Key().Kind() != reflect.String {
				// Handle non-string key types
				switch elemType.Key().Kind() {
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if v, err := strconv.ParseUint(key, 10, 64); err == nil {
						// Check for uint8 overflow
						if elemType.Key().Kind() == reflect.Uint8 && v > 255 {
							return fmt.Errorf("uint8 overflow for key: %s", key)
						}
						keyValue = reflect.ValueOf(v).Convert(elemType.Key())
					} else {
						return fmt.Errorf("invalid uint key: %s", key)
					}
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if v, err := strconv.ParseInt(key, 10, 64); err == nil {
						keyValue = reflect.ValueOf(v).Convert(elemType.Key())
					} else {
						return fmt.Errorf("invalid int key: %s", key)
					}
				case reflect.Float32, reflect.Float64:
					if v, err := strconv.ParseFloat(key, 64); err == nil {
						keyValue = reflect.ValueOf(v).Convert(elemType.Key())
					} else {
						return fmt.Errorf("invalid float key: %s", key)
					}
				case reflect.Struct:
					if elemType.Key() == reflect.TypeOf(time.Time{}) {
						if t, err := time.Parse("2006-01-02", key); err == nil {
							keyValue = reflect.ValueOf(t)
						} else {
							return fmt.Errorf("invalid time key format: %s", key)
						}
					} else {
						return fmt.Errorf("unsupported struct key type: %v", elemType.Key())
					}
				default:
					return fmt.Errorf("unsupported key type: %v", elemType.Key())
				}
			}

			var finalValue reflect.Value

			// Handle different value types based on the string value
			if elemType.Elem().Kind() == reflect.Interface {
				// Try to infer the type for interface{}
				if strings.EqualFold(value, "true") || strings.EqualFold(value, "false") {
					if v, err := strconv.ParseBool(value); err == nil {
						finalValue = reflect.ValueOf(v)
					}
				} else if strings.Contains(value, ".") {
					// Try float first for decimal numbers
					if v, err := strconv.ParseFloat(value, 64); err == nil {
						// Always use float64 for decimal numbers in interface{}
						finalValue = reflect.ValueOf(float64(v))
					}
				} else if v, err := strconv.ParseInt(value, 10, 64); err == nil {
					// Use int for whole numbers
					finalValue = reflect.ValueOf(int(v))
				} else {
					finalValue = reflect.ValueOf(value)
				}
			} else {
				// For specific types
				switch elemType.Elem().Kind() {
				case reflect.String:
					finalValue = reflect.ValueOf(value)
				case reflect.Int:
					if v, err := strconv.ParseInt(value, 10, 64); err == nil {
						finalValue = reflect.ValueOf(int(v))
					}
				case reflect.Int64:
					if v, err := strconv.ParseInt(value, 10, 64); err == nil {
						finalValue = reflect.ValueOf(v)
					}
				case reflect.Float64:
					if v, err := strconv.ParseFloat(value, 64); err == nil {
						finalValue = reflect.ValueOf(v)
					}
				case reflect.Bool:
					if v, err := strconv.ParseBool(value); err == nil {
						finalValue = reflect.ValueOf(v)
					}
				default:
					// Try to convert string to the target type
					targetType := elemType.Elem()
					newValue := reflect.New(targetType).Elem()
					if err := SetRFValue(newValue, value); err == nil {
						finalValue = newValue
					} else {
						finalValue = reflect.ValueOf(value)
					}
				}
			}

			if !finalValue.IsValid() {
				return fmt.Errorf("invalid value for key %s: %s", key, value)
			}

			// For interface{}, we don't need to convert
			if elemType.Elem().Kind() == reflect.Interface {
				mapElem.SetMapIndex(keyValue, finalValue)
			} else if finalValue.Type().ConvertibleTo(elemType.Elem()) {
				mapElem.SetMapIndex(keyValue, finalValue.Convert(elemType.Elem()))
			} else {
				return fmt.Errorf("cannot convert value %v to type %v", finalValue.Interface(), elemType.Elem())
			}
		}

		if !hasValidPair {
			return fmt.Errorf("no valid key-value pairs found in input: %s", mapStr)
		}

		elem.Set(mapElem)
	}

	fld.Set(newSlice)
	return nil
}

func handleStructSliceFromString(fld reflect.Value, str string) error {
	// Get positions slice from pool
	posPtr := slicePool.Get().(*[]int)
	positions := *posPtr
	positions = positions[:0]
	defer func() {
		*posPtr = positions[:0]
		slicePool.Put(posPtr)
	}()

	// First, split the input into individual structs
	var structs []string
	current := ""
	inQuote := false
	quoteChar := byte(0)
	for i := 0; i < len(str); i++ {
		if (str[i] == '"' || str[i] == '\'') && (i == 0 || str[i-1] != '\\') {
			if !inQuote {
				inQuote = true
				quoteChar = str[i]
			} else if str[i] == quoteChar {
				inQuote = false
			}
		}
		if !inQuote && str[i] == ';' {
			if current != "" {
				structs = append(structs, current)
				current = ""
			}
		} else {
			current += string(str[i])
		}
	}
	if current != "" {
		structs = append(structs, current)
	}

	// If no explicit struct separation, treat the whole string as one struct
	if len(structs) == 0 {
		structs = []string{str}
	}

	// Create slice with the right size
	newSlice := reflect.MakeSlice(fld.Type(), len(structs), len(structs))

	// Process each struct
	for i, structStr := range structs {
		elem := newSlice.Index(i)
		if elem.Kind() == reflect.Ptr {
			if elem.IsNil() {
				elem.Set(reflect.New(elem.Type().Elem()))
			}
			elem = elem.Elem()
		}

		// Split into key-value pairs
		pairs := strings.Split(structStr, ",")
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

			field := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])

			// Remove surrounding quotes if present
			if len(value) >= 2 {
				if (value[0] == '"' && value[len(value)-1] == '"') ||
					(value[0] == '\'' && value[len(value)-1] == '\'') {
					value = value[1 : len(value)-1]
				}
			}

			// Set struct field
			if err := SetRFValue(elem, KV{Key: field, Value: value}); err != nil {
				return fmt.Errorf("error setting struct element %d: %v", i, err)
			}
		}
	}

	fld.Set(newSlice)
	return nil
}
