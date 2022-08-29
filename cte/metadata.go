package cte

import (
	"reflect"
)

type metaType string

const (
	metaTypeComputerKey metaType = "key"
	metaTypeComputer    metaType = "computer"
	metaTypeInout       metaType = "inout"
)

type MetadataProvider interface {
	CTEMetadata() interface{}
}

func extractMetadata(mp MetadataProvider, isComputerKey bool) parsedMetadata {
	result := make(map[metaType]reflect.Type)

	metadata := mp.CTEMetadata()
	if metadata == nil {
		panic(ErrNilMetadata.Err(reflect.TypeOf(mp)))
	}

	rt := reflect.TypeOf(metadata)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		result[metaType(field.Name)] = field.Type
	}

	if isComputerKey {
		result[metaTypeComputerKey] = reflect.TypeOf(mp)
	}

	return result
}

type parsedMetadata map[metaType]reflect.Type

func (pm parsedMetadata) getComputerKeyType() (reflect.Type, bool) {
	result, ok := pm[metaTypeComputerKey]
	return result, ok
}

func (pm parsedMetadata) getComputerType() (reflect.Type, bool) {
	result, ok := pm[metaTypeComputer]
	return result, ok
}

func (pm parsedMetadata) getInoutInterface() (reflect.Type, bool) {
	result, ok := pm[metaTypeInout]
	return result, ok
}
