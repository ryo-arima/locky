#!/bin/bash

# LocalStack initialization script for Secrets Manager
# This script creates the initial secret for app.yaml configuration

set -e

echo "Initializing LocalStack Secrets Manager..."

# Wait for LocalStack to be ready
until curl -s http://localhost:4566/_localstack/health | grep -q '"secretsmanager": "available"'; do
  echo "Waiting for LocalStack Secrets Manager to be ready..."
  sleep 2
done

echo "LocalStack is ready. Creating secrets..."

# Path to config file (mounted in container)
CONFIG_FILE="/etc/localstack/init/etc/app.dev.yaml"

# Convert YAML to JSON using Python (yq requires jq which may not be available)
echo "Converting $CONFIG_FILE to JSON..."
CONFIG_JSON=$(python3 -c "
import yaml
import json
with open('$CONFIG_FILE', 'r') as f:
    data = yaml.safe_load(f)
print(json.dumps(data, separators=(',', ':')))
")

awslocal secretsmanager create-secret \
  --name locky/config/app \
  --description "Locky application configuration" \
  --secret-string "$CONFIG_JSON" \
  --region us-east-1 \
  --endpoint-url http://localhost:4566 || echo "Secret may already exist"

echo "✅ Secrets created successfully!"
echo ""
echo "To retrieve the secret:"
echo "  awslocal secretsmanager get-secret-value --secret-id locky/config/app --endpoint-url http://localhost:4566"
