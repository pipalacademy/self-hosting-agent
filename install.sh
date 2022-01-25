#!/usr/bin/env sh
set -eu

VERSION="0.1.0"

# Monschool Agent easy install script.
# See https://github.com/fossunited/monschool-agent/ for detailed installation steps.

check_dependencies() {
	if ! command -v curl > /dev/null; then
		echo "curl is not installed."
		exit 1
	fi
}

setup_dirs() {
    mkdir -p /etc/monschool-agent
}

download_binary() {
    mkdir -p /tmp/monschool-agent
    cd /tmp/monschool-agent ; curl -sL https://github.com/fossunited/monschool-agent/releases/download/v${VERSION}/monschool-agent_${VERSION}_linux_amd64.tar.gz | tar xz
}

binary_path() {
    mv /tmp/monschool-agent/monschool-agent.bin /usr/local/bin
}

setup_systemd() {
    mv /tmp/monschool-agent/monschool-agent.service /etc/systemd/system/monschool-agent.service
    systemctl daemon-reload
    systemctl enable --now monschool-agent
    systemctl status monschool-agent
}



check_dependencies
download_binary
binary_path
setup_systemd
