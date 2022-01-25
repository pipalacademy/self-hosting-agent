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

download_binary() {
    mkdir -p /tmp/monschool-agent
    cd /tmp/monschool-agent ; curl -sL https://github.com/fossunited/monschool-agent/releases/download/v${VERSION}/monschool-agent_${VERSION}_linux_amd64.tar.gz | tar xz
}

binary_path() {
    mv /tmp/monschool-agent/monschool-agent.bin /usr/bin
}

setup_config() {
    mkdir -p /etc/monschool-agent
    mv /tmp/monschool-agent/config.sample.toml /etc/monschool-agent/config.toml
}

setup_systemd() {
    mv /tmp/monschool-agent/deployment/monschool-agent.service /etc/systemd/system/monschool-agent.service
    systemctl daemon-reload
    systemctl enable --now monschool-agent
    systemctl status monschool-agent
}


check_dependencies
download_binary
binary_path
setup_config
setup_systemd
