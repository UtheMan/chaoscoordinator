// Package kill implements `kill` command
package kill

import (
	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/loadbalancer"
	"os"
)

// NewCommand returns a new cobra. This command kills loadbalancer
func NewCommand() *cobra.Command {
	cmdFlags := &loadbalancer.Flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "kill",
		Short: "Kills one of [loadbalancer]",
		Long:  "Kills one of kubernetes loadbalancers (loadbalancer) on azure",
		RunE: func(cmd *cobra.Command, args []string) error {
			return loadbalancer.Kill(os.Getenv("SUBSCRIPTION_ID"), *cmdFlags)
		},
	}
	cmd.PersistentFlags().StringVarP(&cmdFlags.Scope, "scope", "s", "", "Scope for the command")
	cmd.PersistentFlags().StringVarP(&cmdFlags.Mode, "mode", "m", "", "Chaos mode")
	cmd.PersistentFlags().StringVarP(&cmdFlags.LoadBalancerName, "loadbalancer name", "l", "", "Load balancer name")
	cmd.PersistentFlags().StringVarP(&cmdFlags.ResourceGroup, "resourcegroup", "r", "", "Resource group name")
	return cmd
}
