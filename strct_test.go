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

type WeekTimeslotss struct {
	Id         uint      `korm:"pk" json:"id,omitempty"`
	DoctorId   uint      `korm:"fk:doctors.id:cascade:cascade" json:"doctor_id,omitempty"`
	Sunday     []string  `json:"sunday,omitempty"`
	Monday     []string  `json:"monday,omitempty"`
	Tuesday    []string  `json:"tuesday,omitempty"`
	Wednesday  []string  `json:"wednesday,omitempty"`
	Thursday   []string  `json:"thursday,omitempty"`
	Friday     []string  `json:"friday,omitempty"`
	Saturday   []string  `json:"saturday,omitempty"`
	Indx       uint8     `json:"indx,omitempty"`
	LastUpdate time.Time `korm:"update" json:"-"`
}

type Reservation struct {
	Id        uint       `korm:"pk" json:"id,omitempty"`
	DoctorId  uint       `korm:"fk:doctors.id:cascade:cascade" json:"doctor_id,omitempty"`
	PatientId uint       `korm:"fk:patients.id:cascade:cascade" json:"patient_id,omitempty"`
	Day       uint8      `korm:"check: day >= 0 AND day <= 6" json:"day,omitempty"`
	Timeslot  string     `json:"timeslot,omitempty"`
	IsVisio   bool       `json:"is_visio,omitempty"`
	VisioLink string     `json:"visio_link,omitempty"`
	Motif     string     `json:"motif,omitempty" korm:"text"`
	Date      *time.Time `json:"date,omitempty"`
	UpdatedAt time.Time  `korm:"update" json:"updated_at,omitempty"`
	CreatedAt time.Time  `korm:"now" json:"-"`
}

type Doctor struct {
	Name          string
	WeekTimeslots *[]WeekTimeslots
}

type DoctorS struct {
	Name          string
	WeekTimeslots *WeekTimeslots
}

type Docto struct {
	Id               uint           `korm:"pk" json:"id,omitempty"`
	Uuid             string         `korm:"size:40" json:"uuid,omitempty"`
	Email            string         `korm:"iunique" json:"email,omitempty"`
	Number           string         `korm:"iunique" json:"number,omitempty"`
	ExtraNumber      string         `korm:"iunique" json:"extra_number,omitempty"`
	Password         string         `json:"-"`
	Name             string         `json:"name,omitempty"`
	Slug             string         `json:"slug,omitempty"`
	Address          string         `json:"address,omitempty"`
	ExtraAddress     string         `json:"extra_address,omitempty"`
	Prefix           string         `korm:"size:50" json:"prefix,omitempty"`
	City             string         `korm:"fk:cities.name:cascade:cascade" json:"city,omitempty"`
	Image            string         `json:"image,omitempty"`
	Kind             string         `json:"kind,omitempty"`
	Speciality       string         `korm:"fk:specialities.name:cascade:cascade" json:"speciality,omitempty"`
	Link             string         `json:"link,omitempty"`
	RegulationSector string         `json:"regulation_sector,omitempty"`
	Description      string         `json:"description,omitempty"`
	ExtraInfos       string         `json:"extra_infos,omitempty"`
	Languages        []string       `json:"languages,omitempty"`
	IsVisio          bool           `json:"is_visio"`
	IsAvailable      bool           `json:"is_available"`
	BackAt           *time.Time     `json:"back_at,omitempty"`
	IsBlocked        bool           `json:"-"`
	AcceptFirst      bool           `json:"accept_first"`
	Latitude         float64        `json:"latitude,omitempty"`
	Longitude        float64        `json:"longitude,omitempty"`
	WeekTimeslots    WeekTimeslotss `json:"week_timeslots,omitempty"`
	Reservations     []Reservation  `json:"reservations,omitempty"`
	VisitTypes       []string       `json:"visit_types"`
	CreatedAt        time.Time      `korm:"now" json:"-"`
}

