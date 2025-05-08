package kstrct

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

const maxFixedFields = 32 // Maximum number of fields to store in fixed array

type fieldInfoMap struct {
	isFixed    bool
	fixedSize  int
	fixedKeys  [maxFixedFields]string
	fixedInfos [maxFixedFields]fieldInfo
	mapInfo    map[string]fieldInfo
}

func newFieldInfoMap(size int) *fieldInfoMap {
	if size <= maxFixedFields/2 { // Use fixed array if we have room for fields + snake case names
		return &fieldInfoMap{isFixed: true}
	}
	return &fieldInfoMap{
		isFixed: false,
		mapInfo: make(map[string]fieldInfo, size*2),
	}
}

func (m *fieldInfoMap) set(key string, info fieldInfo) {
	if m.isFixed {
		if m.fixedSize < maxFixedFields {
			m.fixedKeys[m.fixedSize] = key
			m.fixedInfos[m.fixedSize] = info
			m.fixedSize++
		}
	} else {
		m.mapInfo[key] = info
	}
}

func (m *fieldInfoMap) get(key string) (fieldInfo, bool) {
	if m.isFixed {
		for i := 0; i < m.fixedSize; i++ {
			if m.fixedKeys[i] == key {
				return m.fixedInfos[i], true
			}
		}
		return fieldInfo{}, false
	}
	info, ok := m.mapInfo[key]
	return info, ok
}

// StructBuilder provides a fluent interface for struct manipulation
type StructBuilder struct {
	ptr        any
	fields     *fieldInfoMap
	structType reflect.Type
	err        error
}

type fieldInfo struct {
	offset uintptr
	typ    reflect.Type
}

var (
	// Global cache for field info
	fieldInfoCache sync.Map // map[reflect.Type]map[string]fieldInfo
	// Common time formats to try in order of preference
	timeFormats = []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006/01/02",
		"2006/01/02 15:04",
		"2006/01/02 15:04:05",
		"02/01/2006",
		"02/01/2006 15:04",
		"02/01/2006 15:04:05",
		"02-01-2006",
		"02-01-2006 15:04",
		"02-01-2006 15:04:05",
		"2006-01-02",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
	}
	// Pool for temporary maps used in operations
	mapPool = sync.Pool{
		New: func() any {
			m := make(map[string]any, 8)
			return &m
		},
	}
	// Pool for string parts
	partsPool = sync.Pool{
		New: func() any {
			s := make([]string, 0, 8)
			return &s
		},
	}
	// Buffer pool for string building
	bufPool = sync.Pool{
		New: func() any {
			return new(strings.Builder)
		},
	}
)

// NewBuilder creates a new StructBuilder for the given struct pointer
func NewBuilder(structPtr any) *StructBuilder {
	typ := reflect.TypeOf(structPtr)
	if typ.Kind() != reflect.Ptr {
		if typ.Kind() != reflect.Chan {
			panic("NewBuilder requires a pointer to struct")
		}
	}
	elemType := typ.Elem()

	// Try to get field info from cache
	var fields *fieldInfoMap
	if cached, ok := fieldInfoCache.Load(elemType); ok {
		fields = cached.(*fieldInfoMap)
	} else {
		// Create new field info map
		numFields := elemType.NumField()
		fields = newFieldInfoMap(numFields)

		// Store field info
		for i := 0; i < numFields; i++ {
			field := elemType.Field(i)
			info := fieldInfo{
				offset: field.Offset,
				typ:    field.Type,
			}

			// Store both original and snake case names
			snakeName := ToSnakeCase(field.Name)
			fields.set(snakeName, info)
			if snakeName != field.Name {
				fields.set(field.Name, info)
			}
		}

		// Store in cache
		fieldInfoCache.Store(elemType, fields)
	}

	return &StructBuilder{
		ptr:        structPtr,
		fields:     fields,
		structType: elemType,
	}
}

func (b *StructBuilder) Error() error {
	return b.err
}

