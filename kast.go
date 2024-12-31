package kstrct

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// ModelInfo contains information about a model
type ModelInfo struct {
	Name   string
	Fields []StructField
}

// StructField represents a field in a dynamic struct
type StructField struct {
	Name         string
	Type         reflect.Type
	Tags         map[string]string
	Value        any
	IsSlice      bool
	SliceType    reflect.Type
	IsStruct     bool
	IsInterface  bool
	Implements   []reflect.Type
	StructFields []StructField
	Anonymous    bool // Whether this is an embedded field
}

// Helper function to get reflect.Type from type name
func getTypeFromString(typeName string) reflect.Type {
	switch typeName {
	case "string":
		return reflect.TypeOf("")
	case "int":
		return reflect.TypeOf(0)
	case "int8":
		return reflect.TypeOf(int8(0))
	case "int16":
		return reflect.TypeOf(int16(0))
	case "int32":
		return reflect.TypeOf(int32(0))
	case "int64":
		return reflect.TypeOf(int64(0))
	case "uint":
		return reflect.TypeOf(uint(0))
	case "uint8":
		return reflect.TypeOf(uint8(0))
	case "uint16":
		return reflect.TypeOf(uint16(0))
	case "uint32":
		return reflect.TypeOf(uint32(0))
	case "uint64":
		return reflect.TypeOf(uint64(0))
	case "float32":
		return reflect.TypeOf(float32(0))
	case "float64":
		return reflect.TypeOf(float64(0))
	case "bool":
		return reflect.TypeOf(false)
	case "byte":
		return reflect.TypeOf(byte(0))
	case "rune":
		return reflect.TypeOf(rune(0))
	case "error":
		return reflect.TypeOf((*error)(nil)).Elem()
	case "time.Time":
		return reflect.TypeOf(time.Time{})
	default:
		// For qualified types (e.g., package.Interface), try to parse them
		if strings.Contains(typeName, ".") {
			parts := strings.SplitN(typeName, ".", 2)
			if len(parts) == 2 {
				// For any interface type, return any and let the struct creation handle the type checking
				// The actual interface type will be verified during struct creation using reflection
				return reflect.TypeOf((*any)(nil)).Elem()
			}
		}
		// For unknown types, return any as a safe default
		return reflect.TypeOf((*any)(nil)).Elem()
	}
}

// Extract fields from AST struct type
func extractFieldsFromASTStruct(structType *ast.StructType) ([]StructField, error) {
	var fields []StructField
	for _, field := range structType.Fields.List {
		// Get field type
		var fieldType reflect.Type
		var isSlice bool
		var sliceType reflect.Type
		var isStruct bool
		var isInterface bool
		var structFields []StructField
		var implements []reflect.Type

		switch t := field.Type.(type) {
		case *ast.Ident:
			fieldType = getTypeFromString(t.Name)
			// Check if it's an interface type
			if fieldType != nil && fieldType.Kind() == reflect.Interface {
				isInterface = true
				implements = append(implements, fieldType)
			}
		case *ast.StructType:
			isStruct = true
			nestedFields, err := extractFieldsFromASTStruct(t)
			if err != nil {
				return nil, fmt.Errorf("error extracting nested struct fields: %w", err)
			}
			structFields = nestedFields
			fieldType = reflect.TypeOf((*any)(nil)).Elem() // Placeholder for struct type
		case *ast.ArrayType:
			isSlice = true
			switch elt := t.Elt.(type) {
			case *ast.Ident:
				sliceType = getTypeFromString(elt.Name)
				fieldType = reflect.SliceOf(sliceType)
			case *ast.StarExpr:
				if ident, ok := elt.X.(*ast.Ident); ok {
					sliceType = reflect.PointerTo(getTypeFromString(ident.Name))
					fieldType = reflect.SliceOf(sliceType)
				}
			case *ast.StructType:
				isStruct = true
				nestedFields, err := extractFieldsFromASTStruct(elt)
				if err != nil {
					return nil, fmt.Errorf("error extracting nested struct fields: %w", err)
				}
				structFields = nestedFields
				fieldType = reflect.SliceOf(reflect.TypeOf((*any)(nil)).Elem()) // Placeholder for struct slice
			}
		case *ast.StarExpr:
			switch x := t.X.(type) {
			case *ast.Ident:
				fieldType = reflect.PointerTo(getTypeFromString(x.Name))
			case *ast.ArrayType:
				if ident, ok := x.Elt.(*ast.Ident); ok {
					sliceType = getTypeFromString(ident.Name)
					fieldType = reflect.PointerTo(reflect.SliceOf(sliceType))
				}
			case *ast.SelectorExpr:
				if ident, ok := x.X.(*ast.Ident); ok {
					fieldType = reflect.PointerTo(getTypeFromString(ident.Name + "." + x.Sel.Name))
				}
			case *ast.StructType:
				isStruct = true
				nestedFields, err := extractFieldsFromASTStruct(x)
				if err != nil {
					return nil, fmt.Errorf("error extracting nested struct fields: %w", err)
				}
				structFields = nestedFields
				fieldType = reflect.PointerTo(reflect.TypeOf((*any)(nil)).Elem()) // Placeholder for pointer to struct
			}
		case *ast.SelectorExpr:
			// This is a qualified type (e.g., any interface type)
			if pkg, ok := t.X.(*ast.Ident); ok {
				fieldType = getTypeFromString(pkg.Name + "." + t.Sel.Name)
				if fieldType != nil && fieldType.Kind() == reflect.Interface {
					isInterface = true
					implements = append(implements, fieldType)
				}
			}
		case *ast.InterfaceType:
			isInterface = true
			fieldType = reflect.TypeOf((*any)(nil)).Elem()
		}

		var fieldName string
		if len(field.Names) > 0 {
			fieldName = field.Names[0].Name
		} else if isInterface {
			// For embedded interfaces, use their type name as the field name
			switch t := field.Type.(type) {
			case *ast.SelectorExpr:
				fieldName = t.Sel.Name
			case *ast.Ident:
				fieldName = t.Name
			}
		}

		// Extract tags if present
		tags := make(map[string]string)
		if field.Tag != nil {
			tag := strings.Trim(field.Tag.Value, "`")
			for _, key := range []string{"json", "korm", "xml", "yaml", "toml", "db"} {
				if v := reflect.StructTag(tag).Get(key); v != "" {
					tags[key] = v
				}
			}
		}

		// Create StructField
		structField := StructField{
			Name:         fieldName,
			Type:         fieldType,
			Tags:         tags,
			IsSlice:      isSlice,
			SliceType:    sliceType,
			IsStruct:     isStruct,
			IsInterface:  isInterface,
			StructFields: structFields,
			Implements:   implements,
			Anonymous:    len(field.Names) == 0 || isInterface, // Make embedded fields and interfaces anonymous
		}

		fields = append(fields, structField)
	}
	return fields, nil
}

