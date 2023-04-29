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
	ErrorExpectedPtr = errors.New("expected structToFill to be a pointer")
	cacheFieldsIndex = kmap.New[string, map[int]string](false)
)

func FillFromMap(structOrChanPtr any, fields_values map[string]any) (err error) {
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
	if rs.Kind() == reflect.Ptr {
		rs = rs.Elem()
		if rs.Kind() == reflect.Ptr {
			rs = reflect.New(rs.Type().Elem()).Elem()
		}
	} else {
		rs = rs.Elem()
	}
	rt := rs.Type()
	strctName := rt.Name()
	indexes, ok := cacheFieldsIndex.Get(strctName)
	if !ok {
		if rs.Kind() == reflect.Ptr {
			indexes = make(map[int]string)
		} else {
			indexes = make(map[int]string)
		}
		cacheFieldsIndex.Set(strctName, indexes)
	}

	for i := 0; i < rs.NumField(); i++ {
		field := rs.Field(i)
		var fname string
		if vf, ok := indexes[i]; !ok {
			fname = ToSnakeCase(rt.Field(i).Name)
			indexes[i] = fname
		} else {
			fname = vf
		}

		if v, ok := fields_values[fname]; ok {
			setErr := SetReflectFieldValue(field, v)
			if setErr != nil {
				fmt.Println("SetError:", err)
				err = errors.Join(err, setErr)
			}
			continue
		}
		if field.Kind() == reflect.Struct || (field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct) {
			cp := make(map[string]any)
			for name, val := range fields_values {
				if sp := strings.Split(name, "."); len(sp) == 2 {
					if sp[0] == fname || sp[0] == strctName {
						cp[sp[1]] = val
					}
				}
			}
			if len(cp) > 0 {
				setErr := SetReflectFieldValue(field, cp)
				err = errors.Join(err, setErr)
			}
		} else if field.Kind() == reflect.Slice {
			cp := make(map[string]any)
			for name, val := range fields_values {
				if sp := strings.Split(name, "."); len(sp) == 2 {
					if sp[0] == fname || sp[0] == strctName {
						cp[sp[1]] = val
					}
				}
			}
			var newElem reflect.Value
			if len(cp) > 0 {
				newElem = reflect.New(field.Type().Elem()).Elem()
				setErr := SetReflectFieldValue(newElem, cp)
				err = errors.Join(err, setErr)
				field.Set(reflect.Append(field, newElem))
			}
		}
	}
	return err
}

