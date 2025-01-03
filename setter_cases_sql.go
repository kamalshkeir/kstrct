package kstrct

import (
	"database/sql"
	"reflect"
	"strconv"
	"time"
)

func InitSetterSQL() {
	// Register SQL types handler
	NewSetterCase(func(fld reflect.Value, value reflect.Value, valueI any) error {
		// Get the field's actual type
		fieldType := fld.Type()

		// Check if it's a SQL type and handle accordingly
		switch fieldType {
		case reflect.TypeOf(sql.NullString{}):
			ns := sql.NullString{Valid: true}
			switch v := valueI.(type) {
			case string:
				ns.String = v
			case *string:
				if v != nil {
					ns.String = *v
				} else {
					ns.Valid = false
				}
			case []byte:
				ns.String = string(v)
			case nil:
				ns.Valid = false
			default:
				if value.IsValid() {
					ns.String = value.String()
				} else {
					ns.Valid = false
				}
			}
			fld.Set(reflect.ValueOf(ns))
			return nil

		case reflect.TypeOf(sql.NullInt64{}):
			ni := sql.NullInt64{Valid: true}
			switch v := valueI.(type) {
			case int64:
				ni.Int64 = v
			case int:
				ni.Int64 = int64(v)
			case string:
				if i, err := strconv.ParseInt(v, 10, 64); err == nil {
					ni.Int64 = i
				} else {
					ni.Valid = false
				}
			case nil:
				ni.Valid = false
			default:
				if value.IsValid() && value.Type().ConvertibleTo(reflect.TypeOf(int64(0))) {
					ni.Int64 = value.Convert(reflect.TypeOf(int64(0))).Interface().(int64)
				} else {
					ni.Valid = false
				}
			}
			fld.Set(reflect.ValueOf(ni))
			return nil

		case reflect.TypeOf(sql.NullFloat64{}):
			nf := sql.NullFloat64{Valid: true}
			switch v := valueI.(type) {
			case float64:
				nf.Float64 = v
			case float32:
				nf.Float64 = float64(v)
			case int:
				nf.Float64 = float64(v)
			case int64:
				nf.Float64 = float64(v)
			case string:
				if f, err := strconv.ParseFloat(v, 64); err == nil {
					nf.Float64 = f
				} else {
					nf.Valid = false
				}
			case nil:
				nf.Valid = false
			default:
				if value.IsValid() && value.Type().ConvertibleTo(reflect.TypeOf(float64(0))) {
					nf.Float64 = value.Convert(reflect.TypeOf(float64(0))).Interface().(float64)
				} else {
					nf.Valid = false
				}
			}
			fld.Set(reflect.ValueOf(nf))
			return nil

		case reflect.TypeOf(sql.NullBool{}):
			nb := sql.NullBool{Valid: true}
			switch v := valueI.(type) {
			case bool:
				nb.Bool = v
			case string:
				if b, err := strconv.ParseBool(v); err == nil {
					nb.Bool = b
				} else {
					nb.Valid = false
				}
			case int:
				nb.Bool = v != 0
			case nil:
				nb.Valid = false
			default:
				if value.IsValid() && value.Type().ConvertibleTo(reflect.TypeOf(true)) {
					nb.Bool = value.Convert(reflect.TypeOf(true)).Interface().(bool)
				} else {
					nb.Valid = false
				}
			}
			fld.Set(reflect.ValueOf(nb))
			return nil

		case reflect.TypeOf(sql.NullTime{}):
			nt := sql.NullTime{Valid: true}
			switch v := valueI.(type) {
			case time.Time:
				nt.Time = v
			case string:
				if t, err := parseTimeString(v); err == nil {
					nt.Time = t
				} else {
					nt.Valid = false
				}
			case int64:
				nt.Time = time.Unix(v, 0)
			case nil:
				nt.Valid = false
			default:
				if value.IsValid() && value.Type().ConvertibleTo(reflect.TypeOf(time.Time{})) {
					nt.Time = value.Convert(reflect.TypeOf(time.Time{})).Interface().(time.Time)
				} else {
					nt.Valid = false
				}
			}
			fld.Set(reflect.ValueOf(nt))
			return nil
		}

		// Not a SQL type, let other handlers deal with it
		return nil
	}, reflect.Struct)
}
