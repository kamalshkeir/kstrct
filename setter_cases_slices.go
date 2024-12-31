package kstrct

import (
	"fmt"
	"reflect"
	"strconv"
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
)

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

			// Pre-allocate slice with exact capacity needed
			newSlice := reflect.MakeSlice(fld.Type(), value.Len(), value.Len())
			elemType := fld.Type().Elem()
			for i := 0; i < value.Len(); i++ {
				elem := newSlice.Index(i)
				sourceElem := value.Index(i).Interface()

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
						return fmt.Errorf("error converting slice element %d: %v", i, err)
					}
				}
			}
			fld.Set(newSlice)
			return nil
		}

		// Handle comma-separated string values for slices
		if str, ok := valueI.(string); ok {
			// Get positions slice from pool
			posPtr := slicePool.Get().(*[]int)
			positions := *posPtr
			positions = positions[:0]
			defer func() {
				*posPtr = positions[:0]
				slicePool.Put(posPtr)
			}()

			// Count non-empty elements and track positions in one pass
			count := 0
			start := 0
			for i := 0; i < len(str); i++ {
				if str[i] == ',' || i == len(str)-1 {
					end := i
					if i == len(str)-1 && str[i] != ',' {
						end = i + 1
					}
					// Check if segment is non-empty without allocating
					segStart := start
					segEnd := end
					for segStart < segEnd && str[segStart] == ' ' {
						segStart++
					}
					for segEnd > segStart && str[segEnd-1] == ' ' {
						segEnd--
					}
					if segEnd > segStart {
						positions = append(positions, segStart, segEnd)
						count++
					}
					start = i + 1
				}
			}

			// Pre-allocate result slice with exact size
			newSlice := reflect.MakeSlice(fld.Type(), count, count)
			elemType := fld.Type().Elem()
			for i := 0; i < count; i++ {
				elem := newSlice.Index(i)
				start := positions[i*2]
				end := positions[i*2+1]
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
						} else {
							return fmt.Errorf("error converting slice element from string: %v", err)
						}
					} else {
						if err := SetRFValue(elem, part); err != nil {
							return fmt.Errorf("error converting slice element from string: %v", err)
						}
					}
				}
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
