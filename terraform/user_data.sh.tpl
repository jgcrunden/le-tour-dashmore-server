#!/bin/bash

set -eux

useradd -m webhook

mkdir -p ${webhook_conf}

curl -OL https://github.com/adnanh/webhook/releases/download/2.8.2/webhook-linux-arm64.tar.gz
tar zxf webhook-linux-arm64.tar.gz
chown root:root webhook-linux-arm64/webhook
mv webhook-linux-arm64/webhook ${webhook_bin}

cat > ${webhook_conf}/hooks.json << EOF
[
	{
		"id": "upgrade",
		"execute-command": "${webhook_bin}/${script_name}",
		"command-working-directory": "${webhook_bin}",
		"trigger-rule": {
			"match": {
				"type": "value",
				"value": "${webhook_token}",
				"parameter": {
					"source": "url",
					"name": "token"
				}
			}
		}
	}
]
EOF

cat > ${webhook_bin}/${script_name} << EOF
#!/bin/bash

sudo dnf upgrade ${app}
EOF

chmod 755 ${webhook_bin}/${script_name}

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
ExecStart=${webhook_bin}/webhook -hooks ${webhook_conf}/hooks.json -verbose -logfile /home/webhook/webhook.log

[Install]
WantedBy=multi-user.target
EOF

semanage fcontext -a -t bin_t "${webhook_bin}/webhook"
restorecon -vR "${webhook_bin}/webhook"

systemctl enable webhook.service
systemctl start webhook.service

cat > /etc/yum.repos.d/joshuacrunden.repo << EOF
[joshuacrunden]
name=Joshua Crunden
baseurl=https://rpm.joshuacrunden.com
enabled=1
gpgcheck=1
gpgkey=https://rpm.joshuacrunden.com/public.key
EOF

rpm --import https://rpm.joshuacrunden.com/public.key

dnf upgrade -y
dnf install ${app} -y
reboot
