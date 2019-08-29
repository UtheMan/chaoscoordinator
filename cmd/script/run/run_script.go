package run

import (
	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/script"
	"os"
)

func NewCommand() *cobra.Command {
	cmdFlags := &script.Flags{}
	cmd := &cobra.Command{
		Args:  cobra.ArbitraryArgs,
		Use:   "run",
		Short: "Execute script on VM",
		Long:  "Execute script on VM on Azure",
		RunE: func(cmd *cobra.Command, args []string) error {
			return script.Run(os.Getenv("SUBSCRIPTION_ID"), *cmdFlags, args)
		},
	}
	cmd.PersistentFlags().IntVarP(&cmdFlags.Timeout, "time out", "t", 0, "Time out - additional time given for scripts to execute before operation times out")
	cmd.PersistentFlags().IntVarP(&cmdFlags.Duration, "duration", "d", 0, "Duration - time for script execution")
	cmd.PersistentFlags().StringVarP(&cmdFlags.ResourceGroup, "resource group", "r", "", "Resource group name")
	cmd.PersistentFlags().StringVarP(&cmdFlags.Path, "path", "p", "", "Path to script")
	cmd.PersistentFlags().StringVarP(&cmdFlags.ScaleSet, "scale set", "s", "", "Scale set name")
	cmd.PersistentFlags().StringVarP(&cmdFlags.Filter, "filter", "f", "", "Filter query for Azure")
	cmd.PersistentFlags().StringVarP(&cmdFlags.Kind, "kind", "k", "", "Short description of executed script")
	cmd.MarkFlagRequired("scale set")
	cmd.MarkFlagRequired("path")
	cmd.MarkFlagRequired("resource group")
	cmd.MarkFlagRequired("time out")
	return cmd
}
