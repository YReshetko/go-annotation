package update

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// UpdateCommand with
// @Cobra(
//
//	usage = "cli update",
//	example = "cli update [-F field -V value] resource",
//	short = "Update command of the application (short)",
//	long = "Update command of the application (long)",
//
// )
type UpdateCommand struct {
	Flag1 string        `flag:"first" short:"f" description:"Flag 1 description"`
	Flag2 int           `flag:"second" short:"s" default:"42" description:"Flag 2 description"`
	Flag3 uint8         `flag:"third" short:"t" default:"8" description:"Flag 3 description"`
	Flag4 float64       `flag:"fourth,persist" default:"3.14" description:"PI"`
	Flag5 bool          `flag:"fifth"`
	Flag6 uint16        `flag:"sixth"`
	Dur   time.Duration `flag:"dur,persist" short:"d" default:"12s" description:"Duration flag description"`
}

// Run - @CobraRun command
func (c UpdateCommand) Run(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Update run command", c)

	return nil
}

// PreRun - @CobraPreRun command
func (c UpdateCommand) PreRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Update post run command", c)
	return nil
}

// PersistPreRun - @CobraPersistPostRun command
/*func (c UpdateCommand) PersistPreRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Update persist post run command", c)
	return nil
}*/

// PostRun - @CobraPostRun command
func (c UpdateCommand) PostRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Update post run command", c)
	return nil
}

// PersistPostRun - @CobraPersistPostRun command
/*func (c UpdateCommand) PersistPostRun(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Update persist post run command", c)
	return nil
}*/
