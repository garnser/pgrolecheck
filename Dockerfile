FROM rockylinux:9

RUN dnf -y update && \
    dnf -y install rpm-build go git && \
    dnf clean all

# Create the rpmbuild directory structure
RUN mkdir -p /root/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}

WORKDIR /root/rpmbuild

COPY SOURCES/ ./SOURCES/
COPY SPECS/pgrolecheck.spec ./SPECS/

COPY build.sh ./build.sh

RUN chmod +x build.sh

CMD ["./build.sh"]

