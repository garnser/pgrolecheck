# PgRoleCheck

PgRoleCheck is a service designed to check whether a PostgreSQL server is running as a primary server or a standby replica. It provides a web interface accessible over HTTP(S), offering an easy way to monitor the role of your PostgreSQL instances.

## Features

- Determines if a PostgreSQL instance is running as a primary or replica.
- Securely accessible over HTTPS.
- Option to run in the foreground and log to STDOUT for debugging.
- Supports logging to a file or syslog for production deployments.
- Configurable via a `.conf` file for easy customization.

## Prerequisites

- Postgresql

## Configuration

The service can be configured via the `pgrolecheck.conf` file, which allows specifying database connection details, the web server's listen address and port, SSL certificate details for HTTPS, and logging configuration.

### Configuration Options

#### Database Configuration

- `dbname`: Name of the database to connect to.
- `user`: The user to connect as.
- `password`: Password for the database user.
- `host`: Hostname or IP address of the database server.
- `port`: Port number of the database server.
- `sslmode`: SSL mode for the database connection.

#### Server Configuration

- `listen_ip`: The IP address the web server listens on.
- `use_ssl`: Enable SSL for the web server.
- `https_port`: The port number for HTTPS connections.
- `cert_file`: Path to the SSL certificate file.
- `key_file`: Path to the SSL private key file.
- `output_format`: Configures the format of the response from the server. Options are `json`, `simple`, or `csv`.

#### Logging Configuration

- `log_file`: Path to the log file. Set to "syslog" to use the system logger, or specify a file path.

#### Security Configuration

- `ip_whitelist`: A comma-separated list of IP addresses that are allowed to access the web interface. If left blank, all IPs are allowed.
- `auth_token`: A token that must be provided in the `Authorization` header of requests for access.

### Flags

PgRoleCheck can be started with the following flags:

- `-f`: Runs the application in the foreground, printing logs directly to STDOUT.
- `-dbconfig`: Specifies a database instance to monitor in addition or instead of those specified in `pgrolecheck.conf` file.

Example usage:

\```bash
pgrolecheck -f -dbconfig --dbconfig '{"Name":"database_db12","DbName":"postgres","User":"user","Password":"password","Host":"localhost","Port":"5432","SslMode":"disable"}'
\```

### Sample Configuration

```ini
[database_name1]
dbname=yourdbname
user=youruser
password=yourpassword
host=localhost
port=5432
sslmode=disable

[server]
listen_ip=0.0.0.0
use_ssl=true
https_port=8443
cert_file=path/to/cert.pem
key_file=path/to/key.pem
output_format=json

[logging]
log_file=/var/log/pgrolecheck.log

[security]
ip_whitelist=192.168.1.100,192.168.1.101
auth_token=someVerySecretToken
```

## Building
PgRoleCheck can be built directly using Go or packaged into an RPM for distribution. Here's how you can do both:

### Building with Go
To compile the service directly:

```bash
go build -o pgrolecheck main.go
```

### Building RPM Package
To package PgRoleCheck as an RPM:

1. Build the Docker image: `make docker`.
2. Then, you can build the RPM using the created Docker image: `make rpm`.

## Installation
After building the RPM package, install it with:

```bash
sudo dnf install ./rpms/x86_64/pgrolecheck-1.0.0-1.el9.x86_64.rpm
```

or using
```bash
make install
```

## Running PgRoleCheck
After installation, PgRoleCheck can be started with:

```bash
systemctl start pgrolecheck
```

If you want to run it in the foreground you can start it with:

```bash
pgrolecheck -f
```

## Interacting with PgRoleCheck

### Query Parameters

- **`?dbname=`**: Specify a database name to check its role. Essential for environments monitoring multiple databases.

### Possible Responses

Responses vary by `output_format`:

- **JSON Format**: `{"db1":{"status":"error","message":"pq: password authentication failed for user \"youruser\""}}`
- **Simple Format**: Returns only the status (`primary`, `replica`, or `error`).
- **CSV Format**: Provides `Name,Status,Message` for each database.

## Security and Access Control

Include the `auth_token` in the request header for API access:

\```bash
curl -H "Authorization: Bearer someVerySecretToken" https://localhost:8443/?dbname=primarydb
\```

PgRoleCheck provides a simple HTTP API that can be interacted with using tools like `curl`. Here's how you can use `curl` to check the role of a PostgreSQL server and the possible responses:

### Checking the Role

To check the role of a PostgreSQL server, send an HTTP GET request to the PgRoleCheck service. For example:

```bash
curl https://localhost:8443/
```


## Example Usage
```bash
# Check the role of the PostgreSQL server
curl -k https://127.0.0.1:8443
{"db1":{"status":"primary",""}}
```

You can use these responses to automate monitoring or integration with other systems.

Ensure you have configured `pgrolecheck.conf` according to your environment before starting the service.

## Contributing
Contributions are welcome! Please feel free to submit pull requests or open issues to discuss potential improvements or features.

## License
PgRoleCheck is released under the MIT License. See the LICENSE file for more details.
