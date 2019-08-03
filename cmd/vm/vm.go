// Package vm implements `vm` command
package vm

import (
	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/cmd/vm/kill"
)

// NewCommand returns a new cobra.Command for vm targeting.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "vm",
		Short: "Targets virtual machines",
		Long:  "Targets virtual machines for the chaos testing",
	}
	cmd.AddCommand(kill.NewCommand())
	return cmd
}
