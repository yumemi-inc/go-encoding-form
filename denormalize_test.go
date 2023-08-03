package form

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDenormalize(t *testing.T) {
	type newType string

	type inner struct {
		Baz bool `form:"baz"`
	}

	type values struct {
		Inner inner

		Foo     string   `form:"foo"`
		Bar     *int64   `form:"bar"`
		NewType newType  `form:"new_type,omitempty"`
		Slice   []string `form:"slice"`
		Array   [2]int   `form:"array"`
	}

	data := url.Values{}
	data.Set("foo", "testing")
	data.Set("bar", "123")
	data.Set("baz", "true")
	data.Set("new_type", "foo")
	data.Add("slice", "abc")
	data.Add("slice", "def")
	data.Add("array", "123")
	data.Add("array", "456")

	v := new(values)
	err := Denormalize(data, v)
	require.NoError(t, err)
	assert.Equal(t, "testing", v.Foo)
	assert.Equal(t, int64(123), *v.Bar)
	assert.Equal(t, true, v.Inner.Baz)
	assert.Equal(t, newType("foo"), v.NewType)
}

func TestDenormalizeFormValue(t *testing.T) {
	b := new(bool)
	require.NoError(t, DenormalizeFormValue("true", b))
	assert.Equal(t, true, *b)

	i := new(int)
	require.NoError(t, DenormalizeFormValue("123", i))
	assert.Equal(t, 123, *i)

	i8 := new(int8)
	require.NoError(t, DenormalizeFormValue("123", i8))
	assert.Equal(t, int8(123), *i8)

	i16 := new(int16)
	require.NoError(t, DenormalizeFormValue("123", i16))
	assert.Equal(t, int16(123), *i16)

	i32 := new(int32)
	require.NoError(t, DenormalizeFormValue("123", i32))
	assert.Equal(t, int32(123), *i32)

	i64 := new(int64)
	require.NoError(t, DenormalizeFormValue("123", i64))
	assert.Equal(t, int64(123), *i64)

	u := new(uint)
	require.NoError(t, DenormalizeFormValue("123", u))
	assert.Equal(t, uint(123), *u)

	u8 := new(uint8)
	require.NoError(t, DenormalizeFormValue("123", u8))
	assert.Equal(t, uint8(123), *u8)

	u16 := new(uint16)
	require.NoError(t, DenormalizeFormValue("123", u16))
	assert.Equal(t, uint16(123), *u16)

	u32 := new(uint32)
	require.NoError(t, DenormalizeFormValue("123", u32))
	assert.Equal(t, uint32(123), *u32)

	u64 := new(uint64)
	require.NoError(t, DenormalizeFormValue("123", u64))
	assert.Equal(t, uint64(123), *u64)

	f32 := new(float32)
	require.NoError(t, DenormalizeFormValue("123.456", f32))
	assert.Equal(t, float32(123.456), *f32)

	f64 := new(float64)
	require.NoError(t, DenormalizeFormValue("123.456", f64))
	assert.Equal(t, 123.456, *f64)

	s := new(string)
	require.NoError(t, DenormalizeFormValue("testing", s))
	assert.Equal(t, "testing", *s)
}