func FillByIndex(structOrChanPtr any, fields_values map[int]any) (err error) {
	rs := reflect.ValueOf(structOrChanPtr)
	if rs.Kind() != reflect.Pointer && rs.Kind() == reflect.Struct {
		return ErrorExpectedPtr
	} else if rs.Kind() == reflect.Chan || reflect.ValueOf(structOrChanPtr).Kind() == reflect.Chan {
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
	rs = rs.Elem()
	for i := range fields_values {
		field := rs.Field(i)
		err = SetReflectFieldValue(field, fields_values[i])
		if err != nil {
			return err
		}
	}
	return err
}

func FillFromMapS[T any](structOrChanPtr *T, fields_values map[string]any) (err error) {
	rs := reflect.ValueOf(structOrChanPtr).Elem()
	if rs.Kind() == reflect.Chan || reflect.ValueOf(structOrChanPtr).Kind() == reflect.Chan {
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
	if rs.Kind() == reflect.Ptr {
		rs = reflect.New(rs.Type().Elem()).Elem()
	}
	rt := rs.Type()
	strctName := rt.Name()
	indexes, ok := cacheFieldsIndex.Get(strctName)
	if !ok {
		indexes = make(map[int]string, rs.NumField())
		cacheFieldsIndex.Set(strctName, indexes)
	}
	for i := 0; i < rs.NumField(); i++ {
		field := rs.Field(i)
		var fname string
		if vf, ok := indexes[i]; !ok {
			fname = ToSnakeCase(rt.Field(i).Name)
			indexes[i] = fname
		} else {
			fname = vf
		}
		if v, ok := fields_values[fname]; ok {
			setErr := SetReflectFieldValue(field, v)
			if setErr != nil {
				err = errors.Join(err, setErr)
			}
			continue
		}
		if field.Kind() == reflect.Struct || field.Kind() == reflect.Ptr {
			cp := make(map[string]any)
			for name, val := range fields_values {
				if sp := strings.Split(name, "."); len(sp) == 2 {
					if sp[0] == fname || sp[0] == strctName {
						cp[sp[1]] = val
					}
				}
			}
			if len(cp) > 0 {
				setErr := SetReflectFieldValue(field, cp)
				err = errors.Join(err, setErr)
			}
		} else if field.Kind() == reflect.Slice {
			cp := make(map[string]any)
			for name, val := range fields_values {
				if sp := strings.Split(name, "."); len(sp) == 2 {
					if sp[0] == fname || sp[0] == strctName {
						cp[sp[1]] = val
					}
				}
			}
			var newElem reflect.Value
			if len(cp) > 0 {
				newElem = reflect.New(field.Type().Elem()).Elem()
				setErr := SetReflectFieldValue(newElem, cp)
				err = errors.Join(err, setErr)
				field.Set(reflect.Append(field, newElem))
			}
		}
	}
	if structOrChanPtr != new(T) {
		return err
	} else {
		return fmt.Errorf("pointer is nil")
	}
}

type FieldCtx struct {
	NumFields int
	Index     int
	Field     reflect.Value
	Name      string
	Value     any
	Type      string
	Tags      []string
}

var fieldCtxPool = sync.Pool{
	New: func() interface{} {
		return &FieldCtx{
			Tags: []string{},
		}
	},
}

func Range[T any](strctPtr *T, fn func(fCtx FieldCtx), tagsToGet ...string) T {
	rs := reflect.ValueOf(strctPtr).Elem()
	typeOfT := rs.Type()
	numFields := rs.NumField()
	for i := 0; i < numFields; i++ {
		f := rs.Field(i)
		fname := ToSnakeCase(typeOfT.Field(i).Name)

		// Get a fieldCtx from the pool
		ctx := fieldCtxPool.Get().(*FieldCtx)
		val := f.Interface()
		ctx.Field = f
		ctx.Name = fname
		ctx.Value = val
		ctx.Type = reflect.TypeOf(val).Name()
		ctx.NumFields = numFields
		ctx.Index = i
		ctx.Tags = ctx.Tags[:0]
		for _, t := range tagsToGet {
			if ftag, ok := typeOfT.Field(i).Tag.Lookup(t); ok {
				ctx.Tags = append(ctx.Tags, ftag)
			}
		}
		fn(*ctx)

		// Put the fieldCtx back into the pool
		fieldCtxPool.Put(ctx)
	}
	return *strctPtr
}

type KV struct {
	Key   string
	Value any
}

func FillFromKV[T any](structOrChanPtr *T, fields_values []KV) (err error) {
	rs := reflect.ValueOf(structOrChanPtr).Elem()
	if rs.Kind() == reflect.Chan || reflect.ValueOf(structOrChanPtr).Kind() == reflect.Chan {
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
	if rs.Kind() == reflect.Ptr {
		rs = reflect.New(rs.Type().Elem()).Elem()
	}
	rt := rs.Type()
	strctName := rt.Name()
	indexes, ok := cacheFieldsIndex.Get(strctName)
	if !ok {
		indexes = make(map[int]string, rs.NumField())
		cacheFieldsIndex.Set(strctName, indexes)
	}

loop:
	for i := 0; i < rs.NumField(); i++ {
		field := rs.Field(i)
		var fname string
		if vf, ok := indexes[i]; !ok {
			fname = ToSnakeCase(rt.Field(i).Name)
			indexes[i] = fname
		} else {
			fname = vf
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
		if field.Kind() == reflect.Struct || field.Kind() == reflect.Ptr {
			cp := make(map[string]any)
			for _, val := range fields_values {
				if sp := strings.Split(val.Key, "."); len(sp) == 2 {
					if sp[0] == fname || sp[0] == strctName {
						cp[sp[1]] = val.Value
					}
				}
			}
			if len(cp) > 0 {
				setErr := SetReflectFieldValue(field, cp)
				err = errors.Join(err, setErr)
			}
		} else if field.Kind() == reflect.Slice {
			cp := make(map[string]any)
			for _, val := range fields_values {
				if sp := strings.Split(val.Key, "."); len(sp) == 2 {
					if sp[0] == fname || sp[0] == strctName {
						cp[sp[1]] = val.Value
					}
				}
			}
			var newElem reflect.Value
			if len(cp) > 0 {
				newElem = reflect.New(field.Type().Elem()).Elem()
				setErr := SetReflectFieldValue(newElem, cp)
				err = errors.Join(err, setErr)
				field.Set(reflect.Append(field, newElem))
			}
		}
	}
	if structOrChanPtr != new(T) {
		return err
	} else {
		return fmt.Errorf("pointer is nil")
	}
}
