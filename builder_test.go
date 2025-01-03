package kstrct

// BenchmarkStructManipulation/Builder/Get-4                2779670               429.9 ns/op             0 B/op          0 allocs/op
// BenchmarkStructManipulation/Builder/Set-4                2767867               433.4 ns/op             0 B/op          0 allocs/op
// BenchmarkStructManipulation/Builder/FromMap-4             923730              1130 ns/op               0 B/op          0 allocs/op
// BenchmarkStructManipulation/Builder/FromKV-4             1203914               995.1 ns/op             0 B/op          0 allocs/op
// BenchmarkStructManipulation/Builder/ToMap-4               805682              1487 ns/op              56 B/op          5 allocs/op
// BenchmarkStructManipulation/Builder/Clone-4             10200862               121.5 ns/op            48 B/op          1 allocs/op
// BenchmarkStructManipulation/Existing/Fill-4               717308              1632 ns/op              48 B/op          1 allocs/op
// BenchmarkStructManipulation/Builder/Filter-4             1716523               695.0 ns/op            32 B/op          3 allocs/op
// BenchmarkStructManipulation/Builder/Each-4               2100909               576.8 ns/op            24 B/op          2 allocs/op
// BenchmarkBuilder/Set-4                                    407602              2569 ns/op             264 B/op         10 allocs/op
// BenchmarkBuilder/FromMap-4                                347602              3453 ns/op             528 B/op          9 allocs/op
// BenchmarkBuilder/FromKV-4                                 351524              3247 ns/op             528 B/op          9 allocs/op
// BenchmarkBuilder/Fill-4                                   351994              3432 ns/op             528 B/op          9 allocs/op
// BenchmarkBuilder/FillOLD-4                                228141              5217 ns/op             272 B/op          8 allocs/op
// BenchmarkSetReflectFieldValue-4                           356673              3372 ns/op             720 B/op         20 allocs/op
// BenchmarkSetReflectFieldValueNew-4                        461209              2458 ns/op             376 B/op         11 allocs/op
// BenchmarkSetDirect-4                                      144570              8230 ns/op             976 B/op         40 allocs/op
// BenchmarkMethodCalls/Direct_Method_Call-4               1000000000               1.196 ns/op           0 B/op          0 allocs/op
// BenchmarkMethodCalls/Kstrct_CallMethod-4                  861462              1247 ns/op             184 B/op          5 allocs/op
// BenchmarkFillFromMap-4                                    519988              2049 ns/op             480 B/op         12 allocs/op
// BenchmarkFillFromKV-4                                     772657              1555 ns/op             120 B/op          9 allocs/op
// BenchmarkFrom-4                                          1592268               756.9 ns/op            56 B/op          4 allocs/op
// BenchmarkRange-4                                         1728423               697.0 ns/op            56 B/op          4 allocs/op
// BenchmarkFill-4                                          1288695               945.3 ns/op             0 B/op          0 allocs/op
// BenchmarkFillM-4                                          823796              1469 ns/op              24 B/op          1 allocs/op
// BenchmarkStructOperations/GetStringField-4              1000000000               0.5978 ns/op          0 B/op          0 allocs/op
// BenchmarkStructOperations/SetStringField-4              500898694                2.380 ns/op           0 B/op          0 allocs/op
// BenchmarkStructOperations/GetIntField-4                 1000000000               0.5945 ns/op          0 B/op          0 allocs/op
// BenchmarkStructOperations/SetIntField-4                 672623871                1.783 ns/op           0 B/op          0 allocs/op
// BenchmarkStructOperations/GetBoolField-4                1000000000               0.5947 ns/op          0 B/op          0 allocs/op
// BenchmarkStructOperations/SetBoolField-4                670106172                1.791 ns/op           0 B/op          0 allocs/op
// BenchmarkStructOperations/CopyStruct-4                  34569778                33.99 ns/op            0 B/op          0 allocs/op
// BenchmarkStructOperations/StructToMap/Unsafe-4           1979080               601.4 ns/op            40 B/op          4 allocs/op
// BenchmarkStructOperations/CompareStructs-4              28991248                40.38 ns/op            0 B/op          0 allocs/op
// BenchmarkStructOperations/GetFieldsByTag-4               1537294               778.4 ns/op            24 B/op          3 allocs/op
// BenchmarkReflectionOperations/GetField/Reflection-4             889956176                1.339 ns/op           0 B/op          0 allocs/op
// BenchmarkReflectionOperations/SetField/Reflection-4             201536242                5.948 ns/op           0 B/op          0 allocs/op
// BenchmarkReflectionOperations/CopyStruct/Reflection-4           28468332                41.30 ns/op            0 B/op          0 allocs/op

// type TestPerson struct {
// 	Name   string  `json:"name"`
// 	Age    int     `json:"age"`
// 	Active bool    `json:"active"`
// 	Score  float64 `json:"score"`
// }

// func BenchmarkStructManipulation(b *testing.B) {
// 	// Test data
// 	person := &TestPerson{
// 		Name:   "John Doe",
// 		Age:    30,
// 		Active: true,
// 		Score:  95.5,
// 	}

// 	// Data for map operations
// 	data := map[string]any{
// 		"Name":   "Jane Doe",
// 		"Age":    25,
// 		"Active": false,
// 		"Score":  98.5,
// 	}

// 	// Test Builder approach
// 	b.Run("Builder/Get", func(b *testing.B) {
// 		builder := NewBuilder(person)
// 		b.ResetTimer()
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			_ = builder.GetString("name")
// 			_ = builder.GetInt("age")
// 			_ = builder.GetBool("active")
// 			_ = builder.GetFloat64("score")
// 		}
// 	})

// 	b.Run("Builder/Set", func(b *testing.B) {
// 		builder := NewBuilder(person)
// 		b.ResetTimer()
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			builder.SetString("name", "Jane Doe").
// 				SetInt("age", 25).
// 				SetBool("active", false).
// 				SetFloat64("score", 98.5)
// 		}
// 	})