// Set sets a field value by name
func (b *StructBuilder) Set(fieldPath string, value any) *StructBuilder {
	// Try direct access first
	if offset, ok := b.getOffset(fieldPath); ok {
		fieldType, _ := b.getFieldType(fieldPath)
		ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)

		if b.handleDirectSet(fieldPath, offset, value, fieldType) {
			return b
		}

		// Fallback to reflection for complex types
		field := reflect.NewAt(fieldType, ptr).Elem()
		if err := SetRFValue(field, value); err != nil {
			if Debug {
				fmt.Printf("Error setting field %s: %v\n", fieldPath, err)
			}
			b.err = fmt.Errorf("%w %w", b.err, err)
		}
		return b
	}

	// Handle nested fields
	if strings.Contains(fieldPath, ".") {
		// Get parts slice from pool
		partsPtr := partsPool.Get().(*[]string)
		*partsPtr = strings.SplitN(fieldPath, ".", 3)
		parts := *partsPtr
		defer partsPool.Put(partsPtr)

		rootField := parts[0]
		if offset, ok := b.getOffset(rootField); ok {
			fieldType, _ := b.getFieldType(rootField)
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			field := reflect.NewAt(fieldType, ptr).Elem()

			// Handle pointer to slice
			if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Slice {
				if field.IsNil() {
					field.Set(reflect.New(field.Type().Elem()))
				}
				sliceVal := field.Elem()
				if sliceVal.Len() == 0 {
					newSlice := reflect.MakeSlice(sliceVal.Type(), 1, 1)
					sliceVal.Set(newSlice)
				}
				firstElem := sliceVal.Index(0)
				if len(parts) > 1 {
					// Get string builder from pool
					buf := bufPool.Get().(*strings.Builder)
					buf.Reset()
					defer bufPool.Put(buf)

					// Build remaining path
					for i := 1; i < len(parts); i++ {
						if i > 1 {
							buf.WriteByte('.')
						}
						buf.WriteString(parts[i])
					}

					// Create KV with remaining path
					kv := KV{Key: buf.String(), Value: value}
					if err := SetRFValue(firstElem, kv); err != nil {
						if Debug {
							fmt.Printf("Error setting nested field %s: %v\n", fieldPath, err)
						}
						b.err = fmt.Errorf("%w %w", b.err, err)
					}
					return b
				}
			}

			// Handle other nested fields
			if len(parts) > 1 {
				// Get string builder from pool
				buf := bufPool.Get().(*strings.Builder)
				buf.Reset()
				defer bufPool.Put(buf)

				// Build remaining path
				for i := 1; i < len(parts); i++ {
					if i > 1 {
						buf.WriteByte('.')
					}
					buf.WriteString(parts[i])
				}

				// Create KV with remaining path
				kv := KV{Key: buf.String(), Value: value}
				if err := SetRFValue(field, kv); err != nil {
					if Debug {
						fmt.Printf("Error setting nested field %s: %v\n", fieldPath, err)
					}
					b.err = fmt.Errorf("%w %w", b.err, err)
				}
				return b
			}
		}
	}

	return b
}

// Get gets a field value by name with zero allocations
func (b *StructBuilder) Get(fieldPath string) any {
	// Handle nested fields first
	if strings.Contains(fieldPath, ".") {
		return b.getComplexField(fieldPath)
	}

	// Try direct access with original name
	if offset, ok := b.getOffset(fieldPath); ok {
		fieldType, _ := b.getFieldType(fieldPath)
		ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
		return b.getFieldValue(ptr, fieldType)
	}

	// Try with snake case
	snakePath := ToSnakeCase(fieldPath)
	if offset, ok := b.getOffset(snakePath); ok {
		fieldType, _ := b.getFieldType(snakePath)
		ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
		return b.getFieldValue(ptr, fieldType)
	}

	return nil
}

