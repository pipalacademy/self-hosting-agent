package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/zerodha/fastglue"
)

type Health struct {
	Addr        string `json:"agent_addr"`
	AgentUptime string `json:"agent_uptime"`
	HostUptime  string `json:"host_uptime"`
	Hostname    string `json:"hostname"`
	PublicIP    string `json:"public_ip"`
}

// handleIndex serves the `index` page.
func handleIndex(r *fastglue.Request) error {
	return r.SendEnvelope("Welcome to monschool-agent.")
}

// handlePing serves a ping response.
func handlePing(r *fastglue.Request) error {
	return r.SendEnvelope("pong")
}

// handleInfo serves the host info page.
func handleInfo(r *fastglue.Request) error {
	var (
		app         = r.Context.(*App)
		agentUptime = time.Since(app.initTime)
	)

	hostname, err := os.Hostname()
	if err != nil {
		app.log.WithError(err).Error("error fetching hostname")
		return r.SendErrorEnvelope(http.StatusInternalServerError, "error fetching hostname", nil, "HostError")
	}

	hostUptime, err := calcHostUptime()
	if err != nil {
		app.log.WithError(err).Error("error fetching host uptime")
		return r.SendErrorEnvelope(http.StatusInternalServerError, "error fetching host uptime", nil, "HostError")
	}

	publicIP, err := getPublicIP()
	if err != nil {
		app.log.WithError(err).Error("error fetching public IP")
		return r.SendErrorEnvelope(http.StatusInternalServerError, "error fetching public IP", nil, "HostError")
	}

	return r.SendEnvelope(&Health{
		Addr:        app.opts.ServerAddr,
		AgentUptime: agentUptime.String(),
		Hostname:    hostname,
		HostUptime:  hostUptime,
		PublicIP:    publicIP,
	})
}

func handleFileListing(r *fastglue.Request) error {
	var (
		app    = r.Context.(*App)
		result = map[string][]string{}
		dirs   = app.opts.WhitelistedDirs
	)
	for _, path := range dirs {
		files, err := getFileNames(path)
		if err != nil {
			app.log.WithError(err).Errorf("error fetching files for directory", path)
			return r.SendErrorEnvelope(http.StatusInternalServerError, "error fetching files", nil, "HostError")
		}
		result[path] = files
	}

	return r.SendEnvelope(result)
}

func handleFileListingByPath(r *fastglue.Request) error {
	var (
		app    = r.Context.(*App)
		result = map[string][]string{}
		path   = r.RequestCtx.UserValue("path").(string)
		dirs   = app.opts.WhitelistedDirs
	)

	// Check if path is whitelisted.
	if !stringInSlice(path, dirs) {
		return r.SendErrorEnvelope(http.StatusBadRequest, "path is not whitelisted", nil, "InputError")
	}

	// List all files under that path.
	files, err := getFileNames(path)
	if err != nil {
		app.log.WithError(err).Errorf("error fetching files for directory", path)
		return r.SendErrorEnvelope(http.StatusInternalServerError, "error fetching files", nil, "HostError")
	}
	result[path] = files

	return r.SendEnvelope(result)
}

func handleUsers(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
	)
	users, err := getUsers()
	if err != nil {
		app.log.WithError(err).Error("error fetching public IP")
		return r.SendErrorEnvelope(http.StatusInternalServerError, "error fetching public IP", nil, "HostError")
	}

	pubKeys, err := parseSSHKeys()
	if err != nil {
		app.log.WithError(err).Error("error fetching public IP")
		return r.SendErrorEnvelope(http.StatusInternalServerError, "error fetching public IP", nil, "HostError")
	}

	return r.SendEnvelope(struct {
		Users   []string `json:"users"`
		PubKeys []string `json:"pub_keys"`
	}{users, pubKeys})
}

func handleVerifyPackages(r *fastglue.Request) error {
	var (
		app      = r.Context.(*App)
		result   = map[string]bool{}
		packages = app.opts.WhitelistedPkgs
	)

	for _, p := range packages {
		ok, err := isPkgInstalled(p)
		if err != nil {
			app.log.WithError(err).Error("error fetching package status")
			return r.SendErrorEnvelope(http.StatusInternalServerError, "error fetching package status", nil, "HostError")
		}
		result[p] = ok
	}

	return r.SendEnvelope(result)
}

func handleVerifyPackageByName(r *fastglue.Request) error {
	var (
		app      = r.Context.(*App)
		result   = map[string]bool{}
		pkg      = r.RequestCtx.UserValue("pkg").(string)
		packages = app.opts.WhitelistedPkgs
	)

	// Check if pkg is whitelisted.
	fmt.Println(packages, pkg)
	if !stringInSlice(pkg, packages) {
		return r.SendErrorEnvelope(http.StatusBadRequest, "package is not whitelisted", nil, "InputError")
	}

	// Check if package is installed.
	ok, err := isPkgInstalled(pkg)
	if err != nil {
		app.log.WithError(err).Error("error fetching package status")
		return r.SendErrorEnvelope(http.StatusInternalServerError, "error fetching package status", nil, "HostError")
	}

	result[pkg] = ok
	return r.SendEnvelope(result)
}