// 	b.Run("Builder/FromMap", func(b *testing.B) {
// 		builder := NewBuilder(person)
// 		data := map[string]any{
// 			"name":   "Jane Doe",
// 			"age":    25,
// 			"active": false,
// 			"score":  98.5,
// 		}
// 		b.ResetTimer()
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			builder.FromMap(data)
// 		}
// 	})
// 	b.Run("Builder/FromKV", func(b *testing.B) {
// 		builder := NewBuilder(person)
// 		kvs := []KV{
// 			{Key: "name", Value: "Jane Doe"},
// 			{Key: "age", Value: 25},
// 			{Key: "active", Value: true},
// 			{Key: "score", Value: 98.5},
// 		}
// 		b.ResetTimer()
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			builder.FromKV(kvs...)
// 		}
// 	})

// 	b.Run("Builder/ToMap", func(b *testing.B) {
// 		builder := NewBuilder(person)
// 		b.ResetTimer()
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			m := builder.ToMap()
// 			PutMap(m)
// 		}
// 	})

// 	b.Run("Builder/Clone", func(b *testing.B) {
// 		builder := NewBuilder(person)
// 		b.ResetTimer()
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			_ = builder.Clone()
// 		}
// 	})

// 	// Test existing Fill approach
// 	b.Run("Existing/Fill", func(b *testing.B) {
// 		b.ResetTimer()
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			p := &TestPerson{}
// 			_ = FillM(p, data)
// 		}
// 	})

