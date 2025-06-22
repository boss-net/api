# Variables
DOCKER_REGISTRY=boss-net
API_IMAGE=$(DOCKER_REGISTRY)/boss-api
VERSION=latest

build-api:
	@echo "Building API Docker image: $(API_IMAGE):$(VERSION)..."
	docker build -t $(API_IMAGE):$(VERSION) .
	@echo "API Docker image built successfully: $(API_IMAGE):$(VERSION)"

push-api:
	@echo "Pushing API Docker image: $(API_IMAGE):$(VERSION)..."
	docker push $(API_IMAGE):$(VERSION)
	@echo "API Docker image pushed successfully: $(API_IMAGE):$(VERSION)"

# Phony targets
.PHONY: build-api push-api
# Default target
all: build-api push-api