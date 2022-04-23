package nested

import (
	"reflect"
	"strings"
)

type unmarshaller struct {
}

func Marshal(v any) ([]byte, error) {
	tags, _ := getTags(v)
	return nil, nil
}

func getTags(v any) (map[string][]string, error) {
	r := make(map[string][]string)
	vType := reflect.TypeOf(v)
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		fieldTag := field.Tag.Get("nested")
		r[field.Name] = strings.Split(fieldTag, ",")
	}
	return r, nil
}
