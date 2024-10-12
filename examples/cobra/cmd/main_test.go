package main

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type Flags struct {
	Output string `flag:"output,required" short:"o" description:"output file name"`
	Num    int    `flag:"num" default:"42" description:"some number for command"`
	IsOK   bool   `flag:"is-ok,persist" short:"i" default:"true" description:"some persistent flag"`
}

func TestParseFlags(t *testing.T) {
	f := Flags{}
	cmd := &cobra.Command{}

	cmd.Flags().StringP("output", "o", "", "output file name")
	cmd.Flags().Int("num", 42, "some number for command")
	cmd.PersistentFlags().BoolP("is-ok", "i", true, "some persistent flag")

	require.NoError(t, cmd.Flags().Set("output", "hello.txt"))
	require.NoError(t, cmd.Flags().Set("num", "15"))
	require.NoError(t, cmd.PersistentFlags().Set("is-ok", "false"))

	assert.NoError(t, parseFlags(cmd, &f))

	assert.Equal(t, "hello.txt", f.Output)
	assert.Equal(t, 15, f.Num)
	assert.Equal(t, false, f.IsOK)
}
