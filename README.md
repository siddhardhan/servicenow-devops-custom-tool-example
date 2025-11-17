# Evidence Service

A microservice that provides evidence information based on control IDs. Built with Go and Gin framework.

## Documentation

- **API Documentation (Swagger UI)**: Available at `http://localhost:8080/swagger/index.html`
- **OpenAPI Specification**: Available at `http://localhost:8080/swagger/doc.json`

## Features

- RESTful API endpoint `/v1/evidences`
- Query evidences by control ID (1234 for datadog, 5678 for sonar, 9012 for gitlab, 3456 for practitest)
- Returns 10-26 evidences per request
- Unique appIDs for each evidence (e.g., "Corpsite", "Cruz Bike Rentals", etc.)
- Evidence type based on control ID:
  - controlId 1234: datadog evidences
  - controlId 5678: sonar evidences
  - controlId 9012: gitlab evidences
  - controlId 3456: practitest evidences
- Random SUCCESS/FAILED status for each evidence
- Randomly generated 32-character system IDs
- Docker support
- Task-based workflow for development and deployment

## Prerequisites

- Go 1.21 or later
- Docker (for containerized deployment)
- Task (for running task commands)
- Azure CLI (for Azure deployments)
- An Azure subscription with permissions to:
  - Create and manage Container Apps
  - Create and manage Container Registries

## Azure Setup

1. Ensure you're logged into Azure CLI:
```bash
az login
```

2. Create an Azure Container Registry:
```bash
task azure:create-acr
```

This will create an ACR with admin access enabled, which is required for Container App authentication.

3. Get ACR credentials (if needed for manual verification):
```bash
task azure:get-acr-creds
```

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd evidence-service
```

2. Install dependencies:
```bash
go mod download
```

3. Install Task (optional, for using Taskfile commands):
```bash
# On macOS
brew install go-task/tap/go-task

# On other systems, visit: https://taskfile.dev/#/installation
```

## Running the Service

### Local Development

Run the service directly on your machine:

```bash
# Using Go
go run main.go

# Using Task
task dev
```

### Docker Deployment

Build and run using Docker:

```bash
# Build and run in one command
task build:all

# Or run individual commands
task docker:build
task docker:run
```

The service will be available at `http://localhost:8080`

### Docker Commands

```bash
# Start container in background
task docker:run

# View container logs
task docker:logs

# Check container status
task docker:status

# Stop and remove container
task docker:stop

# Remove Docker image
task docker:clean
```

## API Usage

### Get Evidences

**Endpoint:** `GET /v1/evidences`

**Query Parameters:**
- `controlId` (required): The ID of the control to fetch evidences for
  - Use "1234" for datadog evidences
  - Use "5678" for sonar evidences
  - Use "9012" for gitlab evidences
  - Use "3456" for practitest evidences

**Response Details:**
- Returns between 10 and 26 evidences per request
- Each evidence has a unique appID representing a specific service (e.g., "Corpsite", "Cruz Bike Rentals")
- Evidence type is determined by controlId
- Evidence status is randomly either "SUCCESS" or "FAILED"
- Each evidence has a unique 32-character sysID with "sys_" prefix

**Example Request for datadog:**
```bash
curl "http://localhost:8080/v1/evidences?controlId=1234"
```

**Example Response for datadog:**
```json
[
  {
    "evidenceId": "sys_1a2b3c4d5e6f7890abcdef0123456789",
    "evidenceType": "datadog",
    "controlId": "1234",
    "evidenceStatus": "SUCCESS",
    "appId": "Corpsite"
  },
  {
    "evidenceId": "sys_0987654321fedcba0123456789abcdef",
    "evidenceType": "datadog",
    "controlId": "1234",
    "evidenceStatus": "FAILED",
    "appId": "Cruz Bike Rentals"
  }
]
```

**Example Request for Sonar:**
```bash
curl "http://localhost:8080/v1/evidences?controlId=5678"
```

**Example Response for Sonar:**
```json
[
  {
    "evidenceId": "sys_abcdef1234567890fedcba0987654321",
    "evidenceType": "sonar",
    "controlId": "5678",
    "evidenceStatus": "SUCCESS",
    "appId": "Hotel Reservation System"
  },
  {
    "evidenceId": "sys_fedcba0987654321abcdef1234567890",
    "evidenceType": "sonar",
    "controlId": "5678",
    "evidenceStatus": "FAILED",
    "appId": "Portfolio"
  }
]
```

## Get Evidences by App ID

**Endpoint:** `GET /v1/evidences/by-app`

**Query Parameters:**
- `app_id` (required): The application identifier. Supported values:
  - `Hogan` — returns 2 evidences (both SUCCESS)
  - `SinglePoint` — returns 4 evidences (one FAILED — the 3rd)
