package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

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

// getFileNames lists all filenames present inside
// the directory.
func getFileNames(dir string) ([]string, error) {
	files := make([]string, 0)
	collection, err := os.ReadDir(dir)
	if err != nil {
		return files, err
	}

	for _, f := range collection {
		files = append(files, f.Name())
	}
	return files, nil
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

func isPkgInstalled(pkg string) (bool, error) {
	cmd := exec.Command(fmt.Sprintf("apt -qq list %s", pkg))
	stdout, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return strings.Contains((string(stdout)), "installed"), nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
