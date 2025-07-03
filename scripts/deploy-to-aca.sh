#!/bin/bash

# Configuration
ACR_NAME="cmgregistry"                  # Your ACR name (without .azurecr.io)
RESOURCE_GROUP="cmg-ss2025"             # Azure Resource Group
AZURE_CMD="az"                          # Azure CLI command
DOCKER_CMD="docker"                     # Docker command
COMPOSE_CMD="docker compose"           # Docker Compose command
FULL_ACR_NAME="$ACR_NAME.azurecr.io"    # Full registry name

# Images
IMAGES=(
    "cmg-ss2025-worker-gateway"
    "cmg-ss2025-worker-registry"
    "cmg-ss2025-job-service"
    "cmg-ss2025-job-scheduler"
    "cmg-ss2025-consumer-gateway"
    "cmg-ss2025-user-management"
    "cmg-ss2025-carbon-intensity-provider"
)

# Apps to restart
APPS_TO_RESTART=(
    "worker-registry"
    "worker-gateway"
    "user-management"
    "job-service"
    "job-scheduler"
    "consumer-gateway"
    "carbon-intensity-provider"
)

# Step 1: Login to ACR
echo "Logging in to ACR: $ACR_NAME..."
$AZURE_CMD acr login --name "$ACR_NAME"

# Step 2: Build and start containers via docker compose
echo "Building and starting services using Docker Compose..."
cd "$(dirname "$0")/.." || exit 1  # Navigate to root directory
$COMPOSE_CMD up --build -d

# Step 3: Push all images with 'latest' tag
echo "Pushing updated images to ACR..."
for IMAGE in "${IMAGES[@]}"; do
    echo "Pushing $IMAGE:latest"
    $DOCKER_CMD tag "$IMAGE" "$FULL_ACR_NAME/$IMAGE:latest"
    $DOCKER_CMD push "$FULL_ACR_NAME/$IMAGE:latest"
done

# Step 4: Restart Container Apps
echo "Updating Azure Container Apps..."
for APP in "${APPS_TO_RESTART[@]}"; do
    IMAGE_NAME="$FULL_ACR_NAME/cmg-ss2025-$APP"
    
    echo "Updating $APP to use latest image from $IMAGE_NAME..."
    
    $AZURE_CMD containerapp update \
      --name "$APP" \
      --resource-group "$RESOURCE_GROUP" \
      --image "$IMAGE_NAME:latest"
done

echo "Deployment complete!"