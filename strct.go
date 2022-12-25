package kstrct

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrorLength = errors.New("error FillSelectedValues: len(values_to_fill) and len(struct fields) should be the same")
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

func FillFromMap(struct_to_fill any, fields_values map[string]any) error {
	rs := reflect.ValueOf(struct_to_fill)
	if rs.Kind() == reflect.Pointer {
		rs = reflect.ValueOf(struct_to_fill).Elem()
	}
	for k, v := range fields_values {
		var field *reflect.Value
		if f := rs.FieldByName(SnakeCaseToTitle(k)); f.IsValid() {
			field = &f
		} else if f := rs.FieldByName(k); f.IsValid() {
			field = &f
		}
		if field == nil {
			return errors.New("FillFromMap error: " + k + " not valid")
		}
		err := SetReflectFieldValue(*field, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func FillFromSelected(structToFill interface{}, fieldsCommaSeparated string, valuesToFill ...interface{}) error {
	rv := reflect.ValueOf(structToFill)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

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
