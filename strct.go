package kstrct

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/kamalshkeir/kmap"
)

var (
	ErrorExpectedPtr = errors.New("expected value to Fill to be a pointer")
	ErrorDoublePtr   = errors.New("got double pointer")
	cacheFieldsIndex = kmap.New[string, *fieldCache]()
	kvPool           = sync.Pool{
		New: func() any {
			s := make([]KV, 0, 50)
			return &s
		},
	}
	Debug = false
)

type KV struct {
	Key   string
	Value any
}

type fieldCache struct {
	names []string
}

func TrySet(ptr, value any) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Pointer {
		return ErrorExpectedPtr
	}
	rv = rv.Elem()
	return SetReflectFieldValue(rv, value)
}

func FillFromMap(structOrChanPtr any, fields_values map[string]any, nested ...bool) (err error) {
	rs := reflect.ValueOf(structOrChanPtr)
	if rs.Kind() != reflect.Pointer && rs.Kind() == reflect.Struct {
		return ErrorExpectedPtr
	} else if rs.Kind() == reflect.Chan || rs.Elem().Kind() == reflect.Chan {
		if rs.Kind() == reflect.Pointer {
			rs = rs.Elem()
		}
		chanType := reflect.New(rs.Type().Elem()).Elem()
		err := SetReflectFieldValue(chanType, fields_values)
		if err != nil {
			return err
		}
		rs.Send(chanType)
		return nil
	}
	if rs.Kind() == reflect.Pointer {
		rs = rs.Elem()
		if rs.Kind() == reflect.Pointer {
			rs = reflect.New(rs.Type().Elem()).Elem()
		}
	} else {
		rs = rs.Elem()
	}
	rt := rs.Type()
	nnested := []string{}
	nn := make(map[string]any)
	for i := 0; i < rs.NumField(); i++ {
		field := rs.Field(i)
		fname := ToSnakeCase(rt.Field(i).Name)

		if v, ok := fields_values[fname]; ok {
			setErr := SetReflectFieldValue(field, v)
			if setErr != nil {
				fmt.Println("SetError:", err)
				err = errors.Join(err, setErr)
			}
			continue
		}
		if field.Kind() == reflect.Struct || (field.Kind() == reflect.Pointer && field.Elem().Kind() == reflect.Struct) {
			if field.Kind() == reflect.Pointer {
				field = field.Elem()
				if !field.IsValid() {
					return fmt.Errorf("field nested slice not valid %v", field.Kind())
				}
			}
			cp := make(map[string]any)
			for name, val := range fields_values {
				if sp := strings.Split(name, "."); len(sp) > 0 {
					if len(sp) == 2 && sp[0] == fname {
						cp[sp[1]] = val
					} else if len(sp) > 2 {
						found := false
						for _, n := range nnested {
							if n == sp[1] {
								found = true
							}
						}
						if !found {
							nnested = append(nnested, sp[1])
						}
						nn[sp[2]] = val
					}
				}
			}
			if len(cp) > 0 {
				setErr := SetReflectFieldValue(field, cp)
				err = errors.Join(err, setErr)
			}
			if len(nnested) > 0 {
				for _, n := range nnested {
					ff := field.FieldByName(SnakeCaseToTitle(n))
					if ff.Kind() == reflect.Struct || ff.Kind() == reflect.Pointer && ff.Elem().Kind() == reflect.Struct {
						setErr := SetReflectFieldValue(ff, nn)
						err = errors.Join(err, setErr)
					} else if len(nn) > 0 {
						if ff.Kind() == reflect.Pointer {
							ff = ff.Elem()
						}
						ne := reflect.New(ff.Type().Elem()).Elem()
						setErr := SetReflectFieldValue(ne, nn)
						err = errors.Join(err, setErr)
						if !ne.IsZero() {
							ff.Set(reflect.Append(ff, ne))
						}
					}
				}
			}
		} else if len(nested) > 0 && nested[0] && field.Kind() == reflect.Slice {
			if field.Kind() == reflect.Pointer {
				field = field.Elem()
				if !field.IsValid() {
					return fmt.Errorf("field nested slice not valid %v", field.Kind())
				}
			}
			cp := make(map[string]any)
			for name, val := range fields_values {
				if sp := strings.Split(name, "."); len(sp) > 0 {
					if len(sp) == 2 && sp[0] == fname {
						cp[sp[1]] = val
					} else if len(sp) > 2 {
						found := false
						for _, n := range nnested {
							if n == sp[1] {
								found = true
							}
						}
						if !found {
							nnested = append(nnested, sp[1])
						}
						nn[sp[2]] = val
					}
				}
			}
			var newElem reflect.Value
			if len(cp) > 0 {
				setErr := SetReflectFieldValue(newElem, cp)
				err = errors.Join(err, setErr)
			}

			// if len(cp) > 0 {
			// 	newElem = reflect.New(field.Type().Elem()).Elem()
			// 	setErr := SetReflectFieldValue(newElem, cp)
			// 	err = errors.Join(err, setErr)
			// 	field.Set(reflect.Append(field, newElem))
			// }

			isStrct := false
			if len(nnested) > 0 {
				for _, n := range nnested {
					ff := newElem.FieldByName(SnakeCaseToTitle(n))
					if ff.Kind() == reflect.Struct || ff.Kind() == reflect.Pointer && ff.Elem().Kind() == reflect.Struct {
						setErr := SetReflectFieldValue(ff, nn)
						err = errors.Join(err, setErr)
						isStrct = true
					} else if len(nn) > 0 {
						if ff.Kind() == reflect.Pointer {
							ff = ff.Elem()
						}
						ne := reflect.New(ff.Type().Elem()).Elem()
						setErr := SetReflectFieldValue(ne, nn)
						err = errors.Join(err, setErr)
						if !ne.IsZero() {
							ff.Set(reflect.Append(ff, ne))
							isStrct = true
						}
					}
				}
			}
			if isStrct {
				field.Set(reflect.Append(field, newElem))
			}
		}
	}
	return err
}

