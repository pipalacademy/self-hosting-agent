package main

import (
	"errors"
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
	PrivateIP   string `json:"private_ip"`
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

	privateIP, err := getPrivateIP()
	if err != nil {
		app.log.WithError(err).Error("error fetching private IP")
		return r.SendErrorEnvelope(http.StatusInternalServerError, "error fetching private IP", nil, "HostError")
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
		PrivateIP:   privateIP,
		PublicIP:    publicIP,
	})
}

func handleVerifyFile(r *fastglue.Request) error {
	var (
		app = r.Context.(*App)
	)
	err := isFileExists(app.opts.PublicFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return r.SendErrorEnvelope(http.StatusBadRequest, fmt.Sprintf("%s does not exists", app.opts.PublicFile), nil, "InputError")
		}
		app.log.WithError(err).Error("error fetching file path")
		return r.SendErrorEnvelope(http.StatusInternalServerError, "error fetching file path", nil, "HostError")
	}
	return r.SendEnvelope(nil)
}