// 	// Test Filter operation
// 	b.Run("Builder/Filter", func(b *testing.B) {
// 		builder := NewBuilder(person)
// 		b.ResetTimer()
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			m := builder.Filter(func(name string, value any) bool {
// 				return name == "Name" || name == "Age"
// 			})
// 			PutMap(m)
// 		}
// 	})

// 	// Test Each operation
// 	b.Run("Builder/Each", func(b *testing.B) {
// 		builder := NewBuilder(person)
// 		var sum int
// 		b.ResetTimer()
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			builder.Each(func(name string, value any) {
// 				if v, ok := value.(int); ok {
// 					sum += v
// 				}
// 			})
// 		}
// 	})
// }

// func ExampleStructBuilder() {
// 	// Create a new person
// 	person := &TestPerson{
// 		Name:   "John Doe",
// 		Age:    30,
// 		Active: true,
// 		Score:  95.5,
// 	}

// 	// Create a builder
// 	builder := NewBuilder(person)

// 	// Modify fields fluently
// 	builder.Set("name", "Jane Doe").
// 		Set("age", 25).
// 		Set("active", false)

// 	// Get field values
// 	name := builder.Get("name")
// 	age := builder.Get("age")

// 	// Print values for example
// 	fmt.Printf("Name: %v, Age: %v\n", name, age)

// 	// Convert to map
// 	m := builder.ToMap()
// 	defer PutMap(m) // Return map to pool when done

// 	// Clone the struct
// 	clone := builder.Clone().(*TestPerson)
// 	fmt.Printf("Clone: %+v\n", clone)

// 	// Filter fields
// 	filtered := builder.Filter(func(name string, value any) bool {
// 		return name == "name" || name == "age"
// 	})
// 	defer PutMap(filtered)

// 	// Iterate over fields
// 	builder.Each(func(name string, value any) {
// 		fmt.Printf("Field: %s = %v\n", name, value)
// 	})
// }

// func TestBuilder(t *testing.T) {
// 	t.Run("Basic types", func(t *testing.T) {
// 		var s BasicTypes
// 		builder := NewBuilder(&s)

// 		// Test direct Set
// 		builder.Set("string", "test").
// 			Set("int", 42).
// 			Set("int8", int8(8)).
// 			Set("int16", int16(16)).
// 			Set("int32", int32(32)).
// 			Set("int64", int64(64)).
// 			Set("uint", uint(42)).
// 			Set("uint8", uint8(8)).
// 			Set("uint16", uint16(16)).
// 			Set("uint32", uint32(32)).
// 			Set("uint64", uint64(64)).
// 			Set("float32", float32(3.14)).
// 			Set("float64", 3.14159).
// 			Set("bool", true).
// 			Set("time", time.Now())

// 		// Test FromMap
// 		builder.FromMap(map[string]any{
// 			"string":  "test2",
// 			"int":     43,
// 			"float64": 3.14,
// 			"bool":    true,
// 			"time":    time.Now(),
// 		})
// 		assert.Equal(t, "test2", s.String)
// 		assert.Equal(t, 43, s.Int)
// 		assert.Equal(t, 3.14, s.Float64)
// 		assert.True(t, s.Bool)
// 	})

// 	t.Run("Pointer types", func(t *testing.T) {
// 		var s BasicTypes
// 		builder := NewBuilder(&s)

// 		str := "pointer"
// 		num := 42
// 		f := 3.14
// 		b := true
// 		tm := time.Now()

// 		builder.Set("string_ptr", &str).
// 			Set("int_ptr", &num).
// 			Set("float_ptr", &f).
// 			Set("bool_ptr", &b).
// 			Set("time_ptr", &tm)

// 		assert.Equal(t, str, *s.StringPtr)
// 		assert.Equal(t, num, *s.IntPtr)
// 		assert.Equal(t, f, *s.FloatPtr)
// 		assert.Equal(t, b, *s.BoolPtr)
// 		assert.Equal(t, tm, *s.TimePtr)
// 	})

// 	t.Run("Slice types", func(t *testing.T) {
// 		var s ComplexTypes
// 		builder := NewBuilder(&s)

// 		// Test string slice
// 		builder.Set("string_slice", []string{"a", "b", "c"})
// 		assert.Equal(t, []string{"a", "b", "c"}, s.StringSlice)

// 		// Test FromMap with comma-separated strings
// 		builder.FromMap(map[string]any{
// 			"string_slice": "x,  y,z",
// 			"int_slice":    "1,2, 3",
// 			"float_slice":  "1.1, 2.2,3.3",
// 			"bool_slice":   "true,false,true",
// 		})

// 		assert.Equal(t, []string{"x", "y", "z"}, s.StringSlice)
// 		assert.Equal(t, []int{1, 2, 3}, s.IntSlice)
// 		assert.Equal(t, []float64{1.1, 2.2, 3.3}, s.FloatSlice)
// 		assert.Equal(t, []bool{true, false, true}, s.BoolSlice)
// 	})

// 	t.Run("Map types", func(t *testing.T) {
// 		var s ComplexTypes
// 		builder := NewBuilder(&s)

// 		// Test direct map setting
// 		stringMap := map[string]string{"key1": "value1", "key2": "value2"}
// 		intMap := map[string]int{"key1": 1, "key2": 2}
// 		floatMap := map[string]float64{"key1": 1.1, "key2": 2.2}
// 		boolMap := map[string]bool{"key1": true, "key2": false}

// 		builder.Set("string_map", stringMap).
// 			Set("int_map", intMap).
// 			Set("float_map", floatMap).
// 			Set("bool_map", boolMap)

// 		assert.Equal(t, stringMap, s.StringMap)
// 		assert.Equal(t, intMap, s.IntMap)
// 		assert.Equal(t, floatMap, s.FloatMap)
// 		assert.Equal(t, boolMap, s.BoolMap)

// 		// Test dot notation
// 		builder.Set("string_map.key3", "value3").
// 			Set("int_map.key3", 3).
// 			Set("float_map.key3", 3.3).
// 			Set("bool_map.key3", true)

// 		assert.Equal(t, "value3", s.StringMap["key3"])
// 		assert.Equal(t, 3, s.IntMap["key3"])
// 		assert.Equal(t, 3.3, s.FloatMap["key3"])
// 		assert.Equal(t, true, s.BoolMap["key3"])
// 	})

// 	t.Run("Nested struct", func(t *testing.T) {
// 		var s ComplexTypes
// 		builder := NewBuilder(&s)

// 		// Test direct nested struct setting
// 		nested := NestedStruct{
// 			Field1: "nested string",
// 			Field2: 456,
// 		}
// 		builder.Set("nested", nested)
// 		assert.Equal(t, nested, s.Nested)

// 		// Test nested fields using FromMap
// 		builder.FromMap(map[string]any{
// 			"nested.field1": "updated string",
// 			"nested.field2": 789,
// 		})

// 		assert.Equal(t, "updated string", s.Nested.Field1)
// 		assert.Equal(t, 789, s.Nested.Field2)
// 	})

// 	t.Run("Nested pointer struct", func(t *testing.T) {
// 		var s ComplexTypes
// 		builder := NewBuilder(&s)

// 		// Test nested pointer fields using FromMap
// 		builder.FromMap(map[string]any{
// 			"nested_ptr.field1": "nested string",
// 			"nested_ptr.field2": 456,
// 		})

// 		assert.NotNil(t, s.NestedPtr)
// 		assert.Equal(t, "nested string", s.NestedPtr.Field1)
// 		assert.Equal(t, 456, s.NestedPtr.Field2)
// 	})

// 	t.Run("Slice of structs", func(t *testing.T) {
// 		var s ComplexTypes
// 		builder := NewBuilder(&s)

// 		// Test setting slice of structs
// 		structs := []NestedStruct{
// 			{Field1: "first", Field2: 1},
// 			{Field1: "second", Field2: 2},
// 		}
// 		builder.Set("struct_slice", structs)
// 		assert.Equal(t, structs, s.StructSlice)

// 		// Test FromMap with indexed notation
// 		builder.FromMap(map[string]any{
// 			"struct_slice": []map[string]any{
// 				{"field1": "third", "field2": 3},
// 				{"field1": "fourth", "field2": 4},
// 			},
// 		})

// 		assert.Len(t, s.StructSlice, 2)
// 		assert.Equal(t, "third", s.StructSlice[0].Field1)
// 		assert.Equal(t, 3, s.StructSlice[0].Field2)
// 		assert.Equal(t, "fourth", s.StructSlice[1].Field1)
// 		assert.Equal(t, 4, s.StructSlice[1].Field2)
// 	})

// 	t.Run("SQL types", func(t *testing.T) {
// 		var s SQLTypes
// 		builder := NewBuilder(&s)

// 		now := time.Now()
// 		builder.Set("null_string", "test").
// 			Set("null_int64", 42).
// 			Set("null_float64", 3.14).
// 			Set("null_bool", true).
// 			Set("null_time", now)

// 		assert.True(t, s.NullString.Valid)
// 		assert.Equal(t, "test", s.NullString.String)
// 		assert.True(t, s.NullInt64.Valid)
// 		assert.Equal(t, int64(42), s.NullInt64.Int64)
// 		assert.True(t, s.NullFloat64.Valid)
// 		assert.Equal(t, 3.14, s.NullFloat64.Float64)
// 		assert.True(t, s.NullBool.Valid)
// 		assert.True(t, s.NullBool.Bool)
// 		assert.True(t, s.NullTime.Valid)
// 		assert.Equal(t, now, s.NullTime.Time)
// 	})

// 	t.Run("KORM-style mapping", func(t *testing.T) {
// 		var user User
// 		builder := NewBuilder(&user)

// 		// Test nested struct creation and setting
// 		builder.FromMap(map[string]any{
// 			"id":       1,
// 			"email":    "test@example.com",
// 			"password": "secret",
// 			"profile": map[string]any{
// 				"first_name": "John",
// 				"last_name":  "Doe",
// 				"age":        30,
// 				"address": map[string]any{
// 					"street":  "123 Main St",
// 					"city":    "New York",
// 					"country": "USA",
// 				},
// 			},
// 			"settings": map[string]string{
// 				"theme":    "dark",
// 				"language": "en",
// 			},
// 		})

// 		assert.Equal(t, 1, user.ID)
// 		assert.Equal(t, "test@example.com", user.Email)
// 		assert.Equal(t, "secret", user.Password)
// 		assert.NotNil(t, user.Profile)
// 		assert.Equal(t, "John", user.Profile.FirstName)
// 		assert.Equal(t, "Doe", user.Profile.LastName)
// 		assert.Equal(t, 30, user.Profile.Age)
// 		assert.NotNil(t, user.Profile.Address)
// 		assert.Equal(t, "123 Main St", user.Profile.Address.Street)
// 		assert.Equal(t, "New York", user.Profile.Address.City)
// 		assert.Equal(t, "USA", user.Profile.Address.Country)
// 		assert.Equal(t, "dark", user.Settings["theme"])
// 		assert.Equal(t, "en", user.Settings["language"])

// 		// Test dot notation for deep updates
// 		builder.Set("profile.first_name", "Jane").
// 			Set("profile.address.city", "Los Angeles").
// 			Set("settings.theme", "light")

// 		assert.Equal(t, "Jane", user.Profile.FirstName)
// 		assert.Equal(t, "Los Angeles", user.Profile.Address.City)
// 		assert.Equal(t, "light", user.Settings["theme"])
// 	})

// 	t.Run("Type conversions", func(t *testing.T) {
// 		var s BasicTypes
// 		builder := NewBuilder(&s)

// 		builder.Set("string", 123). // int to string
// 						Set("int", "123").      // string to int
// 						Set("float64", "3.14"). // string to float
// 						Set("bool", 1).         // int to bool
// 						Set("time", 1640995200) // unix timestamp to time

// 		assert.Equal(t, "123", s.String)
// 		assert.Equal(t, 123, s.Int)
// 		assert.Equal(t, 3.14, s.Float64)
// 		assert.True(t, s.Bool)
// 		assert.Equal(t, time.Unix(1640995200, 0), s.Time)
// 	})

// 	t.Run("Edge cases", func(t *testing.T) {
// 		var s BasicTypes
// 		builder := NewBuilder(&s)

// 		// Test nil pointer
// 		builder.Set("string_ptr", nil)
// 		assert.Nil(t, s.StringPtr)

// 		// Test invalid conversions
// 		builder.Set("int", "invalid")
// 		builder.Set("bool", "invalid")
// 		builder.Set("time", "invalid")

// 		// Values should remain at their zero values
// 		assert.Equal(t, 0, s.Int)
// 		assert.False(t, s.Bool)
// 		assert.Equal(t, time.Time{}, s.Time)
// 	})
// }

// func TestBuilderDocto(t *testing.T) {
// 	u := &Docto{
// 		Languages:    make([]string, 0),
// 		Reservations: make([]Reservation, 0),
// 		VisitTypes:   make([]string, 0),
// 	}
// 	builder := NewBuilder(u)
// 	builder.Set("uuid", "xxx-xxx-xxx").
// 		Set("name", "kamal").
// 		Set("languages", "fr,en,es").
// 		Set("back_at", time.Now()).
// 		Set("is_blocked", true).
// 		Set("week_timeslots.sunday", "8:00,9:00,10:00").
// 		Set("week_timeslots.monday", "10:00,11:00,12:00").
// 		Set("reservations.0.id", 1).
// 		Set("reservations.0.patient_id", 12345).
// 		Set("visit_types", "bla,bla2,bla3,bla4")
// }

// func TestBuilderNestedFieldsSlice(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		expected Doctor
// 	}{
// 		{
// 			name: "nested week timeslots",
// 			expected: Doctor{
// 				Name: "Dr. Smith",
// 				WeekTimeslots: &[]WeekTimeslots{
// 					{
// 						Monday:  []string{"09:00", "10:00", "11:00"},
// 						Tuesday: []string{"14:00", "15:00"},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var got Doctor
// 			builder := NewBuilder(&got)

// 			// Set values using builder
// 			builder.Set("name", "Dr. Smith").
// 				Set("week_timeslots.monday", "09:00,10:00,11:00").
// 				Set("week_timeslots.tuesday", "14:00,15:00")

// 			// Compare results
// 			assert.Equal(t, tt.expected.Name, builder.Get("name"))

// 			weekTimeslots := (*got.WeekTimeslots)[0]
// 			expectedTimeslots := (*tt.expected.WeekTimeslots)[0]
// 			assert.Equal(t, expectedTimeslots.Monday, weekTimeslots.Monday)
// 			assert.Equal(t, expectedTimeslots.Tuesday, weekTimeslots.Tuesday)
// 		})
// 	}
// }

// func TestBuilderNestedFieldsStruct(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		expected DoctorS
// 	}{
// 		{
// 			name: "nested week timeslots",
// 			expected: DoctorS{
// 				Name: "Dr. Smith",
// 				WeekTimeslots: &WeekTimeslots{
// 					Monday:  []string{"09:00", "10:00", "11:00"},
// 					Tuesday: []string{"14:00", "15:00"},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var got DoctorS
// 			builder := NewBuilder(&got)

// 			// Set values using builder
// 			builder.Set("name", "Dr. Smith").
// 				Set("week_timeslots.monday", "09:00,10:00,11:00").
// 				Set("week_timeslots.tuesday", "14:00,15:00")

// 			// Compare results
// 			assert.Equal(t, tt.expected.Name, builder.Get("name"))

// 			weekTimeslots := got.WeekTimeslots
// 			assert.Equal(t, tt.expected.WeekTimeslots.Monday, weekTimeslots.Monday)
// 			assert.Equal(t, tt.expected.WeekTimeslots.Tuesday, weekTimeslots.Tuesday)
// 		})
// 	}
// }

// func TestBuilderFillDocto(t *testing.T) {
// 	u := &Docto{}
// 	builder := NewBuilder(u)

// 	builder.Set("uuid", "xxx-xxx-xxx").
// 		Set("name", "kamal").
// 		Set("languages", "fr,en,es").
// 		Set("back_at", time.Now()).
// 		Set("is_blocked", true).
// 		Set("week_timeslots.sunday", "8:00,9:00,10:00").
// 		Set("week_timeslots.monday", "10:00,11:00,12:00").
// 		Set("reservations.0.id", 1).
// 		Set("reservations.0.patient_id", 12345).
// 		Set("visit_types", "bla,bla2,bla3,bla4")
// }

// func TestBuilderFillNestedFieldsSlice(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		expected Doctor
// 	}{
// 		{
// 			name: "nested week timeslots",
// 			expected: Doctor{
// 				Name: "Dr. Smith",
// 				WeekTimeslots: &[]WeekTimeslots{
// 					{
// 						Monday:  []string{"09:00", "10:00", "11:00"},
// 						Tuesday: []string{"14:00", "15:00"},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var got Doctor
// 			builder := NewBuilder(&got)

// 			builder.Set("name", "Dr. Smith").
// 				Set("week_timeslots.monday", "09:00,10:00,11:00").
// 				Set("week_timeslots.tuesday", "14:00,15:00")

// 			t.Log(got)
// 			if got.Name != tt.expected.Name {
// 				t.Errorf("Name = %v, want %v", got.Name, tt.expected.Name)
// 			}

// 			if !reflect.DeepEqual((*got.WeekTimeslots)[0].Monday, (*tt.expected.WeekTimeslots)[0].Monday) {
// 				t.Errorf("Monday slots = %v, want %v", (*got.WeekTimeslots)[0].Monday, (*tt.expected.WeekTimeslots)[0].Monday)
// 			}

// 			if !reflect.DeepEqual((*got.WeekTimeslots)[0].Tuesday, (*tt.expected.WeekTimeslots)[0].Tuesday) {
// 				t.Errorf("Tuesday slots = %v, want %v", (*got.WeekTimeslots)[0].Tuesday, (*tt.expected.WeekTimeslots)[0].Tuesday)
// 			}
// 		})
// 	}
// }

// func TestBuilderFillNestedFieldsStruct(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		expected DoctorS
// 	}{
// 		{
// 			name: "nested week timeslots",
// 			expected: DoctorS{
// 				Name: "Dr. Smith",
// 				WeekTimeslots: &WeekTimeslots{
// 					Monday:  []string{"09:00", "10:00", "11:00"},
// 					Tuesday: []string{"14:00", "15:00"},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var got DoctorS
// 			builder := NewBuilder(&got)

// 			builder.Set("name", "Dr. Smith").
// 				Set("week_timeslots.monday", "09:00,10:00,11:00").
// 				Set("week_timeslots.tuesday", "14:00,15:00")

// 			t.Log(got)
// 			if got.Name != tt.expected.Name {
// 				t.Errorf("Name = %v, want %v", got.Name, tt.expected.Name)
// 				return
// 			}

// 			if !reflect.DeepEqual(got.WeekTimeslots.Monday, tt.expected.WeekTimeslots.Monday) {
// 				t.Errorf("Monday slots = %v, want %v", got.WeekTimeslots.Monday, tt.expected.WeekTimeslots.Monday)
// 			}

// 			if !reflect.DeepEqual(got.WeekTimeslots.Tuesday, tt.expected.WeekTimeslots.Tuesday) {
// 				t.Errorf("Tuesday slots = %v, want %v", got.WeekTimeslots.Tuesday, tt.expected.WeekTimeslots.Tuesday)
// 			}
// 		})
// 	}
// }

// func TestBuilderDoctoFromMap(t *testing.T) {
// 	u := &Docto{
// 		Languages:    make([]string, 0),
// 		Reservations: make([]Reservation, 0),
// 		VisitTypes:   make([]string, 0),
// 	}
// 	builder := NewBuilder(u)

// 	data := map[string]any{
// 		"uuid":                      "xxx-xxx-xxx",
// 		"name":                      "kamal",
// 		"languages":                 "fr,en,es",
// 		"back_at":                   time.Now(),
// 		"is_blocked":                true,
// 		"week_timeslots.sunday":     "8:00,9:00,10:00",
// 		"week_timeslots.monday":     "10:00,11:00,12:00",
// 		"reservations.0.id":         1,
// 		"reservations.0.patient_id": 12345,
// 		"visit_types":               "bla,bla2,bla3,bla4",
// 	}

// 	builder.FromMap(data)
// }

// func TestBuilderNestedFieldsSliceFromMap(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		expected Doctor
// 	}{
// 		{
// 			name: "nested week timeslots",
// 			expected: Doctor{
// 				Name: "Dr. Smith",
// 				WeekTimeslots: &[]WeekTimeslots{
// 					{
// 						Monday:  []string{"09:00", "10:00", "11:00"},
// 						Tuesday: []string{"14:00", "15:00"},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var got Doctor
// 			builder := NewBuilder(&got)

// 			data := map[string]any{
// 				"name":                   "Dr. Smith",
// 				"week_timeslots.monday":  "09:00,10:00,11:00",
// 				"week_timeslots.tuesday": "14:00,15:00",
// 			}

// 			builder.FromMap(data)

// 			// Compare results
// 			assert.Equal(t, tt.expected.Name, builder.Get("name"))

// 			weekTimeslots := (*got.WeekTimeslots)[0]
// 			expectedTimeslots := (*tt.expected.WeekTimeslots)[0]
// 			assert.Equal(t, expectedTimeslots.Monday, weekTimeslots.Monday)
// 			assert.Equal(t, expectedTimeslots.Tuesday, weekTimeslots.Tuesday)
// 		})
// 	}
// }

// func TestBuilderNestedFieldsStructFromMap(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		expected DoctorS
// 	}{
// 		{
// 			name: "nested week timeslots",
// 			expected: DoctorS{
// 				Name: "Dr. Smith",
// 				WeekTimeslots: &WeekTimeslots{
// 					Monday:  []string{"09:00", "10:00", "11:00"},
// 					Tuesday: []string{"14:00", "15:00"},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var got DoctorS
// 			builder := NewBuilder(&got)

// 			data := map[string]any{
// 				"name":                   "Dr. Smith",
// 				"week_timeslots.monday":  "09:00,10:00,11:00",
// 				"week_timeslots.tuesday": "14:00,15:00",
// 			}

// 			builder.FromMap(data)

// 			// Compare results
// 			assert.Equal(t, tt.expected.Name, builder.Get("name"))

// 			weekTimeslots := got.WeekTimeslots
// 			assert.Equal(t, tt.expected.WeekTimeslots.Monday, weekTimeslots.Monday)
// 			assert.Equal(t, tt.expected.WeekTimeslots.Tuesday, weekTimeslots.Tuesday)
// 		})
// 	}
// }

// func TestBuilderFillDoctoFromMap(t *testing.T) {
// 	u := &Docto{}
// 	builder := NewBuilder(u)

// 	data := map[string]any{
// 		"uuid":                      "xxx-xxx-xxx",
// 		"name":                      "kamal",
// 		"languages":                 "fr,en,es",
// 		"back_at":                   time.Now(),
// 		"is_blocked":                true,
// 		"week_timeslots.sunday":     "8:00,9:00,10:00",
// 		"week_timeslots.monday":     "10:00,11:00,12:00",
// 		"reservations.0.id":         1,
// 		"reservations.0.patient_id": 12345,
// 		"visit_types":               "bla,bla2,bla3,bla4",
// 	}

// 	builder.FromMap(data)
// }

// func TestFillNestedFieldsSliceFromMap(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		expected Doctor
// 	}{
// 		{
// 			name: "nested week timeslots",
// 			expected: Doctor{
// 				Name: "Dr. Smith",
// 				WeekTimeslots: &[]WeekTimeslots{
// 					{
// 						Monday:  []string{"09:00", "10:00", "11:00"},
// 						Tuesday: []string{"14:00", "15:00"},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var got Doctor
// 			builder := NewBuilder(&got)

// 			data := map[string]any{
// 				"name":                   "Dr. Smith",
// 				"week_timeslots.monday":  "09:00,10:00,11:00",
// 				"week_timeslots.tuesday": "14:00,15:00",
// 			}

// 			builder.FromMap(data)

// 			fmt.Printf("%+v\n", got)
// 			assert.Equal(t, tt.expected.Name, got.Name)
// 			weekTimeslots := (*got.WeekTimeslots)[0]
// 			expectedTimeslots := (*tt.expected.WeekTimeslots)[0]
// 			assert.Equal(t, expectedTimeslots.Monday, weekTimeslots.Monday)
// 			assert.Equal(t, expectedTimeslots.Tuesday, weekTimeslots.Tuesday)
// 		})
// 	}
// }

// func TestFillNestedFieldsStructFromMap(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		expected DoctorS
// 	}{
// 		{
// 			name: "nested week timeslots",
// 			expected: DoctorS{
// 				Name: "Dr. Smith",
// 				WeekTimeslots: &WeekTimeslots{
// 					Monday:  []string{"09:00", "10:00", "11:00"},
// 					Tuesday: []string{"14:00", "15:00"},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var got DoctorS
// 			builder := NewBuilder(&got)

// 			data := map[string]any{
// 				"name":                   "Dr. Smith",
// 				"week_timeslots.monday":  "09:00,10:00,11:00",
// 				"week_timeslots.tuesday": "14:00,15:00",
// 			}

// 			builder.FromMap(data)

// 			fmt.Printf("%+v\n", got)
// 			assert.Equal(t, tt.expected.Name, got.Name)
// 			weekTimeslots := got.WeekTimeslots
// 			assert.Equal(t, tt.expected.WeekTimeslots.Monday, weekTimeslots.Monday)
// 			assert.Equal(t, tt.expected.WeekTimeslots.Tuesday, weekTimeslots.Tuesday)
// 		})
// 	}
// }

// func TestBuilderFromKV(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		kvs  []KV
// 		want any
// 	}{
// 		{
// 			name: "Basic types",
// 			kvs: []KV{
// 				{Key: "name", Value: "kamal"},
// 				{Key: "is_blocked", Value: true},
// 				{Key: "languages", Value: "fr,en,es"},
// 			},
// 			want: &Docto{
// 				Name:      "kamal",
// 				IsBlocked: true,
// 				Languages: []string{"fr", "en", "es"},
// 			},
// 		},
// 		{
// 			name: "Nested fields",
// 			kvs: []KV{
// 				{Key: "name", Value: "Dr. Smith"},
// 				{Key: "week_timeslots.monday", Value: "09:00,10:00,11:00"},
// 				{Key: "week_timeslots.tuesday", Value: "14:00,15:00"},
// 			},
// 			want: &Doctor{
// 				Name: "Dr. Smith",
// 				WeekTimeslots: &[]WeekTimeslots{
// 					{
// 						Monday:  []string{"09:00", "10:00", "11:00"},
// 						Tuesday: []string{"14:00", "15:00"},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var got any
// 			switch tt.want.(type) {
// 			case *Docto:
// 				got = &Docto{
// 					Languages:    make([]string, 0),
// 					Reservations: make([]Reservation, 0),
// 					VisitTypes:   make([]string, 0),
// 				}
// 			case *Doctor:
// 				got = &Doctor{}
// 			}

// 			builder := NewBuilder(got)
// 			builder.FromKV(tt.kvs...)

// 			switch want := tt.want.(type) {
// 			case *Docto:
// 				gotDocto := got.(*Docto)
// 				assert.Equal(t, want.Name, gotDocto.Name)
// 				assert.Equal(t, want.IsBlocked, gotDocto.IsBlocked)
// 				assert.Equal(t, want.Languages, gotDocto.Languages)
// 			case *Doctor:
// 				gotDoctor := got.(*Doctor)
// 				assert.Equal(t, want.Name, gotDoctor.Name)
// 				if want.WeekTimeslots != nil {
// 					weekTimeslots := (*gotDoctor.WeekTimeslots)[0]
// 					expectedTimeslots := (*want.WeekTimeslots)[0]
// 					assert.Equal(t, expectedTimeslots.Monday, weekTimeslots.Monday)
// 					assert.Equal(t, expectedTimeslots.Tuesday, weekTimeslots.Tuesday)
// 				}
// 			}
// 		})
// 	}
// }

// func BenchmarkBuilder(b *testing.B) {
// 	b.Run("Set", func(b *testing.B) {
// 		u := &Docto{}
// 		builder := NewBuilder(u)
// 		b.ResetTimer()

// 		for i := 0; i < b.N; i++ {
// 			builder.Set("name", "kamal").
// 				Set("is_blocked", true).
// 				Set("languages", "fr,en,es").
// 				Set("week_timeslots.monday", "10:00,11:00,12:00")
// 		}
// 	})

// 	b.Run("FromMap", func(b *testing.B) {
// 		u := &Docto{}
// 		builder := NewBuilder(u)
// 		data := map[string]any{
// 			"name":                  "kamal",
// 			"is_blocked":            true,
// 			"languages":             "fr,en,es",
// 			"week_timeslots.monday": "10:00,11:00,12:00",
// 		}
// 		b.ResetTimer()

// 		for i := 0; i < b.N; i++ {
// 			builder.FromMap(data)
// 		}
// 	})

// 	b.Run("FromKV", func(b *testing.B) {
// 		u := &Docto{}
// 		builder := NewBuilder(u)
// 		kvs := []KV{
// 			{Key: "name", Value: "kamal"},
// 			{Key: "is_blocked", Value: true},
// 			{Key: "languages", Value: "fr,en,es"},
// 			{Key: "week_timeslots.monday", Value: "10:00,11:00,12:00"},
// 		}
// 		b.ResetTimer()

// 		for i := 0; i < b.N; i++ {
// 			builder.FromKV(kvs...)
// 		}
// 	})
// 	b.Run("Fill", func(b *testing.B) {
// 		u := &Docto{}
// 		kvs := []KV{
// 			{Key: "name", Value: "kamal"},
// 			{Key: "is_blocked", Value: true},
// 			{Key: "languages", Value: "fr,en,es"},
// 			{Key: "week_timeslots.monday", Value: "10:00,11:00,12:00"},
// 		}
// 		b.ResetTimer()

// 		for i := 0; i < b.N; i++ {
// 			Fill(u, kvs, true)
// 		}
// 	})
// 	b.Run("FillOLD", func(b *testing.B) {
// 		u := &Docto{}
// 		kvs := []KV{
// 			{Key: "name", Value: "kamal"},
// 			{Key: "is_blocked", Value: true},
// 			{Key: "languages", Value: "fr,en,es"},
// 			{Key: "week_timeslots.monday", Value: "10:00,11:00,12:00"},
// 		}
// 		b.ResetTimer()

// 		for i := 0; i < b.N; i++ {
// 			FillOLD(u, kvs, true)
// 		}
// 	})
// }

// func BenchmarkSetReflectFieldValue(b *testing.B) {
// 	type WeekTimeslots struct {
// 		Monday    []string
// 		Tuesday   []string
// 		Wednesday []string
// 		Thursday  []string
// 		Friday    []string
// 	}
// 	type Doctor struct {
// 		Name          string
// 		WeekTimeslots *WeekTimeslots
// 	}

// 	d := &Doctor{
// 		Name:          "Dr. Smith",
// 		WeekTimeslots: &WeekTimeslots{},
// 	}

// 	// Get fields for benchmarking
// 	v := reflect.ValueOf(d).Elem()
// 	nameField := v.FieldByName("Name")
// 	weekField := v.FieldByName("WeekTimeslots").Elem()
// 	mondayField := weekField.FieldByName("Monday")
// 	tuesdayField := weekField.FieldByName("Tuesday")
// 	wednesdayField := weekField.FieldByName("Wednesday")
// 	thursdayField := weekField.FieldByName("Thursday")
// 	fridayField := weekField.FieldByName("Friday")

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		// Set simple field
// 		SetReflectFieldValue(nameField, "Dr. Jones")

// 		// Set all week timeslots
// 		SetReflectFieldValue(mondayField, "9:00,10:00,11:00")
// 		SetReflectFieldValue(tuesdayField, "9:00,10:00,11:00")
// 		SetReflectFieldValue(wednesdayField, "9:00,10:00,11:00")
// 		SetReflectFieldValue(thursdayField, "9:00,10:00,11:00")
// 		SetReflectFieldValue(fridayField, "9:00,10:00,11:00")
// 	}
// }

// func BenchmarkSetReflectFieldValueNew(b *testing.B) {
// 	type WeekTimeslots struct {
// 		Monday    []string
// 		Tuesday   []string
// 		Wednesday []string
// 		Thursday  []string
// 		Friday    []string
// 	}
// 	type Doctor struct {
// 		Name          string
// 		WeekTimeslots *WeekTimeslots
// 	}

// 	d := &Doctor{
// 		Name:          "Dr. Smith",
// 		WeekTimeslots: &WeekTimeslots{},
// 	}

// 	// Get fields for benchmarking
// 	v := reflect.ValueOf(d).Elem()
// 	nameField := v.FieldByName("Name")
// 	weekField := v.FieldByName("WeekTimeslots").Elem()
// 	mondayField := weekField.FieldByName("Monday")
// 	tuesdayField := weekField.FieldByName("Tuesday")
// 	wednesdayField := weekField.FieldByName("Wednesday")
// 	thursdayField := weekField.FieldByName("Thursday")
// 	fridayField := weekField.FieldByName("Friday")

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		// Set simple field
// 		SetRFValue(nameField, "Dr. Jones")

// 		// Set all week timeslots
// 		SetRFValue(mondayField, "9:00,10:00,11:00")
// 		SetRFValue(tuesdayField, "9:00,10:00,11:00")
// 		SetRFValue(wednesdayField, "9:00,10:00,11:00")
// 		SetRFValue(thursdayField, "9:00,10:00,11:00")
// 		SetRFValue(fridayField, "9:00,10:00,11:00")
// 	}
// }

// func BenchmarkSetDirect(b *testing.B) {
// 	type WeekTimeslots struct {
// 		Monday    []string
// 		Tuesday   []string
// 		Wednesday []string
// 		Thursday  []string
// 		Friday    []string
// 	}
// 	type Doctor struct {
// 		Name          string
// 		WeekTimeslots *WeekTimeslots
// 	}

// 	d := &Doctor{
// 		Name:          "Dr. Smith",
// 		WeekTimeslots: &WeekTimeslots{},
// 	}
// 	builder := NewBuilder(d)

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		// Set simple field
// 		builder.Set("name", "Dr. Jones")

// 		// Set all week timeslots
// 		builder.Set("week_timeslots.monday", "9:00,10:00,11:00")
// 		builder.Set("week_timeslots.tuesday", "9:00,10:00,11:00")
// 		builder.Set("week_timeslots.wednesday", "9:00,10:00,11:00")
// 		builder.Set("week_timeslots.thursday", "9:00,10:00,11:00")
// 		builder.Set("week_timeslots.friday", "9:00,10:00,11:00")
// 	}
// }

// func TestNestedFieldHandling(t *testing.T) {
// 	type Address struct {
// 		Street string
// 		City   string
// 	}
// 	type User struct {
// 		Name    string
// 		Address *Address
// 		Tags    map[string]string
// 	}

// 	tests := []struct {
// 		name   string
// 		input  map[string]any
// 		expect User
// 	}{
// 		{
// 			name: "nested pointer and map",
// 			input: map[string]any{
// 				"name":           "John",
// 				"address.street": "123 Main St",
// 				"address.city":   "Boston",
// 				"tags.color":     "blue",
// 				"tags.size":      "large",
// 			},
// 			expect: User{
// 				Name: "John",
// 				Address: &Address{
// 					Street: "123 Main St",
// 					City:   "Boston",
// 				},
// 				Tags: map[string]string{
// 					"color": "blue",
// 					"size":  "large",
// 				},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Test FromMap
// 			u1 := User{}
// 			t.Logf("Before FromMap: %+v", u1)
// 			err := NewBuilder(&u1).FromMap(tt.input).Error()
// 			if err != nil {
// 				t.Errorf("FromMap error: %v", err)
// 			}
// 			t.Logf("After FromMap: %+v", u1)
// 			if u1.Address != nil {
// 				t.Logf("Address after FromMap: {Street:%s City:%s}", u1.Address.Street, u1.Address.City)
// 			}
// 			if !reflect.DeepEqual(u1, tt.expect) {
// 				t.Errorf("FromMap got = %+v, want %+v", u1, tt.expect)
// 			}

// 			// Test FromKV
// 			u2 := User{}
// 			t.Logf("Before FromKV: %+v", u2)
// 			var kvs []KV
// 			for k, v := range tt.input {
// 				kvs = append(kvs, KV{k, v})
// 			}
// 			err = NewBuilder(&u2).FromKV(kvs...).Error()
// 			if err != nil {
// 				t.Errorf("FromKV error: %v", err)
// 			}
// 			t.Logf("After FromKV: %+v", u2)
// 			if u2.Address != nil {
// 				t.Logf("Address after FromKV: {Street:%s City:%s}", u2.Address.Street, u2.Address.City)
// 			}
// 			if !reflect.DeepEqual(u2, tt.expect) {
// 				t.Errorf("FromKV got = %+v, want %+v", u2, tt.expect)
// 			}

// 			// Test Fill
// 			u3 := &User{}
// 			t.Logf("Before Fill: %+v", u3)
// 			err = Fill(u3, kvs)
// 			if err != nil {
// 				t.Errorf("Fill error: %v", err)
// 			}
// 			t.Logf("After Fill: %+v", *u3)
// 			if u3.Address != nil {
// 				t.Logf("Address after Fill: {Street:%s City:%s}", u3.Address.Street, u3.Address.City)
// 			}
// 			if !reflect.DeepEqual(*u3, tt.expect) {
// 				t.Errorf("Fill got = %+v, want %+v", *u3, tt.expect)
// 			}
// 		})
// 	}
// }

// func BenchmarkNewBuilder(b *testing.B) {
// 	type TestStruct struct {
// 		String  string
// 		Int     int
// 		Bool    bool
// 		Float   float64
// 		Slice   []string
// 		Map     map[string]interface{}
// 		Pointer *string
// 		Nested  struct {
// 			Field string
// 		}
// 	}

// 	var s TestStruct
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		NewBuilder(&s)
// 	}
// }

// func BenchmarkKMapOperations(b *testing.B) {
// 	m := kmap.New[string, fieldInfo]()
// 	info := fieldInfo{
// 		offset: 123,
// 		typ:    reflect.TypeOf(""),
// 	}

// 	// Benchmark Set operation
// 	b.Run("Set", func(b *testing.B) {
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			key := fmt.Sprintf("field_%d", i)
// 			m.Set(key, info)
// 		}
// 	})

// 	// Benchmark Get operation
// 	b.Run("Get", func(b *testing.B) {
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			key := fmt.Sprintf("field_%d", i%100) // Reuse keys to ensure they exist
// 			_, _ = m.Get(key)
// 		}
// 	})

// 	// Benchmark Range operation
// 	b.Run("Range", func(b *testing.B) {
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			m.Range(func(k string, v fieldInfo) bool {
// 				return true
// 			})
// 		}
// 	})

// 	// Benchmark Delete operation
// 	b.Run("Delete", func(b *testing.B) {
// 		b.ReportAllocs()
// 		for i := 0; i < b.N; i++ {
// 			key := fmt.Sprintf("field_%d", i%100)
// 			m.Delete(key)
// 		}
// 	})
// }
