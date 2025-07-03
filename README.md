# cmg-ss2025

Cloud-native Microservices mit Go - SoSe 2025

## Deploy to Azure Container Apps

The script `/scripts/deploy-to-aca.sh` was created to deploy the newest images to the Azure Container Apps automatically.

## Execution

1. Navigate from root to `/scripts`

2. Make file executable:

   ```bash
   chmod +x deploy-to-aca.sh
   ```

3. Execute file:
   ```bash
    ./deploy-to-aca.sh
   ```

## Prerequisites

1. **Azure CLI** installed and configured:
   ```bash
   az login
   ```
2. **Docker** installed and Docker Daemon running
3. **Permissions** to access the ACR (Contributor or Owner Role). To do so you have to be added to the azure resource group by authority.

## How It Works

1. Uses `az acr login` to authenticate Docker with your ACR
2. Builds all images
3. Pushes images to the Azure Container Registry if they aren't up to date and tags images with version `latest`
4. Updates all App Containers with the newest image from the Azure Container Registry

---

## Azure Container Registry (ACR) Management Script

The script `/scripts/manage-acr-manually.sh` was created to manually configure the Azure Container Registry.

## Prerequisites

1. **Azure CLI** installed and configured:
   ```bash
   az login
   ```
2. **Docker** installed and Docker Daemon running
3. **Permissions** to access the ACR (Contributor or Owner Role). To do so you have to be added to the azure resource group by authority.

## Execution

1. Navigate from root to `/scripts`

2. Make it executable:

   ```bash
   chmod +x manage-acr-manually.sh
   ```

3. Execute:
   ```bash
    ./acr_manamanage-acr-manuallyger.sh
   ```

## Usage

```
1 - Login to ACR          # Authenticates with your registry
2 - Push specific image   # Tags and uploads a local image
3 - Push all images       # Pushes all images and tags them as `latest`
4 - Delete image          # Removes tags or purges entire repository
5 - Exit
```

## How It Works

### Authentication

- Uses `az acr login` to authenticate Docker with your ACR
- Required only once per session

### Pushing specific Images

- Enter local image name (e.g., `cmg-ss2025-job-service`)
- Target tag (e.g., `v1`)
- Actions:
  ```bash
  docker tag my-app:latest cmgss2025.azurecr.io/my-app:v1
  docker push cmgss2025.azurecr.io/my-app:v1
  ```

### Deleting Images

1. Safe Workflow:
   - First lists all existing tags
   - Chose between:
     - Single tag deletion (`untag`)
     - Full purge (`--purge`) with confirmation

## Examples

### Push an Image:

```bash
$ ./acr_manager.sh
Choose action (1-4): 2
Local image name: cmg-ss2025-job-service
Tag: v1
✅ Successfully pushed: cmgss2025.azurecr.io/my-app:v1
```

### Delete an Image:

```bash
$ ./acr_manager.sh
Choose action (1-4): 3
Image name: my-app
Existing tags: latest, v1, v2
Tag to delete: v1
✅ Deleted tag: my-app:v1
```

---

## 🔐 Authentication & Secrets Overview

This document outlines the authentication flows, JWT usage, registration restrictions, and secure secret handling for the Green Load Shifting Platform.

---

### 1. ✅ Identity Provider

We use **Auth0** as our centralized identity provider.

- **Token URL:** `https://dev-jqhwcu7xuwgdqi56.eu.auth0.com/oauth/token`
- **Audience:** `https://green-load-shifting-platform/`

---

### 2. 🔁 Auth Flows

#### a. **Consumers (external users)**

- Not permitted to self-register.
- Must contact the system admin for access.
- Receive JWTs via internal issuance if authorized.

#### b. **User Management Provider (Internal)**

- The only service allowed to call the `/auth/login` endpoint.
- Uses `client_id.client_secret` flow to obtain JWT.
- Role-based access ensures only the provider can manage users.

#### c. **Other Internal Services**

- Use **client credentials** grant to authenticate.
- Tokens are attached as `Authorization: Bearer <token>` in requests.

---

### 3. 📦 JWT Payload Example

A sample JWT issued to an internal service using the Client Credentials flow:

