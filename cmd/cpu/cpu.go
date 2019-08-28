// Package cpu implements `cpu` command
package cpu

import (
	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/cmd/cpu/stress"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "cpu",
		Short: "Targets CPU resources",
		Long:  "Targets CPU resources for certain amount of time",
	}
	cmd.AddCommand(stress.NewCommand())
	return cmd
}
