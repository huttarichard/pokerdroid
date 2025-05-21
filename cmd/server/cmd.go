package cmdserver

import (
	"github.com/spf13/cobra"
)

func init() {
	CMD.AddCommand(serverCMD)
	CMD.AddCommand(testCMD)
}

var CMD = &cobra.Command{
	Use:   "server",
	Short: "will run server that will serve solutions over HTTP",
}
