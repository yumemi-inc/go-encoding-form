package form

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
)

var (
	ErrNonNilPointerRequired = errors.New("a non-nil pointer is required")
	ErrUnknownField          = errors.New("unknown field")
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
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrNonNilPointerRequired
	}

	unmarshaler, ok := v.(Unmarshaler)
	if ok {
		return unmarshaler.UnmarshalForm(bytes)
	}

	rv = rv.Elem()
	ty := rv.Type()

	if ty.Kind() != reflect.Struct {
		return ErrUnknownType
	}

	name := ""
	nameBuf := make([]byte, 0, len(bytes))
	valueBuf := make([]byte, 0, len(bytes))

	for _, b := range append(bytes, 0) {
		if b == '=' {
			name = string(nameBuf)

			continue
		}

		if b == '&' || b == 0 {
			fieldValue, tag := findFieldInStruct(ty, rv, name)
			if fieldValue == nil {
				return ErrUnknownField
			}

			if len(valueBuf) > 0 || !tag.OmitEmpty {
				if fieldValue.Kind() == reflect.Pointer && fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}

				if err := UnmarshalFormValue(valueBuf, fieldValue.Addr().Interface()); err != nil {
					return err
				}
			}

			name = ""
			nameBuf = make([]byte, 0, len(bytes))
			valueBuf = make([]byte, 0, len(bytes))

			if b == 0 {
				break
			} else {
				continue
			}
		}

		if name == "" {
			nameBuf = append(nameBuf, b)
		} else {
			valueBuf = append(valueBuf, b)
		}
	}

	return nil
}

func UnmarshalFormValue(bytes []byte, v any) error {
	unmarshaler, ok := v.(ValueUnmarshaler)
	if ok {
		return unmarshaler.UnmarshalFormValue(bytes)
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrNonNilPointerRequired
	}

	switch reflect.Indirect(rv).Kind() {
	case reflect.Pointer:
		return UnmarshalFormValue(bytes, rv.Elem().Interface())

	case reflect.Bool:
		b, err := strconv.ParseBool(string(bytes))
		if err != nil {
			return err
		}

		rv.Elem().SetBool(b)

	case reflect.Int:
		i, err := strconv.Atoi(string(bytes))
		if err != nil {
			return err
		}

		rv.Elem().SetInt(int64(i))

	case reflect.Int8:
		i, err := strconv.ParseInt(string(bytes), 10, 8)
		if err != nil {
			return err
		}

		rv.Elem().SetInt(i)

	case reflect.Int16:
		i, err := strconv.ParseInt(string(bytes), 10, 16)
		if err != nil {
			return err
		}

		rv.Elem().SetInt(i)

	case reflect.Int32:
		i, err := strconv.ParseInt(string(bytes), 10, 32)
		if err != nil {
			return err
		}

		rv.Elem().SetInt(i)

	case reflect.Int64:
		i, err := strconv.ParseInt(string(bytes), 10, 64)
		if err != nil {
			return err
		}

		rv.Elem().SetInt(i)

	case reflect.Uint:
		i, err := strconv.ParseUint(string(bytes), 10, 64)
		if err != nil {
			return err
		}

		rv.Elem().SetUint(i)

	case reflect.Uint8:
		i, err := strconv.ParseUint(string(bytes), 10, 8)
		if err != nil {
			return err
		}

		rv.Elem().SetUint(i)

	case reflect.Uint16:
		i, err := strconv.ParseUint(string(bytes), 10, 16)
		if err != nil {
			return err
		}

		rv.Elem().SetUint(i)

	case reflect.Uint32:
		i, err := strconv.ParseUint(string(bytes), 10, 32)
		if err != nil {
			return err
		}

		rv.Elem().SetUint(i)

	case reflect.Uint64:
		i, err := strconv.ParseUint(string(bytes), 10, 64)
		if err != nil {
			return err
		}

		rv.Elem().SetUint(i)

	case reflect.Float32:
		f, err := strconv.ParseFloat(string(bytes), 32)
		if err != nil {
			return err
		}

		rv.Elem().SetFloat(f)

	case reflect.Float64:
		f, err := strconv.ParseFloat(string(bytes), 64)
		if err != nil {
			return err
		}

		rv.Elem().SetFloat(f)

	case reflect.String:
		s, err := url.QueryUnescape(string(bytes))
		if err != nil {
			return err
		}

		rv.Elem().SetString(s)

	default:
		return ErrUnknownType
	}

	return nil
}