func FillFromKV(structOrChanPtr any, fields_values []KV, nested ...bool) (err error) {
	rs := reflect.ValueOf(structOrChanPtr)
	if rs.Kind() != reflect.Pointer && rs.Kind() == reflect.Struct {
		return ErrorExpectedPtr
	} else if rs.Kind() == reflect.Chan || rs.Elem().Kind() == reflect.Chan {
		if rs.Kind() == reflect.Pointer {
			rs = rs.Elem()
		}
		chanType := reflect.New(rs.Type().Elem()).Elem()
		err := SetReflectFieldValue(chanType, fields_values)
		if err != nil {
			return err
		}
		rs.Send(chanType)
		return nil
	}
	if rs.Kind() == reflect.Pointer {
		rs = rs.Elem()
	}
	rt := rs.Type()
	nnested := []string{}
	nn := make(map[string]any)
loop:
	for i := 0; i < rs.NumField(); i++ {
		field := rs.Field(i)
		fname := ToSnakeCase(rt.Field(i).Name)

		for _, v := range fields_values {
			v.Key = strings.TrimSpace(v.Key)
			if v.Key == fname {
				setErr := SetReflectFieldValue(field, v.Value)
				if setErr != nil {
					err = errors.Join(err, setErr)
				}
				continue loop
			}
		}

		if field.Kind() == reflect.Pointer {
			if !field.Elem().IsValid() {
				fieldNew := reflect.New(field.Type().Elem())
				field.Set(fieldNew)
			}
		}
		if field.Kind() == reflect.Struct || (field.Kind() == reflect.Pointer && field.Elem().Kind() == reflect.Struct) {
			if field.Kind() == reflect.Pointer {
				field = field.Elem()
				if !field.IsValid() {
					return fmt.Errorf("field nested slice not valid %v", field.Kind())
				}
			}
			cp := make(map[string]any)
			for _, val := range fields_values {
				if sp := strings.Split(val.Key, "."); len(sp) > 0 {
					if len(sp) == 2 && sp[0] == fname {
						cp[sp[1]] = val.Value
					} else if len(sp) > 2 {
						found := false
						for _, n := range nnested {
							if n == sp[1] {
								found = true
							}
						}
						if !found {
							nnested = append(nnested, sp[1])
						}
						nn[sp[2]] = val.Value
					}
				}
			}
			if len(cp) > 0 {
				setErr := SetReflectFieldValue(field, cp)
				err = errors.Join(err, setErr)
			}
			if len(nnested) > 0 {
				for _, n := range nnested {
					ff := field.FieldByName(SnakeCaseToTitle(n))
					if ff.Kind() == reflect.Struct || ff.Kind() == reflect.Pointer && ff.Elem().Kind() == reflect.Struct {
						setErr := SetReflectFieldValue(ff, nn)
						err = errors.Join(err, setErr)
					} else if len(nn) > 0 {
						if ff.Kind() == reflect.Pointer {
							ff = ff.Elem()
						}
						ne := reflect.New(ff.Type().Elem()).Elem()
						setErr := SetReflectFieldValue(ne, nn)
						err = errors.Join(err, setErr)
						if !ne.IsZero() {
							ff.Set(reflect.Append(ff, ne))
						}
					}
				}
			}
		} else if (len(nested) > 0 && nested[0]) && (field.Kind() == reflect.Slice || field.Kind() == reflect.Pointer) {
			if field.Kind() == reflect.Pointer {
				field = field.Elem()
				if !field.IsValid() {
					return fmt.Errorf("field nested slice not valid %v", field.Kind())
				}
			}
			cp := make(map[string]any)
			for _, val := range fields_values {
				if sp := strings.Split(val.Key, "."); len(sp) > 0 {
					if len(sp) == 2 && sp[0] == fname {
						cp[sp[1]] = val.Value
					} else if len(sp) > 2 {
						found := false
						for _, n := range nnested {
							if n == sp[1] {
								found = true
							}
						}
						if !found {
							nnested = append(nnested, sp[1])
						}
						nn[sp[2]] = val.Value
					}
				}
			}
			newElem := reflect.New(field.Type().Elem()).Elem()
			if len(cp) > 0 {
				setErr := SetReflectFieldValue(newElem, cp)
				err = errors.Join(err, setErr)
				field.Set(reflect.Append(field, newElem))
				continue loop
			}
			isStrct := false
			if len(nnested) > 0 {
				for _, n := range nnested {
					ff := newElem.FieldByName(SnakeCaseToTitle(n))
					if ff.Kind() == reflect.Struct || ff.Kind() == reflect.Pointer && ff.Elem().Kind() == reflect.Struct {
						setErr := SetReflectFieldValue(ff, nn)
						err = errors.Join(err, setErr)
						isStrct = true
					} else if len(nn) > 0 {
						if ff.Kind() == reflect.Pointer {
							ff = ff.Elem()
						}
						ne := reflect.New(ff.Type().Elem()).Elem()
						setErr := SetReflectFieldValue(ne, nn)
						err = errors.Join(err, setErr)
						if !ne.IsZero() {
							ff.Set(reflect.Append(ff, ne))
							isStrct = true
						}
					}
				}
			}
			if isStrct {
				field.Set(reflect.Append(field, newElem))
			}
		}
	}
	return err
}

