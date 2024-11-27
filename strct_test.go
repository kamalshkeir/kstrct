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
