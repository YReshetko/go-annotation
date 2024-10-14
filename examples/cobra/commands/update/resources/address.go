package resources

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

// AddressCommand with
// @Cobra(
//
//	usage = "cli update address",
//	example = "cli update address -a 'Some street, 14'",
//	short = "Update address command of the application (short)",
//	long = "Update address command of the application (long)",
//
// )
type AddressCommand struct {
	Address Address `flag:"addr,required" short:"a" description:"Address separated by comma (Example: Some street, 14)"`
}

type Address struct {
	Street   string
	Building int
}

func (a *Address) MarshalFlag(value string) error {
	add := strings.Split(value, ",")
	if len(add) != 2 {
		return fmt.Errorf("expected to get address in format 'street, building', but got: '%s'", value)
	}

	b, err := strconv.Atoi(add[1])
	if err != nil {
		return err
	}

	*a = Address{
		Street:   add[0],
		Building: b,
	}
	return nil
}

// Run - @CobraRun command
func (c AddressCommand) Run(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Update address run command", c)
	return nil
}
