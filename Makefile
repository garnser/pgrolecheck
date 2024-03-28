# Variables
DOCKER_IMAGE_NAME=pgrolecheck
DOCKER_TAG=latest
RPM_BUILD_DIR=./rpms

# Target for building the Docker image
.PHONY: docker
docker:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .

# Target for building the RPM package using the Docker container
.PHONY: rpm
rpm:
	mkdir -p $(RPM_BUILD_DIR)
	docker run --rm -v "$(PWD)/$(RPM_BUILD_DIR):/build/rpms" $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

# Target for installing the RPM package
.PHONY: install
install:
	sudo dnf install $(RPM_BUILD_DIR)/x86_64/*.rpm -y

# Clean up build artifacts
.PHONY: clean
clean:
	rm -rf $(RPM_BUILD_DIR)/*
	docker rmi $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)
