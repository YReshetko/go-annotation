package tag_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/YReshetko/go-annotation/internal/tag"
)

type TestTypeAll struct {
	F1 string `annotation:"default=hello,name=field1,required"`
	F2 int    `annotation:"default=1,name=field2,required"`
	F3 int32  `annotation:"required"`
	F4 int64  `annotation:"name=field4"`
	F5 uint   `annotation:""`
	F6 uint32 `annotation:"default=10"`
	F7 bool   `annotation:"default=true"`
}

func TestParse_TypeAll_Success(t *testing.T) {
	toTest := TestTypeAll{}
	actual, err := tag.Parse(toTest, map[string]string{
		"field1": "world",
		"f3":     "10",
		"field4": "160",
		"f5":     "30",
		"f6":     "11",
		"f7":     "false",
	})
	require.NoError(t, err)
	a, ok := actual.(TestTypeAll)
	require.True(t, ok)
	assert.Equal(t, "world", a.F1)
	assert.Equal(t, 1, a.F2)
	assert.Equal(t, int32(10), a.F3)
	assert.Equal(t, int64(160), a.F4)
	assert.Equal(t, uint(30), a.F5)
	assert.Equal(t, uint32(11), a.F6)
	assert.Equal(t, false, a.F7)
}

func TestParse_TypeAll_NoRequiredField(t *testing.T) {
	toTest := TestTypeAll{}
	_, err := tag.Parse(toTest, map[string]string{})
	require.EqualError(t, err, "unable to get value for annotation TestTypeAll: parameter f3 is required, but not set")
}

type TestTypeWrongTag struct {
	F1 string `annotation:"default=hello=world"`
}

func TestParse_WrongTag_Fail(t *testing.T) {
	toTest := TestTypeWrongTag{}
	require.PanicsWithError(t, "tag violates key=value format", func() {
		_, _ = tag.Parse(toTest, map[string]string{})
	})
}

type TestTypeWrongRequired struct {
	F1 string `annotation:"default=hello,require"`
}

func TestParse_WrongRequired_Fail(t *testing.T) {
	toTest := TestTypeWrongRequired{}
	require.PanicsWithError(t, "tag has only one option without values: 'required'", func() {
		_, _ = tag.Parse(toTest, map[string]string{})
	})
}
