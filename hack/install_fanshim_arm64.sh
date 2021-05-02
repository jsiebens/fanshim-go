#!/bin/bash
set -e

sudo curl -L -o /usr/local/bin/fanshim_linux_arm64 https://github.com/jsiebens/fanshim-go/releases/download/v0.3.0/fanshim_linux_arm64
sudo chmod 755 /usr/local/bin/fanshim_linux_arm64


sudo tee /etc/systemd/system/fanshim.service >/dev/null <<EOF
[Unit]
Description="FanShim Controller"

[Service]
Type=exec
ExecStart=/usr/local/bin/fanshim_linux_arm64 -delay 5 -verbose --on-threshold 80 --off-threshold 65
KillMode=process
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF
sudo chmod 0600 /etc/systemd/system/fanshim.service

sudo systemctl enable fanshim.service
sudo systemctl start fanshim.service
