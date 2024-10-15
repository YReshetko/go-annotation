package delete

import (
	"fmt"
	"github.com/spf13/cobra"
)

// DeleteCommand with
// @Cobra(
//
//	usage = "cli delete",
//	example = "cli delete [-R strategy] resource",
//	short = "Delete command of the application (short)",
//	long = "Delete command of the application (long)",
//
// )
type DeleteCommand struct{}

// Run - @CobraRun command
func (c DeleteCommand) Run(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Delete run command", c)
	return nil
}

// PostRun - @CobraPostRun command
func (c DeleteCommand) PostRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Delete post run command", c)
	return nil
}

// PersistPostRun - @CobraPersistPostRun command
func (c DeleteCommand) PersistPostRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Delete persist post run command", c)
	return nil
}
