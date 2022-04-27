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
	return Marshaler{input: v, reshaped: make(map[string]any)}
}

type Marshaler struct {
	input    any
	reshaped map[string]any
	tags     map[string]tags.TagInfo
}

func (m *Marshaler) MarshalJSON() ([]byte, error) {
	if !isStruct(m.input) {
		return json.Marshal(m.input)
	}
	structTags, err := tags.GetTags(m.input)
	if err != nil {
		return nil, fmt.Errorf("could not parse the struct tags: %v", err)
	}
	m.tags = structTags
	return nil, nil
}

func addFieldToReshape(reshaped map[string]any, path []string, v any) error {
	fieldName := path[0]
	field := reshaped[fieldName]
	if len(path) == 1 {
		if field != nil {
			return errors.New(fieldName + " is already occupied")
		}
		reshaped[fieldName] = v
		return nil
	}
	mapV, isMap := field.(map[string]any)
	if field != nil && isMap {
		return addFieldToReshape(mapV, path[1:], v)
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
