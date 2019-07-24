// Package loadbalancer implements `loadbalancer` command
package loadbalancer

import (
	"github.com/UtheMan/chaosCoordinator/cmd/loadbalancer/kill"
	"github.com/spf13/cobra"
)

var Scope string
var Mode string

// NewCommand returns a new cobra.Command for loadbalancer targeting.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "loadbalancer",
		Short: "Targets load balancers",
		Long:  "Targets load balancers for the chaos testing",
	}
	cmd.PersistentFlags().StringVarP(&Scope, "scope", "s", "", "Scope for the command")
	// cmd.MarkFlagRequired("scope")
	cmd.PersistentFlags().StringVarP(&Mode, "mode", "m", "", "Chaos mode")
	// cmd.MarkFlagRequired("mode")
	cmd.AddCommand(kill.NewCommand())
	return cmd
}
