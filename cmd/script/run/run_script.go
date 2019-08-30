package run

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/script"
	"os"
)

func NewCommand() *cobra.Command {
	cmdFlags := &script.Flags{}
	cmd := &cobra.Command{
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args)%2 != 0 {
				return errors.New("every argument needs a name and value")
			}
			return nil
		},
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
