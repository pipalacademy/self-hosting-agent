package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	// Version of the build. This is injected at build-time.
	buildString = "unknown"
)

func main() {
	// Create a new context which is cancelled when `SIGINT`/`SIGTERM` is received.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// Initialise and load the config.
	ko, err := initConfig("config.sample.toml", "MONSCHOOL_AGENT_")
	if err != nil {
		fmt.Println("error initialising config", err)
		os.Exit(1)
	}

	var (
		log    = initLogger(ko)
		server = initHTTPServer(ko)
		opts   = initOpts(ko)
	)

	// Initialise a new instance of app.
	app := App{
		log:    log,
		opts:   opts,
		server: server,
	}

	// Start an instance of app.
	app.log.WithField("version", buildString).Info("booting monschool agent")
	app.Start(ctx)
	// Listen on the close channel indefinitely until a
	// `SIGINT` or `SIGTERM` is received.
	<-ctx.Done()
	// Cancel the context to gracefully shutdown and perform
	// any cleanup tasks.
	cancel()
	// Shutdown routines.
	app.Stop()
}
