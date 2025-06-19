# Le Tour d'Ashmore Server

## Setup
Install postgres
```bash
sudo dnf install -y postgres postgres-contrib
```

Initialise postgres
```bash
sudo su - postgres
sudo su - postgres -c "postgresql-setup --initdb"
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

