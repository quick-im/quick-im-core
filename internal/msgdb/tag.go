package msgdb

import (
	"errors"
	"reflect"
	"strings"
)

var (
	errNonStruct = errors.New("non-struct type")
)

type structInfo struct {
	Field string
	Tags  []string
}

func parseStruct(data any) (string, []structInfo, error) {
	if !isStruct(data) {
		return "", nil, errNonStruct
	}
	t := reflect.TypeOf(data)
	ret := make([]structInfo, t.NumField())
	var name = strings.ToLower(t.Name())
	for i := 0; i < t.NumField(); i++ {
		tmp := structInfo{}
		field := t.Field(i)
		tag := field.Tag.Get("imdb")
		tmp.Field = strings.ToLower(field.Name)
		// 如果设置了rethinkdb字段名，则以rethinkdb为准
		rtags := strings.Split(field.Tag.Get("rethinkdb"), ",")
		if len(rtags) != 0 {
			if rtags[0] != "-" && rtags[0] != "" {
				tmp.Field = rtags[0]
			}
		}
		if tag != "" {
			tmp.Tags = strings.Split(tag, ",")
		}
		ret[i] = tmp
	}
	return name, ret, nil
}

func isStruct(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv.Kind() == reflect.Struct
}

func (s *structInfo) ContainsTag(tag string) bool {
	for _, x := range s.Tags {
		if x == tag {
			return true
		}
	}
	return false
}

func (s *structInfo) IsPrimaryKey() bool {
	if s.ContainsTag("pk") || s.ContainsTag("primary_key") {
		return true
	}
	return false
}

func (s *structInfo) IsSimpleIndex() bool {
	if s.ContainsTag("index") {
		return true
	}
	return false
}
