package disk

import (
	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/cmd/disk/fill"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "disk",
		Short: "Targets disk resources",
		Long:  "Targets disk resources for specified amount of time",
	}
	cmd.AddCommand(fill.NewCommand())
	return cmd
}
