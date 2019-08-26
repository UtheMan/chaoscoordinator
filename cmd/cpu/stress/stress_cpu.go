// Package cpu implements `cpu` command
package stress

import (
	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/hardware/cmdutil"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/hardware/cpu"
	"os"
)

func NewCommand() *cobra.Command {
	cmdFlags := &cmdutil.Flags{}

	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "stress",
		Short: "Performs stress test on cpu",
		Long:  "Performs stress test on cpu for specified amount of time",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cpu.BeginStress(os.Getenv("SUBSCRIPTION_ID"), *cmdFlags)
		},
	}
	cmd.PersistentFlags().StringVarP(&cmdFlags.ResourceGroup, "resource group", "r", "", "Resource group name")
	cmd.PersistentFlags().StringVarP(&cmdFlags.ScaleSetName, "scale set", "s", "", "Scale set name")
	cmd.PersistentFlags().IntVarP(&cmdFlags.Duration, "duration", "d", 0, "Stress test duration")
	cmd.PersistentFlags().IntVarP(&cmdFlags.TimeOut, "timeout", "t", 0, "Stress test operation timeout")

	return cmd
}
