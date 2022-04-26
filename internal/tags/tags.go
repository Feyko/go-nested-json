package tags

import (
	"reflect"
	"strings"
)

type TagInfo struct {
	path      []string
	omitempty bool
}

func GetTags(v any) (map[string]TagInfo, error) {
	r := make(map[string]TagInfo)
	vType := reflect.TypeOf(v)
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		fieldTag := field.Tag.Get("nested")
		r[field.Name] = parseTag(fieldTag)
	}
	return r, nil
}

func parseTag(tag string) (r TagInfo) {
	parts := strings.Split(tag, ",")
	r.path, r.omitempty = chopOmitempty(parts)
	return
}

func chopOmitempty(tagParts []string) (chopped []string, omitempty bool) {
	if len(tagParts) == 0 {
		return tagParts, false
	}
	if tagParts[len(tagParts)-1] == "omitempty" {
		omitempty = true
		chopped = tagParts[:len(tagParts)-1]
	}
	return
}
