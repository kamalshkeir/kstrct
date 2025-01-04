package kstrct

import (
	"fmt"
	"reflect"
)

func InitSetterSQL() {
	// Register SQL types handler
	NewSetterCase(func(fld reflect.Value, value reflect.Value, valueI any) error {
		if Debug {
			fmt.Printf("\n=== SqlNull HANDLER DEBUG ===\n")
			fmt.Printf("SQL type: %s\n", fld.Type().String())
			fmt.Printf("Value type: %T\n", valueI)
			fmt.Printf("Value: %+v\n", valueI)
		}
		return handleSqlNull(fld, value)
	}, reflect.Struct)
}
