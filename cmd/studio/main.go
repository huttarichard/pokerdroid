package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	absp "github.com/pokerdroid/poker/abs/pack"
	"github.com/pokerdroid/poker/studio"
	studiotree "github.com/pokerdroid/poker/studio/tree"
	"github.com/pokerdroid/poker/tree"
	"github.com/spf13/cobra"
)

func init() {
	CMD.AddCommand(devCMD)
}

var CMD = &cobra.Command{
	Use:   "pokerdroid",
	Short: "pokerdroid studio",
}

func main() {
	err := CMD.Execute()
	if err != nil {
		os.Exit(1)
	}
}

type studioArgs struct {
	dir string
	abs string
}

var tf = studioArgs{}

func init() {
	flags := devCMD.Flags()

	str, _ := os.UserHomeDir()
	dir := filepath.Join(str, ".pokerdroid")

	p := os.Getenv("POKERDROID_DIR")
	if p != "" {
		dir = p
	}

	flags.StringVarP(&tf.dir, "dir", "d", dir, "directory containing the pokerdroid solutions")
	cobra.MarkFlagRequired(flags, "dir")

	flags.StringVarP(&tf.abs, "abs", "a", "pack_600_600_1000.bin", "default abstraction")
}

var devCMD = &cobra.Command{
	Use:   "dev",
	Short: "Run development server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGTERM, os.Interrupt)
		defer cancel()

		// flags := cmd.Flags()
		logger := log.Default()

		dev, err := studio.NewDevelopment()
		if err != nil {
			logger.Fatal("failed to create development server:", err)
		}

		os.MkdirAll(tf.dir, 0755)

		dev.WebV.Bind("rpc_log", func(msg string) error {
			logger.Print(msg)
			return nil
		})

		rxs, err := tree.NewFileRootsFromDir(tf.dir)
		if err != nil {
			logger.Fatal(err)
		}
		defer rxs.Close()

		logger.Printf("%s", rxs.String())

		pack := filepath.Join(tf.dir, "clustering", tf.abs)
		logger.Printf("loading clustering: %s", pack)

		abs, err := absp.NewFromFile(pack)
		if err != nil {
			logger.Fatalf("failed to create abstraction: %s", err)
		}

		err = studiotree.Bind(studiotree.BindParams{
			Abs:     abs,
			Roots:   rxs.Roots(),
			WebView: dev.WebV,
		})
		if err != nil {
			logger.Fatal("failed to bind tree:", err)
		}

		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			logger.Fatal("failed to create listener:", err)
		}
		defer listener.Close()

		logger.Printf("Listening on http://%s", listener.Addr().String())

		err = dev.Run(ctx, listener)
		if err != nil {
			logger.Fatal("failed to run development server:", err)
		}
	},
}
