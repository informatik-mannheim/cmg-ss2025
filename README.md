# cmg-ss2025

Cloud-native Microservices mit Go - SoSe 2025

## Azure Container Registry (ACR) Management Script

A Bash script for managing Docker images in Azure Container Registry with easy push/delete operations.

## 🛠️ Prerequisites
1. **Azure CLI** installed and configured:
   ```bash
   az login
   ```
2. **Docker** installed and Docker Daemon running
3. **Permissions** to access the ACR (Contributor or Owner Role). To do so you have to be added to the azure resource group by authority.

## 📥 Execution

1. Make it executable:
   ```bash
   chmod +x acr_manager.sh
   ``` 

2. Execute from **project root**:
   ```bash
    ./acr_manager.sh
   ```

## 🚀 Usage
```
1 - Login to ACR          # Authenticates with your registry
2 - Push specific image   # Tags and uploads a local image
3 - Push all images       # Pushes all images and tags them as `latest`
4 - Delete image          # Removes tags or purges entire repository
5 - Exit
```

## 🛠️ How It Works
### 🔐 Authentication
- Uses `az acr login` to authenticate Docker with your ACR
- Required only once per session

### 📤 Pushing specific Images
- Enter local image name (e.g., `cmg-ss2025-job-service`)
- Target tag (e.g., `v1`)
- Actions:
   ```bash
   docker tag my-app:latest cmgss2025.azurecr.io/my-app:v1
   docker push cmgss2025.azurecr.io/my-app:v1
   ``` 

### 🗑️ Deleting Images
1. Safe Workflow:
    - First lists all existing tags
    - Chose between:
        - Single tag deletion (`untag`)
        - Full purge (`--purge`) with confirmation

## 📋 Examples
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

* **Token URL:** `https://dev-jqhwcu7xuwgdqi56.eu.auth0.com/oauth/token`
* **Audience:** `https://green-load-shifting-platform/`

---

### 2. 🔁 Auth Flows

#### a. **Consumers (external users)**

* Not permitted to self-register.
* Must contact the system admin for access.
* Receive JWTs via internal issuance if authorized.

#### b. **User Management Provider (Internal)**

* The only service allowed to call the `/login` endpoint.
* Uses `client_credentials` flow to obtain JWT.
* Role-based access ensures only the provider can manage users.

#### c. **Other Internal Services**

* Use **client credentials** grant to authenticate.
* Tokens are attached as `Authorization: Bearer <token>` in requests.

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

* Only the **User Management Provider** is allowed to access the `/register` endpoint.
* Manual onboarding: External users must **contact the administrator** to request access.
* Unauthorized clients will be **rejected at gateway level**.

---

### 5. 🔐 JWT Issuance

* **Endpoint:** `POST /login`
* **Flow:** Resource Owner Password (for users) or Client Credentials (for services)
* **Response:** JWT token with custom claims used for RBAC

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

   * Role: `Key Vault Secrets User`
   * Assign access to: `Managed identity`
   * Select the identity of your **Container App**

#### 🔗 Link Secrets to Container App

1. Go to **Container App → Security → Secrets → + Add**
2. Set:

   * **Name**: Internal reference key (e.g. `db-password`)
   * **Type**: `Key Vault reference`
   * **Key Vault Secret URI**: From your secret’s **Current Version**
   * **Identity**: `System-assigned`
3. Add the secret.

#### 🌱 Use Secrets in Environment Variables

1. Go to **Application → Containers**
2. Scroll to **Environment Variables → + Add**
3. Set:

   * **Name**: Final env var (e.g. `DB_PASSWORD`)
   * **Value Source**: `Secret`
   * **Value**: Select the secret key (e.g. `db-password`)

---

### 7. 📄 Summary Table

| Component          | Description                                                             |
| ------------------ | ----------------------------------------------------------------------- |
| Identity Provider  | Auth0 (`client_credentials` + resource owner password)                  |
| JWT Role Mapping   | Custom claims under `https://green-load-shifting-platform/*` namespace  |
| Registration Flow  | Only via internal provider; no public access                            |
| JWT Auth           | All services authenticate using JWT headers                             |
| Secret Management  | Azure Key Vault + system-assigned managed identity                      |
| Secrets to Runtime | Secrets referenced in Container App → injected as environment variables |

---

Let me know if you want this as a downloadable `.md` file or formatted for Swagger/OpenAPI annotations too!

# List of approved packages

| Package                              | Description                                  |
|--------------------------------------|----------------------------------------------|
| `github.com/gorilla/mux`             | HTTP request router and dispatcher           |
| `github.com/google/uuid`             | UUID generation (e.g., for user identifiers) |


