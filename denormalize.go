package form

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
)

var (
	ErrNonNilPointerRequired = errors.New("a non-nil pointer is required")
	ErrSliceNotSupported     = errors.New("slices are not supported yet")
	ErrUnknownField          = errors.New("unknown field")
)

type ValueDenormalizer interface {
	DenormalizeFormValue(data string) error
}

type Denormalizer interface {
	DenormalizeForm(data url.Values) error
}

func Denormalize(data url.Values, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrNonNilPointerRequired
	}

	if denormalizer, ok := v.(Denormalizer); ok {
		return denormalizer.DenormalizeForm(data)
	}

	rv = rv.Elem()
	ty := rv.Type()

	if ty.Kind() != reflect.Struct {
		return ErrUnknownType
	}

	for name, values := range data {
		fieldValue, tag := findFieldInStruct(ty, rv, name)
		if fieldValue == nil {
			return ErrUnknownField
		}

		if len(values) == 1 {
			value := values[0]

			if len(value) > 0 || !tag.OmitEmpty {
				if fieldValue.Kind() == reflect.Pointer && fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}

				if err := DenormalizeFormValue(value, fieldValue.Addr().Interface()); err != nil {
					return err
				}
			}
		} else {
			return ErrSliceNotSupported
		}
	}

	return nil
}

func DenormalizeFormValue(value string, v any) error {
	unmarshaler, ok := v.(ValueDenormalizer)
	if ok {
		return unmarshaler.DenormalizeFormValue(value)
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrNonNilPointerRequired
	}

	switch reflect.Indirect(rv).Kind() {
	case reflect.Pointer:
		return DenormalizeFormValue(value, rv.Elem().Interface())

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		rv.Elem().SetBool(b)

	case reflect.Int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		rv.Elem().SetInt(int64(i))

	case reflect.Int8:
		i, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return err
		}

		rv.Elem().SetInt(i)

	case reflect.Int16:
		i, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return err
		}

		rv.Elem().SetInt(i)

	case reflect.Int32:
		i, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}

		rv.Elem().SetInt(i)

	case reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		rv.Elem().SetInt(i)

	case reflect.Uint:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}

		rv.Elem().SetUint(i)

	case reflect.Uint8:
		i, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return err
		}

		rv.Elem().SetUint(i)

	case reflect.Uint16:
		i, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return err
		}

		rv.Elem().SetUint(i)

	case reflect.Uint32:
		i, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}

		rv.Elem().SetUint(i)

	case reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}

		rv.Elem().SetUint(i)

	case reflect.Float32:
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}

		rv.Elem().SetFloat(f)

	case reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		rv.Elem().SetFloat(f)

	case reflect.String:
		rv.Elem().SetString(value)

	default:
		return ErrUnknownType
	}

	return nil
}