- `controlIds` (required): Comma-separated list of control IDs to include. Example: `1234,5678`.
- `version` (optional): Optional version string to include in the evidence objects (e.g. `v1.2.3`).

Version-specific behavior
- If no `version` is provided, the endpoint will return all evidences with `evidenceStatus: "SUCCESS"`.
- If `version=R22-2.0`, the endpoint forces all returned evidences to `evidenceStatus: "SUCCESS"`.
- If `version=R22-1.0`, the endpoint guarantees at least one returned evidence has `evidenceStatus: "FAILED"` (if none are failed, the first evidence will be marked FAILED).

**Example Request (Hogan):**
```bash
curl "http://localhost:8080/v1/evidences/by-app?app_id=Hogan&controlIds=1234,5678"
```

**Example Response (Hogan):**
```json
[
  {
    "evidenceId": "sys_abcdef1234567890abcdef1234567890",
    "evidenceType": "datadog",
    "controlId": "1234",
    "evidenceStatus": "SUCCESS",
    "appId": "Hogan"
  },
  {
    "evidenceId": "sys_0123456789abcdef0123456789abcdef",
    "evidenceType": "sonar",
    "controlId": "5678",
    "evidenceStatus": "SUCCESS",
    "appId": "Hogan"
  }
]
```

**Example Request (SinglePoint):**
```bash
curl "http://localhost:8080/v1/evidences/by-app?app_id=SinglePoint&controlIds=1234,5678,9012,3456&version=v1.2.3"
```

**Example Response (SinglePoint):**
```json
[
  { "evidenceId": "sys_a1b2c3d4e5f60123456789abcdefabcd", "evidenceType": "datadog", "controlId": "1234", "evidenceStatus": "SUCCESS", "appId": "SinglePoint" },
  { "evidenceId": "sys_b1c2d3e4f5a60123456789abcdefabcd", "evidenceType": "sonar",   "controlId": "5678", "evidenceStatus": "SUCCESS", "appId": "SinglePoint" },
  { "evidenceId": "sys_c1d2e3f4a5b60123456789abcdefabcd", "evidenceType": "gitlab",  "controlId": "9012", "evidenceStatus": "FAILED",  "appId": "SinglePoint" },
  { "evidenceId": "sys_d1e2f3a4b5c60123456789abcdefabcd", "evidenceType": "practitest","controlId": "3456", "evidenceStatus": "SUCCESS", "appId": "SinglePoint" }
]
```

## Project Structure

```
evidence-service/
├── main.go           # Main application file
├── Dockerfile        # Docker configuration
├── .dockerignore     # Docker ignore file
├── Taskfile.yml      # Task runner configuration
├── go.mod           # Go module file
├── go.sum           # Go module checksum
└── README.md        # This file
```

## Development

### Available Tasks

Development:
- `task build` - Build the Go application locally
- `task test` - Run tests
- `task dev` - Run the application locally for development

Docker Operations:
- `task docker:build` - Build Docker image with AMD64 platform support
- `task docker:tag` - Tag Docker image for registry
- `task docker:push` - Push Docker image to registry
- `task docker:run` - Run Docker container in background
- `task docker:stop` - Stop and remove container
- `task docker:logs` - View container logs
- `task docker:status` - Check container status
- `task docker:clean` - Remove Docker image
- `task build:all` - Build and run Docker container locally with logs

Azure Operations:
- `task azure:login` - Login to Azure
- `task azure:acr-login` - Login to Azure Container Registry
- `task azure:create-rg` - Create Azure Resource Group
- `task azure:create-acr` - Create Azure Container Registry
- `task azure:get-acr-creds` - Get ACR credentials
- `task azure:create-env` - Create Container Apps Environment
- `task azure:deploy` - Deploy to Azure Container Apps (low-level create/update)
- `task deploy:prepare` - Validate environment and prerequisites (runs `scripts/validate_env.sh`)
- `task deploy:infra` - Ensure resource group, ACR and Container Apps env (idempotent)
- `task deploy:auth` - Authenticate to Azure and ACR (`azure:login` + `azure:acr-login`)
- `task deploy:build` - Run local build and tests
- `task deploy:image` - Build Docker image
- `task deploy:publish` - Tag and push image to ACR
- `task deploy:app` - Create / update Container App (wraps `azure:deploy`)
- `task deploy:verify` - Run smoke tests against deployed app (`scripts/smoke_test.sh`)
- `task deploy:all` - Orchestrated deploy (prepare -> infra -> auth -> build -> image -> publish -> app -> verify)
- `task azure:cleanup` - Clean up all Azure resources

### Environment Variables

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Update the following variables in `.env`:

Azure Configuration:
- `AZURE_SUBSCRIPTION_ID`: Your Azure subscription ID
- `AZURE_TENANT_ID`: Your Azure tenant ID
- `AZURE_LOCATION`: Azure region (default: eastus)

