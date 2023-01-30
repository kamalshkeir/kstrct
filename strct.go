package kstrct

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var (
	ErrorLength      = errors.New("error FillSelectedValues: len(values_to_fill) and len(struct fields) should be the same")
	ErrorExpectedPtr = errors.New("expected structToFill to be a pointer")
)

func FillFromValues(structToFill interface{}, valuesToFill ...interface{}) error {
	rs := reflect.ValueOf(structToFill)
	if rs.Kind() == reflect.Pointer {
		rs = reflect.ValueOf(structToFill).Elem()
	}
	typeOfT := rs.Type()
	fieldsIndexes := make([]int, 0, rs.NumField())

	for i := 0; i < rs.NumField(); i++ {
		field := typeOfT.Field(i)
		kormTag, kormOk := field.Tag.Lookup("korm")
		kstrctTag, kstrctOk := field.Tag.Lookup("kstrct")

		if kormOk {
			switch kormTag {
			case "pk", "autoinc", "-":
				continue
			case "m2m":
				if !strings.Contains(kormTag, "m2m") {
					fieldsIndexes = append(fieldsIndexes, i)
				}
			default:
				fieldsIndexes = append(fieldsIndexes, i)
			}
		} else if kstrctOk && kstrctTag != "-" {
			fieldsIndexes = append(fieldsIndexes, i)
		} else if len(valuesToFill) < rs.NumField() {
			if i != 0 {
				fieldsIndexes = append(fieldsIndexes, i)
			}
		}
	}

	for i, fi := range fieldsIndexes {
		idx := i
		if ftag, ok := typeOfT.Field(fi).Tag.Lookup("korm"); ok {
			if ftag == "pk" || ftag == "autoinc" || ftag == "-" || strings.Contains(ftag, "m2m") {
				if len(valuesToFill) < len(fieldsIndexes) {
					continue
				}
			}
		} else if ftag, ok := typeOfT.Field(fi).Tag.Lookup("kstrct"); ok {
			if ftag == "-" && len(valuesToFill) < len(fieldsIndexes) {
				continue
			}
		}

		field := rs.Field(fi)
		if field.IsValid() {
			if len(valuesToFill) < len(fieldsIndexes) {
				idx = i - 1
			}
			err := SetReflectFieldValue(field, valuesToFill[idx])
			if err != nil {
				return err
			}
		} else {
			return errors.New("FillFromValues error: " + ToSnakeCase(typeOfT.Field(fi).Name) + " not valid")
		}
	}
	return nil
}

func FillFromMap(structToFill any, fields_values map[string]any) error {
	rs := reflect.ValueOf(structToFill)
	if rs.Kind() != reflect.Pointer {
		return ErrorExpectedPtr
	}
	rs = rs.Elem()
	for k, v := range fields_values {
		var field *reflect.Value
		if f := rs.FieldByName(SnakeCaseToTitle(k)); f.IsValid() {
			field = &f
		} else if f := rs.FieldByName(k); f.IsValid() {
			field = &f
		}
		if field == nil {
			return fmt.Errorf("fillFromMap error: %s not valid", k)
		}
		err := SetReflectFieldValue(*field, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func FillFromMapS[T any](fields_values map[string]any, model ...T) (T, error) {
	ptr := new(T)
	rs := reflect.ValueOf(ptr).Elem()

	for k, v := range fields_values {
		var field *reflect.Value
		if f := rs.FieldByName(SnakeCaseToTitle(k)); f.IsValid() {
			field = &f
		} else if f := rs.FieldByName(k); f.IsValid() {
			field = &f
		}
		if field == nil {
			return *new(T), fmt.Errorf("fillFromMap error: %s not valid", k)
		}
		err := SetReflectFieldValue(*field, v)
		if err != nil {
			return *new(T), err
		}
	}
	if ptr != new(T) {
		return *ptr, nil
	} else {
		return *new(T), fmt.Errorf("pointer is nil")
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

func FillFromSelected(structToFill interface{}, fieldsCommaSeparated string, valuesToFill ...interface{}) error {
	rv := reflect.ValueOf(structToFill)
	if rv.Kind() != reflect.Pointer {
		return ErrorExpectedPtr
	}
	rv = rv.Elem()
	rt := rv.Type()

	skipped := 0
	for i := 0; i < rv.NumField(); i++ {
		fieldName := ToSnakeCase(rt.Field(i).Name)
		if fi := strings.Index(fieldsCommaSeparated, fieldName); fi == -1 {
			skipped++
			continue
		}

		field := rv.Field(i)
		if !field.IsValid() {
			return fmt.Errorf("FillFromValues error: %s not valid", fieldName)
		}

		if err := SetReflectFieldValue(field, valuesToFill[i-skipped]); err != nil {
			return err
		}
	}

	return nil
}
