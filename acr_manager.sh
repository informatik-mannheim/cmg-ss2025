#!/bin/bash

# Configuration
ACR_NAME="cmgregistry"                  # Your ACR name (without .azurecr.io)
RESOURCE_GROUP="cmg-ss2025"           # Azure Resource Group
AZURE_CMD="az"                        # Azure CLI command
DOCKER_CMD="docker"                   # Docker command
FULL_ACR_NAME="$ACR_NAME.azurecr.io"  # Full registry name

# Images
WORKER_GATEWAY="cmg-ss2025-worker-gateway"
WORKER_REGISTRY="cmg-ss2025-worker-registry"
JOB_SERVICE="cmg-ss2025-job-service"
JOB_SCHEDULER="cmg-ss2025-job-scheduler"
CONSUMER_GATEWAY="cmg-ss2025-consumer-gateway"
USER_MANAGEMENT="cmg-ss2025-user-management"
CARBON_INTENSITY_PROVIDER="cmg-ss2025-carbon-intensity-provider"

# Functions
login_to_acr() {
    echo "### Login to Azure Container Registry ###"
    $AZURE_CMD acr login --name $ACR_NAME
}

push_image() {
    echo "### Push Image to ACR ###"
    read -p "Local image name (e.g., cmg-ss2025-job-service): " LOCAL_IMAGE
    read -p "Tag (e.g., v1): " TAG

    # Remove existing tag if it exists
    $DOCKER_CMD rmi "$FULL_ACR_NAME/$LOCAL_IMAGE:$TAG" 2>/dev/null || true

    # Tag and push
    $DOCKER_CMD tag "$LOCAL_IMAGE" "$FULL_ACR_NAME/$LOCAL_IMAGE:$TAG"
    $DOCKER_CMD push "$FULL_ACR_NAME/$LOCAL_IMAGE:$TAG"

    echo "Successfully pushed: $FULL_ACR_NAME/$LOCAL_IMAGE:$TAG"
}

push_all_images() {
    echo "### Push ALL images with 'latest' tag to ACR ###"

    IMAGES=(
        "$WORKER_GATEWAY"
        "$WORKER_REGISTRY"
        "$JOB_SERVICE"
        "$JOB_SCHEDULER"
        "$CONSUMER_GATEWAY"
        "$USER_MANAGEMENT"
        "$CARBON_INTENSITY_PROVIDER"
    )

    for IMAGE in "${IMAGES[@]}"; do
        echo "🚀 Pushing image: $IMAGE:latest"
        $DOCKER_CMD tag "$IMAGE" "$FULL_ACR_NAME/$IMAGE:latest"
        $DOCKER_CMD push "$FULL_ACR_NAME/$IMAGE:latest"
        echo "✅ Pushed: $FULL_ACR_NAME/$IMAGE:latest"
    done

    echo "🎉 All images have been pushed."
}

delete_image() {
    echo "### Delete Image from ACR ###"
    read -p "Image name in ACR (without tag, e.g., cmg-ss2025-job-service): " IMAGE_NAME
    
    # Show existing tags first
    echo "🔍 Existing tags for $IMAGE_NAME:"
    $AZURE_CMD acr repository show-tags --name $ACR_NAME --repository "$IMAGE_NAME" --output tsv
    
    read -p "Tag to delete (e.g., v1) or '--purge' to delete ALL tags: " TAG

    if [[ "$TAG" == "--purge" ]]; then
        echo "⚠️ WARNING: This will delete ALL tags and manifests for $IMAGE_NAME!"
        read -p "Are you sure? (y/n): " CONFIRM
        if [[ "$CONFIRM" == "y" ]]; then
            $AZURE_CMD acr repository delete \
                --name $ACR_NAME \
                --repository "$IMAGE_NAME" \
                --yes
            echo "✅ Deleted ALL tags for $IMAGE_NAME"
        else
            echo "🚫 Deletion cancelled"
        fi
    else
        echo "ℹ️ This will only delete the specific tag '$TAG'"
        $AZURE_CMD acr repository untag \
            --name $ACR_NAME \
            --image "$IMAGE_NAME:$TAG"
        echo "✅ Deleted tag: $IMAGE_NAME:$TAG"
    fi
}

# Menu
while true; do
    echo ""
    echo "1 - Login to ACR"
    echo "2 - Push specific image"
    echo "3 - Push all images (tagged as latest)"
    echo "4 - Delete image"
    echo "5 - Exit"
    read -p "Choose action (1-4): " CHOICE

    case $CHOICE in
        1) login_to_acr ;;
        2) push_image ;;
        3) push_all_images ;;
        4) delete_image ;;
        5) break ;;
        *) echo "❌ Invalid input" ;;
    esac
done