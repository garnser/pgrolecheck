# PgRoleCheck

PgRoleCheck is a service designed to check whether a PostgreSQL server is running as a primary server or a standby replica. It provides a web interface accessible over HTTPS, offering an easy way to monitor the role of your PostgreSQL instances.

## Features

- Determines if a PostgreSQL instance is running as a primary or replica.
- Securely accessible over HTTPS.
- Option to run in the foreground and log to STDOUT for debugging.
- Supports logging to a file or syslog for production deployments.
- Configurable via a `.conf` file for easy customization.

## Prerequisites

- Go (1.15 or later recommended)
- PostgreSQL
- RPM Build tools (if building RPM packages)

## Configuration

The service can be configured via the `pgrolecheck.conf` file, which allows specifying database connection details, the web server's listen address and port, SSL certificate details for HTTPS, and logging configuration.

A sample configuration looks like this:

```ini
[database]
dbname=yourdbname
user=youruser
password=yourpassword
host=localhost
port=5432
sslmode=disable

[server]
listen_ip=0.0.0.0
https_port=8443
cert_file=path/to/cert.pem
key_file=path/to/key.pem

[logging]
log_file=/var/log/pgrolecheck.log
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

1. Ensure you have RPM build tools installed.
2. Run `make rpm` from the root of the repository. This will generate an RPM in the ./rpms directory.

### Docker Build
You can also build PgRoleCheck using Docker:

1. Build the Docker image: make docker.
2. Then, you can build the RPM using the created Docker image: make rpm.

## Installation
After building the RPM package, install it with:

```bash
sudo dnf install ./rpms/x86_64/pgrolecheck-1.0.0-1.el9.x86_64.rpm
```

## Running PgRoleCheck
After installation, PgRoleCheck can be started with:

```bash
systemctl start pgrolecheck
```

Ensure you have configured `pgrolecheck.conf` according to your environment before starting the service.

## Contributing
Contributions are welcome! Please feel free to submit pull requests or open issues to discuss potential improvements or features.

## License
PgRoleCheck is released under the MIT License. See the LICENSE file for more details.