// getFieldValue gets a field value based on its type
func (b *StructBuilder) getFieldValue(ptr unsafe.Pointer, fieldType reflect.Type) any {
	// If it's a pointer, handle pointer types
	if fieldType.Kind() == reflect.Ptr {
		ptrVal := *(*unsafe.Pointer)(ptr)
		if ptrVal == nil {
			return nil
		}
		// Create a new pointer of the correct type
		val := reflect.NewAt(fieldType.Elem(), ptrVal)
		return val.Interface()
	}

	// Now handle the actual type
	switch fieldType.Kind() {
	case reflect.String:
		return *(*string)(ptr)
	case reflect.Int:
		return *(*int)(ptr)
	case reflect.Int64:
		return *(*int64)(ptr)
	case reflect.Int32:
		return *(*int32)(ptr)
	case reflect.Int16:
		return *(*int16)(ptr)
	case reflect.Int8:
		return *(*int8)(ptr)
	case reflect.Uint:
		return *(*uint)(ptr)
	case reflect.Uint64:
		return *(*uint64)(ptr)
	case reflect.Uint32:
		return *(*uint32)(ptr)
	case reflect.Uint16:
		return *(*uint16)(ptr)
	case reflect.Uint8:
		return *(*uint8)(ptr)
	case reflect.Bool:
		return *(*bool)(ptr)
	case reflect.Float64:
		return *(*float64)(ptr)
	case reflect.Float32:
		return *(*float32)(ptr)
	case reflect.Map:
		mapVal := reflect.NewAt(fieldType, ptr).Elem()
		if mapVal.IsNil() {
			return nil
		}
		return mapVal.Interface()
	case reflect.Slice:
		sliceVal := reflect.NewAt(fieldType, ptr).Elem()
		if sliceVal.IsNil() {
			return nil
		}
		// For slices, just return the interface directly to preserve the type
		return sliceVal.Interface()
	case reflect.Struct:
		if fieldType == reflect.TypeOf(time.Time{}) {
			t := *(*time.Time)(ptr)
			return t.Format("2006-01-02 15:04:05")
		}
		structVal := reflect.NewAt(fieldType, ptr).Elem()
		return structVal.Interface()
	}
	return nil
}

// getComplexField handles complex field types and nested fields
func (b *StructBuilder) getComplexField(fieldPath string) any {
	// Split the path and convert each part to snake case
	parts := strings.Split(fieldPath, ".")
	for i := range parts {
		parts[i] = ToSnakeCase(parts[i])
	}

	// Get root field
	rootField := parts[0]
	if offset, ok := b.getOffset(rootField); ok {
		fieldType, ok := b.getFieldType(rootField)
		if !ok {
			return nil
		}
		ptr := GetFieldPointer(b.ptr, offset)

		// Handle pointer to struct
		if fieldType.Kind() == reflect.Ptr {
			ptrVal := *(*unsafe.Pointer)(ptr)
			if ptrVal == nil {
				return nil
			}
			return b.getFieldValue(ptrVal, fieldType.Elem())
		}

		// Handle different nested types
		return b.getFieldValue(ptr, fieldType)
	}

	return nil
}

func FindFirstNumber(s string) (hasNumber bool, position int) {
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			return true, i
		}
	}
	return false, -1
}

// SplitByDot splits a string by dots with zero allocations using a callback
func SplitByDot(s string, callback func(part string) bool) {
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			if !callback(s[start:i]) {
				return
			}
			start = i + 1
		}
	}
	if start < len(s) {
		callback(s[start:])
	}
}

