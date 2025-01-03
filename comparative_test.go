package kstrct

// // Helper functions
// func strPtr(s string) *string {
// 	return &s
// }

// func intPtr(i int) *int {
// 	return &i
// }

// // Test structures
// type ComparativeTestStruct struct {
// 	// Basic types
// 	String    string
// 	Int       int
// 	Float64   float64
// 	Bool      bool
// 	Time      time.Time
// 	StringPtr *string
// 	IntPtr    *int

// 	// Slices
// 	StringSlice []string
// 	IntSlice    []int
// 	FloatSlice  []float64

// 	// Maps
// 	StringMap map[string]string
// 	IntMap    map[string]int
// 	FloatMap  map[string]float64

// 	// Nested struct
// 	Nested struct {
// 		Name   string
// 		Value  int
// 		Active bool
// 	}

// 	// Pointer to nested struct
// 	NestedPtr *struct {
// 		Name   string
// 		Value  int
// 		Active bool
// 	}

// 	// Slice of structs
// 	StructSlice []struct {
// 		Name   string
// 		Value  int
// 		Active bool
// 	}

// 	// SQL types
// 	NullString  sql.NullString
// 	NullInt64   sql.NullInt64
// 	NullFloat64 sql.NullFloat64
// 	NullBool    sql.NullBool
// 	NullTime    sql.NullTime

// 	// Deep nesting
// 	DeepNested struct {
// 		Level1 struct {
// 			Level2 struct {
// 				Name   string
// 				Value  int
// 				Active bool
// 			}
// 		}
// 	}

// 	// Map with complex values
// 	ComplexMap map[string]struct {
// 		Name   string
// 		Value  int
// 		Active bool
// 	}
// }

