#!/bin/bash

set -eux

useradd -m webhook

export webhook_bin="/home/webhook/bin"
mkdir -p ${webhook_bin}

curl -OL https://github.com/adnanh/webhook/releases/download/2.8.2/webhook-linux-arm64.tar.gz
tar zxf webhook-linux-arm64.tar.gz
mv webhook-linux-arm64/webhook ${webhook_bin}

export script_name="upgrade.sh"

cat > ${webhook_bin}/hooks.json << EOF
[
	{
		"id": "upgrade",
		"execute-command": "${webhook_bin}/${script_name}",
		"command-working-directory": "${webhook_bin}",
		"trigger-rule": {
			"match": {
				"type": "value",
				"value": "<PASSWORD>",
				"parameter": {
					"source": "url",
					"name": "token"
				}
			}
		}
	}
]
EOF

export app=le-tour-dashmore-server
cat > ${webhook_bin}/${script_name} << EOF
#!/bin/bash

sudo dnf upgrade ${app}
EOF

export password=$(tr -dc 'A-Za-z0-9' < /dev/urandom | head -c 20)
sed -i "s/<PASSWORD>/${password}/" ${webhook_bin}/hooks.json

chown webhook:webhook -R ${webhook_bin}
chmod 700 ${webhook_bin}/${script_name}

echo "webhook ALL=(ROOT) NOPASSWD: /usr/bin/dnf upgrade ${app}" >> /etc/sudoers

cat > /etc/systemd/system/webhook.service << EOF
[Unit]
Description=Webhook
After=network.target
StartLimitIntervalSec=0
[Service]
Type=simple
Restart=always
RestartSec=1
User=webhook
ExecStart=${webhook_bin}/webhook -hooks ${webhook_bin}/hooks.json -verbose -logfile ${webhook_bin}/webhook.log

[Install]
WantedBy=multi-user.target
EOF

systemctl enable webhook.service
systemctl start webhook.service

cat > /etc/yum.repos.d/joshuacrunden.repo << EOF
[joshuacrunden]
name=Joshua Crunden
baseurl=https://rpm.joshuacrunden.com
enabled=1
gpgcheck=1
gpgkey=https://rpm.joshuacrunden.com/pgp-key.public
EOF

dnf upgrade -y

rpm --import https://rpm.joshuacrunden.com/pgp-key.public
dnf install ${app} -y
reboot
