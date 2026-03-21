package models

import (
	"net/url"
	"reflect"
	"strconv"
)

type ExternalAddQueryRequest struct {
	Quiet             *bool   `json:"quiet,omitempty"`
	Quieter           *bool   `json:"quieter,omitempty"`
	Silent            *bool   `json:"silent,omitempty"`
	Progress          *bool   `json:"progress,omitempty"`
	Trickle           *bool   `json:"trickle,omitempty"`
	OnlyHash          *bool   `json:"only-hash,omitempty"`
	WrapWithDirectory *bool   `json:"wrap-with-directory,omitempty"`
	Chunker           *string `json:"chunker,omitempty"`
	RawLeaves         *bool   `json:"raw-leaves,omitempty"`
	NoCopy            *bool   `json:"nocopy,omitempty"`
	FSCache           *bool   `json:"fscache,omitempty"`
	CIDVersion        *int    `json:"cid-version,omitempty"`
	Hash              *string `json:"hash,omitempty"`
	Inline            *bool   `json:"inline,omitempty"`
	InlineLimit       *int    `json:"inline-limit,omitempty"`
	Pin               *bool   `json:"pin,omitempty"`
	ToFiles           *string `json:"to-files,omitempty"`
	PreserveMode      *bool   `json:"preserve-mode,omitempty"`
	PreserveMtime     *bool   `json:"preserve-mtime,omitempty"`
	Mode              *uint   `json:"mode,omitempty"`
	Mtime             *int64  `json:"mtime,omitempty"`
	MtimeNsecs        *uint   `json:"mtime-nsecs,omitempty"`
}

func (req *ExternalAddQueryRequest) ToQueryParam() (string, error) {
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

		if !field.IsNil() {
			switch field.Kind() {
			case reflect.Ptr:
				switch field.Elem().Kind() {
				case reflect.Bool:
					values.Set(tagName, strconv.FormatBool(field.Elem().Bool()))
				case reflect.Int, reflect.Int64:
					values.Set(tagName, strconv.FormatInt(field.Elem().Int(), 10))
				case reflect.Uint:
					values.Set(tagName, strconv.FormatUint(field.Elem().Uint(), 10))
				case reflect.String:
					values.Set(tagName, field.Elem().String())
				}
			}
		}
	}

	return values.Encode(), nil
}
