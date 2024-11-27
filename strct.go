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
			s := make([]KV, 0, 32) // Most structs have < 32 fields
			return &s              // Return pointer to slice
		},
	}
)

type fieldCache struct {
	names []string
}
type KV struct {
	Key   string
	Value any
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
	strctName := rt.Name()
	indexes, ok := cacheFieldsIndex.Get(strctName)
	if !ok {
		indexes = &fieldCache{
			names: make([]string, rs.NumField()),
		}
		cacheFieldsIndex.Set(strctName, indexes)
	}
	nnested := []string{}
	nn := make(map[string]any)
	for i := 0; i < rs.NumField(); i++ {
		field := rs.Field(i)
		var fname string
		if indexes.names[i] == "" {
			fname = ToSnakeCase(rt.Field(i).Name)
			indexes.names[i] = fname
		} else {
			fname = indexes.names[i]
		}

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
	strctName := rt.Name()
	numFields := rs.NumField()

	// Get or create field index cache
	cache, ok := cacheFieldsIndex.Get(strctName)
	if !ok {
		cache = &fieldCache{
			names: make([]string, numFields),
		}
		cacheFieldsIndex.Set(strctName, cache)
	}

	nnested := []string{}
	nn := make(map[string]any)
loop:
	for i := 0; i < numFields; i++ {
		field := rs.Field(i)
		var fname string
		if cache.names[i] == "" {
			fname = ToSnakeCase(rt.Field(i).Name)
			cache.names[i] = fname
		} else {
			fname = cache.names[i]
		}

		for _, v := range fields_values {
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
	strctName := rt.Name()
	numFields := rs.NumField()

	indexes, ok := cacheFieldsIndex.Get(strctName)
	if !ok {
		indexes = &fieldCache{
			names: make([]string, numFields),
		}
		cacheFieldsIndex.Set(strctName, indexes)
	}

	// Process non-nested fields first
	for _, v := range fields_values {
		if !strings.Contains(v.Key, ".") {
			i := 0
			for ; i < numFields; i++ {
				if indexes.names[i] == "" {
					indexes.names[i] = ToSnakeCase(rt.Field(i).Name)
				}
				if indexes.names[i] == v.Key {
					if err := SetReflectFieldValue(rs.Field(i), v.Value); err != nil {
						return err
					}
					break
				}
			}
		}
	}

	// Process nested fields if needed
	if len(nested) > 0 && nested[0] {
		nnested := []string{}
		nn := make(map[string]any)
		for i := 0; i < numFields; i++ {
			field := rs.Field(i)
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
						if len(sp) == 2 && sp[0] == indexes.names[i] {
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
						if len(sp) == 2 && sp[0] == indexes.names[i] {
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
