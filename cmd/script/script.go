package script

import (
	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/cmd/script/run"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "script",
		Short: "Handles scripts that are launched on VMs",
		Long:  "Handles scripts that are launched on VMs running on Azure cloud",
	}
	cmd.AddCommand(run.NewCommand())
	return cmd
}
