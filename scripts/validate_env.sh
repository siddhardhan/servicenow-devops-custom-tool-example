#!/usr/bin/env bash
set -euo pipefail

# Load .env if present
if [ -f .env ]; then
  # shellcheck disable=SC1091
  source .env
fi

MISSING=0
check() {
  local name=$1
  if [ -z "${!name:-}" ]; then
    echo "Error: $name is not set"
    MISSING=1
  fi
}

check AZURE_SUBSCRIPTION_ID
check AZURE_TENANT_ID
check AZURE_LOCATION
check AZURE_RESOURCE_GROUP
check AZURE_CONTAINER_REGISTRY
check CONTAINER_APP_NAME
check CONTAINER_APP_ENV
check DOCKER_IMAGE
check DOCKER_TAG

if [ "$MISSING" -ne 0 ]; then
  echo "One or more required environment variables are missing. Please update .env or export the variables."
  exit 1
fi

echo "Environment variables validated."
