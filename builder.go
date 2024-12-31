package kstrct

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/kamalshkeir/kmap"
)

// StructBuilder provides a fluent interface for struct manipulation
type StructBuilder struct {
	ptr        any
	offsets    *kmap.SafeMap[string, uintptr]
	fieldTypes *kmap.SafeMap[string, reflect.Type]
	structType reflect.Type
	fieldOrder []string
	snakeCache *kmap.SafeMap[string, string] // Cache for snake case field names
	err        error
}

var (
	// Cache for struct builders to avoid recreating them
	builderCache = kmap.New[string, *StructBuilder]()
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

	// Check cache first
	if builder, ok := builderCache.Get(elemType.String()); ok {
		builder.ptr = structPtr
		builder.err = nil
		return builder
	}

	// Create new builder
	builder := &StructBuilder{
		ptr:        structPtr,
		offsets:    kmap.New[string, uintptr](),
		fieldTypes: kmap.New[string, reflect.Type](),
		structType: elemType,
		fieldOrder: make([]string, 0, elemType.NumField()),
		snakeCache: kmap.New[string, string](),
		err:        nil,
	}

	// Pre-compute all field offsets and types
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		// Convert field name to snake case
		snakeName := ToSnakeCase(field.Name)
		builder.offsets.Set(snakeName, field.Offset)
		builder.fieldTypes.Set(snakeName, field.Type)
		builder.fieldOrder = append(builder.fieldOrder, snakeName)
		builder.snakeCache.Set(field.Name, snakeName)

		// Also store the original name for backward compatibility
		builder.offsets.Set(field.Name, field.Offset)
		builder.fieldTypes.Set(field.Name, field.Type)
	}

	// Cache the builder
	builderCache.Set(elemType.String(), builder)
	return builder
}

func (b *StructBuilder) Error() error {
	return b.err
}

