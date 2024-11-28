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

// cpu: Intel(R) Core(TM) i5-7300HQ CPU @ 2.50GHz
// BenchmarkFillFromMap-4                   1536951               745.6 ns/op           408 B/op          4 allocs/op
// BenchmarkFillFromKV-4                    3356922               355.5 ns/op            48 B/op          1 allocs/op
// BenchmarkFrom-4                          2882827               434.9 ns/op            56 B/op          4 allocs/op
// BenchmarkRange-4                         3131361               379.2 ns/op            56 B/op          4 allocs/op
// BenchmarkFill-4                          3929871               306.1 ns/op             0 B/op          0 allocs/op
// BenchmarkFillM-4                         2054829               579.4 ns/op            24 B/op          1 allocs/op
// BenchmarkMapstructure-4                   356590              3218 ns/op            1496 B/op         31 allocs/op
// BenchmarkMapstructureDecoder-4            404091              2980 ns/op            1344 B/op         28 allocs/op
// PASS
// ok      github.com/kamalshkeir/kstrct   12.692s

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
func BenchmarkFillFromKV(b *testing.B) {
	t := time.Now()
	a := Something{}
	b.ResetTimer()
	kv := []KV{}
	kv = append(kv, KV{"id", 1}, KV{"email", "something"}, KV{"is_admin", true}, KV{"created_at", t})
	for i := 0; i < b.N; i++ {
		err := FillFromKV(&a, kv)
		if err != nil {
			b.Error(err)
		}
		if a.Id != 1 || !a.IsAdmin || a.CreatedAt != t {
			b.Errorf("something wrong %v", a)
		}
	}
}

func BenchmarkFrom(b *testing.B) {
	t := time.Now()
	s := Something{
		Id:        1,
		Email:     "something",
		IsAdmin:   true,
		CreatedAt: t,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var count int
		for _, ctx := range From(&s) {
			if ctx.Value != nil {
				count++
			}
		}
		if count != 4 {
			b.Errorf("expected 4 fields, got %d", count)
		}
	}
}

func BenchmarkRange(b *testing.B) {
	t := time.Now()
	s := Something{
		Id:        1,
		Email:     "something",
		IsAdmin:   true,
		CreatedAt: t,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var count int
		Range(&s, func(ctx FieldCtx) bool {
			if ctx.Value != nil {
				count++
			}
			return true
		})
		if count != 4 {
			b.Errorf("expected 4 fields, got %d", count)
		}
	}
}

func BenchmarkFill(b *testing.B) {
	t := time.Now()
	a := Something{}
	b.ResetTimer()
	kv := []KV{}
	kv = append(kv, KV{"id", 1}, KV{"email", "something"}, KV{"is_admin", true}, KV{"created_at", t})
	for i := 0; i < b.N; i++ {
		err := Fill(&a, kv)
		if err != nil {
			b.Error(err)
		}
		if a.Id != 1 || !a.IsAdmin || a.CreatedAt != t {
			b.Errorf("something wrong %v", a)
		}
	}
}

func BenchmarkFillM(b *testing.B) {
	t := time.Now()
	a := Something{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := FillM(&a, map[string]any{
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
