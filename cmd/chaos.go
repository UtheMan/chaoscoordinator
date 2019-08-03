package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/utheman/chaoscoordinator/cmd/loadbalancer"
	"github.com/utheman/chaoscoordinator/cmd/vm"
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
