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
	rv := reflect.Indirect(reflect.ValueOf(data))
	if rv.Kind() != reflect.Struct {
		return "", nil, errNonStruct
	}
	t := rv.Type()
	ret := make([]structInfo, 0, t.NumField())
	var name = strings.ToLower(t.Name())
	for i := 0; i < t.NumField(); i++ {
		var tmp structInfo
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
		ret = append(ret, tmp)
	}
	return name, ret, nil
}

// func isStruct(v any) bool {
// 	rv := reflect.ValueOf(v)
// 	if rv.Kind() == reflect.Ptr {
// 		rv = rv.Elem()
// 	}
// 	return rv.Kind() == reflect.Struct
// }

// func isStruct(v any) bool {
// 	rv := reflect.Indirect(reflect.ValueOf(v))
// 	return rv.Kind() == reflect.Struct
// }

func (s *structInfo) ContainsTag(tag string) bool {
	for i := range s.Tags {
		if s.Tags[i] == tag {
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
	return s.ContainsTag("index")
}