func TestFillDocto(t *testing.T) {
	u := &Docto{}
	err := Fill(u, []KV{
		{Key: "uuid", Value: "xxx-xxx-xxx"},
		{Key: "name", Value: "kamal"},
		{Key: "languages", Value: "fr,en,es"},
		{Key: "back_at", Value: time.Now()},
		{Key: "is_blocked", Value: true},
		{Key: "week_timeslots.sunday", Value: "8:00,9:00,10:00"},
		{Key: "week_timeslots.monday", Value: "10:00,11:00,12:00"},
		{Key: "reservations.id", Value: 1},
		{Key: "reservations.patient_id", Value: 12345},
		{Key: "visit_types", Value: "bla,bla2,bla3,bla4"},
	})
	if err != nil {
		t.Error(err)
		return
	}
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
			err := Fill(&got, tt.input)
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
			err := Fill(&got, tt.input)
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

func TestCreateStruct(t *testing.T) {
	// Define fields for our dynamic struct
	fields := []StructField{
		{
			Name: "ID",
			Type: reflect.TypeOf(0),
			Tags: map[string]string{
				"json": "id,omitempty",
				"korm": "pk",
			},
			Value: 1,
		},
		{
			Name: "Name",
			Type: reflect.TypeOf(""),
			Tags: map[string]string{
				"json": "name,omitempty",
			},
			Value: "John Doe",
		},
		{
			Name: "CreatedAt",
			Type: reflect.TypeOf(time.Time{}),
			Tags: map[string]string{
				"json": "-",
				"korm": "update",
			},
			Value: time.Now(),
		},
	}

	// Create the struct
	dynamicStruct, err := CreateStruct(fields)
	if err != nil {
		t.Fatalf("Failed to create struct: %v", err)
	}

	// Verify the struct was created correctly
	val := reflect.ValueOf(dynamicStruct).Elem()
	typ := val.Type()

	// Check number of fields
	if typ.NumField() != len(fields) {
		t.Errorf("Expected %d fields, got %d", len(fields), typ.NumField())
	}

	// Check field types and tags
	for i, field := range fields {
		structField := typ.Field(i)

		// Check name
		if structField.Name != field.Name {
			t.Errorf("Field %d: expected name %s, got %s", i, field.Name, structField.Name)
		}

		// Check type
		if structField.Type != field.Type {
			t.Errorf("Field %d: expected type %v, got %v", i, field.Type, structField.Type)
		}

		// Check tags
		for key, expectedValue := range field.Tags {
			if gotValue, ok := structField.Tag.Lookup(key); !ok || gotValue != expectedValue {
				t.Errorf("Field %d: expected tag %s:\"%s\", got %s:\"%s\"", i, key, expectedValue, key, gotValue)
			}
		}

		// Check value
		if field.Value != nil {
			fieldValue := val.Field(i).Interface()
			if !reflect.DeepEqual(fieldValue, field.Value) {
				t.Errorf("Field %d: expected value %v, got %v", i, field.Value, fieldValue)
			}
		}
	}
}

func TestExtractStructFields(t *testing.T) {
	// Create a test struct
	type TestStruct struct {
		ID        uint      `json:"id,omitempty" korm:"pk"`
		Name      string    `json:"name,omitempty"`
		CreatedAt time.Time `json:"-" korm:"update"`
	}

	testValue := TestStruct{
		ID:        1,
		Name:      "Test",
		CreatedAt: time.Now(),
	}

	// Extract fields
	fields, err := ExtractStructFields(&testValue)
	if err != nil {
		t.Fatalf("Failed to extract struct fields: %v", err)
	}

	// Verify number of fields
	if len(fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(fields))
	}

	// Verify field details
	expectedFields := map[string]struct {
		fieldType reflect.Type
		tags      map[string]string
		value     any
	}{
		"ID": {
			fieldType: reflect.TypeOf(uint(0)),
			tags: map[string]string{
				"json": "id,omitempty",
				"korm": "pk",
			},
			value: uint(1),
		},
		"Name": {
			fieldType: reflect.TypeOf(""),
			tags: map[string]string{
				"json": "name,omitempty",
			},
			value: "Test",
		},
		"CreatedAt": {
			fieldType: reflect.TypeOf(time.Time{}),
			tags: map[string]string{
				"json": "-",
				"korm": "update",
			},
			value: testValue.CreatedAt,
		},
	}

	for _, field := range fields {
		expected, ok := expectedFields[field.Name]
		if !ok {
			t.Errorf("Unexpected field: %s", field.Name)
			continue
		}

		// Check type
		if field.Type != expected.fieldType {
			t.Errorf("Field %s: expected type %v, got %v", field.Name, expected.fieldType, field.Type)
		}

		// Check tags
		for key, expectedValue := range expected.tags {
			if gotValue, ok := field.Tags[key]; !ok || gotValue != expectedValue {
				t.Errorf("Field %s: expected tag %s:\"%s\", got %s:\"%s\"", field.Name, key, expectedValue, key, gotValue)
			}
		}

		// Check value
		if !reflect.DeepEqual(field.Value, expected.value) {
			t.Errorf("Field %s: expected value %v, got %v", field.Name, expected.value, field.Value)
		}
	}

	// Test error case - non-struct input
	_, err = ExtractStructFields(42)
	if err == nil {
		t.Error("Expected error when passing non-struct value, got nil")
	}
}
