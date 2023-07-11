package form

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
)

var (
	ErrUnknownType = errors.New("could not determine how to marshal this type")
)

type ValueMarshaler interface {
	MarshalFormValue() ([]byte, error)
}

type Marshaler interface {
	MarshalForm() ([]byte, error)
}

func MarshalForm(v any) ([]byte, error) {
	marshaler, ok := v.(Marshaler)
	if ok {
		return marshaler.MarshalForm()
	}

	value := reflect.ValueOf(v)
	for value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	ty := value.Type()
	if ty.Kind() != reflect.Struct {
		return nil, ErrUnknownType
	}

	bytes := make([]byte, 0)

	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)

		name := field.Name
		tag := ParseTag(field.Tag.Get("form"))
		if tag.Key != "" {
			name = tag.Key
		}

		rv := value.Field(i)
		if rv.IsZero() && tag.OmitEmpty {
			continue
		}

		if len(bytes) > 0 {
			bytes = append(bytes, '&')
		}

		bytes = append(bytes, []byte(name)...)
		bytes = append(bytes, byte('='))

		valueBytes, err := MarshalFormValue(rv.Interface())
		if err != nil {
			return nil, err
		}

		bytes = append(bytes, valueBytes...)
	}

	return bytes, nil
}

func MarshalFormValue(v any) ([]byte, error) {
	marshaler, ok := v.(ValueMarshaler)
	if ok {
		return marshaler.MarshalFormValue()
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() {
			return []byte{}, nil
		}

		return MarshalFormValue(reflect.ValueOf(v).Elem().Interface())

	case reflect.Bool:
		return []byte(strconv.FormatBool(v.(bool))), nil

	case reflect.Int:
		return []byte(strconv.Itoa(v.(int))), nil

	case reflect.Int8:
		return []byte(strconv.Itoa(int(v.(int8)))), nil

	case reflect.Int16:
		return []byte(strconv.Itoa(int(v.(int16)))), nil

	case reflect.Int32:
		return []byte(strconv.FormatInt(int64(v.(int32)), 10)), nil

	case reflect.Int64:
		return []byte(strconv.FormatInt(v.(int64), 10)), nil

	case reflect.Uint:
		return []byte(strconv.FormatUint(uint64(v.(uint)), 10)), nil

	case reflect.Uint8:
		return []byte(strconv.Itoa(int(v.(uint8)))), nil

	case reflect.Uint16:
		return []byte(strconv.Itoa(int(v.(uint16)))), nil

	case reflect.Uint32:
		return []byte(strconv.FormatUint(uint64(v.(uint32)), 10)), nil

	case reflect.Uint64:
		return []byte(strconv.FormatUint(v.(uint64), 10)), nil

	case reflect.Float32:
		return []byte(strconv.FormatFloat(float64(v.(float32)), 'G', -1, 32)), nil

	case reflect.Float64:
		return []byte(strconv.FormatFloat(v.(float64), 'G', -1, 64)), nil

	case reflect.String:
		return []byte(url.QueryEscape(v.(string))), nil
	}

	return nil, ErrUnknownType
}
