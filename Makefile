# Top-level Makefile for Multi-Service Go Project
# 
# This Makefile orchestrates building, testing, and deploying multiple Go services.
# It assumes each service has its own Makefile located in services/<service>/.
#
# Key Features:
# - Intelligent change detection: Only tests services with modified files
# - Support for both PR and direct builds via BUILD_REASON environment variable
# - Parallel service operations where possible
# - Comprehensive error handling and logging

# List of services that support full deployment (containerization)
# Add new services here when they're ready for deployment
DEPLOYMENT_SERVICES = carbon-intensity-provider consumer-gateway job job-scheduler user-management worker-gateway worker-registry

# Base directory containing all microservices
SERVICE_DIR = services

.PHONY: integrationcheck deployment

# Target: integrationcheck
# Purpose: Run integration checks (clean, test, build) only for services with changed files
# Usage: Called by Azure Pipelines for PR validation
# Logic: 
#   1. Detect changed files using git diff
#   2. Extract service names from changed file paths
#   3. Run integration checks only for affected services
#   4. Skip execution if no service changes detected
integrationcheck:
	@echo "Detecting changed files since target branch..."; \
	if [ "$$BUILD_REASON" = "PullRequest" ]; then \
		echo "PR build detected, using Azure DevOps variables"; \
		CHANGED_FILES=$$(git diff --name-only $$(git merge-base HEAD origin/main) HEAD); \
	else \
		echo "Direct build, comparing with main"; \
		CHANGED_FILES=$$(git diff --name-only origin/main...HEAD 2>/dev/null || git diff --name-only HEAD~1 HEAD); \
	fi; \
	echo "Changed files: $$CHANGED_FILES"; \
	CHANGED_SERVICES=""; \
	for file in $$CHANGED_FILES; do \
		if echo $$file | grep -q "^$(SERVICE_DIR)/"; then \
			SERVICE=$$(echo $$file | cut -d'/' -f2); \
			if [ ! -z "$$SERVICE" ] && [ -d "$(SERVICE_DIR)/$$SERVICE" ]; then \
				if ! echo "$$CHANGED_SERVICES" | grep -q "$$SERVICE"; then \
					CHANGED_SERVICES="$$CHANGED_SERVICES $$SERVICE"; \
				fi; \
			fi; \
		fi; \
	done; \
	CHANGED_SERVICES=$$(echo $$CHANGED_SERVICES | xargs); \
	if [ -z "$$CHANGED_SERVICES" ]; then \
		echo "No service changes detected. Skipping integration check."; \
	else \
		echo "Running integration checks for changed services: $$CHANGED_SERVICES"; \
		for service in $$CHANGED_SERVICES; do \
			echo "Testing and building $$service..."; \
			if [ -f "$(SERVICE_DIR)/$$service/Makefile" ]; then \
				$(MAKE) -C $(SERVICE_DIR)/$$service integrationcheck || exit 1; \
			else \
				echo "Warning: No Makefile found for service $$service"; \
			fi; \
		done; \
	fi

# Target: deployment
# Purpose: Run deployment only for changed and deployable services
# Usage: Called during production deployment after merge to main
# Logic:
#   - Detects changed files using git diff between HEAD and HEAD~1
#   - Extracts affected services from changed paths
#   - Filters to services listed in DEPLOYMENT_SERVICES
#   - Verifies the presence of a Makefile and a 'deployment' target per service
#   - Executes deployment only for qualified services
deployment:
	@echo "Detecting changed files for deployment..."; \
	git fetch --all --prune; \
	if [ "$$BUILD_REASON" = "PullRequest" ]; then \
		echo "PR build – skipping deployment."; \
		exit 0; \
	else \
		CHANGED_FILES=$$(git diff --name-only HEAD~1 HEAD); \
	fi; \
	echo "Changed files: $$CHANGED_FILES"; \
	CHANGED_SERVICES=""; \
	for file in $$CHANGED_FILES; do \
		if echo $$file | grep -q "^$(SERVICE_DIR)/"; then \
			SERVICE=$$(echo $$file | cut -d'/' -f2); \
			if [ -n "$$SERVICE" ] && [ -d "$(SERVICE_DIR)/$$SERVICE" ]; then \
				if ! echo "$$CHANGED_SERVICES" | grep -qw "$$SERVICE"; then \
					CHANGED_SERVICES="$$CHANGED_SERVICES $$SERVICE"; \
				fi; \
			fi; \
		fi; \
	done; \
	CHANGED_SERVICES=$$(echo $$CHANGED_SERVICES | xargs); \
	if [ -z "$$CHANGED_SERVICES" ]; then \
		echo "No service changes detected. Skipping deployment."; \
	else \
		echo "Changed services: $$CHANGED_SERVICES"; \
		for service in $$CHANGED_SERVICES; do \
			if echo "$(DEPLOYMENT_SERVICES)" | grep -qw "$$service"; then \
				echo "Deploying $$service..."; \
				if [ -f "$(SERVICE_DIR)/$$service/Makefile" ]; then \
					if grep -q "^deployment:" "$(SERVICE_DIR)/$$service/Makefile"; then \
						$(MAKE) -C $(SERVICE_DIR)/$$service deployment || exit 1; \
					else \
						echo "No 'deployment' target in Makefile for $$service – skipping."; \
					fi; \
				else \
					echo "No Makefile found for $$service – skipping."; \
				fi; \
			else \
				echo "$$service is not in deployable service list – skipping."; \
			fi; \
		done; \
	fi