// GetStructFromName finds a struct by name in the codebase and returns its schema and actual type
func GetStructFromName(rootDir string, structName string) (*ModelInfo, any, error) {
	p, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, nil, fmt.Errorf("expect Absolute Path: %w", err)
	}
	fset := token.NewFileSet()

	// Walk through all directories recursively
	var foundModel *ModelInfo
	var foundType reflect.Type

	// Split the struct name into package and type parts if it contains a dot
	var wantPkg, wantType string
	parts := strings.Split(structName, ".")
	if len(parts) == 2 {
		wantPkg, wantType = parts[0], parts[1]
	} else {
		wantType = structName
	}

	err = filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip vendor directory and hidden directories
		if info.IsDir() {
			if info.Name() == "vendor" || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Only parse .go files
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Parse the file
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil // Skip files with parse errors
		}

		// If package was specified, check if this is the package we want
		if wantPkg != "" && file.Name.Name != wantPkg {
			return nil
		}

		// Look for the struct declaration
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if typeSpec.Name.Name == wantType {
						if structType, ok := typeSpec.Type.(*ast.StructType); ok {
							// Found the struct, extract its fields
							fields, err := extractFieldsFromASTStruct(structType)
							if err != nil {
								return fmt.Errorf("error extracting fields: %w", err)
							}

							// Get package path
							dirPath := filepath.Dir(path)
							relPath, err := filepath.Rel(p, dirPath)
							if err != nil {
								return fmt.Errorf("error getting relative path for %s: %w", dirPath, err)
							}

							// Convert path separators to package format
							pkgPath := strings.ReplaceAll(relPath, string(filepath.Separator), "/")
							modelName := pkgPath + "." + wantType

							foundModel = &ModelInfo{
								Name:   modelName,
								Fields: fields,
							}

							// Create struct fields for reflection
							structFields := make([]reflect.StructField, len(fields))
							for i, field := range fields {
								var tags string
								if len(field.Tags) > 0 {
									var tagParts []string
									for key, value := range field.Tags {
										tagParts = append(tagParts, fmt.Sprintf(`%s:"%s"`, key, value))
									}
									tags = strings.Join(tagParts, " ")
								}

								structFields[i] = reflect.StructField{
									Name: field.Name,
									Type: field.Type,
									Tag:  reflect.StructTag(tags),
								}
							}

							// Create the struct type
							foundType = reflect.StructOf(structFields)

							return filepath.SkipDir // Stop walking once we find the struct
						}
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("error walking directory: %w", err)
	}

	if foundModel == nil {
		return nil, nil, fmt.Errorf("struct %s not found", structName)
	}

	instance := reflect.New(foundType).Interface()
	return foundModel, instance, nil
}
