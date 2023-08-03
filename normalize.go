package form

import (
	"net/url"
	"reflect"
	"strconv"
)

type ValueNormalizer interface {
	NormalizeFormValue() (string, error)
}

type Normalizer interface {
	NormalizeForm() (url.Values, error)
}

func Normalize(v any) (url.Values, error) {
	if normalizer, ok := v.(Normalizer); ok {
		return normalizer.NormalizeForm()
	}

	value := reflect.ValueOf(v)
	for value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	ty := value.Type()
	if ty.Kind() != reflect.Struct {
		return nil, ErrUnknownType
	}

	values := url.Values{}

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

		switch rv.Kind() {
		case reflect.Struct:
			nestedValues, err := Normalize(rv.Interface())
			if err != nil {
				return nil, err
			}

			for name, value := range nestedValues {
				for _, v := range value {
					values.Add(name, v)
				}
			}

		case reflect.Slice, reflect.Array:
			for i := 0; i < rv.Len(); i++ {
				value, err := NormalizeFormValue(rv.Index(i).Interface())
				if err != nil {
					return nil, err
				}

				values.Add(name, value)
			}

		default:
			value, err := NormalizeFormValue(rv.Interface())
			if err != nil {
				return nil, err
			}

			values.Add(name, value)
		}
	}

	return values, nil
}

func NormalizeFormValue(v any) (string, error) {
	if normalizer, ok := v.(ValueNormalizer); ok {
		return normalizer.NormalizeFormValue()
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() {
			return "", nil
		}

		return NormalizeFormValue(reflect.ValueOf(v).Elem().Interface())

	case reflect.Bool:
		return strconv.FormatBool(rv.Bool()), nil

	case reflect.Int:
		return strconv.Itoa(int(rv.Int())), nil

	case reflect.Int8:
		return strconv.Itoa(int(rv.Int())), nil

	case reflect.Int16:
		return strconv.Itoa(int(rv.Int())), nil

	case reflect.Int32:
		return strconv.FormatInt(rv.Int(), 10), nil

	case reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10), nil

	case reflect.Uint:
		return strconv.FormatUint(rv.Uint(), 10), nil

	case reflect.Uint8:
		return strconv.Itoa(int(rv.Uint())), nil

	case reflect.Uint16:
		return strconv.Itoa(int(rv.Uint())), nil

	case reflect.Uint32:
		return strconv.FormatUint(rv.Uint(), 10), nil

	case reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10), nil

	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'G', -1, 32), nil

	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'G', -1, 64), nil

	case reflect.String:
		return rv.String(), nil
	}

	return "", ErrUnknownType
}
