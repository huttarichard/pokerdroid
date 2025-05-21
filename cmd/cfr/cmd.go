package cmdcfr

import (
	"github.com/spf13/cobra"
)

func init() {
	CMD.AddCommand(trainCMD)
	CMD.AddCommand(exploitCMD)
	CMD.AddCommand(analyzeCMD)
	CMD.AddCommand(testCMD)
}

var CMD = &cobra.Command{
	Use:   "cfr",
	Short: "will run training & analysis for tree",
}
