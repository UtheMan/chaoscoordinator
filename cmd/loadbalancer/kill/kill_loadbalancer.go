// Package kill implements `kill` command
package kill

import (
	"github.com/spf13/cobra"
)

// NewCommand returns a new cobra. This command kills loadbalancer
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "kill",
		Short: "Kills one of [loadbalancer]",
		Long:  "Kills one of kubernetes loadbalancers (loadbalancer) on azure",
	}
	return cmd
}
