package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/cdrivex4/agy-plus-plus/internal/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	code := cli.Execute(ctx, os.Args)
	os.Exit(code)
}
