package create

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CreateCommand with
// @Cobra(
//
//	usage = "cli create",
//	example = "cli create [-F file] resource",
//	short = "Create command of the application (short)",
//	long = "Create command of the application (long)",
//
// )
type CreateCommand struct{}

// Run - @CobraRun command
func (c CreateCommand) Run(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Create run command", c)
	return nil
}

// PostRun - @CobraPostRun command
func (c CreateCommand) PostRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Create post run command", c)
	return nil
}

// PersistPostRun - @CobraPersistPostRun command
func (c CreateCommand) PersistPostRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Create persist post run command", c)
	return nil
}