// func TestComparativeFunctionality(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		input    map[string]any // Base input that will be converted for each function
// 		expected ComparativeTestStruct
// 	}{
// 		{
// 			name: "basic types",
// 			input: map[string]any{
// 				"string":  "test",
// 				"int":     42,
// 				"float64": 3.14,
// 				"bool":    true,
// 				"time":    "2023-01-02 15:04:05",
// 			},
// 			expected: ComparativeTestStruct{
// 				String:  "test",
// 				Int:     42,
// 				Float64: 3.14,
// 				Bool:    true,
// 				Time:    time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC),
// 			},
// 		},
// 		{
// 			name: "pointer types",
// 			input: map[string]any{
// 				"string_ptr": "pointer test",
// 				"int_ptr":    123,
// 			},
// 			expected: ComparativeTestStruct{
// 				StringPtr: strPtr("pointer test"),
// 				IntPtr:    intPtr(123),
// 			},
// 		},
// 		{
// 			name: "slice types",
// 			input: map[string]any{
// 				"string_slice": []string{"a", "b", "c"},
// 				"int_slice":    []int{1, 2, 3},
// 				"float_slice":  []float64{1.1, 2.2, 3.3},
// 			},
// 			expected: ComparativeTestStruct{
// 				StringSlice: []string{"a", "b", "c"},
// 				IntSlice:    []int{1, 2, 3},
// 				FloatSlice:  []float64{1.1, 2.2, 3.3},
// 			},
// 		},
// 		{
// 			name: "map types",
// 			input: map[string]any{
// 				"string_map": map[string]string{"key1": "value1", "key2": "value2"},
// 				"int_map":    map[string]int{"key1": 1, "key2": 2},
// 				"float_map":  map[string]float64{"key1": 1.1, "key2": 2.2},
// 			},
// 			expected: ComparativeTestStruct{
// 				StringMap: map[string]string{"key1": "value1", "key2": "value2"},
// 				IntMap:    map[string]int{"key1": 1, "key2": 2},
// 				FloatMap:  map[string]float64{"key1": 1.1, "key2": 2.2},
// 			},
// 		},
// 		{
// 			name: "nested struct",
// 			input: map[string]any{
// 				"nested.name":   "nested test",
// 				"nested.value":  42,
// 				"nested.active": true,
// 			},
// 			expected: ComparativeTestStruct{
// 				Nested: struct {
// 					Name   string
// 					Value  int
// 					Active bool
// 				}{
// 					Name:   "nested test",
// 					Value:  42,
// 					Active: true,
// 				},
// 			},
// 		},
// 		{
// 			name: "nested pointer struct",
// 			input: map[string]any{
// 				"nested_ptr.name":   "ptr test",
// 				"nested_ptr.value":  24,
// 				"nested_ptr.active": true,
// 			},
// 			expected: ComparativeTestStruct{
// 				NestedPtr: &struct {
// 					Name   string
// 					Value  int
// 					Active bool
// 				}{
// 					Name:   "ptr test",
// 					Value:  24,
// 					Active: true,
// 				},
// 			},
// 		},
// 		{
// 			name: "slice of structs",
// 			input: map[string]any{
// 				"struct_slice.0.name":   "first",
// 				"struct_slice.0.value":  1,
// 				"struct_slice.0.active": true,
// 				"struct_slice.1.name":   "second",
// 				"struct_slice.1.value":  2,
// 				"struct_slice.1.active": false,
// 			},
// 			expected: ComparativeTestStruct{
// 				StructSlice: []struct {
// 					Name   string
// 					Value  int
// 					Active bool
// 				}{
// 					{Name: "first", Value: 1, Active: true},
// 					{Name: "second", Value: 2, Active: false},
// 				},
// 			},
// 		},
// 		{
// 			name: "sql null types",
// 			input: map[string]any{
// 				"null_string":  "null test",
// 				"null_int64":   int64(42),
// 				"null_float64": 3.14,
// 				"null_bool":    true,
// 				"null_time":    "2023-01-02 15:04:05",
// 			},
// 			expected: ComparativeTestStruct{
// 				NullString:  sql.NullString{String: "null test", Valid: true},
// 				NullInt64:   sql.NullInt64{Int64: 42, Valid: true},
// 				NullFloat64: sql.NullFloat64{Float64: 3.14, Valid: true},
// 				NullBool:    sql.NullBool{Bool: true, Valid: true},
// 				NullTime:    sql.NullTime{Time: time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC), Valid: true},
// 			},
// 		},
// 		{
// 			name: "deep nested struct",
// 			input: map[string]any{
// 				"deep_nested.level1.level2.name":   "deep",
// 				"deep_nested.level1.level2.value":  42,
// 				"deep_nested.level1.level2.active": true,
// 			},
// 			expected: ComparativeTestStruct{
// 				DeepNested: struct {
// 					Level1 struct {
// 						Level2 struct {
// 							Name   string
// 							Value  int
// 							Active bool
// 						}
// 					}
// 				}{
// 					Level1: struct {
// 						Level2 struct {
// 							Name   string
// 							Value  int
// 							Active bool
// 						}
// 					}{
// 						Level2: struct {
// 							Name   string
// 							Value  int
// 							Active bool
// 						}{
// 							Name:   "deep",
// 							Value:  42,
// 							Active: true,
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			name: "complex map",
// 			input: map[string]any{
// 				"complex_map": map[string]any{
// 					"key1": map[string]any{
// 						"name":   "map struct 1",
// 						"value":  1,
// 						"active": true,
// 					},
// 					"key2": map[string]any{
// 						"name":   "map struct 2",
// 						"value":  2,
// 						"active": false,
// 					},
// 				},
// 			},
// 			expected: ComparativeTestStruct{
// 				ComplexMap: map[string]struct {
// 					Name   string
// 					Value  int
// 					Active bool
// 				}{
// 					"key1": {Name: "map struct 1", Value: 1, Active: true},
// 					"key2": {Name: "map struct 2", Value: 2, Active: false},
// 				},
// 			},
// 		},
// 		{
// 			name: "empty_values",
// 			input: map[string]any{
// 				"string":       "",
// 				"int":          0,
// 				"float64":      0.0,
// 				"bool":         false,
// 				"string_slice": []string{},
// 				"int_slice":    []int{},
// 				"float_slice":  []float64{},
// 				"string_map":   map[string]string{},
// 				"int_map":      map[string]int{},
// 				"float_map":    map[string]float64{},
// 			},
// 			expected: ComparativeTestStruct{
// 				String:      "",
// 				Int:         0,
// 				Float64:     0.0,
// 				Bool:        false,
// 				StringSlice: []string{},
// 				IntSlice:    []int{},
// 				FloatSlice:  []float64{},
// 				StringMap:   map[string]string{},
// 				IntMap:      map[string]int{},
// 				FloatMap:    map[string]float64{},
// 			},
// 		},
// 		{
// 			name: "null_sql_types",
// 			input: map[string]any{
// 				"null_string":  nil,
// 				"null_int64":   nil,
// 				"null_float64": nil,
// 				"null_bool":    nil,
// 				"null_time":    nil,
// 			},
// 			expected: ComparativeTestStruct{
// 				NullString:  sql.NullString{String: "", Valid: false},
// 				NullInt64:   sql.NullInt64{Int64: 0, Valid: false},
// 				NullFloat64: sql.NullFloat64{Float64: 0, Valid: false},
// 				NullBool:    sql.NullBool{Bool: false, Valid: false},
// 				NullTime:    sql.NullTime{Time: time.Time{}, Valid: false},
// 			},
// 		},
// 		{
// 			name: "type_conversions",
// 			input: map[string]any{
// 				"int":     42.5,  // float64 to int
// 				"float64": 42,    // int to float64
// 				"string":  123,   // number to string
// 				"int_ptr": "456", // string to number pointer
// 			},
// 			expected: ComparativeTestStruct{
// 				Int:     42,
// 				Float64: 42.0,
// 				String:  "123",
// 				IntPtr:  intPtr(456),
// 			},
// 		},
// 		{
// 			name: "special_chars",
// 			input: map[string]any{
// 				"string":      "Hello 世界 🌍",
// 				"nested.name": "Test 测试",
// 			},
// 			expected: ComparativeTestStruct{
// 				String: "Hello 世界 🌍",
// 				Nested: struct {
// 					Name   string
// 					Value  int
// 					Active bool
// 				}{
// 					Name: "Test 测试",
// 				},
// 			},
// 		},
// 		{
// 			name: "alternative_time_formats",
// 			input: map[string]any{
// 				"time":      "2023-01-02T15:04:05Z", // RFC3339
// 				"null_time": "2023-01-02 15:04:05",  // Standard format
// 			},
// 			expected: ComparativeTestStruct{
// 				Time: time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC),
// 				NullTime: sql.NullTime{
// 					Time:  time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC),
// 					Valid: true,
// 				},
// 			},
// 		},
// 		{
// 			name: "partial_nested",
// 			input: map[string]any{
// 				"nested.name":                    "partial",
// 				"deep_nested.level1.level2.name": "deep partial",
// 			},
// 			expected: ComparativeTestStruct{
// 				Nested: struct {
// 					Name   string
// 					Value  int
// 					Active bool
// 				}{
// 					Name: "partial",
// 				},
// 				DeepNested: struct {
// 					Level1 struct {
// 						Level2 struct {
// 							Name   string
// 							Value  int
// 							Active bool
// 						}
// 					}
// 				}{
// 					Level1: struct {
// 						Level2 struct {
// 							Name   string
// 							Value  int
// 							Active bool
// 						}
// 					}{
// 						Level2: struct {
// 							Name   string
// 							Value  int
// 							Active bool
// 						}{
// 							Name: "deep partial",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Test builder.FromMap
// 			t.Run("FromMap", func(t *testing.T) {
// 				var got ComparativeTestStruct
// 				builder := NewBuilder(&got)
// 				builder.FromMap(tt.input)
// 				assert.Equal(t, tt.expected, got)
// 			})

