package main

import (
	"os"

	cmdbench "github.com/pokerdroid/poker/cmd/bench"
	cmdcfr "github.com/pokerdroid/poker/cmd/cfr"
	cmdclus "github.com/pokerdroid/poker/cmd/clus"
	cmddeep "github.com/pokerdroid/poker/cmd/deep"
	cmdserver "github.com/pokerdroid/poker/cmd/server"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pokerdroid",
	Short: "solver based on vrmccfr with search capabilities",
	Long:  ``,
}

func main() {
	rootCmd.AddCommand(cmdbench.CMD)
	rootCmd.AddCommand(cmdcfr.CMD)
	rootCmd.AddCommand(cmdserver.CMD)
	rootCmd.AddCommand(cmddeep.CMD)
	rootCmd.AddCommand(cmdclus.CMD)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
