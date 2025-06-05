# Makefile for building and testing all services in the project
# This Makefile assumes that each service has its own Makefile in the services directory

SERVICES = carbon-intensity-provider cli consumer-gateway job job-scheduler user-management worker-deamon worker-gateway worker-registry
SERVICE_DIR = services

.PHONY: all build test containerize

all: build test containerize

# Targets for building and testing all services
# You can run `make` to build and test all services
# or run `make build` to only build them, or `make test` to only test them.
# Each service's Makefile should define its own build and test targets

build:
	@for service in $(SERVICES); do \
		echo "Building $$service ..."; \
		$(MAKE) -C $(SERVICE_DIR)/$$service build; \
	done

test:
	@for service in $(SERVICES); do \
		echo "Testing $$service ..."; \
		$(MAKE) -C $(SERVICE_DIR)/$$service test; \
	done

containerize:
	@for service in $(SERVICES); do \
		echo "Containerizing $$service ..."; \
		$(MAKE) -C $(SERVICE_DIR)/$$service containerize; \
	done