// 			// Test builder.FromKV
// 			t.Run("FromKV", func(t *testing.T) {
// 				var got ComparativeTestStruct
// 				builder := NewBuilder(&got)
// 				kvs := mapToKVs(tt.input)
// 				builder.FromKV(kvs...)
// 				assert.Equal(t, tt.expected, got)
// 			})

// 			// Test Fill
// 			t.Run("Fill", func(t *testing.T) {
// 				var got ComparativeTestStruct
// 				kvs := mapToKVs(tt.input)
// 				err := Fill(&got, kvs)
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.expected, got)
// 			})
// 		})
// 	}
// }

// func mapToKVs(m map[string]any) []KV {
// 	kvs := make([]KV, 0, len(m))
// 	for k, v := range m {
// 		kvs = append(kvs, KV{Key: k, Value: v})
// 	}
// 	return kvs
// }

// // Edge cases and error handling tests
// func TestComparativeEdgeCases(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		input         map[string]any
// 		expectError   bool
// 		errorContains string
// 	}{
// 		{
// 			name: "invalid number conversion",
// 			input: map[string]any{
// 				"int": "not a number",
// 			},
// 			expectError:   true,
// 			errorContains: "cannot convert string",
// 		},
// 		{
// 			name: "invalid bool conversion",
// 			input: map[string]any{
// 				"bool": "not a bool",
// 			},
// 			expectError:   true,
// 			errorContains: "cannot convert string",
// 		},
// 		{
// 			name: "invalid time format",
// 			input: map[string]any{
// 				"time": "invalid time",
// 			},
// 			expectError:   true,
// 			errorContains: "could not parse time",
// 		},
// 		{
// 			name: "nil pointer dereference",
// 			input: map[string]any{
// 				"nested_ptr.name": "test",
// 			},
// 			expectError: false, // Should handle nil pointer by creating new instance
// 		},
// 		{
// 			name: "invalid slice index",
// 			input: map[string]any{
// 				"struct_slice.-1.name": "negative index",
// 			},
// 			expectError:   true,
// 			errorContains: "negative slice index",
// 		},
// 		{
// 			name: "overflow_number",
// 			input: map[string]any{
// 				"int": "9999999999999999999999999999999999999999", // Number too large for int
// 			},
// 			expectError:   true,
// 			errorContains: "cannot convert",
// 		},
// 		{
// 			name: "invalid_field_type",
// 			input: map[string]any{
// 				"int": "not_a_number", // String that can't be converted to int
// 			},
// 			expectError:   true,
// 			errorContains: "cannot convert",
// 		},
// 		{
// 			name: "invalid_slice_format",
// 			input: map[string]any{
// 				"string_slice": struct{ foo string }{foo: "bar"}, // Cannot convert struct to string slice
// 			},
// 			expectError:   true,
// 			errorContains: "cannot assign value",
// 		},
// 		{
// 			name: "invalid_struct_field",
// 			input: map[string]any{
// 				"nested.invalid": "value", // Invalid nested field
// 			},
// 			expectError:   true,
// 			errorContains: "field invalid not found",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Test builder.FromMap
// 			t.Run("FromMap", func(t *testing.T) {
// 				var got ComparativeTestStruct
// 				builder := NewBuilder(&got)
// 				builder.FromMap(tt.input)
// 				err := builder.Error()
// 				if tt.expectError {
// 					assert.Error(t, err)
// 					if tt.errorContains != "" && err != nil {
// 						assert.Contains(t, err.Error(), tt.errorContains)
// 					}
// 				} else {
// 					assert.NoError(t, err)
// 				}
// 			})

