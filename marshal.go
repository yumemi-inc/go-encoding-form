package form

import (
	"errors"
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
	values, err := Normalize(v)
	if err != nil {
		return nil, err
	}

	return []byte(values.Encode()), nil
}

func MarshalFormValue(v any) ([]byte, error) {
	marshaler, ok := v.(ValueMarshaler)
	if ok {
		return marshaler.MarshalFormValue()
	}

	value, err := NormalizeFormValue(v)
	if err != nil {
		return nil, err
	}

	return []byte(value), nil
}
