package kstrct

import (
	"reflect"
	"time"
	"unsafe"
)

// emptyInterface is the header of an any value.
// This matches the internal representation of any in the runtime.
type emptyInterface struct {
	typ  unsafe.Pointer // points to the type information
	word unsafe.Pointer // points to the data
}

// FieldOffset returns the offset of a field in a struct type.
// This is useful for direct memory access to struct fields.
// Returns -1 if field not found.
func FieldOffset(structType reflect.Type, fieldName string) uintptr {
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	if field, ok := structType.FieldByName(fieldName); ok {
		return field.Offset
	}
	return ^uintptr(0)
}

// GetFieldPointer returns a pointer to a field in a struct.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetFieldPointer(structPtr any, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr((*emptyInterface)(unsafe.Pointer(&structPtr)).word) + offset)
}

// GetStringField returns a string field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetStringField(structPtr any, offset uintptr) string {
	strPtr := (*string)(GetFieldPointer(structPtr, offset))
	return *strPtr
}

// SetStringField sets a string field value in a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func SetStringField(structPtr any, offset uintptr, value string) {
	strPtr := (*string)(GetFieldPointer(structPtr, offset))
	*strPtr = value
}

// GetIntField returns an int field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetIntField(structPtr any, offset uintptr) int {
	intPtr := (*int)(GetFieldPointer(structPtr, offset))
	return *intPtr
}

// SetIntField sets an int field value in a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func SetIntField(structPtr any, offset uintptr, value int) {
	intPtr := (*int)(GetFieldPointer(structPtr, offset))
	*intPtr = value
}
func SetInt64Field(structPtr any, offset uintptr, value int64) {
	intPtr := (*int64)(GetFieldPointer(structPtr, offset))
	*intPtr = value
}
func SetInt32Field(structPtr any, offset uintptr, value int32) {
	intPtr := (*int32)(GetFieldPointer(structPtr, offset))
	*intPtr = value
}
func SetUInt32Field(structPtr any, offset uintptr, value uint32) {
	intPtr := (*uint32)(GetFieldPointer(structPtr, offset))
	*intPtr = value
}
func SetUIntField(structPtr any, offset uintptr, value uint) {
	intPtr := (*uint)(GetFieldPointer(structPtr, offset))
	*intPtr = value
}
func SetUInt8Field(structPtr any, offset uintptr, value uint8) {
	intPtr := (*uint8)(GetFieldPointer(structPtr, offset))
	*intPtr = value
}

func SetTimeField(structPtr any, offset uintptr, value time.Time) {
	timePtr := (*time.Time)(GetFieldPointer(structPtr, offset))
	*timePtr = value
}

// GetBoolField returns a bool field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetBoolField(structPtr any, offset uintptr) bool {
	boolPtr := (*bool)(GetFieldPointer(structPtr, offset))
	return *boolPtr
}

// SetBoolField sets a bool field value in a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func SetBoolField(structPtr any, offset uintptr, value bool) {
	boolPtr := (*bool)(GetFieldPointer(structPtr, offset))
	*boolPtr = value
}

// CopyStruct makes a shallow copy of a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func CopyStruct(dst, src any) {
	dstPtr := (*emptyInterface)(unsafe.Pointer(&dst)).word
	srcPtr := (*emptyInterface)(unsafe.Pointer(&src)).word
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	size := typ.Size()
	typAlign := uintptr(typ.Align())

	// Copy memory with proper alignment
	for offset := uintptr(0); offset < size; offset += typAlign {
		remaining := size - offset
		if remaining < typAlign {
			// Handle the last chunk that might be smaller than alignment
			*(*byte)(unsafe.Pointer(uintptr(dstPtr) + offset)) = *(*byte)(unsafe.Pointer(uintptr(srcPtr) + offset))
		} else {
			// Copy aligned chunks
			*(*uintptr)(unsafe.Pointer(uintptr(dstPtr) + offset)) = *(*uintptr)(unsafe.Pointer(uintptr(srcPtr) + offset))
		}
	}
}

// StructToMap converts a struct to a map without allocations using a preallocated map.
// The caller must provide a pre-allocated map to avoid allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func StructToMap(structPtr any, result map[string]any) {
	val := reflect.ValueOf(structPtr)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		offset := field.Offset
		ptr := GetFieldPointer(structPtr, offset)

		switch field.Type.Kind() {
		case reflect.String:
			result[field.Name] = *(*string)(ptr)
		case reflect.Int:
			result[field.Name] = *(*int)(ptr)
		case reflect.Bool:
			result[field.Name] = *(*bool)(ptr)
			// Add more types as needed
		}
	}
}

