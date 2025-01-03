package kstrct

import (
	"reflect"
	"testing"
	"time"
)

func TestMapFromString(t *testing.T) {
	tests := []struct {
		name     string
		target   interface{}
		input    string
		expected interface{}
	}{
		{
			name:     "map_string_string_basic",
			target:   &map[string]string{},
			input:    "name:John,age:25",
			expected: map[string]string{"name": "John", "age": "25"},
		},
		{
			name:     "map_string_string_with_quotes",
			target:   &map[string]string{},
			input:    `name:"John Doe",description:'Software Engineer'`,
			expected: map[string]string{"name": "John Doe", "description": "Software Engineer"},
		},
		{
			name:     "map_uint_string",
			target:   &map[uint]string{},
			input:    "1:first,2:second,3:third",
			expected: map[uint]string{1: "first", 2: "second", 3: "third"},
		},
		{
			name:     "map_uint8_string",
			target:   &map[uint8]string{},
			input:    "1:one,2:two,255:max",
			expected: map[uint8]string{1: "one", 2: "two", 255: "max"},
		},
		{
			name:     "map_string_interface",
			target:   &map[string]interface{}{},
			input:    "name:John,age:25,active:true,salary:50000.50",
			expected: map[string]interface{}{"name": "John", "age": int64(25), "active": true, "salary": 50000.50},
		},
		{
			name:     "map_string_int",
			target:   &map[string]int{},
			input:    "one:1,two:2,three:3",
			expected: map[string]int{"one": 1, "two": 2, "three": 3},
		},
		{
			name:     "map_string_float64",
			target:   &map[string]float64{},
			input:    "pi:3.14,e:2.718",
			expected: map[string]float64{"pi": 3.14, "e": 2.718},
		},
		{
			name:     "map_string_bool",
			target:   &map[string]bool{},
			input:    "valid:true,active:false",
			expected: map[string]bool{"valid": true, "active": false},
		},
		{
			name:   "map_time_string",
			target: &map[time.Time]string{},
			input:  "2023-01-01:new year,2023-12-25:christmas",
			expected: map[time.Time]string{
				mustParseTime("2023-01-01"): "new year",
				mustParseTime("2023-12-25"): "christmas",
			},
		},
		{
			name:   "map_uint_interface",
			target: &map[uint]interface{}{},
			input:  "1:first,2:true,3:123,4:45.67",
			expected: map[uint]interface{}{
				1: "first",
				2: true,
				3: int64(123),
				4: 45.67,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetRFValue(reflect.ValueOf(tt.target).Elem(), tt.input)
			if err != nil {
				t.Errorf("SetRFValue() error = %v", err)
				return
			}

			got := reflect.ValueOf(tt.target).Elem().Interface()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SetRFValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMapFromStringEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		target      interface{}
		input       string
		expectError bool
	}{
		{
			name:        "empty_string",
			target:      &map[string]string{},
			input:       "",
			expectError: false,
		},
		{
			name:        "invalid_format",
			target:      &map[string]string{},
			input:       "invalid:format:here",
			expectError: false, // Should skip invalid entries
		},
		{
			name:        "missing_value",
			target:      &map[string]string{},
			input:       "key:",
			expectError: false,
		},
		{
			name:        "missing_key",
			target:      &map[string]string{},
			input:       ":value",
			expectError: false,
		},
		{
			name:        "invalid_uint_value",
			target:      &map[uint]string{},
			input:       "invalid:value",
			expectError: false, // Should skip invalid entries
		},
		{
			name:        "uint8_overflow",
			target:      &map[uint8]string{},
			input:       "256:overflow",
			expectError: false, // Should skip invalid entries
		},
		{
			name:        "invalid_time_format",
			target:      &map[time.Time]string{},
			input:       "invalid-date:value",
			expectError: false, // Should skip invalid entries
		},
		{
			name:        "mixed_valid_and_invalid_uint_keys",
			target:      &map[uint]string{},
			input:       "1:valid,invalid:value",
			expectError: false, // Should keep valid entries and skip invalid ones
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetRFValue(reflect.ValueOf(tt.target).Elem(), tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("SetRFValue() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestMapFromMap(t *testing.T) {
	tests := []struct {
		name     string
		target   interface{}
		input    interface{}
		expected interface{}
	}{
		{
			name:     "map_string_any_to_map_string_string",
			target:   &map[string]string{},
			input:    map[string]interface{}{"name": "John", "age": "25"},
			expected: map[string]string{"name": "John", "age": "25"},
		},
		{
			name:     "map_string_string_to_map_string_string",
			target:   &map[string]string{},
			input:    map[string]string{"name": "John", "age": "25"},
			expected: map[string]string{"name": "John", "age": "25"},
		},
		{
			name:     "map_string_int_to_map_string_string",
			target:   &map[string]string{},
			input:    map[string]int{"age": 25, "id": 123},
			expected: map[string]string{"age": "25", "id": "123"},
		},
		{
			name:     "map_string_float_to_map_string_interface",
			target:   &map[string]interface{}{},
			input:    map[string]float64{"pi": 3.14, "e": 2.718},
			expected: map[string]interface{}{"pi": 3.14, "e": 2.718},
		},
		{
			name:   "map_time_string_to_map_time_string",
			target: &map[time.Time]string{},
			input: map[time.Time]string{
				mustParseTime("2023-01-01"): "new year",
				mustParseTime("2023-12-25"): "christmas",
			},
			expected: map[time.Time]string{
				mustParseTime("2023-01-01"): "new year",
				mustParseTime("2023-12-25"): "christmas",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetRFValue(reflect.ValueOf(tt.target).Elem(), tt.input)
			if err != nil {
				t.Errorf("SetRFValue() error = %v", err)
				return
			}

			got := reflect.ValueOf(tt.target).Elem().Interface()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SetRFValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Helper function to parse time strings
func mustParseTime(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}
