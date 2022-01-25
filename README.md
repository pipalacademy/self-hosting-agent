# monschool-agent
A deamon program that runs on the user's node for the Self Hosting 101 Course

## Install

### Easy Install

```bash
curl --proto '=https' --tlsv1.2 -sSf https://raw.githubusercontent.com/fossunited/monschool-agent/main/install.sh | bash
```

### Manual

```bash
$ mkdir ~/monschool-agent ; cd ~/monschool-agent
$ curl -sL https://github.com/fossunited/monschool-agent/releases/download/v0.1.0/monschool-agent_0.1.0_linux_amd64.tar.gz | tar xz
$ sudo mv ./monschool-agent.bin /usr/bin
$ mkdir -p /etc/monschool-agent
$ mv ./config.sample.toml /etc/monschool-agent/config.toml
```

### Running as a service

Save this service file at `/etc/systemd/system/monschool.service`

```
[Unit]
Description=Monschool Agent
Documentation=https://github.com/fossunited/monschool-agent/
After=network-online.target
Requires=network-online.target

StartLimitIntervalSec=500
StartLimitBurst=5

[Service]
ExecStart=/usr/bin/monschool-agent.bin --config /etc/monschool-agent/config.toml
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

#### Enable the service

```bash
$ sudo systemctl daemon-reload
$ sudo systemctl enable --now monschool
$ sudo systemctl status monschool
```
