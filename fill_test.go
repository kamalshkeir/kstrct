package kstrct

// // Test structs
// type BasicTypes struct {
// 	String    string
// 	Int       int
// 	Int8      int8
// 	Int16     int16
// 	Int32     int32
// 	Int64     int64
// 	Uint      uint
// 	Uint8     uint8
// 	Uint16    uint16
// 	Uint32    uint32
// 	Uint64    uint64
// 	Float32   float32
// 	Float64   float64
// 	Bool      bool
// 	Time      time.Time
// 	StringPtr *string
// 	IntPtr    *int
// 	FloatPtr  *float64
// 	BoolPtr   *bool
// 	TimePtr   *time.Time
// }

// type NestedStruct struct {
// 	Field1 string
// 	Field2 int
// }

// type ComplexTypes struct {
// 	StringSlice    []string
// 	IntSlice       []int
// 	FloatSlice     []float64
// 	BoolSlice      []bool
// 	InterfaceSlice []any
// 	StringMap      map[string]string
// 	IntMap         map[string]int
// 	FloatMap       map[string]float64
// 	BoolMap        map[string]bool
// 	InterfaceMap   map[string]any
// 	Nested         NestedStruct
// 	NestedPtr      *NestedStruct
// 	StructSlice    []NestedStruct
// }

// type SQLTypes struct {
// 	NullString  sql.NullString
// 	NullInt64   sql.NullInt64
// 	NullFloat64 sql.NullFloat64
// 	NullBool    sql.NullBool
// 	NullTime    sql.NullTime
// }

// type User struct {
// 	ID       int
// 	Email    string
// 	Password string
// 	Profile  *UserProfile
// 	Settings map[string]string
// }

// type UserProfile struct {
// 	FirstName string
// 	LastName  string
// 	Age       int
// 	Address   *Address
// }

// type Address struct {
// 	Street  string
// 	City    string
// 	Country string
// }

// func TestFill(t *testing.T) {
// 	t.Run("Basic types", func(t *testing.T) {
// 		var s BasicTypes
// 		kvs := []KV{
// 			{Key: "string", Value: "test"},
// 			{Key: "int", Value: 42},
// 			{Key: "int8", Value: int8(8)},
// 			{Key: "int16", Value: int16(16)},
// 			{Key: "int32", Value: int32(32)},
// 			{Key: "int64", Value: int64(64)},
// 			{Key: "uint", Value: uint(42)},
// 			{Key: "uint8", Value: uint8(8)},
// 			{Key: "uint16", Value: uint16(16)},
// 			{Key: "uint32", Value: uint32(32)},
// 			{Key: "uint64", Value: uint64(64)},
// 			{Key: "float32", Value: float32(3.14)},
// 			{Key: "float64", Value: 3.14159},
// 			{Key: "bool", Value: true},
// 			{Key: "time", Value: time.Now()},
// 		}
// 		err := Fill(&s, kvs)
// 		assert.NoError(t, err)
// 	})

// 	t.Run("Pointer types", func(t *testing.T) {
// 		var s BasicTypes
// 		str := "pointer"
// 		num := 42
// 		f := 3.14
// 		b := true
// 		tm := time.Now()
// 		kvs := []KV{
// 			{Key: "string_ptr", Value: &str},
// 			{Key: "int_ptr", Value: &num},
// 			{Key: "float_ptr", Value: &f},
// 			{Key: "bool_ptr", Value: &b},
// 			{Key: "time_ptr", Value: &tm},
// 		}
// 		err := Fill(&s, kvs)
// 		assert.NoError(t, err)
// 		assert.Equal(t, str, *s.StringPtr)
// 		assert.Equal(t, num, *s.IntPtr)
// 		assert.Equal(t, f, *s.FloatPtr)
// 		assert.Equal(t, b, *s.BoolPtr)
// 		assert.Equal(t, tm, *s.TimePtr)
// 	})

