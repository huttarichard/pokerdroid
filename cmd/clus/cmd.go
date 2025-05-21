package cmdclus

import "github.com/spf13/cobra"

func init() {
	CMD.AddCommand(equitiesCMD)
	CMD.AddCommand(riverCMD)
	CMD.AddCommand(turnCMD)
	CMD.AddCommand(flopCMD)
	CMD.AddCommand(packCMD)
}

var CMD = &cobra.Command{
	Use:   "clustering",
	Short: "build poker domain clustering",
}
