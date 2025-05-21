package cmddeep

import "github.com/spf13/cobra"

func init() {
	CMD.AddCommand(baselineCMD)
}

var CMD = &cobra.Command{
	Use:   "deep",
	Short: "deep learning",
}
