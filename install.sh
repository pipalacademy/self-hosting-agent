#!/usr/bin/env sh
set -eu

VERSION="0.2.0"

REPO_URL="https://github.com/pipalacademy/self-hosting-agent"
DOWNLOAD_URL="$REPO_URL/releases/download/v${VERSION}/self-hosting-agent-${VERSION}_linux_amd64.tar.gz"

# self-hosting-agent install script.
# See https://github.com/pipalacademy/self-hosting-agent/ for detailed installation steps.

check_dependencies() {
	if ! command -v curl > /dev/null; then
		echo "curl is not installed."
		exit 1
	fi
}

download_binary() {
    mkdir -p /tmp/self-hosting-agent
    cd /tmp/self-hosting-agent ; curl -sL $DOWNLOAD_URL | tar xz
    mv /tmp/self-hosting-agent/self-hosting-agent.bin /usr/bin
}

setup_config() {
    mkdir -p /etc/self-hosting-agent
    mv /tmp/self-hosting-agent/config.sample.toml /etc/self-hosting-agent/config.toml
}

setup_systemd() {
    mv /tmp/self-hosting-agent/deployment/self-hosting-agent.service /etc/systemd/system/self-hosting-agent.service
    systemctl daemon-reload
    systemctl enable --now self-hosting-agent
    systemctl status self-hosting-agent
}


check_dependencies
download_binary
setup_config
setup_systemd
