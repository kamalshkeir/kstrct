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


func FillFromValues(struct_to_fill any, values_to_fill ...any) error {
	rs := reflect.ValueOf(struct_to_fill)
	if rs.Kind() == reflect.Pointer {
		rs = reflect.ValueOf(struct_to_fill).Elem()
	}
	typeOfT := rs.Type()
	ignored := []int{}
	
	loop:
		for i := 0; i < rs.NumField(); i++ {	
			if ftag, ok := typeOfT.Field(i).Tag.Lookup("korm"); ok {
				if ftag == "-" {
					ignored = append(ignored, i)
					continue
				}
				if strings.Contains(ftag,"pk") || strings.Contains(ftag,"autoinc") {
					if len(values_to_fill) != rs.NumField() {
						ignored = append(ignored, i)
						continue loop
					}
				}
			}

			field := rs.Field(i)
			if field.IsValid() {
				index := i
				if len(values_to_fill) < rs.NumField() {
					index= i-len(ignored)
				}
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovered. Error:\n", r)
						fmt.Println("ignored:",ignored)
						fmt.Println("values:",values_to_fill)
						fmt.Println(ToSnakeCase(typeOfT.Field(index).Name),":",values_to_fill[index])
					}
				}()
				
				SetReflectFieldValue(field, values_to_fill[index])
			} else {
				return errors.New("FillFromValues error: "+ToSnakeCase(typeOfT.Field(i).Name)+" not valid")
			}
		}
	return nil
}

func FillFromMap(struct_to_fill any,fields_values map[string]any) error {	
	rs := reflect.ValueOf(struct_to_fill)
	if rs.Kind() == reflect.Pointer {
		rs = reflect.ValueOf(struct_to_fill).Elem()
	}
	for k,v := range fields_values {
		var field *reflect.Value
		if f := rs.FieldByName(SnakeCaseToTitle(k));f.IsValid() && f.CanSet() {
			field=&f
		} else if f := rs.FieldByName(k);f.IsValid() && f.CanSet() {
			field=&f
		}
		if field != nil {
			SetReflectFieldValue(*field, v)
		} else {
			return errors.New("FillFromValues error: "+k+" not valid")
		}
	}
	return nil
}

func FillFromSelected(struct_to_fill any, fields_comma_separated string, values_to_fill ...any) error {
	rs := reflect.ValueOf(struct_to_fill)
	if rs.Kind() == reflect.Pointer {
		rs = reflect.ValueOf(struct_to_fill).Elem()
	}
	typeOfT := rs.Type()
	skipped := 0
	for i := 0; i < rs.NumField(); i++ {
		if len(values_to_fill) < rs.NumField() && (!strings.Contains(fields_comma_separated,ToSnakeCase(typeOfT.Field(i).Name)) || !strings.Contains(fields_comma_separated,typeOfT.Field(i).Name)){
			skipped++
			continue
		} 
		field := rs.Field(i)
		if field.IsValid() {
			SetReflectFieldValue(field, values_to_fill[i-skipped])
		} else {
			return errors.New("FillFromValues error: "+ToSnakeCase(typeOfT.Field(i).Name)+" not valid")
		}
	}
	return nil
}




