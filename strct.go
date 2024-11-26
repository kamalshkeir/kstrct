package kstrct

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"iter"

	"github.com/kamalshkeir/kmap"
)

var (
	ErrorExpectedPtr = errors.New("expected structToFill to be a pointer")
)

type fieldCache struct {
	names []string
}

var cacheFieldsIndex = kmap.New[string, *fieldCache]()

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

type FieldCtx struct {
	NumFields int
	Index     int
	Field     reflect.Value
	Name      string
	Value     any
	Type      string
	Tags      []string
}

func (ctx *FieldCtx) Reset() {
	ctx.NumFields = 0
	ctx.Index = 0
	ctx.Field = reflect.Value{}
	ctx.Name = ""
	ctx.Value = nil
	ctx.Type = ""
	ctx.Tags = ctx.Tags[:0]
}

var fieldCtxPool = sync.Pool{
	New: func() any {
		return &FieldCtx{
			Tags: []string{},
		}
	},
}

// From returns an iterator for struct fields
func From(strctPtr any, tagsToGet ...string) iter.Seq2[int, FieldCtx] {
	return func(yield func(int, FieldCtx) bool) {
		rs := reflect.ValueOf(strctPtr).Elem()
		rt := rs.Type()
		numFields := rs.NumField()
		strctName := rt.Name()

		// Get or create field index cache
		cache, ok := cacheFieldsIndex.Get(strctName)
		if !ok {
			cache = &fieldCache{
				names: make([]string, numFields),
			}
			cacheFieldsIndex.Set(strctName, cache)
		}

		for i := 0; i < numFields; i++ {
			f := rs.Field(i)
			var fname string
			if cache.names[i] == "" {
				fname = ToSnakeCase(rt.Field(i).Name)
				cache.names[i] = fname
			} else {
				fname = cache.names[i]
			}

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
				if ftag, ok := rt.Field(i).Tag.Lookup(t); ok {
					ctx.Tags = append(ctx.Tags, ftag)
				}
			}

			if !yield(i, *ctx) {
				ctx.Reset()
				fieldCtxPool.Put(ctx)
				return
			}
			ctx.Reset()
			fieldCtxPool.Put(ctx)
		}
	}
}

type KV struct {
	Key   string
	Value any
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

func Range[T any](strctPtr *T, fn func(c FieldCtx) bool, tagsToGet ...string) T {
	rs := reflect.ValueOf(strctPtr).Elem()
	rt := rs.Type()
	numFields := rs.NumField()
	strctName := rt.Name()

	// Get or create field index cache
	cache, ok := cacheFieldsIndex.Get(strctName)
	if !ok {
		cache = &fieldCache{
			names: make([]string, numFields),
		}
		cacheFieldsIndex.Set(strctName, cache)
	}

	for i := 0; i < numFields; i++ {
		f := rs.Field(i)
		var fname string
		if cache.names[i] == "" {
			fname = ToSnakeCase(rt.Field(i).Name)
			cache.names[i] = fname
		} else {
			fname = cache.names[i]
		}

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
			if ftag, ok := rt.Field(i).Tag.Lookup(t); ok {
				ctx.Tags = append(ctx.Tags, ftag)
			}
		}
		ok := fn(*ctx)
		if !ok {
			ctx.Reset()
			fieldCtxPool.Put(ctx)
			break
		}
		ctx.Reset()
		fieldCtxPool.Put(ctx)
	}
	return *strctPtr
}