// FromKV fills the struct from variadic KV structs using direct pointer manipulation
func (b *StructBuilder) FromKV(kvs ...KV) *StructBuilder {
	// First, collect all fields with the same prefix
	prefixMap := make(map[string]map[string]any)

	for _, vv := range kvs {
		k := vv.Key
		v := vv.Value

		// Try direct field first
		if _, ok := b.getOffset(k); ok {
			// Special handling for []KV type
			if kvs, ok := v.([]KV); ok {
				// Convert KVs to map[string]any
				kvMap := make(map[string]any)
				for _, kv := range kvs {
					kvMap[kv.Key] = kv.Value
				}
				b.Set(k, kvMap)
				continue
			}
			b.Set(k, v)
			continue
		}

		// Handle nested fields
		if strings.Contains(k, ".") {
			parts := strings.Split(k, ".")
			prefix := parts[0]

			// Check if it's a numeric index
			if len(parts) > 1 {
				if _, err := strconv.Atoi(parts[1]); err == nil {
					// If it's a map value for an index, convert it to proper KV format
					if mapValue, ok := v.(map[string]any); ok && len(parts) == 2 {
						for mk, mv := range mapValue {
							indexedKey := fmt.Sprintf("%s.%s.%s", prefix, parts[1], mk)
							b.Set(indexedKey, mv)
						}
						continue
					}
					// Special handling for []KV type
					if kvs, ok := v.([]KV); ok && len(parts) == 2 {
						// Convert KVs to map[string]any
						kvMap := make(map[string]any)
						for _, kv := range kvs {
							kvMap[kv.Key] = kv.Value
						}
						for mk, mv := range kvMap {
							indexedKey := fmt.Sprintf("%s.%s.%s", prefix, parts[1], mk)
							b.Set(indexedKey, mv)
						}
						continue
					}
					b.Set(k, v)
					continue
				}
			}

			// Collect fields with same prefix
			if _, ok := prefixMap[prefix]; !ok {
				prefixMap[prefix] = make(map[string]any)
			}

			// For deeply nested fields, create nested maps
			current := prefixMap[prefix]
			for i := 1; i < len(parts)-1; i++ {
				key := parts[i]
				if _, ok := current[key]; !ok {
					current[key] = make(map[string]any)
				}
				current = current[key].(map[string]any)
			}

			// Set the final value
			current[parts[len(parts)-1]] = v
			continue
		}

		b.Set(k, v)
	}

	// Now handle all collected prefixes
	for prefix, fields := range prefixMap {
		b.Set(prefix, fields)
	}

	return b
}

// FromMap fills the struct from a map with zero allocations where possible
func (b *StructBuilder) FromMap(data map[string]any) *StructBuilder {
	// First, collect all fields with the same prefix
	prefixMap := make(map[string]map[string]any)

	for k, v := range data {
		// Try direct field first
		if _, ok := b.getOffset(k); ok {
			// Special handling for []KV type
			if kvs, ok := v.([]KV); ok {
				// Convert KVs to map[string]any
				kvMap := make(map[string]any)
				for _, kv := range kvs {
					kvMap[kv.Key] = kv.Value
				}
				b.Set(k, kvMap)
				continue
			}
			b.Set(k, v)
			continue
		}

		// Handle nested fields
		if strings.Contains(k, ".") {
			parts := strings.Split(k, ".")
			prefix := parts[0]

			// Check if it's a numeric index
			if len(parts) > 1 {
				if _, err := strconv.Atoi(parts[1]); err == nil {
					// If it's a map value for an index, convert it to proper KV format
					if mapValue, ok := v.(map[string]any); ok && len(parts) == 2 {
						for mk, mv := range mapValue {
							indexedKey := fmt.Sprintf("%s.%s.%s", prefix, parts[1], mk)
							b.Set(indexedKey, mv)
						}
						continue
					}
					// Special handling for []KV type
					if kvs, ok := v.([]KV); ok && len(parts) == 2 {
						// Convert KVs to map[string]any
						kvMap := make(map[string]any)
						for _, kv := range kvs {
							kvMap[kv.Key] = kv.Value
						}
						for mk, mv := range kvMap {
							indexedKey := fmt.Sprintf("%s.%s.%s", prefix, parts[1], mk)
							b.Set(indexedKey, mv)
						}
						continue
					}
					b.Set(k, v)
					continue
				}
			}

			// Collect fields with same prefix
			if _, ok := prefixMap[prefix]; !ok {
				prefixMap[prefix] = make(map[string]any)
			}

			// For deeply nested fields, create nested maps
			current := prefixMap[prefix]
			for i := 1; i < len(parts)-1; i++ {
				key := parts[i]
				if _, ok := current[key]; !ok {
					current[key] = make(map[string]any)
				}
				current = current[key].(map[string]any)
			}

			// Set the final value
			current[parts[len(parts)-1]] = v
			continue
		}

		b.Set(k, v)
	}

	// Now handle all collected prefixes
	for prefix, fields := range prefixMap {
		b.Set(prefix, fields)
	}

	return b
}