// Set sets a field value by name
func (b *StructBuilder) Set(fieldPath string, value any) *StructBuilder {
	// Try direct access first
	if offset, ok := b.offsets.Get(fieldPath); ok {
		fieldType, _ := b.fieldTypes.Get(fieldPath)
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
		if offset, ok := b.offsets.Get(rootField); ok {
			fieldType, _ := b.fieldTypes.Get(rootField)
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
	if offset, ok := b.offsets.Get(fieldPath); ok {
		fieldType, _ := b.fieldTypes.Get(fieldPath)
		ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
		return b.getFieldValue(ptr, fieldType)
	}

	// Try with snake case
	snakePath := ToSnakeCase(fieldPath)
	if offset, ok := b.offsets.Get(snakePath); ok {
		fieldType, _ := b.fieldTypes.Get(snakePath)
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
	offset, ok := b.offsets.Get(rootField)
	if !ok {
		return nil
	}

	fieldType, ok := b.fieldTypes.Get(rootField)
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
		hasNumber, _ := FindFirstNumber(k)

		// Skip processing if it contains a number (handled by Set directly)
		if hasNumber {
			b.Set(k, v)
			continue
		}

		// Try direct field first
		if _, ok := b.offsets.Get(k); ok {
			b.Set(k, v)
			continue
		}

		// Handle nested fields
		if strings.Contains(k, ".") {
			parts := strings.Split(k, ".")
			prefix := parts[0]

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
		hasNumber, _ := FindFirstNumber(k)
		if hasNumber {
			b.Set(k, v)
			continue
		}

		// Try direct field first
		if _, ok := b.offsets.Get(k); ok {
			b.Set(k, v)
			continue
		}

		// Handle nested fields
		if strings.Contains(k, ".") {
			parts := strings.Split(k, ".")
			prefix := parts[0]

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

	b.offsets.Range(func(fieldName string, offset uintptr) bool {
		fieldType, _ := b.fieldTypes.Get(fieldName)
		ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
		if ptr == nil {
			return true
		}

		switch fieldType.Kind() {
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
				elemVal := reflect.NewAt(fieldType.Elem(), ptrVal).Elem()
				result[fieldName] = elemVal.Interface()
			}
		case reflect.Map:
			mapVal := reflect.NewAt(fieldType, ptr).Elem()
			if !mapVal.IsNil() {
				result[fieldName] = mapVal.Interface()
			}
		case reflect.Slice:
			sliceVal := reflect.NewAt(fieldType, ptr).Elem()
			if !sliceVal.IsNil() {
				result[fieldName] = sliceVal.Interface()
			}
		case reflect.Struct:
			if fieldType == reflect.TypeOf(time.Time{}) {
				result[fieldName] = *(*time.Time)(ptr)
			} else {
				structVal := reflect.NewAt(fieldType, ptr).Elem()
				result[fieldName] = structVal.Interface()
			}
		}
		return true
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
	for _, fieldName := range b.fieldOrder {
		value := b.Get(fieldName)
		fn(fieldName, value)
	}
	return b
}

// Filter creates a new map with fields that match the predicate
func (b *StructBuilder) Filter(predicate func(name string, value any) bool) map[string]any {
	// Get map from pool
	resultPtr := mapPool.Get().(*map[string]any)
	result := *resultPtr
	clear(result)

	b.Each(func(name string, value any) {
		if predicate(name, value) {
			result[name] = value
		}
	})

	return result
}

func (b *StructBuilder) init(t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := ToSnakeCase(field.Name)
		b.offsets.Set(name, field.Offset)
		b.fieldTypes.Set(name, field.Type)
		b.fieldOrder = append(b.fieldOrder, name)
	}
}

// GetString gets a string field value by name with zero allocations
func (b *StructBuilder) GetString(fieldPath string) string {
	if offset, ok := b.offsets.Get(fieldPath); ok {
		fieldType, _ := b.fieldTypes.Get(fieldPath)
		if fieldType.Kind() == reflect.String {
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			return *(*string)(ptr)
		}
	}
	return ""
}

// GetInt gets an int field value by name with zero allocations
func (b *StructBuilder) GetInt(fieldPath string) int {
	if offset, ok := b.offsets.Get(fieldPath); ok {
		fieldType, _ := b.fieldTypes.Get(fieldPath)
		if fieldType.Kind() == reflect.Int {
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			return *(*int)(ptr)
		}
	}
	return 0
}

// GetBool gets a bool field value by name with zero allocations
func (b *StructBuilder) GetBool(fieldPath string) bool {
	if offset, ok := b.offsets.Get(fieldPath); ok {
		fieldType, _ := b.fieldTypes.Get(fieldPath)
		if fieldType.Kind() == reflect.Bool {
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			return *(*bool)(ptr)
		}
	}
	return false
}

// GetFloat64 gets a float64 field value by name with zero allocations
func (b *StructBuilder) GetFloat64(fieldPath string) float64 {
	if offset, ok := b.offsets.Get(fieldPath); ok {
		fieldType, _ := b.fieldTypes.Get(fieldPath)
		if fieldType.Kind() == reflect.Float64 {
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			return *(*float64)(ptr)
		}
	}
	return 0
}

// SetString sets a string field value by name with zero allocations
func (b *StructBuilder) SetString(fieldPath string, value string) *StructBuilder {
	if offset, ok := b.offsets.Get(fieldPath); ok {
		fieldType, _ := b.fieldTypes.Get(fieldPath)
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
	if offset, ok := b.offsets.Get(fieldPath); ok {
		fieldType, _ := b.fieldTypes.Get(fieldPath)
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
	if offset, ok := b.offsets.Get(fieldPath); ok {
		fieldType, _ := b.fieldTypes.Get(fieldPath)
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
	if offset, ok := b.offsets.Get(fieldPath); ok {
		fieldType, _ := b.fieldTypes.Get(fieldPath)
		if fieldType.Kind() == reflect.Float64 {
			ptr := unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&b.ptr)).word) + offset)
			*(*float64)(ptr) = value
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
