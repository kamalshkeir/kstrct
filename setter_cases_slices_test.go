package kstrct

import (
	"reflect"
	"testing"
	"time"
)

func TestSliceFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		target   interface{}
		expected interface{}
	}{
		{
			name:   "slice of maps basic",
			input:  "name:John,age:25",
			target: &[]map[string]string{},
			expected: []map[string]string{
				{"name": "John", "age": "25"},
			},
		},
		{
			name:   "slice of maps with quotes",
			input:  `name:"John Doe",description:'Software Engineer'`,
			target: &[]map[string]string{},
			expected: []map[string]string{
				{"name": "John Doe", "description": "Software Engineer"},
			},
		},
		{
			name:  "slice of structs basic",
			input: "name:John,age:25",
			target: &[]struct {
				Name string
				Age  int
			}{},
			expected: []struct {
				Name string
				Age  int
			}{
				{Name: "John", Age: 25},
			},
		},
		{
			name:  "slice of structs with nested fields",
			input: "name:John,address.street:Main St,address.number:123",
			target: &[]struct {
				Name    string
				Address struct {
					Street string
					Number int
				}
			}{},
			expected: []struct {
				Name    string
				Address struct {
					Street string
					Number int
				}
			}{
				{
					Name: "John",
					Address: struct {
						Street string
						Number int
					}{
						Street: "Main St",
						Number: 123,
					},
				},
			},
		},
		{
			name:   "slice of maps with mixed types",
			input:  "name:John,age:25,active:true,salary:50000.50",
			target: &[]map[string]interface{}{},
			expected: []map[string]interface{}{
				{
					"name":   "John",
					"age":    25,
					"active": true,
					"salary": 50000.50,
				},
			},
		},
		{
			name:   "slice of maps with uint keys",
			input:  "1:first,2:second,3:third",
			target: &[]map[uint]string{},
			expected: []map[uint]string{
				{
					1: "first",
					2: "second",
					3: "third",
				},
			},
		},
		{
			name:   "slice of maps with uint8 keys",
			input:  "1:one,2:two,255:max",
			target: &[]map[uint8]string{},
			expected: []map[uint8]string{
				{
					1:   "one",
					2:   "two",
					255: "max",
				},
			},
		},
		{
			name:   "slice of maps with time.Time keys",
			input:  "2023-01-01:new year,2023-12-25:christmas",
			target: &[]map[time.Time]string{},
			expected: []map[time.Time]string{
				{
					time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC):   "new year",
					time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC): "christmas",
				},
			},
		},
		{
			name:   "slice of maps with mixed value types and uint keys",
			input:  "1:first,2:true,3:123,4:45.67",
			target: &[]map[uint]interface{}{},
			expected: []map[uint]interface{}{
				{
					uint(1): "first",
					uint(2): true,
					uint(3): 123,
					uint(4): 45.67,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetValue := reflect.ValueOf(tt.target).Elem()
			err := SetRFValue(targetValue, tt.input)
			if err != nil {
				t.Errorf("SetRFValue() error = %v", err)
				return
			}

			// Compare the results
			got := targetValue.Interface()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("SetRFValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSliceFromStringEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		target      interface{}
		shouldError bool
	}{
		{
			name:        "empty string",
			input:       "",
			target:      &[]map[string]string{},
			shouldError: false,
		},
		{
			name:        "invalid format",
			input:       "invalid:format:here",
			target:      &[]map[string]string{},
			shouldError: true,
		},
		{
			name:        "missing value",
			input:       "key:",
			target:      &[]map[string]string{},
			shouldError: true,
		},
		{
			name:        "missing key",
			input:       ":value",
			target:      &[]map[string]string{},
			shouldError: true,
		},
		{
			name:  "invalid nested path",
			input: "address..street:Main St",
			target: &[]struct {
				Address struct {
					Street string
				}
			}{},
			shouldError: true,
		},
		{
			name:        "invalid uint value",
			input:       "invalid:value",
			target:      &[]map[uint]string{},
			shouldError: true,
		},
		{
			name:        "uint8 overflow",
			input:       "256:overflow",
			target:      &[]map[uint8]string{},
			shouldError: true,
		},
		{
			name:        "invalid time format",
			input:       "invalid-date:value",
			target:      &[]map[time.Time]string{},
			shouldError: true,
		},
		{
			name:        "mixed valid and invalid uint keys",
			input:       "1:valid,invalid:value",
			target:      &[]map[uint]string{},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetValue := reflect.ValueOf(tt.target).Elem()
			err := SetRFValue(targetValue, tt.input)
			if (err != nil) != tt.shouldError {
				t.Errorf("SetRFValue() error = %v, shouldError = %v", err, tt.shouldError)
			}
		})
	}
}
