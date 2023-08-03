package form

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeForm(t *testing.T) {
	type newType string

	type inner struct {
		Baz bool `form:"baz"`
	}

	type values struct {
		Inner inner

		Foo     string   `form:"foo"`
		Bar     *int64   `form:"bar,omitempty"`
		NewType newType  `form:"new_type"`
		Slice   []string `form:"slice"`
		Array   [2]int   `form:"array"`
	}

	v := values{
		Foo:     "testing",
		Bar:     nil,
		NewType: "foo",
		Slice:   []string{},
	}

	actual, err := Normalize(v)
	require.NoError(t, err)
	assert.Equal(t, "false", actual.Get("baz"))
	assert.Equal(t, "testing", actual.Get("foo"))
	assert.Equal(t, "foo", actual.Get("new_type"))
	assert.Equal(t, []string{"0", "0"}, actual["array"])
	assert.Empty(t, actual.Get("bar"))
	assert.Empty(t, actual["slice"])

	i := int64(123)
	v.Inner.Baz = true
	v.Bar = &i
	v.Slice = []string{"abc", "def"}
	v.Array = [2]int{123, 456}
	actual, err = Normalize(v)

	require.NoError(t, err)
	assert.Equal(t, "true", actual.Get("baz"))
	assert.Equal(t, "testing", actual.Get("foo"))
	assert.Equal(t, "123", actual.Get("bar"))
	assert.Equal(t, "foo", actual.Get("new_type"))
	assert.Equal(t, []string{"abc", "def"}, actual["slice"])
	assert.Equal(t, []string{"123", "456"}, actual["array"])
}

func TestNormalizeFormValue(t *testing.T) {
	bytes, err := NormalizeFormValue(true)
	require.NoError(t, err)
	assert.Equal(t, "true", bytes)

	bytes, err = NormalizeFormValue(123)
	require.NoError(t, err)
	assert.Equal(t, "123", bytes)

	bytes, err = NormalizeFormValue(int8(123))
	require.NoError(t, err)
	assert.Equal(t, "123", bytes)

	bytes, err = NormalizeFormValue(int16(123))
	require.NoError(t, err)
	assert.Equal(t, "123", bytes)

	bytes, err = NormalizeFormValue(int32(123))
	require.NoError(t, err)
	assert.Equal(t, "123", bytes)

	bytes, err = NormalizeFormValue(int64(123))
	require.NoError(t, err)
	assert.Equal(t, "123", bytes)

	bytes, err = NormalizeFormValue(uint(123))
	require.NoError(t, err)
	assert.Equal(t, "123", bytes)

	bytes, err = NormalizeFormValue(uint8(123))
	require.NoError(t, err)
	assert.Equal(t, "123", bytes)

	bytes, err = NormalizeFormValue(uint16(123))
	require.NoError(t, err)
	assert.Equal(t, "123", bytes)

	bytes, err = NormalizeFormValue(uint32(123))
	require.NoError(t, err)
	assert.Equal(t, "123", bytes)

	bytes, err = NormalizeFormValue(uint64(123))
	require.NoError(t, err)
	assert.Equal(t, "123", bytes)

	bytes, err = NormalizeFormValue(float32(123.456))
	require.NoError(t, err)
	assert.Equal(t, "123.456", bytes)

	bytes, err = NormalizeFormValue(123.456)
	require.NoError(t, err)
	assert.Equal(t, "123.456", bytes)

	bytes, err = NormalizeFormValue("abc!def")
	require.NoError(t, err)
	assert.Equal(t, "abc%21def", bytes)
}