```json
{
  "https://green-load-shifting-platform/role": "dummy_role",
  "https://green-load-shifting-platform/client_id": "QgXJrkSv5Z5dF8hc8wrfODv2VOHeWBj9",
  "iss": "https://dev-jqhwcu7xuwgdqi56.eu.auth0.com/",
  "sub": "QgXJrkSv5Z5dF8hc8wrfODv2VOHeWBj9@clients",
  "aud": "https://green-load-shifting-platform/",
  "iat": 1746911116,
  "exp": 1746997516,
  "gty": "client-credentials",
  "azp": "QgXJrkSv5Z5dF8hc8wrfODv2VOHeWBj9",
  "permissions": []
}
```

---

### 4. 📝 User Registration

- Only the **User Management Provider** is allowed to access the `/register` endpoint.
- Manual onboarding: External users must **contact the administrator** to request access.
- Unauthorized clients will be **rejected at gateway level**.

---

### 5. 🔐 JWT Issuance

- **Endpoint:** `POST /login`
- **Flow:** Resource Owner Password (for users) or Client Credentials (for services)
- **Response:** JWT token with custom claims used for RBAC
- **Type of Roles:** consumer / provider / job_scheduler

---

### 6. 🧠 Secret Management with Azure Key Vault

Secrets such as DB credentials, certificates, or external API keys are **never stored in code**. They are securely managed using Azure Key Vault and injected into container apps via managed identity.

#### 🔑 Key Vault Setup

1. Go to **Key Vault → Secrets → + Generate/Import**
2. Add your secret name and value.
3. Save the secret.

#### 🆔 Enable Managed Identity

1. Go to your **Container App → Security → Identity**
2. Turn **System-assigned managed identity** to `On` and save.

#### 🛡️ Assign Key Vault Permissions

1. Go back to **Key Vault → Access Control (IAM)**
2. **Add Role Assignment**

   - Role: `Key Vault Secrets User`
   - Assign access to: `Managed identity`
   - Select the identity of your **Container App**

#### 🔗 Link Secrets to Container App

1. Go to **Container App → Security → Secrets → + Add**
2. Set:

   - **Name**: Internal reference key (e.g. `db-password`)
   - **Type**: `Key Vault reference`
   - **Key Vault Secret URI**: From your secret’s **Current Version**
   - **Identity**: `System-assigned`

3. Add the secret.

#### 🌱 Use Secrets in Environment Variables

1. Go to **Application → Containers**
2. Scroll to **Environment Variables → + Add**
3. Set:

   - **Name**: Final env var (e.g. `DB_PASSWORD`)
   - **Value Source**: `Secret`
   - **Value**: Select the secret key (e.g. `db-password`)

---

### 7. 📄 Summary Table

| Component          | Description                                                             |
| ------------------ |-------------------------------------------------------------------------|
| Identity Provider  | Auth0 (`client_id.client_secret` + resource owner password)             |
| JWT Role Mapping   | Custom claims under `https://green-load-shifting-platform/*` namespace  |
| Registration Flow  | Only via internal provider; no public access                            |
| JWT Auth           | All services authenticate using JWT headers                             |
| Secret Management  | Azure Key Vault + system-assigned managed identity                      |
| Secrets to Runtime | Secrets referenced in Container App → injected as environment variables |

---

## 🚀 CI/CD Pipeline & Build System

This project uses **Azure DevOps Pipelines** for automated testing, building, and containerization of all microservices. The pipeline system is designed to work with **Makefiles** that are specifically configured for the CI/CD environment.

### 📋 Pipeline Overview

The [`azure-pipelines.yml`](azure-pipelines.yml) file defines our CI/CD workflow:

1. **Dependency Installation**: Installs Go, make, Docker, and required testing tools (`go-junit-report`, `gocover-cobertura`)
2. **Change Detection**: Identifies which services have been modified in pull requests
3. **Integration Testing**: Runs tests, coverage analysis, builds, and containerization for affected services
4. **Test Reporting**: Publishes JUnit test results and code coverage to Azure DevOps
5. **Image Building**: Creates Docker images for changed services and pushes them to Azure Container Registry

### 🔧 Makefile Structure

The project uses a two-level Makefile system:

#### Root Makefile
The root [`Makefile`](Makefile) orchestrates operations across all microservices:
- Detects which services have been modified
- Runs integration checks on affected services
- Coordinates containerization of all services

#### Service-Level Makefiles
Each microservice has its own Makefile (e.g., [`services/job/Makefile`](services/job/Makefile)):

