package form

import (
	"net/url"
	"reflect"
)

type ValueUnmarshaler interface {
	UnmarshalFormValue(bytes []byte) error
}

type Unmarshaler interface {
	UnmarshalForm(bytes []byte) error
}

func findFieldInStruct(ty reflect.Type, rv reflect.Value, name string) (*reflect.Value, Tag) {
	for i := 0; i < rv.NumField(); i++ {
		field := ty.Field(i)

		if field.Type.Kind() == reflect.Struct {
			if v, t := findFieldInStruct(field.Type, rv.Field(i), name); v != nil {
				return v, t
			}
		}

		tag := ParseTag(field.Tag.Get("form"))
		if tag.Key != "" {
			if tag.Key == name {
				f := rv.Field(i)

				return &f, tag
			}

			continue
		}

		if field.Name == name {
			f := rv.Field(i)

			return &f, tag
		}
	}

	return nil, Tag{}
}

func UnmarshalForm(bytes []byte, v any) error {
	values, err := url.ParseQuery(string(bytes))
	if err != nil {
		return err
	}

	return Denormalize(values, v)
}

func UnmarshalFormValue(bytes []byte, v any) error {
	s, err := url.QueryUnescape(string(bytes))
	if err != nil {
		return err
	}

	return DenormalizeFormValue(s, v)
}