// 			// Test builder.FromKV
// 			t.Run("FromKV", func(t *testing.T) {
// 				var got ComparativeTestStruct
// 				builder := NewBuilder(&got)
// 				kvs := mapToKVs(tt.input)
// 				builder.FromKV(kvs...)
// 				err := builder.Error()
// 				if tt.expectError {
// 					assert.Error(t, err)
// 					if tt.errorContains != "" && err != nil {
// 						assert.Contains(t, err.Error(), tt.errorContains)
// 					}
// 				} else {
// 					assert.NoError(t, err)
// 				}
// 			})

// 			// Test Fill
// 			t.Run("Fill", func(t *testing.T) {
// 				var got ComparativeTestStruct
// 				kvs := mapToKVs(tt.input)
// 				err := Fill(&got, kvs)
// 				if tt.expectError {
// 					assert.Error(t, err)
// 					if tt.errorContains != "" && err != nil {
// 						assert.Contains(t, err.Error(), tt.errorContains)
// 					}
// 				} else {
// 					assert.NoError(t, err)
// 				}
// 			})
// 		})
// 	}
// }

// // Performance comparison tests
// func BenchmarkComparative(b *testing.B) {
// 	input := map[string]any{
// 		"string":               "test",
// 		"int":                  42,
// 		"float64":              3.14,
// 		"bool":                 true,
// 		"time":                 "2023-01-02 15:04:05",
// 		"string_slice":         []string{"a", "b", "c"},
// 		"int_slice":            []int{1, 2, 3},
// 		"float_slice":          []float64{1.1, 2.2, 3.3},
// 		"nested.name":          "nested test",
// 		"nested.value":         42,
// 		"nested.active":        true,
// 		"nested_ptr.name":      "ptr test",
// 		"nested_ptr.value":     24,
// 		"nested_ptr.active":    true,
// 		"struct_slice.0.name":  "first",
// 		"struct_slice.0.value": 1,
// 		"struct_slice.1.name":  "second",
// 		"struct_slice.1.value": 2,
// 	}

