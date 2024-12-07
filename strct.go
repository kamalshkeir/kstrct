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

func Fill(structOrChanPtr any, fields_values []KV, nested ...bool) (err error) {
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
				fmt.Println("err set struct", field.Interface(), nestedKVs)
			}
			continue loop
		} else if field.Kind() == reflect.Slice || (field.Kind() == reflect.Pointer && field.Elem().Kind() == reflect.Slice) {
			if field.Kind() == reflect.Pointer {
				field = field.Elem()
			}
			err = SetReflectFieldValue(field, nestedKVs)
			if err != nil {
				fmt.Println("err set slice", field.Interface(), nestedKVs)
			}
		}
	}
	return err
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
		kvs = append(kvs, KV{k, v})
	}

	// Use our zero-allocation optimized function
	err = Fill(structOrChanPtr, kvs, nested...)

	// Return pointer to slice to pool
	kvPool.Put(kvsp)

	return err
}

// func FillPP(structOrChanPtr any, fields_values []KV, nested ...bool) (err error) {
// 	rs := reflect.ValueOf(structOrChanPtr)
// 	if rs.Kind() != reflect.Pointer && rs.Kind() == reflect.Struct {
// 		return ErrorExpectedPtr
// 	} else if rs.Kind() == reflect.Chan || rs.Elem().Kind() == reflect.Chan {
// 		if rs.Kind() == reflect.Pointer {
// 			rs = rs.Elem()
// 		}
// 		chanType := reflect.New(rs.Type().Elem()).Elem()
// 		err := SetReflectFieldValue(chanType, fields_values)
// 		if err != nil {
// 			return err
// 		}
// 		rs.Send(chanType)
// 		return nil
// 	}
// 	if rs.Kind() == reflect.Pointer {
// 		rs = rs.Elem()
// 	}
// 	rt := rs.Type()
// 	numFields := rs.NumField()

// 	// Use rt.String() as cache key - simpler and good enough
// 	cacheKey := rt.String()
// 	cache, ok := cacheFieldsIndex.Get(cacheKey)
// 	if !ok {
// 		cache = make(map[string]string, numFields)
// 		cacheFieldsIndex.Set(cacheKey, cache)
// 	}

// loop:
// 	for i := 0; i < numFields; i++ {
// 		field := rs.Field(i)
// 		fieldName := rt.Field(i).Name
// 		fname, ok := cache[fieldName]
// 		if !ok {
// 			fname = ToSnakeCase(fieldName)
// 			cache[fieldName] = fname
// 		}

// 		for _, v := range fields_values {
// 			if v.Key == fname {
// 				if err := SetReflectFieldValue(field, v.Value); err != nil {
// 					return err
// 				}
// 				continue loop
// 			}
// 		}
// 	}
// 	return err
// }
