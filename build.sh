#!/bin/bash

APP_NAME="pgrolecheck"
VERSION="1.0.0"
ARCHIVE_NAME="${APP_NAME}-${VERSION}.tar.gz"

# Ensure the current directory is where the script and source code reside
cd /root/rpmbuild/SOURCES/

# Initialize Go module if go.mod does not exist
if [ ! -f go.mod ]; then
    echo "Initializing Go module"
    go mod init ${APP_NAME}
fi

# Tidy up dependencies and ensure go.sum is up to date
echo "Tidying Go module dependencies"
go mod tidy

# Create a directory matching the expected structure and copy files
mkdir -p ${APP_NAME}-${VERSION}
cp main.go pgrolecheck.conf pgrolecheck.service pgrolecheck.1 go.mod go.sum ${APP_NAME}-${VERSION}/

# Create the source archive
echo "Creating source archive"
tar czf ${ARCHIVE_NAME} ${APP_NAME}-${VERSION}

# Build the RPM package
echo "Building RPM package"
rpmbuild -ba /root/rpmbuild/SPECS/pgrolecheck.spec \
         --define "_topdir /root/rpmbuild" \
         --define "version ${VERSION}" \
         --define "name ${APP_NAME}" \
         --define "_rpmdir /build/rpms"