```makefile
BINARY_NAME=job-binary
DOCKER_IMAGE=job

.PHONY: all build test containerize clean integrationcheck deployment junit coverage

integrationcheck: clean test build

deployment: clean test build containerize

clean:
	rm -f $(BINARY_NAME)
	rm -f report.xml
	rm -f coverage.out
	rm -f coverage.xml

test: junit coverage

junit:
	go test -v ./... | go-junit-report > report.xml

coverage:
	go test -coverprofile=coverage.out ./...
	gocover-cobertura < coverage.out > coverage.xml

build:
	go build -o $(BINARY_NAME) .

containerize:
	docker build -t $(DOCKER_IMAGE) .
```

### ⚠️ Important: Pipeline-Only Execution

**The Makefiles are designed exclusively for the CI/CD pipeline environment and are NOT intended for local development.**

**Why you shouldn't run `make` commands locally:**

1. **Missing Dependencies**: The pipeline installs specific tools (`go-junit-report`, `gocover-cobertura`) that may not be available locally
2. **Environment Differences**: The pipeline runs on Ubuntu with specific configurations and Docker daemon
3. **File Path Assumptions**: Some commands assume the pipeline's directory structure
4. **Azure Credentials**: Image pushing requires Azure Container Registry credentials only available in the pipeline

### 📊 Pipeline Targets

| Target | Purpose | When Used |
|--------|---------|-----------|
| `integrationcheck` | Run tests, coverage, build, and containerize | Every PR and commit |
| `containerize` | Build Docker images for the service | Part of integration check |
| `clean` | Remove build artifacts and old images | Before each build |
| `test` | Run unit tests with coverage reporting | Part of integration check |

### 🔄 Workflow Examples

#### Pull Request Workflow:
1. Developer creates PR → Pipeline detects changed services
2. Runs `make integrationcheck` on affected services
3. Runs go test and build binary on selected Microservice
4. Reports test results and coverage to Azure DevOps
5. Blocks merge if tests fail

#### Main Branch Workflow:
1. Code merged to `main` → Pipeline runs full integration check
2. Builds Docker images for all changed services
3. Pushes images to Azure Container Registry with `latest` tag
4. **Manual deployment still required** using [`/scripts/deploy-to-aca.sh`](scripts/deploy-to-aca.sh)

### 🚨 Important: Manual Deployment Required

**The pipeline only builds Docker. Actual deployment to Azure Container Apps must be done manually using the deployment script.**

After the pipeline completes successfully:

1. Navigate to `/scripts`
2. Run `./deploy-to-aca.sh` to deploy the newly built images to Azure Container Apps


---

# List of approved packages

| Package                                                           | Description                                                                |
| ----------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `github.com/gorilla/mux`                                          | HTTP request router and dispatcher                                         |
| `github.com/google/uuid`                                          | UUID generation (e.g., for user identifiers)                               |
| `go.opentelemetry.io/otel`                                        | Core OpenTelemetry API for tracing and metrics                             |
| `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp` | Sends trace data to an OTLP-compatible backend over HTTP                   |
| `go.opentelemetry.io/otel/sdk`                                    | SDK implementation for OpenTelemetry (tracer provider, processors, etc.)   |
| `go.opentelemetry.io/otel/trace`                                  | Trace-related types and interfaces (e.g., Tracer, Span) from OpenTelemetry |
| `github.com/jackc/pgx/v5`                                         | PostgreSQL Driver and Toolkit                                              |
| `github.com/cenkalti/backoff/v5`                                  | Transitive dependency                                                      |
| `github.com/go-logr/logr`                                         | Transitive dependency                                                      |
| `github.com/go-logr/stdr`                                         | Transitive dependency                                                      |
| `github.com/grpc-ecosystem/grpc-gateway/v2`                       | Transitive dependency                                                      |
| `go.opentelemetry.io/auto/sdk`                                    | Transitive dependency                                                      |
| `go.opentelemetry.io/otel/exporters/otlp/otlptrace`               | Transitive dependency                                                      |
| `go.opentelemetry.io/otel/metric`                                 | Transitive dependency                                                      |
| `go.opentelemetry.io/proto/otlp`                                  | Transitive dependency                                                      |
| `golang.org/x/net`                                                | Transitive dependency                                                      |
| `golang.org/x/sys`                                                | Transitive dependency                                                      |
| `golang.org/x/text`                                               | Transitive dependency                                                      |
| `google.golang.org/genproto/googleapis/api`                       | Transitive dependency                                                      |
| `google.golang.org/genproto/googleapis/rpc`                       | Transitive dependency                                                      |
| `google.golang.org/grpc`                                          | Transitive dependency                                                      |
| `google.golang.org/protobuf`                                      | Transitive dependency                                                      |
