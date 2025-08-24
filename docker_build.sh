# cmd to grant permission to run
# chmod +x docker_build.sh
# Set the base image name
IMAGE_NAME=social-todo-image
IMAGE_CACHED_NAME=social-todo-lib-cached

# First argument to script (optional) – if provided, we will build cached image
CACHED_BUILD=$1

# If a cached build flag/argument is passed in
if [[ "$MODE" == "cache" || "$MODE" == "1" ]]; then
    echo "Docker building cached image..."

    # Remove old cached image if it exists
    docker rmi ${IMAGE_CACHED_NAME} ${IMAGE_NAME}

    # Build a new cached image using Dockerfile-cache
    # This is usually used to cache Go modules (go mod download)
    docker build -t ${IMAGE_CACHED_NAME} -f Dockerfile-cache .
fi

echo "Docker building main image..."

# Build the main multi-stage Docker image using Dockerfile-multi-stage
docker build -t ${IMAGE_NAME}:latest .

echo "Done!!"
