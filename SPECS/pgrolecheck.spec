Name:           pgrolecheck
Version:        %{version}
Release:        1%{?dist}
Summary:        A service to check the role of a PostgreSQL server.
License:        MIT
URL:            https://example.com/pgrolecheck
Source0:        %{name}-%{version}.tar.gz
BuildRequires:  go, rpm-build
Requires:       postgresql

%global debug_package %{nil}
%global __debug_install_post %{nil}

%description
pgrolecheck is a service that checks whether a PostgreSQL server is running as a primary server or a standby replica. It provides a web interface accessible over HTTPS.

%prep
%setup -q

%build
go mod tidy
go build -o %{name} .

%install
mkdir -p %{buildroot}/usr/local/bin
install -m 0755 %{name} %{buildroot}/usr/local/bin/%{name}
mkdir -p %{buildroot}/etc/pgrolecheck
install -m 0644 pgrolecheck.conf %{buildroot}/etc/pgrolecheck/%{name}.conf
mkdir -p %{buildroot}/usr/lib/systemd/system
install -m 0644 pgrolecheck.service %{buildroot}/usr/lib/systemd/system/%{name}.service
install -Dm644 pgrolecheck.1 %{buildroot}/usr/share/man/man1/pgrolecheck.1

%files
/usr/local/bin/%{name}
/etc/pgrolecheck/%{name}.conf
/usr/lib/systemd/system/%{name}.service
%{_mandir}/man1/pgrolecheck.1*

%pre
getent group postgres >/dev/null || groupadd -r -g 26 postgres
if ! getent passwd postgres >/dev/null; then
    useradd -r -u 26 -g postgres -d /var/lib/pgsql -s /bin/bash -c "PostgreSQL Server" postgres
    mkdir -p /var/lib/pgsql
    chown postgres:postgres /var/lib/pgsql
fi

%post
systemctl daemon-reload
systemctl enable %{name}.service

%preun
if [ $1 -eq 0 ]; then
    systemctl stop %{name}.service
    systemctl disable %{name}.service
    systemctl daemon-reload
fi

%changelog
* Tue Mar 28 2024 Jonathan Petersson <jpetersson@garnser.se> - 1.0.0-1
- Initial RPM package for pgrolecheck.
