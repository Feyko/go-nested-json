package nested

import (
	"encoding/json"
	"fmt"
	"nested-json/internal/tags"
	"reflect"
)

func Marshal(v any) ([]byte, error) {
	marshaler := NewMarshaler(v)
	return marshaler.MarshalJSON()
}

func NewMarshaler(v any) Marshaler {
	return Marshaler{input: v}
}

type Marshaler struct {
	input    any
	reshaped map[string]any
	tags     map[string]tags.TagInfo
}

func (m Marshaler) MarshalJSON() ([]byte, error) {
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

func isStruct(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Struct
}
