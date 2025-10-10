# Evidence Service

A microservice that provides evidence information based on control IDs. Built with Go and Gin framework.

## Documentation

- **API Documentation (Swagger UI)**: Available at `http://localhost:8080/swagger/index.html`
- **OpenAPI Specification**: Available at `http://localhost:8080/swagger/doc.json`

## Features

- RESTful API endpoint `/v1/evidences`
- Query evidences by control ID (1234 for DataDog, 5678 for Sonar)
- Returns 10-26 evidences per request
- Unique appIDs (A through Z) for each evidence
- Evidence type based on control ID:
  - controlId 1234: DataDog evidences
  - controlId 5678: Sonar evidences
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
  - Use "1234" for DataDog evidences
  - Use "5678" for Sonar evidences

**Response Details:**
- Returns between 10 and 26 evidences per request
- Each evidence has a unique appID (A through Z)
- Evidence type is determined by controlId
- Evidence status is randomly either "SUCCESS" or "FAILED"
- Each evidence has a unique 32-character sysID with "sys_" prefix

**Example Request for DataDog:**
```bash
curl "http://localhost:8080/v1/evidences?controlId=1234"
```

**Example Response for DataDog:**
```json
[
  {
    "evidenceId": "sys_1a2b3c4d5e6f7890abcdef0123456789",
    "evidenceType": "dataDog",
    "controlId": "1234",
    "evidenceStatus": "SUCCESS",
    "appId": "A"
  },
  {
    "evidenceId": "sys_0987654321fedcba0123456789abcdef",
    "evidenceType": "dataDog",
    "controlId": "1234",
    "evidenceStatus": "FAILED",
    "appId": "B"
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
    "appId": "X"
  },
  {
    "evidenceId": "sys_fedcba0987654321abcdef1234567890",
    "evidenceType": "sonar",
    "controlId": "5678",
    "evidenceStatus": "FAILED",
    "appId": "Y"
  }
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

- `task build` - Build the Go application locally
- `task test` - Run tests
- `task dev` - Run the application locally
- `task docker:build` - Build Docker image
- `task docker:run` - Run Docker container in background
- `task docker:stop` - Stop Docker container
- `task docker:logs` - View container logs
- `task docker:status` - Check container status
- `task docker:clean` - Remove Docker image
- `task all` - Build and run Docker container

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

You can deploy the entire application with a single command that will set up all required Azure resources, build and push the Docker image, and deploy the application:

```bash
task deploy:all
```

This command will:
1. Log in to Azure and set up Azure resources:
   - Create Resource Group
   - Create Azure Container Registry (ACR)
   - Log in to ACR
   - Create Container Apps Environment

2. Build and push the Docker image:
   - Build image with AMD64 platform support
   - Tag image for ACR
   - Push to Azure Container Registry

3. Deploy the application:
   - Create Container App
   - Set up external ingress
   - Configure ACR authentication
   - Deploy the application

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