// 	t.Run("Slice types", func(t *testing.T) {
// 		var s ComplexTypes
// 		kvs := []KV{
// 			{Key: "string_slice", Value: "a, b,c"},
// 			{Key: "int_slice", Value: "1, 2, 3"},
// 			{Key: "float_slice", Value: "1.1,2.2,3.3"},
// 			{Key: "bool_slice", Value: "true,false,true"},
// 		}
// 		err := Fill(&s, kvs)
// 		assert.NoError(t, err)
// 		assert.Equal(t, []string{"a", "b", "c"}, s.StringSlice)
// 		assert.Equal(t, []int{1, 2, 3}, s.IntSlice)
// 		assert.Equal(t, []float64{1.1, 2.2, 3.3}, s.FloatSlice)
// 		assert.Equal(t, []bool{true, false, true}, s.BoolSlice)
// 	})

// 	t.Run("Map types", func(t *testing.T) {
// 		var s ComplexTypes
// 		kvs := []KV{
// 			{Key: "string_map.key1", Value: "value1"},
// 			{Key: "string_map.key2", Value: "value2"},
// 			{Key: "int_map.key1", Value: 1},
// 			{Key: "int_map.key2", Value: 2},
// 			{Key: "float_map.key1", Value: 1.1},
// 			{Key: "float_map.key2", Value: 2.2},
// 			{Key: "bool_map.key1", Value: true},
// 			{Key: "bool_map.key2", Value: false},
// 		}
// 		err := Fill(&s, kvs)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "value1", s.StringMap["key1"])
// 		assert.Equal(t, 1, s.IntMap["key1"])
// 		assert.Equal(t, 1.1, s.FloatMap["key1"])
// 		assert.Equal(t, true, s.BoolMap["key1"])
// 	})

// 	t.Run("Nested struct", func(t *testing.T) {
// 		var s ComplexTypes
// 		kvs := []KV{
// 			{Key: "nested.field1", Value: "nested string"},
// 			{Key: "nested.field2", Value: 456},
// 		}
// 		err := Fill(&s, kvs)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "nested string", s.Nested.Field1)
// 		assert.Equal(t, 456, s.Nested.Field2)
// 	})

// 	t.Run("Nested pointer struct", func(t *testing.T) {
// 		var s ComplexTypes
// 		kvs := []KV{
// 			{Key: "nested_ptr.field1", Value: "nested string"},
// 			{Key: "nested_ptr.field2", Value: 456},
// 		}
// 		err := Fill(&s, kvs)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, s.NestedPtr)
// 		assert.Equal(t, "nested string", s.NestedPtr.Field1)
// 		assert.Equal(t, 456, s.NestedPtr.Field2)
// 	})

// 	t.Run("Slice of structs", func(t *testing.T) {
// 		var s ComplexTypes
// 		kvs := []KV{
// 			{Key: "struct_slice.0.field1", Value: "first"},
// 			{Key: "struct_slice.0.field2", Value: 1},
// 			{Key: "struct_slice.1.field1", Value: "second"},
// 			{Key: "struct_slice.1.field2", Value: 2},
// 		}
// 		err := Fill(&s, kvs)
// 		assert.NoError(t, err)
// 		assert.Len(t, s.StructSlice, 2)
// 		assert.Equal(t, "first", s.StructSlice[0].Field1)
// 		assert.Equal(t, 1, s.StructSlice[0].Field2)
// 		assert.Equal(t, "second", s.StructSlice[1].Field1)
// 		assert.Equal(t, 2, s.StructSlice[1].Field2)
// 	})

// 	t.Run("SQL types", func(t *testing.T) {
// 		var s SQLTypes
// 		kvs := []KV{
// 			{Key: "null_string", Value: "test"},
// 			{Key: "null_int64", Value: 42},
// 			{Key: "null_float64", Value: 3.14},
// 			{Key: "null_bool", Value: true},
// 			{Key: "null_time", Value: time.Now()},
// 		}
// 		err := Fill(&s, kvs)
// 		assert.NoError(t, err)
// 		assert.True(t, s.NullString.Valid)
// 		assert.Equal(t, "test", s.NullString.String)
// 		assert.True(t, s.NullInt64.Valid)
// 		assert.Equal(t, int64(42), s.NullInt64.Int64)
// 		assert.True(t, s.NullFloat64.Valid)
// 		assert.Equal(t, 3.14, s.NullFloat64.Float64)
// 		assert.True(t, s.NullBool.Valid)
// 		assert.True(t, s.NullBool.Bool)
// 		assert.True(t, s.NullTime.Valid)
// 	})

