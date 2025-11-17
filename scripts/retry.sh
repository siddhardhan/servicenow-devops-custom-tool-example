#!/usr/bin/env bash
set -euo pipefail

# Simple retry wrapper: retries a command up to N times with exponential backoff
# Usage: ./scripts/retry.sh 3 2 <command> [args...]
#  - first arg: max attempts (default 3)
#  - second arg: base delay seconds (default 2)

MAX_ATTEMPTS=3
BASE_DELAY=2
if [ $# -ge 1 ]; then
  MAX_ATTEMPTS=$1
fi
if [ $# -ge 2 ]; then
  BASE_DELAY=$2
fi
shift 2 || true

if [ $# -lt 1 ]; then
  echo "Usage: $0 <max_attempts> <base_delay> <command> [args...]"
  exit 2
fi

attempt=1
while true; do
  if "$@"; then
    exit 0
  fi
  if [ $attempt -ge $MAX_ATTEMPTS ]; then
    echo "Command failed after $attempt attempts"
    exit 1
  fi
  sleep_seconds=$((BASE_DELAY * 2 ** (attempt - 1)))
  echo "Attempt $attempt failed — retrying in ${sleep_seconds}s..."
  sleep $sleep_seconds
  attempt=$((attempt + 1))
done
