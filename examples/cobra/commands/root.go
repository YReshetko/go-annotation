package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RootCommand with
// @Cobra(
//
//	usage = "cli",
//	example = "cli [-F file | -D dir] ... [-f format] profile",
//	short = "Root command of the application (short)",
//	long = "Root command of the application (long)",
//
// )
type RootCommand struct {
	Output string `flag:"output,required" short:"o" description:"output file name"`
	Num    int    `flag:"num" default:"42" description:"some number for command"`
	IsOK   bool   `flag:"is-ok,persist" short:"i" default:"true" description:"some persistent flag"`
}

// Run - @CobraRun command
func (c RootCommand) Run(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Root run command", c)
	return nil
}

// PreRun - @CobraPreRun command
func (c RootCommand) PreRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Root pre run command", c)
	return nil
}

// PersistPreRun - @CobraPersistPreRun command
//func (c RootCommand) PersistPreRun(cmd *cobra.Command, agrs []string) error {
//	fmt.Println("Root persist pre run command", c)
//	return nil
//}