// ToMap converts the struct to a map with minimal allocations
func (b *StructBuilder) ToMap() map[string]any {
	// Get map from pool
	resultPtr := mapPool.Get().(*map[string]any)
	result := *resultPtr
	clear(result)

	b.fields.each(func(fieldName string, info fieldInfo) {
		ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + info.offset)
		if ptr == nil {
			return
		}

		switch info.typ.Kind() {
		case reflect.String:
			result[fieldName] = *(*string)(ptr)
		case reflect.Int:
			result[fieldName] = *(*int)(ptr)
		case reflect.Bool:
			result[fieldName] = *(*bool)(ptr)
		case reflect.Float64:
			result[fieldName] = *(*float64)(ptr)
		case reflect.Ptr:
			ptrVal := *(*unsafe.Pointer)(ptr)
			if ptrVal != nil {
				elemVal := reflect.NewAt(info.typ.Elem(), ptrVal).Elem()
				result[fieldName] = elemVal.Interface()
			}
		case reflect.Map:
			mapVal := reflect.NewAt(info.typ, ptr).Elem()
			if !mapVal.IsNil() {
				result[fieldName] = mapVal.Interface()
			}
		case reflect.Slice:
			sliceVal := reflect.NewAt(info.typ, ptr).Elem()
			if !sliceVal.IsNil() {
				result[fieldName] = sliceVal.Interface()
			}
		case reflect.Struct:
			if info.typ == reflect.TypeOf(time.Time{}) {
				result[fieldName] = *(*time.Time)(ptr)
			} else {
				structVal := reflect.NewAt(info.typ, ptr).Elem()
				result[fieldName] = structVal.Interface()
			}
		}
	})

	return result
}

// PutMap returns a map to the pool
func PutMap(m map[string]any) {
	mapPool.Put(&m)
}

// Clone creates a deep copy of the struct with zero allocations
func (b *StructBuilder) Clone() any {
	// Create new instance of same type
	newVal := reflect.New(b.structType)
	newPtr := newVal.Interface()

	// Copy memory directly
	srcPtr := (*emptyInterface)(unsafe.Pointer(&b.ptr)).word
	dstPtr := (*emptyInterface)(unsafe.Pointer(&newPtr)).word
	size := b.structType.Size()

	// Copy in 8-byte chunks
	chunks := size / 8
	remainder := size % 8

	for i := uintptr(0); i < chunks; i++ {
		offset := i * 8
		*(*uint64)(unsafe.Pointer(uintptr(dstPtr) + offset)) = *(*uint64)(unsafe.Pointer(uintptr(srcPtr) + offset))
	}

	// Copy remaining bytes
	for i := uintptr(0); i < remainder; i++ {
		offset := chunks*8 + i
		*(*uint8)(unsafe.Pointer(uintptr(dstPtr) + offset)) = *(*uint8)(unsafe.Pointer(uintptr(srcPtr) + offset))
	}

	return newPtr
}

// Each iterates over struct fields with a callback
func (b *StructBuilder) Each(fn func(name string, value any)) *StructBuilder {
	b.fields.each(func(fieldName string, _ fieldInfo) {
		value := b.Get(fieldName)
		fn(fieldName, value)
	})
	return b
}

