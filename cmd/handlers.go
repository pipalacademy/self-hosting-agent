package main

import (
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
		return r.SendErrorEnvelope(http.StatusBadRequest, "error fetching hostname", nil, "HostError")
	}

	hostUptime, err := calcHostUptime()
	if err != nil {
		app.log.WithError(err).Error("error fetching host uptime")
		return r.SendErrorEnvelope(http.StatusBadRequest, "error fetching host uptime", nil, "HostError")
	}

	privateIP, err := getPrivateIP()
	if err != nil {
		app.log.WithError(err).Error("error fetching private IP")
		return r.SendErrorEnvelope(http.StatusBadRequest, "error fetching private IP", nil, "HostError")
	}

	publicIP, err := getPublicIP()
	if err != nil {
		app.log.WithError(err).Error("error fetching public IP")
		return r.SendErrorEnvelope(http.StatusBadRequest, "error fetching public IP", nil, "HostError")
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
