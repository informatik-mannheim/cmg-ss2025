# Makefile for building and testing all services in the project
# This Makefile assumes that each service has its own Makefile in the services directory

INTEGRATION_SERVICES = carbon-intensity-provider cli consumer-gateway job job-scheduler user-management worker-deamon worker-gateway worker-registry
DEPLOYMENT_SERVICES = carbon-intensity-provider consumer-gateway job job-scheduler user-management worker-gateway worker-registry
SERVICE_DIR = services

.PHONY: integrationcheck deployment


# Calls the integration check for each service
# Only services that have an integration check defined in their Makefile will be processed
integrationcheck: 
	@for service in $(INTEGRATION_SERVICES); do \
		echo "Test and Building $$service ..."; \
		$(MAKE) -C $(SERVICE_DIR)/$$service integrationcheck; \
	done

# Calls the deployment for each service
# Only services that have a deployment defined in their Makefile will be processed
# That excludes the cli and worker-deamon services
deployment:
	@for service in $(DEPLOYMENT_SERVICES); do \
		echo "Containerizing $$service ..."; \
		$(MAKE) -C $(SERVICE_DIR)/$$service deployment; \
	done
	