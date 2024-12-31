package kstrct

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	IsActive bool   `json:"is_active"`
}

func BenchmarkStructOperations(b *testing.B) {
	s := &TestStruct{
		Name:     "John Doe",
		Age:      30,
		IsActive: true,
	}
	s2 := &TestStruct{}

	// Cache field offsets
	nameOffset := FieldOffset(reflect.TypeOf(s), "Name")
	ageOffset := FieldOffset(reflect.TypeOf(s), "Age")
	activeOffset := FieldOffset(reflect.TypeOf(s), "IsActive")

	// Pre-allocate map for StructToMap
	m := make(map[string]any, 3)
	tagMap := make(map[string]uintptr)

	b.Run("GetStringField", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = GetStringField(s, nameOffset)
		}
	})

	b.Run("SetStringField", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			SetStringField(s, nameOffset, "Jane Doe")
		}
	})

	b.Run("GetIntField", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = GetIntField(s, ageOffset)
		}
	})

	b.Run("SetIntField", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			SetIntField(s, ageOffset, 31)
		}
	})

	b.Run("GetBoolField", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = GetBoolField(s, activeOffset)
		}
	})

	b.Run("SetBoolField", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			SetBoolField(s, activeOffset, false)
		}
	})

	b.Run("CopyStruct", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			CopyStruct(s2, s)
		}
	})

	b.Run("StructToMap/Unsafe", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			clear(m)
			StructToMap(s, m)
		}
	})

	b.Run("CompareStructs", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = CompareStructs(s, s2)
		}
	})

	b.Run("GetFieldsByTag", func(b *testing.B) {
		typ := reflect.TypeOf(s)

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			clear(tagMap)
			GetFieldsByTag(typ, "json", tagMap)
		}
	})
}

// Regular reflection-based implementations for comparison
func BenchmarkReflectionOperations(b *testing.B) {
	s := &TestStruct{
		Name:     "John Doe",
		Age:      30,
		IsActive: true,
	}
	s2 := &TestStruct{}

	b.Run("GetField/Reflection", func(b *testing.B) {
		b.ReportAllocs()
		val := reflect.ValueOf(s).Elem()
		field := val.FieldByName("Name")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = field.String()
		}
	})

	b.Run("SetField/Reflection", func(b *testing.B) {
		b.ReportAllocs()
		val := reflect.ValueOf(s).Elem()
		field := val.FieldByName("Name")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			field.SetString("Jane Doe")
		}
	})

	b.Run("CopyStruct/Reflection", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			reflect.ValueOf(s2).Elem().Set(reflect.ValueOf(s).Elem())
		}
	})
}