// Filter creates a new map with fields that match the predicate
func (b *StructBuilder) Filter(predicate func(name string, value any) bool) map[string]any {
	// Get map from pool
	resultPtr := mapPool.Get().(*map[string]any)
	result := *resultPtr
	clear(result)

	b.fields.each(func(fieldName string, _ fieldInfo) {
		value := b.Get(fieldName)
		if predicate(fieldName, value) {
			result[fieldName] = value
		}
	})

	return result
}

// func (b *StructBuilder) init(t reflect.Type) {
// 	fields := newFieldInfoMap(t.NumField())
// 	for i := 0; i < t.NumField(); i++ {
// 		field := t.Field(i)
// 		name := ToSnakeCase(field.Name)
// 		fields.set(name, fieldInfo{
// 			offset: field.Offset,
// 			typ:    field.Type,
// 		})
// 	}
// 	b.fields = fields
// }

// GetString gets a string field value by name with zero allocations
func (b *StructBuilder) GetString(fieldPath string) string {
	if offset, ok := b.getOffset(fieldPath); ok {
		fieldType, _ := b.getFieldType(fieldPath)
		if fieldType.Kind() == reflect.String {
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			return *(*string)(ptr)
		}
	}
	return ""
}

// GetInt gets an int field value by name with zero allocations
func (b *StructBuilder) GetInt(fieldPath string) int {
	if offset, ok := b.getOffset(fieldPath); ok {
		fieldType, _ := b.getFieldType(fieldPath)
		if fieldType.Kind() == reflect.Int {
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			return *(*int)(ptr)
		}
	}
	return 0
}

// GetBool gets a bool field value by name with zero allocations
func (b *StructBuilder) GetBool(fieldPath string) bool {
	if offset, ok := b.getOffset(fieldPath); ok {
		fieldType, _ := b.getFieldType(fieldPath)
		if fieldType.Kind() == reflect.Bool {
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			return *(*bool)(ptr)
		}
	}
	return false
}

// GetFloat64 gets a float64 field value by name with zero allocations
func (b *StructBuilder) GetFloat64(fieldPath string) float64 {
	if offset, ok := b.getOffset(fieldPath); ok {
		fieldType, _ := b.getFieldType(fieldPath)
		if fieldType.Kind() == reflect.Float64 {
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			return *(*float64)(ptr)
		}
	}
	return 0
}

// SetString sets a string field value by name with zero allocations
func (b *StructBuilder) SetString(fieldPath string, value string) *StructBuilder {
	if offset, ok := b.getOffset(fieldPath); ok {
		fieldType, _ := b.getFieldType(fieldPath)
		if fieldType.Kind() == reflect.String {
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			*(*string)(ptr) = value
			return b
		}
	}
	return b
}

// SetInt sets an int field value by name with zero allocations
func (b *StructBuilder) SetInt(fieldPath string, value int) *StructBuilder {
	if offset, ok := b.getOffset(fieldPath); ok {
		fieldType, _ := b.getFieldType(fieldPath)
		if fieldType.Kind() == reflect.Int {
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			*(*int)(ptr) = value
			return b
		}
	}
	return b
}

// SetBool sets a bool field value by name with zero allocations
func (b *StructBuilder) SetBool(fieldPath string, value bool) *StructBuilder {
	if offset, ok := b.getOffset(fieldPath); ok {
		fieldType, _ := b.getFieldType(fieldPath)
		if fieldType.Kind() == reflect.Bool {
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			*(*bool)(ptr) = value
			return b
		}
	}
	return b
}

// SetFloat64 sets a float64 field value by name with zero allocations
func (b *StructBuilder) SetFloat64(fieldPath string, value float64) *StructBuilder {
	if offset, ok := b.getOffset(fieldPath); ok {
		fieldType, _ := b.getFieldType(fieldPath)
		if fieldType.Kind() == reflect.Float64 {
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			*(*float64)(ptr) = value
			return b
		}
	}
	return b
}

