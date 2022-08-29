package cte

import (
	"reflect"
	"strings"
)

func swallowErrPlanExecutionEndingEarly(err error) error {
	// Execution was intentionally ended by clients
	if err == ErrPlanExecutionEndingEarly || err == ErrRootPlanExecutionEndingEarly {
		return nil
	}

	return err
}

func extractFullNameFromValue(v interface{}) string {
	rt := extractUnderlyingType(reflect.ValueOf(v))

	return extractFullNameFromType(rt)
}

func extractFullNameFromType(t reflect.Type) string {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return t.PkgPath() + "/" + t.Name()
}

func extractShortName(fullName string) string {
	shortNameIdx := strings.LastIndex(fullName, "/")
	return fullName[shortNameIdx+1:]
}

func extractFieldTypes(field reflect.StructField) (isPointerType bool, valueType reflect.Type, pointerType reflect.Type) {
	rawFieldType := field.Type
	isPointerType = rawFieldType.Kind() == reflect.Pointer

	valueType = rawFieldType
	if isPointerType {
		valueType = rawFieldType.Elem()
	}

	pointerType = reflect.PointerTo(valueType)

	return
}

func extractUnderlyingType(v reflect.Value) reflect.Type {
	if v.Kind() == reflect.Pointer {
		return v.Elem().Type()
	}

	return v.Type()
}
