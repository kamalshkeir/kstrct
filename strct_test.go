package kstrct

import (
	"reflect"
	"testing"
	"time"
)

type Something struct {
	Id        int
	Email     string
	IsAdmin   bool
	CreatedAt time.Time
}

type WeekTimeslots struct {
	Monday    []string
	Tuesday   []string
	Wednesday []string
	Thursday  []string
	Friday    []string
	Saturday  []string
	Sunday    []string
}

type Doctor struct {
	Name          string
	WeekTimeslots *[]WeekTimeslots
}

type DoctorS struct {
	Name          string
	WeekTimeslots *WeekTimeslots
}

func TestFillNestedFieldsSlice(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		input    []KV
		expected Doctor
	}{
		{
			name: "nested week timeslots",
			input: []KV{
				{"name", "Dr. Smith"},
				{"week_timeslots.monday", "09:00,10:00,11:00"},
				{"week_timeslots.tuesday", "14:00,15:00"},
			},
			expected: Doctor{
				Name: "Dr. Smith",
				WeekTimeslots: &[]WeekTimeslots{
					{
						Monday:  []string{"09:00", "10:00", "11:00"},
						Tuesday: []string{"14:00", "15:00"},
					},
				},
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Doctor
			err := Fill(&got, tt.input, true)
			if err != nil {
				t.Errorf("Fill() error = %v", err)
				return
			}
			t.Log(got)
			// Compare results
			if got.Name != tt.expected.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.expected.Name)
			}

			// Compare Monday slots
			if !reflect.DeepEqual((*got.WeekTimeslots)[0].Monday, (*tt.expected.WeekTimeslots)[0].Monday) {
				t.Errorf("Monday slots = %v, want %v", (*got.WeekTimeslots)[0].Monday, (*tt.expected.WeekTimeslots)[0].Monday)
			}

			// Compare Tuesday slots
			if !reflect.DeepEqual((*got.WeekTimeslots)[0].Tuesday, (*tt.expected.WeekTimeslots)[0].Tuesday) {
				t.Errorf("Tuesday slots = %v, want %v", (*got.WeekTimeslots)[0].Tuesday, (*tt.expected.WeekTimeslots)[0].Tuesday)
			}
		})
	}
}

func TestFillNestedFieldsStruct(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		input    []KV
		expected DoctorS
	}{
		{
			name: "nested week timeslots",
			input: []KV{
				{"name", "Dr. Smith"},
				{"week_timeslots.monday", "09:00,10:00,11:00"},
				{"week_timeslots.tuesday", "14:00,15:00"},
			},
			expected: DoctorS{
				Name: "Dr. Smith",
				WeekTimeslots: &WeekTimeslots{
					Monday:  []string{"09:00", "10:00", "11:00"},
					Tuesday: []string{"14:00", "15:00"},
				},
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got DoctorS
			err := Fill(&got, tt.input, true)
			if err != nil {
				t.Errorf("Fill() error = %v", err)
				return
			}
			t.Log(got)
			// Compare results
			if got.Name != tt.expected.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.expected.Name)
				return
			}

			// if got.WeekTimeslots == nil {
			// 	t.Error("WeekTimeslots is nil")
			// 	return
			// }

			// Compare Monday slots
			if !reflect.DeepEqual(got.WeekTimeslots.Monday, tt.expected.WeekTimeslots.Monday) {
				t.Errorf("Monday slots = %v, want %v", got.WeekTimeslots.Monday, tt.expected.WeekTimeslots.Monday)
			}

			// Compare Tuesday slots
			if !reflect.DeepEqual(got.WeekTimeslots.Tuesday, tt.expected.WeekTimeslots.Tuesday) {
				t.Errorf("Tuesday slots = %v, want %v", got.WeekTimeslots.Tuesday, tt.expected.WeekTimeslots.Tuesday)
			}
		})
	}
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
