package get

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

// GetCommand with
// @Cobra(
//
//	usage = "cli get",
//	example = "cli get [-F format] resource",
//	short = "Get command of the application (short)",
//	long = "Get command of the application (long)",
//
// )
type GetCommand struct {
	Flag1 string        `flag:"first" short:"f" description:"Flag 1 description"`
	Flag2 int           `flag:"second" short:"s" default:"42" description:"Flag 2 description"`
	Flag3 uint8         `flag:"third" short:"t" default:"8" description:"Flag 3 description"`
	Flag4 float64       `flag:"fourth" default:"3.14" description:"PI"`
	Flag5 bool          `flag:"fifth"`
	Flag6 uint16        `flag:"sixth"`
	Dur   time.Duration `flag:"dur,persist" short:"d" default:"12s" description:"Duration flag description"`
}

// Run - @CobraRun command
func (c GetCommand) Run(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Get run command", c)

	return nil
}

// PreRun - @CobraPreRun command
func (c GetCommand) PreRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Get post run command", c)
	return nil
}

// PersistPreRun - @CobraPersistPreRun command
/*func (c GetCommand) PersistPreRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Get persist post run command", c)
	return nil
}*/

// PostRun - @CobraPostRun command
func (c GetCommand) PostRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Get post run command", c)
	return nil
}

// PersistPostRun - @CobraPersistPostRun command
func (c GetCommand) PersistPostRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Get persist post run command", c)
	return nil
}
