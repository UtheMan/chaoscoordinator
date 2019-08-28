package latency

import (
	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/hardware/network"
	"os"
)

func NewCommand() *cobra.Command {
	cmdFlags := &network.Flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "latency",
		Short: "Increases latency for a vm",
		Long:  "Increases latency for a vm by specified amount",
		RunE: func(cmd *cobra.Command, args []string) error {
			return network.BeginLatencyIncrease(os.Getenv("SUBSCRIPTION_ID"), *cmdFlags)
		},
	}
	cmd.PersistentFlags().IntVarP(&cmdFlags.Duration, "duration", "d", 0, "Latency test duration")
	cmd.PersistentFlags().StringVarP(&cmdFlags.ResourceGroup, "resource group", "r", "", "Resource group name")
	cmd.PersistentFlags().StringVarP(&cmdFlags.ScaleSet, "scale set", "s", "", "Scale set name")
	cmd.PersistentFlags().IntVarP(&cmdFlags.Latency, "latency", "l", 0, "Latency (in ms) we want to add to a vm")
	cmd.PersistentFlags().IntVarP(&cmdFlags.Timeout, "time out", "t", 0, "Time out - additional time given for scripts to execute")
	cmd.MarkFlagRequired("scale set")
	cmd.MarkFlagRequired("resource group")
	return cmd
}

//network latency -d 120 -l 200 -t 60 -r chaoscoordinatorresourcegroup -s controlplane
