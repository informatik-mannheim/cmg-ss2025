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
2 - Push image            # Tags and uploads a local image
3 - Delete image          # Removes tags or purges entire repository
4 - Exit
```

## 🛠️ How It Works
### 🔐 Authentication
- Uses `az acr login` to authenticate Docker with your ACR
- Required only once per session

### 📤 Pushing Images
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

# List of approved packages

| Package                              | Description                                  |
|--------------------------------------|----------------------------------------------|
| `github.com/gorilla/mux`             | HTTP request router and dispatcher           |
| `github.com/google/uuid`             | UUID generation (e.g., for user identifiers) |


