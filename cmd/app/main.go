package main

import (
	"avito/internal/app"
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

var (
	configPath string
	isDebug    bool
)

func init() {
	flag.StringVar(&configPath, "config", "./config/config.env", "path to config file")
	flag.BoolVar(&isDebug, "is-debug", true, "enable debug mode")
}

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	defer stop()

	a, err := app.New(ctx, configPath, isDebug)
	if err != nil {
		os.Exit(1)
	}

	go func() {
		if err := a.Run(ctx); err != nil {
			stop()
		}
	}()

	<-ctx.Done()
	if err = a.Stop(context.Background()); err != nil {
		os.Exit(1)
	}
}
