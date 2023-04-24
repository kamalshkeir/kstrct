package kstrct

import (
	"testing"
	"time"
)

type Something struct {
	Id        int
	Email     string
	IsAdmin   bool
	CreatedAt time.Time
}

// BenchmarkFillFromMap-4           1757280               662.4 ns/op           360 B/op          3 allocs/op
// BenchmarkFillFromMapS-4          1893105               618.1 ns/op           360 B/op          3 allocs/op
// BenchmarkFillByIndex-4           2991378               384.6 ns/op            24 B/op          1 allocs/op

func BenchmarkFillFromMap(b *testing.B) {
	t := time.Now()
	a := Something{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := FillFromMap(&a, map[string]any{
			"id":         1,
			"email":      "something",
			"is_admin":   true,
			"created_at": t,
		})
		if err != nil {
			b.Error(err)
		}
		if a.Id != 1 || !a.IsAdmin || a.CreatedAt != t {
			b.Errorf("something wrong %v", a)
		}
	}
}

func BenchmarkFillFromMapS(b *testing.B) {
	t := time.Now()
	a := Something{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := FillFromMapS(&a, map[string]any{
			"id":         1,
			"email":      "something",
			"is_admin":   true,
			"created_at": t,
		})
		if err != nil {
			b.Error(err)
		}
		if a.Id != 1 || !a.IsAdmin || a.CreatedAt != t {
			b.Errorf("something wrong %v", a)
		}
	}
}

func BenchmarkFillByIndex(b *testing.B) {
	t := time.Now()
	a := Something{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := FillByIndex(&a, map[int]any{
			0: 1,
			1: "something",
			2: true,
			3: t,
		})
		if err != nil {
			b.Error(err)
		}
		if a.Id != 1 || !a.IsAdmin || a.CreatedAt != t {
			b.Errorf("something wrong %v", a)
		}
	}
}
