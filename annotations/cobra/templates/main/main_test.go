package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MyEnum string

const (
	MyEnumValue1 MyEnum = "value-1"
	MyEnumValue2 MyEnum = "value-2"
	MyEnumValue3 MyEnum = "value-3"
)

func (e *MyEnum) MarshalFlag(value string) error {
	values := map[MyEnum]MyEnum{
		MyEnumValue1: MyEnumValue1,
		MyEnumValue2: MyEnumValue2,
		MyEnumValue3: MyEnumValue3,
	}
	res, ok := values[MyEnum(value)]
	if !ok {
		return fmt.Errorf("expected one of %v, but got %s", values, value)
	}
	*e = res
	return nil
}

type MyStruct struct {
	val1 string
	val2 string
}

func (e *MyStruct) MarshalFlag(value string) error {
	vals := strings.Split(value, ".")
	if len(vals) != 2 {
		return fmt.Errorf("expeted 2 values separated by dot, but got %d: %s", len(vals), value)
	}
	*e = MyStruct{
		val1: vals[0],
		val2: vals[1],
	}
	return nil
}

type Flags struct {
	Output string        `flag:"output,required" short:"o" description:"output file name"`
	Num    int           `flag:"num" default:"42" description:"some number for command"`
	IsOK   bool          `flag:"is-ok,persist" short:"i" default:"true" description:"some persistent flag"`
	Dur    time.Duration `flag:"dur" default:"2s" description:"some duration flag"`
	Enum   MyEnum        `flag:"enum" default:"value-2" description:"some enum with custom flag marshaling"`
	Custom MyStruct      `flag:"struct" default:"hello.world" description:"some custom structure with marshalling"`
}

func TestParseFlags(t *testing.T) {
	f := Flags{}
	cmd := &cobra.Command{}

	cmd.Flags().StringP("output", "o", "", "output file name")
	cmd.Flags().Int("num", 42, "some number for command")
	cmd.PersistentFlags().BoolP("is-ok", "i", true, "some persistent flag")
	cmd.Flags().String("dur", "10s", "some duration flag")
	cmd.Flags().String("enum", "value-1", "some enum with custom flag marshaling")
	cmd.Flags().String("struct", "hello.world", "some custom structure with marshalling")

	require.NoError(t, cmd.Flags().Set("output", "hello.txt"))
	require.NoError(t, cmd.Flags().Set("num", "15"))
	require.NoError(t, cmd.PersistentFlags().Set("is-ok", "false"))
	require.NoError(t, cmd.Flags().Set("dur", "2s"))
	require.NoError(t, cmd.Flags().Set("enum", "value-2"))
	require.NoError(t, cmd.Flags().Set("struct", "greeting.john"))

	require.NoError(t, parseFlags(cmd, &f))

	assert.Equal(t, "hello.txt", f.Output)
	assert.Equal(t, 15, f.Num)
	assert.Equal(t, false, f.IsOK)
	assert.Equal(t, time.Second*2, f.Dur)
	assert.Equal(t, MyEnumValue2, f.Enum)
	assert.Equal(t, MyStruct{
		val1: "greeting",
		val2: "john",
	}, f.Custom)
}
