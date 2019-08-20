package fill

import (
	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/disk"
	"os"
)

func NewCommand() *cobra.Command {
	cmdFlags := &disk.Flags{}

	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "fill",
		Short: "Fills disk resources with data",
		Long:  "Fills disk resources with data for specified amount of time",
		RunE: func(cmd *cobra.Command, args []string) error {
			return disk.Fill(os.Getenv("SUBSCRIPTION_ID"), *cmdFlags)
		},
	}
	cmd.PersistentFlags().StringVarP(&cmdFlags.Time, "time", "t", "", "Stress test duration")
	cmd.PersistentFlags().IntVarP(&cmdFlags.FillPercentage, "fill percentage", "f", 0, "Disk fill percentage (0-100)")
	return cmd
}
