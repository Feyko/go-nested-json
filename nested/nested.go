package nested

import (
	"encoding/json"
	"errors"
	"fmt"
	"nested-json/internal/tags"
	"reflect"
)

func Marshal(v any) ([]byte, error) {
	marshaler := NewMarshaler(v)
	return marshaler.MarshalJSON()
}

func NewMarshaler(v any) Marshaler {
	m := Marshaler{input: v, reshaped: make(map[string]any)}
	m.reflect()
	return m
}

type Marshaler struct {
	input      any
	inputType  reflect.Type
	inputValue reflect.Value
	reshaped   map[string]any
	tags       map[string]tags.TagInfo
}

func (m *Marshaler) MarshalJSON() ([]byte, error) {
	if !isStruct(m.input) {
		return json.Marshal(m.input)
	}
	structTags, err := tags.GetTags(m.inputType)
	if err != nil {
		return nil, fmt.Errorf("could not parse the struct tags: %v", err)
	}
	m.tags = structTags
	err = m.reshape()
	if err != nil {
		return nil, err
	}
	return json.Marshal(m.reshaped)
}

func (m *Marshaler) reflect() {
	m.inputType = reflect.TypeOf(m.input)
	m.inputValue = reflect.ValueOf(m.input)
}

func (m *Marshaler) reshape() error {
	for fieldName, tag := range m.tags {
		field := m.inputValue.FieldByName(fieldName)
		if isEmpty(field) {
			continue
		}
		err := addFieldToReshape(m.reshaped, tag.Path, field.Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

func addFieldToReshape(reshaped map[string]any, path []string, v any) error {
	fieldName := path[0]
	field := reshaped[fieldName]
	if len(path) == 1 {
		if field != nil {
			return NewFieldCollisionError(fieldName)
		}
		reshaped[fieldName] = v
		return nil
	}
	mapV, isMap := field.(map[string]any)
	if field != nil && isMap {
		err := addFieldToReshape(mapV, path[1:], v)
		if collision, ok := err.(FieldCollisionError); ok {
			return collision
		}
		return err
	}
	if field == nil {
		reshaped[fieldName] = make(map[string]any)
		return addFieldToReshape(reshaped[fieldName].(map[string]any), path[1:], v)
	}
	return errors.New("invalid field state")
}

func isStruct(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Struct
}

type FieldCollisionError struct {
	path []string
}

func (err FieldCollisionError) Error() string {
	pathLength := len(err.path)
	pathstring := ""
	for i := 0; i < pathLength; i++ {
		pathstring += "."
		pathstring += err.path[pathLength-1-i]
	}
	pathstring = pathstring[1:] // Cut the first dot. Clearer that way
	return fmt.Sprintf("multiple struct fields are trying to be put into the following json field: %v")
}

func (err FieldCollisionError) AddField(field string) error {
	return FieldCollisionError{path: append(err.path, field)}
}

func NewFieldCollisionError(field string) FieldCollisionError {
	return FieldCollisionError{path: []string{field}}
}

func isEmpty(v reflect.Value) bool {
	kind := v.Kind()
	if v.CanInt() {
		return v.Int() == 0
	}
	switch kind {
	case reflect.Slice:
		return v.Len() == 0
	case reflect.Pointer:
		return v.IsNil()
	}
	return v.String() == ""
}