Resource Names:
- `AZURE_RESOURCE_GROUP`: Resource group name
- `AZURE_CONTAINER_REGISTRY`: Your ACR name
- `CONTAINER_APP_NAME`: Container App name
- `CONTAINER_APP_ENV`: Container App environment name

Docker Configuration:
- `DOCKER_REGISTRY`: Will be automatically set based on ACR name
- `DOCKER_IMAGE`: Docker image name (default: evidence-service)
- `DOCKER_TAG`: Docker image tag (default: latest)

Application Configuration:
- `PORT`: Application port (default: 8080)
- `GIN_MODE`: Gin framework mode (default: release)

## Cloud Deployment

### Prerequisites

1. Azure CLI installed
2. Azure subscription
3. Azure Container Registry (ACR)
4. Proper permissions to create resources in Azure

### Setup Azure Resources

1. Login to Azure:
```bash
task azure:login
```

2. Create Resource Group and Container Apps Environment:
```bash
task azure:create-rg
```

3. Create Azure Container Registry:
```bash
task azure:create-acr
```
This creates an ACR with admin access enabled for Container App authentication.

4. Login to Azure Container Registry:
```bash
task azure:acr-login
```

5. Create Container Apps Environment:
```bash
task azure:create-env
```

The deployment process will automatically handle ACR authentication using admin credentials.

### Deploy to Azure

You can deploy the application using a modular flow. `task deploy:all` still exists but now orchestrates a set of smaller tasks. This gives you the ability to re-run a specific stage if it fails.

Run the full orchestrated deploy:

```bash
task deploy:all
```

Or run stages individually when troubleshooting:

- `task deploy:prepare` — validate environment and prerequisites (`scripts/validate_env.sh`)
- `task deploy:infra` — ensure Resource Group, ACR and Container Apps Environment (idempotent)
- `task deploy:auth` — authenticate to Azure and ACR
- `task deploy:build` — run unit tests and build the binary
- `task deploy:image` — build the Docker image
- `task deploy:publish` — tag and push the image to ACR
- `task deploy:app` — create or update the Container App
- `task deploy:verify` — run smoke tests against the live app (`scripts/smoke_test.sh`)

Why this helps:

- Easier debugging: re-run only the failing stage (for example `task deploy:publish`) instead of repeating the full flow.
- Safer infra changes: `deploy:infra` is idempotent and will not recreate resources unnecessarily.
- Better CI mapping: each stage can be a separate CI job (lint/test/build/publish/deploy/verify).

If you prefer the old one-shot behavior, `task deploy:all` will run the full orchestrated flow.

### Clean Up Azure Resources

To remove all Azure resources created for this project:

```bash
task azure:cleanup
```

This will:
1. Show all resources that will be deleted
2. Ask for confirmation
3. Delete resources in order:
   - Container App
   - Container Apps Environment
   - Container Registry
   - Resource Group

### Individual Commands

You can also run each step individually if needed:

```bash
# Azure Setup
task azure:login           # Login to Azure
task azure:create-rg       # Create Resource Group
task azure:create-acr      # Create Container Registry
task azure:acr-login       # Login to Container Registry
task azure:create-env      # Create Container Apps Environment

# Docker Operations
task docker:build         # Build the container image
task docker:push          # Push to Container Registry

# Deploy helper tasks (new, modular)
task deploy:prepare      # Validate env and prerequisites
task deploy:infra        # Ensure RG/ACR/Container Apps env
task deploy:auth         # Azure + ACR login
task deploy:build        # Run build and tests
task deploy:image        # Build Docker image
task deploy:publish      # Tag and push image
task deploy:app          # Create/update Container App
task deploy:verify       # Run smoke tests against deployed app

# Deployment
task azure:deploy        # Deploy to Container Apps
```

### Configuration

Update the following variables in `.env` before deployment:

Required variables in `.env`:

```properties
# Azure Configuration
AZURE_SUBSCRIPTION_ID=<your-subscription-id>
AZURE_TENANT_ID=<your-tenant-id>
AZURE_LOCATION=eastus                  # Azure region to deploy to

# Resource Names
AZURE_RESOURCE_GROUP=evidence-service-rg    # Name for your resource group
AZURE_CONTAINER_REGISTRY=evidenceserviceacr # Must be unique across Azure
CONTAINER_APP_NAME=evidence-service         # Name for your container app
CONTAINER_APP_ENV=evidence-env             # Name for your container app environment

# Docker Configuration
DOCKER_IMAGE=evidence-service              # Name of your Docker image
DOCKER_TAG=latest                         # Tag for your Docker image

# Application Configuration
PORT=8080                                 # Port your application listens on
GIN_MODE=release                          # Gin framework mode
```

Note: The ACR name (`AZURE_CONTAINER_REGISTRY`) must be:
- Globally unique across Azure
- 5-50 characters long
- Contain only lowercase letters and numbers

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

[Add your license here]
