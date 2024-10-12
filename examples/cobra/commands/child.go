package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

// ChildCommand with
// @Cobra(
//
//	build = "user,ci",
//	usage = "cli get",
//	example = "cli get [-F file | -D dir] ... [-f format] profile",
//	short = "Child command of the application (short)",
//	long = "Child command of the application (long)",
//
// )
type ChildCommand struct {
	Output string `flag:"output,required" short:"o" description:"output file name"`
	Num    int    `flag:"num" short:"n" default:"42" description:"some number for command"`
	IsOK   bool   `flag:"is-ok,inherited"`
}

// Run - @CobraRun command
func (c ChildCommand) Run(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Child run command", c)
	return nil
}

// PostRun - @CobraPostRun command
func (c ChildCommand) PostRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Child post run command", c)
	return nil
}

// PersistPostRun - @CobraPersistPostRun command
func (c ChildCommand) PersistPostRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Child persist post run command", c)
	return nil
}
