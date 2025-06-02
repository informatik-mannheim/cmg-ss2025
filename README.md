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

Here’s a clean and structured README section based on your steps, titled **🔐 Using Azure Key Vault Secrets in Container Apps**:

---

## 🔐 Using Azure Key Vault Secrets in Container Apps

Follow these steps to securely provide secrets from Azure Key Vault to your Container App using managed identity:

### 1. Create and Add Secrets in Key Vault

1. Go to your **Azure Key Vault** in the portal.
2. Under **Objects → Secrets**, click **+ Generate/Import**.
3. Provide:

   * A **Name** for your secret.
   * The **Value** (actual secret content).
4. Click **Create**.

---

### 2. Enable Managed Identity on Container App

1. Go to your **Container App** in the Azure portal.
2. Navigate to **Security → Identity**.
3. Under **System-assigned managed identity**, switch **Status** to **On**.
4. Click **Save**.

> ✅ This identity will be used to securely access the Key Vault without hardcoding credentials.

---

### 3. Assign Key Vault Access to the Managed Identity

1. Go back to your **Key Vault**.
2. Navigate to the **Secrets** section, click on the specific secret you want to share.
3. On the left panel, go to **Access Control (IAM)**.
4. Click **+ Add → Add role assignment**.
5. Choose the role **Key Vault Secrets User**.
6. Under **Assign access to**, select **Managed identity**.
7. Click **+ Select members**, choose your **Container App’s managed identity**, and click **Select**.
8. Save the role assignment.

---

### 4. Reference the Key Vault Secret in Container App

1. In the **Container App**, go to **Security → Secrets**.
2. Click **+ Add**.
3. Set:

   * **Name** → a secret key name (used for referencing; not your final env var name).
   * **Type** → Key Vault reference.
   * **Key Vault Secret URI** → found in the Key Vault under your secret, click it, then copy the **Secret Identifier** (Current Version).
   * **Identity** → select **System-assigned**.
4. Click **Add**.

---

### 5. Set Environment Variable for Container

1. Go to **Application → Containers** in your Container App.

2. Scroll to **Environment variables**, click **+ Add**.

3. Set:

   * **Name** → the actual environment variable name your container app will use.
   * **Value Source** → **Secret**.
   * **Value** → select the secret key name you defined in step 4.

4. Save the configuration.

---


# List of approved packages

| Package                              | Description                                  |
|--------------------------------------|----------------------------------------------|
| `github.com/gorilla/mux`             | HTTP request router and dispatcher           |
| `github.com/google/uuid`             | UUID generation (e.g., for user identifiers) |


