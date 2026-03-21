package models

import (
	"net/url"
	"reflect"
	"strconv"
)

type ExternalRoutingProvideRequest struct {
	Arg       string `json:"arg"`
	Verbose   *bool  `json:"verbose,omitempty"`
	Recursive *bool  `json:"recursive,omitempty"`
}

func (req *ExternalRoutingProvideRequest) ToQueryParam() (string, error) {
	values := url.Values{}
	v := reflect.ValueOf(req).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		tagName := tag
		if commaIdx := len(tag) - len(",omitempty"); commaIdx > 0 && tag[commaIdx:] == ",omitempty" {
			tagName = tag[:commaIdx]
		}

		if field.Kind() != reflect.Pointer || !field.IsNil() {
			switch field.Kind() {
			case reflect.Ptr:
				switch field.Elem().Kind() {
				case reflect.Bool:
					values.Set(tagName, strconv.FormatBool(field.Elem().Bool()))
				case reflect.String:
					values.Set(tagName, field.Elem().String())
				}
			case reflect.String:
				values.Set(tagName, field.String())
			}
		}
	}

	return values.Encode(), nil
}
