#!/usr/bin/env bash
set -euo pipefail

# Load .env if present
if [ -f .env ]; then
  # shellcheck disable=SC1091
  source .env
fi

if [ -z "${CONTAINER_APP_NAME:-}" ] || [ -z "${AZURE_RESOURCE_GROUP:-}" ]; then
  echo "CONTAINER_APP_NAME and AZURE_RESOURCE_GROUP must be set in environment or .env"
  exit 1
fi

# Fetch FQDN of the deployed Container App
FQDN=$(az containerapp show --name "$CONTAINER_APP_NAME" --resource-group "$AZURE_RESOURCE_GROUP" --query properties.configuration.ingress.fqdn -o tsv)

if [ -z "$FQDN" ]; then
  echo "Could not determine app FQDN. Is the Container App deployed?"
  exit 1
fi

echo "Found app at: https://$FQDN"

# Run a simple health check
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "https://$FQDN/health" || true)

if [ "$HTTP_STATUS" = "200" ]; then
  echo "Health endpoint returned 200 OK"
else
  echo "Health endpoint returned $HTTP_STATUS (may not exist) — trying smoke endpoint"
  # Try smoke endpoint for by-app Hogan
  PAYLOAD=$(curl -sS "https://$FQDN/v1/evidences/by-app?app_id=Hogan" || true)
  if [ -z "$PAYLOAD" ]; then
    echo "Smoke endpoint returned empty response or timed out"
    exit 2
  fi
  echo "Sample response (truncated):"
  echo "$PAYLOAD" | head -c 1000
fi

echo "Smoke tests completed successfully."