// SetFloat32 sets a float32 field value by name with zero allocations
func (b *StructBuilder) SetFloat32(fieldPath string, value float32) *StructBuilder {
	if offset, ok := b.getOffset(fieldPath); ok {
		fieldType, _ := b.getFieldType(fieldPath)
		if fieldType.Kind() == reflect.Float32 {
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			*(*float32)(ptr) = value
			return b
		}
	}
	return b
}

// parseTimeString attempts to parse a time string using multiple formats
func parseTimeString(v string) (time.Time, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return time.Time{}, fmt.Errorf("empty time string")
	}

	// Try exact match first
	for _, format := range timeFormats {
		if t, err := time.Parse(format, v); err == nil {
			return t, nil
		}
	}

	// If no exact match, try truncating for longer formats
	for _, format := range timeFormats {
		if len(v) >= len(format) {
			input := v[:len(format)]
			if t, err := time.Parse(format, input); err == nil {
				return t, nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("could not parse time '%s'", v)
}

// Fill fills the struct with values from a map or KV pairs
func (b *StructBuilder) Fill(data any) *StructBuilder {
	switch v := data.(type) {
	case map[string]any:
		return b.FromMap(v)
	case []KV:
		return b.FromKV(v...)
	case KV:
		return b.FromKV(v)
	case map[string]string:
		// Get map from pool
		mapPtr := mapPool.Get().(*map[string]any)
		m := *mapPtr
		clear(m)
		defer mapPool.Put(mapPtr)

		// Convert string map to any map
		for k, v := range v {
			m[k] = v
		}
		return b.FromMap(m)
	case map[string]int:
		// Get map from pool
		mapPtr := mapPool.Get().(*map[string]any)
		m := *mapPtr
		clear(m)
		defer mapPool.Put(mapPtr)

		// Convert int map to any map
		for k, v := range v {
			m[k] = v
		}
		return b.FromMap(m)
	case map[string]bool:
		// Get map from pool
		mapPtr := mapPool.Get().(*map[string]any)
		m := *mapPtr
		clear(m)
		defer mapPool.Put(mapPtr)

		// Convert bool map to any map
		for k, v := range v {
			m[k] = v
		}
		return b.FromMap(m)
	case map[string]float64:
		// Get map from pool
		mapPtr := mapPool.Get().(*map[string]any)
		m := *mapPtr
		clear(m)
		defer mapPool.Put(mapPtr)

		// Convert float64 map to any map
		for k, v := range v {
			m[k] = v
		}
		return b.FromMap(m)
	default:
		// Try to convert to map using reflection
		val := reflect.ValueOf(data)
		if val.Kind() == reflect.Map && val.Type().Key().Kind() == reflect.String {
			// Get map from pool
			mapPtr := mapPool.Get().(*map[string]any)
			m := *mapPtr
			clear(m)
			defer mapPool.Put(mapPtr)

			// Convert using reflection
			iter := val.MapRange()
			for iter.Next() {
				k := iter.Key().String()
				v := iter.Value().Interface()
				m[k] = v
			}
			return b.FromMap(m)
		}
	}
	return b
}

func (b *StructBuilder) getOffset(name string) (uintptr, bool) {
	info, ok := b.fields.get(name)
	if ok {
		return info.offset, true
	}
	return 0, false
}

func (b *StructBuilder) getFieldType(name string) (reflect.Type, bool) {
	info, ok := b.fields.get(name)
	if ok {
		return info.typ, true
	}
	return nil, false
}

// Add methods to fieldInfoMap for iteration
func (m *fieldInfoMap) each(fn func(name string, info fieldInfo)) {
	if m.isFixed {
		for i := 0; i < m.fixedSize; i++ {
			fn(m.fixedKeys[i], m.fixedInfos[i])
		}
	} else {
		for k, v := range m.mapInfo {
			fn(k, v)
		}
	}
}
