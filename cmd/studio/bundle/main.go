package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pokerdroid/poker/studio"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		os.Interrupt,
	)
	defer cancel()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		logger.Fatal(err)
	}
	defer listener.Close()

	dev := studio.New()

	err = dev.Run(ctx, listener)
	if errors.Is(err, context.Canceled) {
		logger.Println("Done.")
		return
	}

	if err != nil {
		logger.Fatal(err)
	}
}
