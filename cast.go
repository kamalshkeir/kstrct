package kstrct

import (
	"reflect"
)

func Cast[T any](value any, toType T) (T, error) {
	ret := new(T)
	err := SetReflectFieldValue(reflect.ValueOf(ret), value)
	if err != nil {
		return *new(T), err
	}
	return *ret, nil
}
