package tags

import (
	"reflect"
	"strings"
)

type TagInfo struct {
	Path      []string
	OmitEmpty bool
}

func GetTags(v reflect.Type) (map[string]TagInfo, error) {
	r := make(map[string]TagInfo)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !isValidField(field) {
			continue
		}
		fieldTag := field.Tag.Get("nested")
		r[field.Name] = parseTag(fieldTag)
	}
	return r, nil
}

func parseTag(tag string) (r TagInfo) {
	parts := strings.Split(tag, ",")
	r.Path, r.OmitEmpty = chopOmitempty(parts)
	return
}

func chopOmitempty(tagParts []string) (chopped []string, omitempty bool) {
	chopped = tagParts
	if len(tagParts) == 0 {
		return tagParts, false
	}
	if tagParts[len(tagParts)-1] == "omitempty" {
		omitempty = true
		chopped = tagParts[:len(tagParts)-1]
	}
	return
}

func isValidField(field reflect.StructField) bool {
	return field.IsExported()
}