// 	t.Run("KORM-style mapping", func(t *testing.T) {
// 		var user User
// 		kvs := []KV{
// 			{Key: "id", Value: 1},
// 			{Key: "email", Value: "test@example.com"},
// 			{Key: "password", Value: "secret"},
// 			{Key: "profile.first_name", Value: "John"},
// 			{Key: "profile.last_name", Value: "Doe"},
// 			{Key: "profile.age", Value: 30},
// 			{Key: "profile.address.street", Value: "123 Main St"},
// 			{Key: "profile.address.city", Value: "New York"},
// 			{Key: "profile.address.country", Value: "USA"},
// 			{Key: "settings.theme", Value: "dark"},
// 			{Key: "settings.language", Value: "en"},
// 		}
// 		err := Fill(&user, kvs)
// 		assert.NoError(t, err)
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
// 	})

// 	t.Run("Type conversions", func(t *testing.T) {
// 		var s BasicTypes
// 		kvs := []KV{
// 			{Key: "string", Value: 123},      // int to string
// 			{Key: "int", Value: "123"},       // string to int
// 			{Key: "float64", Value: "3.14"},  // string to float
// 			{Key: "bool", Value: 1},          // int to bool
// 			{Key: "time", Value: 1640995200}, // unix timestamp to time
// 		}
// 		err := Fill(&s, kvs)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "123", s.String)
// 		assert.Equal(t, 123, s.Int)
// 		assert.Equal(t, 3.14, s.Float64)
// 		assert.True(t, s.Bool)
// 		assert.Equal(t, time.Unix(1640995200, 0), s.Time)
// 	})

// 	t.Run("Edge cases", func(t *testing.T) {
// 		var s BasicTypes
// 		kvs := []KV{
// 			{Key: "int", Value: "invalid"},  // invalid number
// 			{Key: "bool", Value: "invalid"}, // invalid bool
// 			{Key: "time", Value: "invalid"}, // invalid time
// 		}
// 		err := Fill(&s, kvs)
// 		assert.Error(t, err) // Should return error for invalid conversions
// 	})

// 	// Additional test cases specific to old package
// 	t.Run("Map filling", func(t *testing.T) {
// 		var s ComplexTypes
// 		m := map[string]any{
// 			"string_map": map[string]string{
// 				"key1": "value1",
// 				"key2": "value2",
// 			},
// 			"int_map": map[string]int{
// 				"key1": 1,
// 				"key2": 2,
// 			},
// 		}
// 		err := FillM(&s, m)
// 		assert.NoError(t, err)
// 		assert.Equal(t, "value1", s.StringMap["key1"])
// 		assert.Equal(t, 1, s.IntMap["key1"])
// 	})

// 	t.Run("Channel filling", func(t *testing.T) {
// 		ch := make(chan BasicTypes, 1)
// 		kvs := []KV{
// 			{Key: "string", Value: "test"},
// 			{Key: "int", Value: 42},
// 		}
// 		err := Fill(ch, kvs)
// 		assert.NoError(t, err)
// 		s := <-ch
// 		assert.Equal(t, "test", s.String)
// 		assert.Equal(t, 42, s.Int)
// 	})

// 	t.Run("Pointer to channel filling", func(t *testing.T) {
// 		ch := make(chan BasicTypes, 1)
// 		chPtr := &ch
// 		kvs := []KV{
// 			{Key: "string", Value: "test"},
// 			{Key: "int", Value: 42},
// 		}
// 		err := Fill(chPtr, kvs)
// 		assert.NoError(t, err)
// 		s := <-ch
// 		assert.Equal(t, "test", s.String)
// 		assert.Equal(t, 42, s.Int)
// 	})
// }
