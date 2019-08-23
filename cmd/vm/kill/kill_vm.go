// Package kill implements `kill` command
package kill

import (
	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/vm"
	"os"
)

// NewCommand returns a new cobra. This command kills VM
func NewCommand() *cobra.Command {
	cmdFlags := &vm.Flags{}

	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "kill",
		Short: "Kills one of [vm]",
		Long:  "Kills one of kubernetes virtual machines (vm) on azure",
		RunE: func(cmd *cobra.Command, args []string) error {
			return vm.Kill(os.Getenv("SUBSCRIPTION_ID"), *cmdFlags)
		},
	}
	cmd.PersistentFlags().StringVarP(&cmdFlags.Scope, "scope", "s", "", "Scope for the command")
	cmd.PersistentFlags().StringVarP(&cmdFlags.Mode, "mode", "m", "", "Chaos mode")
	cmd.PersistentFlags().StringVarP(&cmdFlags.ResourceName, "name", "n", "", "resource name")
	cmd.PersistentFlags().StringVarP(&cmdFlags.ResourceGroup, "resource group", "r", "", "resource group name")

	return cmd
}
