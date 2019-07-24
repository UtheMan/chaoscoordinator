// Package kill implements `kill` command
package kill

import (
	"github.com/spf13/cobra"
)

var (
	Scope string
	Mode  string
)

// NewCommand returns a new cobra. This command kills VM
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "kill",
		Short: "Kills one of [vm]",
		Long:  "Kills one of kubernetes virtual machines (vm) on azure",
		// RunE: func(cmd *cobra.Command, args []string) error {
		//command impl
		// },
	}
	cmd.PersistentFlags().StringVarP(&Scope, "scope", "s", "", "Scope for the command")
	cmd.PersistentFlags().StringVarP(&Mode, "mode", "m", "", "Chaos mode")
	return cmd
}
