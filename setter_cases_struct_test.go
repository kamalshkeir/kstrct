package kstrct

// type TestBasicStruct struct {
// 	String    string
// 	Int       int
// 	Uint      uint
// 	Float     float64
// 	Bool      bool
// 	Time      time.Time
// 	StringPtr *string
// 	IntPtr    *int
// 	UintPtr   *uint
// 	FloatPtr  *float64
// 	BoolPtr   *bool
// 	TimePtr   *time.Time
// }

// type TestNestedStruct struct {
// 	Name     string
// 	Basic    TestBasicStruct
// 	BasicPtr *TestBasicStruct
// 	Map      map[string]interface{}
// 	Slice    []string
// }

// type TestComplexStruct struct {
// 	Name     string
// 	Nested   TestNestedStruct
// 	MapSlice []map[string]interface{}
// 	TimeMap  map[time.Time]string
// 	UintMap  map[uint]string
// }

// func TestStructFromKV(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		target   interface{}
// 		input    []KV
// 		expected interface{}
// 	}{
// 		{
// 			name:   "basic struct fields",
// 			target: &TestBasicStruct{},
// 			input: []KV{
// 				{Key: "string", Value: "test"},
// 				{Key: "int", Value: 42},
// 				{Key: "uint", Value: uint(123)},
// 				{Key: "float", Value: 3.14},
// 				{Key: "bool", Value: true},
// 				{Key: "time", Value: "2023-01-01"},
// 			},
// 			expected: &TestBasicStruct{
// 				String: "test",
// 				Int:    42,
// 				Uint:   123,
// 				Float:  3.14,
// 				Bool:   true,
// 				Time:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
// 			},
// 		},
// 		{
// 			name:   "pointer fields",
// 			target: &TestBasicStruct{},
// 			input: []KV{
// 				{Key: "string_ptr", Value: "test"},
// 				{Key: "int_ptr", Value: 42},
// 				{Key: "uint_ptr", Value: uint(123)},
// 				{Key: "float_ptr", Value: 3.14},
// 				{Key: "bool_ptr", Value: true},
// 				{Key: "time_ptr", Value: "2023-01-01"},
// 			},
// 			expected: &TestBasicStruct{
// 				StringPtr: strPtr("test"),
// 				IntPtr:    intPtr(42),
// 				UintPtr:   uintPtr(123),
// 				FloatPtr:  float64Ptr(3.14),
// 				BoolPtr:   boolPtr(true),
// 				TimePtr:   timePtr(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
// 			},
// 		},
// 		{
// 			name: "nested struct",
// 			target: &TestNestedStruct{
// 				Map: make(map[string]interface{}),
// 			},
// 			input: []KV{
// 				{Key: "name", Value: "test"},
// 				{Key: "basic.string", Value: "nested"},
// 				{Key: "basic.int", Value: 42},
// 				{Key: "map", Value: map[string]interface{}{
// 					"key1": "value1",
// 					"key2": "123",
// 				}},
// 				{Key: "slice", Value: []string{"a", "b", "c"}},
// 			},
// 			expected: &TestNestedStruct{
// 				Name: "test",
// 				Basic: TestBasicStruct{
// 					String: "nested",
// 					Int:    42,
// 				},
// 				Map: map[string]interface{}{
// 					"key1": "value1",
// 					"key2": "123",
// 				},
// 				Slice: []string{"a", "b", "c"},
// 			},
// 		},
// 		{
// 			name:   "nested pointer struct",
// 			target: &TestNestedStruct{},
// 			input: []KV{
// 				{Key: "basic_ptr.string", Value: "nested ptr"},
// 				{Key: "basic_ptr.int", Value: 42},
// 			},
// 			expected: &TestNestedStruct{
// 				BasicPtr: &TestBasicStruct{
// 					String: "nested ptr",
// 					Int:    42,
// 				},
// 			},
// 		},
// 		{
// 			name: "complex struct",
// 			target: &TestComplexStruct{
// 				TimeMap: make(map[time.Time]string),
// 				UintMap: make(map[uint]string),
// 			},
// 			input: []KV{
// 				{Key: "name", Value: "complex"},
// 				{Key: "nested.name", Value: "nested"},
// 				{Key: "nested.basic.string", Value: "deep nested"},
// 				{Key: "map_slice", Value: []map[string]interface{}{
// 					{"key1": "value1"},
// 					{"key2": "123"},
// 				}},
// 				{Key: "time_map", Value: map[time.Time]string{
// 					time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC): "new year",
// 				}},
// 				{Key: "uint_map", Value: map[uint]string{
// 					1: "one",
// 				}},
// 			},
// 			expected: &TestComplexStruct{
// 				Name: "complex",
// 				Nested: TestNestedStruct{
// 					Name: "nested",
// 					Basic: TestBasicStruct{
// 						String: "deep nested",
// 					},
// 				},
// 				MapSlice: []map[string]interface{}{
// 					{"key1": "value1"},
// 					{"key2": "123"},
// 				},
// 				TimeMap: map[time.Time]string{
// 					time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC): "new year",
// 				},
// 				UintMap: map[uint]string{
// 					1: "one",
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			for _, kv := range tt.input {
// 				err := SetRFValue(reflect.ValueOf(tt.target).Elem(), kv)
// 				if err != nil {
// 					t.Errorf("SetRFValue() error = %v", err)
// 					return
// 				}
// 			}

// 			if !reflect.DeepEqual(tt.target, tt.expected) {
// 				t.Errorf("SetRFValue() = %v, want %v", tt.target, tt.expected)
// 			}
// 		})
// 	}
// }

// func TestStructEdgeCases(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		target      interface{}
// 		input       KV
// 		shouldError bool
// 	}{
// 		{
// 			name:   "invalid field name",
// 			target: &TestBasicStruct{},
// 			input: KV{
// 				Key:   "nonexistent",
// 				Value: "test",
// 			},
// 			shouldError: true,
// 		},
// 		{
// 			name:   "invalid nested field",
// 			target: &TestNestedStruct{},
// 			input: KV{
// 				Key:   "basic.nonexistent",
// 				Value: "test",
// 			},
// 			shouldError: true,
// 		},
// 		{
// 			name:   "invalid type conversion",
// 			target: &TestBasicStruct{},
// 			input: KV{
// 				Key:   "int",
// 				Value: "not a number",
// 			},
// 			shouldError: true,
// 		},
// 		{
// 			name:   "invalid time format",
// 			target: &TestBasicStruct{},
// 			input: KV{
// 				Key:   "time",
// 				Value: "invalid date",
// 			},
// 			shouldError: true,
// 		},
// 		{
// 			name:   "nil pointer struct",
// 			target: &TestNestedStruct{},
// 			input: KV{
// 				Key:   "basic_ptr.string",
// 				Value: "test",
// 			},
// 			shouldError: false,
// 		},
// 		{
// 			name: "invalid map key type",
// 			target: &TestComplexStruct{
// 				UintMap: make(map[uint]string),
// 			},
// 			input: KV{
// 				Key: "uint_map",
// 				Value: map[string]string{
// 					"invalid": "test",
// 				},
// 			},
// 			shouldError: true,
// 		},
// 		{
// 			name: "invalid slice index",
// 			target: &TestComplexStruct{
// 				MapSlice: make([]map[string]interface{}, 0),
// 			},
// 			input: KV{
// 				Key:   "map_slice.invalid.key",
// 				Value: "test",
// 			},
// 			shouldError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := SetRFValue(reflect.ValueOf(tt.target).Elem(), tt.input)
// 			if (err != nil) != tt.shouldError {
// 				t.Errorf("SetRFValue() error = %v, shouldError = %v", err, tt.shouldError)
// 			}
// 		})
// 	}
// }

// // Helper functions for creating pointers
// func uintPtr(u uint) *uint           { return &u }
// func float64Ptr(f float64) *float64  { return &f }
// func boolPtr(b bool) *bool           { return &b }
// func timePtr(t time.Time) *time.Time { return &t }