// 	kvs := mapToKVs(input)

// 	b.Run("FromMap", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			var s ComparativeTestStruct
// 			builder := NewBuilder(&s)
// 			builder.FromMap(input)
// 		}
// 	})

// 	b.Run("FromKV", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			var s ComparativeTestStruct
// 			builder := NewBuilder(&s)
// 			builder.FromKV(kvs...)
// 		}
// 	})

// 	b.Run("Fill", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			var s ComparativeTestStruct
// 			Fill(&s, kvs)
// 		}
// 	})
// }

// type KormUser struct {
// 	ID        uint      `korm:"pk" json:"id,omitempty"`
// 	Email     string    `korm:"iunique" json:"email,omitempty"`
// 	Password  string    `json:"-"`
// 	Profile   *Profile  `json:"profile,omitempty"`
// 	Settings  Settings  `json:"settings,omitempty"`
// 	CreatedAt time.Time `korm:"now" json:"-"`
// }

// type Profile struct {
// 	FirstName string       `json:"first_name,omitempty"`
// 	LastName  string       `json:"last_name,omitempty"`
// 	Age       int          `json:"age,omitempty"`
// 	Address   *UserAddress `json:"address,omitempty"`
// }

// type UserAddress struct {
// 	Street  string `json:"street,omitempty"`
// 	City    string `json:"city,omitempty"`
// 	Country string `json:"country,omitempty"`
// }

// type Settings struct {
// 	Theme    string   `json:"theme,omitempty"`
// 	Language string   `json:"language,omitempty"`
// 	Tags     []string `json:"tags,omitempty"`
// }

// func TestKormStyleMapping(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		input    map[string]any
// 		expected KormUser
// 	}{
// 		{
// 			name: "korm_style_basic",
// 			input: map[string]any{
// 				"id":       uint(1),
// 				"email":    "test@example.com",
// 				"password": "hashed_password",
// 				"profile": map[string]any{
// 					"first_name": "John",
// 					"last_name":  "Doe",
// 					"age":        30,
// 					"address": map[string]any{
// 						"street":  "123 Main St",
// 						"city":    "New York",
// 						"country": "USA",
// 					},
// 				},
// 				"settings": map[string]any{
// 					"theme":    "dark",
// 					"language": "en",
// 					"tags":     []string{"user", "premium"},
// 				},
// 				"created_at": time.Now(),
// 			},
// 			expected: KormUser{
// 				ID:       1,
// 				Email:    "test@example.com",
// 				Password: "hashed_password",
// 				Profile: &Profile{
// 					FirstName: "John",
// 					LastName:  "Doe",
// 					Age:       30,
// 					Address: &UserAddress{
// 						Street:  "123 Main St",
// 						City:    "New York",
// 						Country: "USA",
// 					},
// 				},
// 				Settings: Settings{
// 					Theme:    "dark",
// 					Language: "en",
// 					Tags:     []string{"user", "premium"},
// 				},
// 			},
// 		},
// 		{
// 			name: "korm_style_dot_notation",
// 			input: map[string]any{
// 				"id":                      uint(2),
// 				"email":                   "jane@example.com",
// 				"password":                "another_hash",
// 				"profile.first_name":      "Jane",
// 				"profile.last_name":       "Smith",
// 				"profile.age":             25,
// 				"profile.address.street":  "456 Oak Ave",
// 				"profile.address.city":    "Los Angeles",
// 				"profile.address.country": "USA",
// 				"settings.theme":          "light",
// 				"settings.language":       "es",
// 				"settings.tags":           []string{"user", "trial"},
// 			},
// 			expected: KormUser{
// 				ID:       2,
// 				Email:    "jane@example.com",
// 				Password: "another_hash",
// 				Profile: &Profile{
// 					FirstName: "Jane",
// 					LastName:  "Smith",
// 					Age:       25,
// 					Address: &UserAddress{
// 						Street:  "456 Oak Ave",
// 						City:    "Los Angeles",
// 						Country: "USA",
// 					},
// 				},
// 				Settings: Settings{
// 					Theme:    "light",
// 					Language: "es",
// 					Tags:     []string{"user", "trial"},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Test FromMap
// 			t.Run("FromMap", func(t *testing.T) {
// 				var got KormUser
// 				builder := NewBuilder(&got)
// 				builder.FromMap(tt.input)
// 				assert.Equal(t, tt.expected.ID, got.ID)
// 				assert.Equal(t, tt.expected.Email, got.Email)
// 				assert.Equal(t, tt.expected.Password, got.Password)
// 				assert.Equal(t, tt.expected.Profile.FirstName, got.Profile.FirstName)
// 				assert.Equal(t, tt.expected.Profile.LastName, got.Profile.LastName)
// 				assert.Equal(t, tt.expected.Profile.Age, got.Profile.Age)
// 				assert.Equal(t, tt.expected.Profile.Address.Street, got.Profile.Address.Street)
// 				assert.Equal(t, tt.expected.Profile.Address.City, got.Profile.Address.City)
// 				assert.Equal(t, tt.expected.Profile.Address.Country, got.Profile.Address.Country)
// 				assert.Equal(t, tt.expected.Settings.Theme, got.Settings.Theme)
// 				assert.Equal(t, tt.expected.Settings.Language, got.Settings.Language)
// 				assert.Equal(t, tt.expected.Settings.Tags, got.Settings.Tags)
// 			})

