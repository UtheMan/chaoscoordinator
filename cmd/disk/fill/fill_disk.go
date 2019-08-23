package fill

import (
	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/hardware/cmdutil"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/hardware/disk"
	"os"
)

func NewCommand() *cobra.Command {
	cmdFlags := &cmdutil.Flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "fill",
		Short: "Fills disk resources with data",
		Long:  "Fills disk resources with data for specified amount of time",
		RunE: func(cmd *cobra.Command, args []string) error {
			return disk.BeginFill(os.Getenv("SUBSCRIPTION_ID"), *cmdFlags)
		},
	}
	cmd.PersistentFlags().IntVarP(&cmdFlags.Duration, "duration", "d", 0, "Stress test duration")
	cmd.PersistentFlags().StringVarP(&cmdFlags.ResourceGroup, "resource group", "r", "", "Resource group name")
	cmd.PersistentFlags().StringVarP(&cmdFlags.ScaleSetName, "scale set", "s", "", "Scale set name")
	cmd.PersistentFlags().StringVarP(&cmdFlags.Amount, "amount", "a", "", "Amount of data (in MB) we want to load onto vms")
	cmd.PersistentFlags().IntVarP(&cmdFlags.TimeOut, "time out", "t", 0, "Time out - additional time given for scripts to execute")
	cmd.MarkFlagRequired("scale set")
	cmd.MarkFlagRequired("resource group")
	return cmd
}
