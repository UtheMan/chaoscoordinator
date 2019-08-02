package main

import (
	"os"

	"github.com/UtheMan/chaosCoordinator/cmd/loadbalancer"
	"github.com/UtheMan/chaosCoordinator/cmd/vm"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Args:         cobra.NoArgs,
		Use:          "chaos",
		Short:        "chaos coordinator is a tool for cluster infrastructure chaos testing",
		Long:         `chaos coordinator is a tool for cluster infrastructure chaos testing. It targets underlying infrastructure with kubernetes cron jobs`,
		SilenceUsage: true,
	}
	rootCmd.AddCommand(vm.NewCommand())
	rootCmd.AddCommand(loadbalancer.NewCommand())
	return rootCmd
}

func Run() error {
	return NewCommand().Execute()
}
func main() {
	if err := Run(); err != nil {
		os.Exit(1)
	}
}
