#!/bin/bash
set -e

sudo systemctl stop fanshim || true

sudo curl -sL -o /usr/local/bin/fanshim_linux_arm64 https://github.com/jsiebens/fanshim-go/releases/download/v0.7.0/fanshim_linux_arm64
sudo chmod 755 /usr/local/bin/fanshim_linux_arm64

sudo mkdir -p /etc/fanshim.d

sudo tee /etc/fanshim.d/env >/dev/null <<EOF
OFF_THRESHOLD=50
ON_THRESHOLD=65
DELAY=5
VERBOSE=false
BRIGHTNESS=50
EOF

sudo tee /etc/systemd/system/fanshim.service >/dev/null <<EOF
[Unit]
Description="FanShim Controller"

[Service]
Type=exec
EnvironmentFile=/etc/fanshim.d/env
ExecStart=/usr/local/bin/fanshim_linux_arm64
KillMode=process
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF
sudo chmod 0600 /etc/systemd/system/fanshim.service

sudo systemctl daemon-reload
sudo systemctl enable fanshim.service
sudo systemctl restart fanshim.service
