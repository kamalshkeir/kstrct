package kstrct

import (
	"fmt"
	"reflect"

	"github.com/kamalshkeir/kmap"
)

// MethodFunc represents a method function that can be added to a struct
type MethodFunc any

// methodCache stores methods for each type
var methodCache = kmap.New[string, *kmap.SafeMap[string, MethodFunc]]()

// AddMethod adds a method to a struct type. The method must be of the form:
// func(receiver T) methodName([args ...]) [results...]
// where T is the struct type or *T for pointer receiver
func AddMethod(structPtr any, methodName string, methodFunc MethodFunc) error {
	// Get the type of the struct
	structValue := reflect.ValueOf(structPtr)
	if structValue.Kind() != reflect.Ptr {
		return fmt.Errorf("structPtr must be a pointer to a struct")
	}

	structType := structValue.Type()
	methodValue := reflect.ValueOf(methodFunc)
	methodType := methodValue.Type()

	// Validate method signature
	if methodType.NumIn() < 1 {
		return fmt.Errorf("method must have at least one input parameter (the receiver)")
	}

	receiverType := methodType.In(0)
	if receiverType != structType && receiverType != structType.Elem() {
		return fmt.Errorf("method receiver type %v does not match struct type %v", receiverType, structType)
	}

	// Get or create method map for this type
	typeKey := structType.String()
	var methodMap *kmap.SafeMap[string, MethodFunc]
	if v, ok := methodCache.Get(typeKey); ok {
		methodMap = v
	} else {
		methodMap = kmap.New[string, MethodFunc]()
		methodCache.Set(typeKey, methodMap)
	}

	// Store the method
	methodMap.Set(methodName, methodFunc)

	return nil
}

// CallMethod calls a previously added method on a struct
func CallMethod(structPtr any, methodName string, args ...any) ([]any, error) {
	// Get the type of the struct
	structValue := reflect.ValueOf(structPtr)
	if structValue.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("structPtr must be a pointer to a struct")
	}

	// Get method map for this type
	typeKey := structValue.Type().String()
	methodMap, ok := methodCache.Get(typeKey)
	if !ok {
		return nil, fmt.Errorf("no methods found for type %s", typeKey)
	}

	// Get the method
	methodI, ok := methodMap.Get(methodName)
	if !ok {
		return nil, fmt.Errorf("method %s not found", methodName)
	}
	method := reflect.ValueOf(methodI)

	// Prepare arguments
	methodType := method.Type()
	numIn := methodType.NumIn()
	if len(args)+1 != numIn {
		return nil, fmt.Errorf("wrong number of arguments: got %d, want %d", len(args), numIn-1)
	}

	callArgs := make([]reflect.Value, numIn)
	callArgs[0] = structValue // receiver

	for i, arg := range args {
		argValue := reflect.ValueOf(arg)
		if !argValue.Type().AssignableTo(methodType.In(i + 1)) {
			return nil, fmt.Errorf("argument %d has wrong type: got %v, want %v", i+1, argValue.Type(), methodType.In(i+1))
		}
		callArgs[i+1] = argValue
	}

	// Call the method
	results := method.Call(callArgs)

	// Convert results to interface slice
	returnValues := make([]any, len(results))
	for i, result := range results {
		returnValues[i] = result.Interface()
	}

	return returnValues, nil
}
