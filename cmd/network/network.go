package network

import (
	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/cmd/network/latency"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "network",
		Short: "Targets network of a vm",
		Long:  "Performs specific actions like increasing the latency for a vm",
	}
	cmd.AddCommand(latency.NewCommand())
	return cmd
}
