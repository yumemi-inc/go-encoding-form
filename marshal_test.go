package form

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalForm(t *testing.T) {
	type values struct {
		Foo string `form:"foo"`
		Bar *int64 `form:"bar,omitempty"`
	}

	v := values{
		Foo: "testing",
		Bar: nil,
	}

	bytes, err := MarshalForm(v)
	require.NoError(t, err)
	assert.Equal(t, []byte("foo=testing"), bytes)

	i := int64(123)
	v.Bar = &i
	bytes, err = MarshalForm(v)

	require.NoError(t, err)
	assert.Equal(t, []byte("foo=testing&bar=123"), bytes)
}

func TestMarshalFormValue(t *testing.T) {
	bytes, err := MarshalFormValue(true)
	require.NoError(t, err)
	assert.Equal(t, []byte("true"), bytes)

	bytes, err = MarshalFormValue(123)
	require.NoError(t, err)
	assert.Equal(t, []byte("123"), bytes)

	bytes, err = MarshalFormValue(int8(123))
	require.NoError(t, err)
	assert.Equal(t, []byte("123"), bytes)

	bytes, err = MarshalFormValue(int16(123))
	require.NoError(t, err)
	assert.Equal(t, []byte("123"), bytes)

	bytes, err = MarshalFormValue(int32(123))
	require.NoError(t, err)
	assert.Equal(t, []byte("123"), bytes)

	bytes, err = MarshalFormValue(int64(123))
	require.NoError(t, err)
	assert.Equal(t, []byte("123"), bytes)

	bytes, err = MarshalFormValue(uint(123))
	require.NoError(t, err)
	assert.Equal(t, []byte("123"), bytes)

	bytes, err = MarshalFormValue(uint8(123))
	require.NoError(t, err)
	assert.Equal(t, []byte("123"), bytes)

	bytes, err = MarshalFormValue(uint16(123))
	require.NoError(t, err)
	assert.Equal(t, []byte("123"), bytes)

	bytes, err = MarshalFormValue(uint32(123))
	require.NoError(t, err)
	assert.Equal(t, []byte("123"), bytes)

	bytes, err = MarshalFormValue(uint64(123))
	require.NoError(t, err)
	assert.Equal(t, []byte("123"), bytes)

	bytes, err = MarshalFormValue(float32(123.456))
	require.NoError(t, err)
	assert.Equal(t, []byte("123.456"), bytes)

	bytes, err = MarshalFormValue(123.456)
	require.NoError(t, err)
	assert.Equal(t, []byte("123.456"), bytes)

	bytes, err = MarshalFormValue("abc!def")
	require.NoError(t, err)
	assert.Equal(t, []byte("abc%21def"), bytes)
}
