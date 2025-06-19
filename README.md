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

## Development
Serve example html pages to avoid continually fetching from the internet
```bash
python -m http.serve -d examples/html
```

Run server
```bash
go run main.go
```
