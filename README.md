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

### Gotchas
Edit `/var/lib/pgsql/data/pg_hba.conf` to ensure it allows the access for users:
e.g.
```bash
# "local" is for Unix domain socket connections only
local   all             all                                     md5
# IPv4 local connections:
host    all             all             127.0.0.1/32            md5
# IPv6 local connections:
host    all             all             ::1/128                 md5
```

Set up database, user and tables:
```bash
CREATE DATABASE letour;
CREATE USER letour WITH ENCRYPTED PASSWORD '<PASSWORD_HERE>';
GRANT ALL PRIVILEGES ON DATABASE letour TO letour;
\c letour postgres
GRANT ALL ON SCHEMA public TO letour;
\c letour letour
```

Create table:
```bash
CREATE TABLE rider(
  id SERIAL PRIMARY KEY,
  name VARCHAR,
  team VARCHAR,
  points INT
);
````



## Development
Serve example html pages to avoid continually fetching from the internet
```bash
python -m http.serve -d examples/html
```

Run server
```bash
go run main.go
```

