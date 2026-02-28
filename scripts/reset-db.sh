#!/usr/bin/env bash
set -euo pipefail

CONTAINER_NAME="edh-tracker-db"
IMAGE_NAME="edh-tracker-db"
PORT="3306:3306"

run_container() {
  docker run --detach --name="$CONTAINER_NAME" --publish "$PORT" "$IMAGE_NAME"
}

echo "Starting fresh $CONTAINER_NAME..."

if output=$(run_container 2>&1); then
  echo "DB running: $output"
  exit 0
fi

# Extract existing container ID from conflict error
existing_id=$(echo "$output" | grep -oP 'container "\K[^"]+')

if [[ -z "$existing_id" ]]; then
  echo "Unexpected error:" >&2
  echo "$output" >&2
  exit 1
fi

echo "Removing existing container ${existing_id:0:12}..."
docker stop "$existing_id" > /dev/null
docker rm "$existing_id" > /dev/null

echo "Starting fresh $CONTAINER_NAME..."
new_id=$(run_container)
echo "DB running: $new_id"
