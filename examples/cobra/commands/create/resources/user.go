package resources

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

// UserCommand with
// @Cobra(
//
//	usage = "cli create user",
//	example = "cli create user -g male -n John -l Doe -s 5",
//	short = "Create user command of the application (short)",
//	long = "Create user command of the application (long)",
//
// )
type UserCommand struct {
	Gender    Gender `flag:"gender,required" short:"g" description:"User gender, one of: male, female"`
	FirstName string `flag:"name,required" short:"n" description:"User name"`
	LastName  string `flag:"lastname" short:"l" description:"User last name"`
	State     State  `flag:"state,required" short:"s" description:"State is a value from 1 to 5"`
}

type State int
type Gender string

func (g *Gender) MarshalFlag(value string) error {
	m := map[string]struct{}{
		"male":   {},
		"female": {},
	}
	_, ok := m[value]
	if !ok {
		return fmt.Errorf("expected gender one of [male, female], but got %s", value)
	}
	*g = Gender(value)
	return nil
}

func (s *State) MarshalFlag(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	if i < 1 || i > 5 {
		return fmt.Errorf("expected state from 1 to 5, but got %d", i)
	}
	*s = State(i)
	return nil
}

// Run - @CobraRun command
func (c UserCommand) Run(cmd *cobra.Command, agrs []string) error {
	fmt.Println("Create user run command", c)
	return nil
}
