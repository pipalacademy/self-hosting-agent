package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/shirou/gopsutil/v3/host"
)

func calcHostUptime() (string, error) {
	uptime, err := host.Uptime()
	if err != nil {
		return "", err
	}
	days := uptime / (60 * 60 * 24)
	hours := (uptime - (days * 60 * 60 * 24)) / (60 * 60)
	minutes := ((uptime - (days * 60 * 60 * 24)) - (hours * 60 * 60)) / 60
	return fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes), nil
}

// Get preferred outbound ip of this machine.
// https://stackoverflow.com/a/37382208
func getPrivateIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://ipinfo.io/ip")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func isFileExists(fPath string) error {
	if _, err := os.Stat(fPath); err != nil {
		return err

	}
	return nil
}

func getUsers() ([]string, error) {
	users := make([]string, 0)
	info, err := host.Users()
	if err != nil {
		return nil, err
	}
	for _, u := range info {
		users = append(users, u.User)
	}
	return users, nil
}

func parseSSHKeys() ([]string, error) {
	keys := make([]string, 0)

	authKeyFile := os.Getenv("HOME") + "/.ssh/authorized_keys"
	file, err := os.Open(authKeyFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return keys, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		keys = append(keys, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}
