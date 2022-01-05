package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"
)

type Opts struct {
	ServerAddr      string
	WhitelistedDirs []string
	WhitelistedPkgs []string
}

// App is the global container that holds
// objects of various routines that run on boot.
type App struct {
	log      *logrus.Logger
	server   *fasthttp.Server
	opts     Opts
	initTime time.Time
}

// Start initialises the sink workers and
// subscription stream in background and waits
// for context to be cancelled to exit.
func (app *App) Start(ctx context.Context) {
	// Set the initialising time. This is
	// used to calculate uptime of the agent.
	app.initTime = time.Now()

	// Initialise a new instance of `fastglue`.
	g := fastglue.NewGlue()

	// Inject an instance of `app` as context.
	g.SetContext(app)

	// Initialise HTTP handlers.
	g.GET("/", handleIndex)
	g.GET("/ping", handlePing)
	g.GET("/host", handleInfo)
	g.GET("/users", handleUsers)
	g.GET("/files/", handleFileListing)
	g.GET("/files/{path}", handleFileListingByPath)
	g.GET("/packages", handleVerifyPackages)
	g.GET("/packages/{pkg}", handleVerifyPackageByName)

	// Start HTTP server.
	app.log.WithField("address", app.opts.ServerAddr).Info("starting server")
	go func() {
		if err := g.ListenAndServe(app.opts.ServerAddr, "", app.server); err != nil {
			app.log.WithError(err).Fatal("error starting server")
		}
	}()
}

func (app *App) Stop() {
	// Perform a graceful shutdown of server.
	err := app.server.Shutdown()
	if err != nil {
		app.log.WithError(err).Fatal("error shutting down server")
	}
	app.log.Info("front server shutdown. bye")
}