func FillOLD(structOrChanPtr any, fields_values []KV, nested ...bool) (err error) {
	rs := reflect.ValueOf(structOrChanPtr)
	if rs.Kind() != reflect.Pointer && rs.Kind() == reflect.Struct {
		return ErrorExpectedPtr
	} else if rs.Kind() == reflect.Chan || rs.Elem().Kind() == reflect.Chan {
		if rs.Kind() == reflect.Pointer {
			rs = rs.Elem()
		}
		chanType := reflect.New(rs.Type().Elem()).Elem()
		err := SetReflectFieldValue(chanType, fields_values)
		if err != nil {
			return err
		}
		rs.Send(chanType)
		return nil
	}
	if rs.Kind() == reflect.Pointer {
		rs = rs.Elem()
	}
	rt := rs.Type()
	numFields := rs.NumField()

	cacheKey := rt.String()
	cache, ok := cacheFieldsIndex.Get(cacheKey)
	if !ok {
		cache = &fieldCache{
			names: make([]string, numFields),
		}
		// Pre-fill all names to avoid race conditions
		for i := 0; i < numFields; i++ {
			cache.names[i] = ToSnakeCase(rt.Field(i).Name)
		}
		cacheFieldsIndex.Set(cacheKey, cache)
	}

	// Safety check - recreate cache if size mismatch
	if len(cache.names) != numFields {
		cache = &fieldCache{
			names: make([]string, numFields),
		}
		for i := 0; i < numFields; i++ {
			cache.names[i] = ToSnakeCase(rt.Field(i).Name)
		}
		cacheFieldsIndex.Set(cacheKey, cache)
	}

loop:
	for i := 0; i < numFields; i++ {
		field := rs.Field(i)
		fname := cache.names[i]
		nestedKVs := []KV{}
		for _, v := range fields_values {
			v.Key = strings.TrimSpace(v.Key)
			if v.Key == fname {
				if err := SetReflectFieldValue(field, v.Value); err != nil {
					return err
				}
				continue loop
			} else if strings.HasPrefix(v.Key, fname) {
				// nested
				nestedKVs = append(nestedKVs, KV{
					Key:   strings.TrimPrefix(v.Key, fname+"."),
					Value: v.Value,
				})
			}
		}
		if len(nestedKVs) == 0 {
			continue loop
		}
		if field.Kind() == reflect.Pointer && field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		if field.Kind() == reflect.Struct || (field.Kind() == reflect.Pointer && field.Elem().Kind() == reflect.Struct) {
			err := SetReflectFieldValue(field, nestedKVs)
			if err != nil {
				if Debug {
					fmt.Println("err set struct", field.Interface(), nestedKVs)
				}
			}
			continue loop
		} else if field.Kind() == reflect.Slice || (field.Kind() == reflect.Pointer && field.Elem().Kind() == reflect.Slice) {
			if field.Kind() == reflect.Pointer {
				field = field.Elem()
				if !field.IsValid() {
					return fmt.Errorf("field nested slice not valid %v", field.Kind())
				}
			}
			err = SetReflectFieldValue(field, nestedKVs)
			if err != nil {
				if Debug {
					fmt.Println("err set slice", field.Interface(), nestedKVs)
				}
			}
		}
	}
	return err
}