// CompareStructs compares two structs of the same type without allocations.
// Returns true if all fields are equal.
// UNSAFE: This function provides direct memory access for maximum performance.
func CompareStructs(a, b any) bool {
	aPtr := (*emptyInterface)(unsafe.Pointer(&a)).word
	bPtr := (*emptyInterface)(unsafe.Pointer(&b)).word
	typ := reflect.TypeOf(a)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	size := typ.Size()

	// Compare memory directly
	for offset := uintptr(0); offset < size; offset++ {
		if *(*byte)(unsafe.Pointer(uintptr(aPtr) + offset)) != *(*byte)(unsafe.Pointer(uintptr(bPtr) + offset)) {
			return false
		}
	}
	return true
}

// GetFieldsByTag returns field offsets for fields with a specific tag without allocations.
// The result is stored in the provided map to avoid allocations.
func GetFieldsByTag(structType reflect.Type, tagName string, result map[string]uintptr) {
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if tag := field.Tag.Get(tagName); tag != "" {
			result[tag] = field.Offset
		}
	}
}

// GetFloat64Field returns a float64 field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetFloat64Field(structPtr any, offset uintptr) float64 {
	ptr := (*float64)(GetFieldPointer(structPtr, offset))
	return *ptr
}

// GetFloat32Field returns a float32 field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetFloat32Field(structPtr any, offset uintptr) float32 {
	ptr := (*float32)(GetFieldPointer(structPtr, offset))
	return *ptr
}

// GetUintField returns a uint field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetUintField(structPtr any, offset uintptr) uint {
	ptr := (*uint)(GetFieldPointer(structPtr, offset))
	return *ptr
}

// GetUint64Field returns a uint64 field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetUint64Field(structPtr any, offset uintptr) uint64 {
	ptr := (*uint64)(GetFieldPointer(structPtr, offset))
	return *ptr
}

// GetUint32Field returns a uint32 field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetUint32Field(structPtr any, offset uintptr) uint32 {
	ptr := (*uint32)(GetFieldPointer(structPtr, offset))
	return *ptr
}

// GetUint16Field returns a uint16 field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetUint16Field(structPtr any, offset uintptr) uint16 {
	ptr := (*uint16)(GetFieldPointer(structPtr, offset))
	return *ptr
}

// GetUint8Field returns a uint8 field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetUint8Field(structPtr any, offset uintptr) uint8 {
	ptr := (*uint8)(GetFieldPointer(structPtr, offset))
	return *ptr
}

// GetInt64Field returns an int64 field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetInt64Field(structPtr any, offset uintptr) int64 {
	ptr := (*int64)(GetFieldPointer(structPtr, offset))
	return *ptr
}

// GetInt32Field returns an int32 field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetInt32Field(structPtr any, offset uintptr) int32 {
	ptr := (*int32)(GetFieldPointer(structPtr, offset))
	return *ptr
}

// GetInt16Field returns an int16 field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetInt16Field(structPtr any, offset uintptr) int16 {
	ptr := (*int16)(GetFieldPointer(structPtr, offset))
	return *ptr
}

// GetInt8Field returns an int8 field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetInt8Field(structPtr any, offset uintptr) int8 {
	ptr := (*int8)(GetFieldPointer(structPtr, offset))
	return *ptr
}

// GetTimeField returns a time.Time field value from a struct without allocations.
// UNSAFE: This function provides direct memory access for maximum performance.
func GetTimeField(structPtr any, offset uintptr) time.Time {
	ptr := (*time.Time)(GetFieldPointer(structPtr, offset))
	return *ptr
}

// NumberSegment holds information about a found number
type NumberSegment struct {
	Value   int // The parsed number value
	Start   int // Start position in string
	End     int // End position in string
	Segment int // Which dot-separated segment (0-based)
}

// FindNumbersInPath scans a path string for numbers with zero allocations
func FindNumbersInPath(s string, fn func(NumberSegment) bool) {
	segment := 0
	start := 0

	for i := 0; i < len(s); i++ {
		// Track segments
		if s[i] == '.' {
			segment++
			start = i + 1
			continue
		}

		// If we find a digit
		if s[i] >= '0' && s[i] <= '9' {
			numStart := i
			numEnd := i

			// Find the complete number
			for numEnd+1 < len(s) && s[numEnd+1] >= '0' && s[numEnd+1] <= '9' {
				numEnd++
			}

			// Only consider it a valid number if it's at the start of a segment
			if numStart == start {
				// Check if this is a complete segment
				isComplete := numEnd+1 == len(s) || s[numEnd+1] == '.'

				if isComplete {
					// Parse number without allocations
					num := 0
					for j := numStart; j <= numEnd; j++ {
						num = num*10 + int(s[j]-'0')
					}

					// Call the callback - if it returns false, stop processing
					if !fn(NumberSegment{
						Value:   num,
						Start:   numStart,
						End:     numEnd + 1,
						Segment: segment,
					}) {
						return
					}
				}
			}

			// Skip to end of number
			i = numEnd
		}
	}
}
