// Package cpu implements `cpu` command
package stress

import (
	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/cpu"
	"os"
)

func NewCommand() *cobra.Command {
	cmdFlags := &cpu.Flags{}

	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "stress",
		Short: "Performs stress test on cpu",
		Long:  "Performs stress test on cpu for specified amount of time",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cpu.Stress(os.Getenv("SUBSCRIPTION_ID"), *cmdFlags)
		},
	}
	cmd.PersistentFlags().StringVarP(&cmdFlags.Time, "time", "t", "", "Stress test duration")

	return cmd
}