var errNil = errors.New("got nil structOrChanPtr")

func Fill(structOrChanPtr any, fields_values []KV, nested ...bool) (err error) {
	if structOrChanPtr == nil {
		return errNil
	}
	rs := reflect.ValueOf(structOrChanPtr)
	if rs.Kind() != reflect.Pointer && rs.Kind() == reflect.Struct {
		return ErrorExpectedPtr
	}

	// Handle channels first
	if rs.Kind() == reflect.Chan || (rs.Kind() == reflect.Pointer && rs.Elem().Kind() == reflect.Chan) {
		if rs.Kind() == reflect.Pointer {
			rs = rs.Elem()
		}
		chanType := reflect.New(rs.Type().Elem()).Elem()
		err := SetReflectFieldValue(chanType, fields_values)
		if err != nil {
			return err
		}
		rs.Send(chanType)
		return nil
	}

	// For non-channel types, use builder with handleNestedField
	builder := NewBuilder(structOrChanPtr)
	builder.FromKV(fields_values...)
	return builder.Error()
}

func FillM(structOrChanPtr any, fields_values map[string]any, nested ...bool) (err error) {
	l := len(fields_values)
	if l == 0 {
		return nil
	}

	// Get pointer to slice from pool
	kvsp := kvPool.Get().(*[]KV)
	kvs := (*kvsp)[:0] // Reset length but keep capacity

	// Fill slice directly
	for k, v := range fields_values {
		switch mv := v.(type) {
		case map[string]string:
			for mk, mv := range mv {
				kvs = append(kvs, KV{k + "." + mk, mv})
			}
		case map[string]int:
			for mk, mv := range mv {
				kvs = append(kvs, KV{k + "." + mk, mv})
			}
		case map[string]any:
			for mk, mv := range mv {
				kvs = append(kvs, KV{k + "." + mk, mv})
			}
		default:
			kvs = append(kvs, KV{k, v})
		}
	}

	// Use our zero-allocation optimized function
	err = Fill(structOrChanPtr, kvs, nested...)

	// Return pointer to slice to pool
	kvPool.Put(kvsp)

	return err
}

// CreateStruct dynamically creates a struct type with the given fields and returns a pointer to a new instance
func CreateStruct(fields []StructField) (any, error) {
	if len(fields) == 0 {
		return nil, errors.New("no fields provided")
	}

	// Create struct fields
	structFields := make([]reflect.StructField, len(fields))
	for i, field := range fields {
		// Build tag string
		var tags string
		if len(field.Tags) > 0 {
			var tagParts []string
			for key, value := range field.Tags {
				tagParts = append(tagParts, fmt.Sprintf(`%s:"%s"`, key, value))
			}
			tags = strings.Join(tagParts, " ")
		}

		structFields[i] = reflect.StructField{
			Name: field.Name,
			Type: field.Type,
			Tag:  reflect.StructTag(tags),
		}
	}

	// Create the struct type
	structType := reflect.StructOf(structFields)

	// Create a new instance of the struct
	structValue := reflect.New(structType)

	// Set initial values if provided
	for i, field := range fields {
		if field.Value != nil {
			err := SetReflectFieldValue(structValue.Elem().Field(i), field.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to set field %s: %w", field.Name, err)
			}
		}
	}

	return structValue.Interface(), nil
}

// ExtractStructFields takes a struct or pointer to struct and returns its fields as []StructField
func ExtractStructFields(structOrPtr any) ([]StructField, error) {
	val := reflect.ValueOf(structOrPtr)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct or pointer to struct, got %v", val.Kind())
	}

	typ := val.Type()
	numFields := typ.NumField()
	fields := make([]StructField, numFields)

	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Extract tags
		tags := make(map[string]string)
		if field.Tag != "" {
			// Get all possible tags from the field
			for _, key := range []string{"json", "korm", "xml", "yaml", "toml", "db"} {
				if tag, ok := field.Tag.Lookup(key); ok {
					tags[key] = tag
				}
			}
		}

		// Create StructField
		fields[i] = StructField{
			Name:  field.Name,
			Type:  field.Type,
			Tags:  tags,
			Value: fieldValue.Interface(),
		}
	}

	return fields, nil
}