// 			// Test FromKV
// 			t.Run("FromKV", func(t *testing.T) {
// 				var got KormUser
// 				builder := NewBuilder(&got)
// 				kvs := mapToKVs(tt.input)
// 				builder.FromKV(kvs...)
// 				assert.Equal(t, tt.expected.ID, got.ID)
// 				assert.Equal(t, tt.expected.Email, got.Email)
// 				assert.Equal(t, tt.expected.Password, got.Password)
// 				assert.Equal(t, tt.expected.Profile.FirstName, got.Profile.FirstName)
// 				assert.Equal(t, tt.expected.Profile.LastName, got.Profile.LastName)
// 				assert.Equal(t, tt.expected.Profile.Age, got.Profile.Age)
// 				assert.Equal(t, tt.expected.Profile.Address.Street, got.Profile.Address.Street)
// 				assert.Equal(t, tt.expected.Profile.Address.City, got.Profile.Address.City)
// 				assert.Equal(t, tt.expected.Profile.Address.Country, got.Profile.Address.Country)
// 				assert.Equal(t, tt.expected.Settings.Theme, got.Settings.Theme)
// 				assert.Equal(t, tt.expected.Settings.Language, got.Settings.Language)
// 				assert.Equal(t, tt.expected.Settings.Tags, got.Settings.Tags)
// 			})

// 			// Test Fill
// 			t.Run("Fill", func(t *testing.T) {
// 				var got KormUser
// 				got.Profile = &Profile{}
// 				got.Profile.Address = &UserAddress{}
// 				got.Settings = Settings{
// 					Tags: make([]string, 0),
// 				}
// 				kvs := mapToKVs(tt.input)
// 				err := Fill(&got, kvs)
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.expected.ID, got.ID)
// 				assert.Equal(t, tt.expected.Email, got.Email)
// 				assert.Equal(t, tt.expected.Password, got.Password)
// 				assert.Equal(t, tt.expected.Profile.FirstName, got.Profile.FirstName)
// 				assert.Equal(t, tt.expected.Profile.LastName, got.Profile.LastName)
// 				assert.Equal(t, tt.expected.Profile.Age, got.Profile.Age)
// 				assert.Equal(t, tt.expected.Profile.Address.Street, got.Profile.Address.Street)
// 				assert.Equal(t, tt.expected.Profile.Address.City, got.Profile.Address.City)
// 				assert.Equal(t, tt.expected.Profile.Address.Country, got.Profile.Address.Country)
// 				assert.Equal(t, tt.expected.Settings.Theme, got.Settings.Theme)
// 				assert.Equal(t, tt.expected.Settings.Language, got.Settings.Language)
// 				assert.Equal(t, tt.expected.Settings.Tags, got.Settings.Tags)
// 			})
// 		})
// 	}
// }
