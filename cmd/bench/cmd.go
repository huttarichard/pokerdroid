package cmdbench

import (
	"github.com/spf13/cobra"
)

func init() {
	CMD.AddCommand(mcSlumbotCMD)
	CMD.AddCommand(mcCFR_CMD)
	CMD.AddCommand(slumbotCFR_CMD)
}

var CMD = &cobra.Command{
	Use:   "bench",
	Short: "bench existing poker agents",
}
