# Evidence Service

A microservice that provides evidence information based on control IDs. Built with Go and Gin framework.

## Features

- RESTful API endpoint `/getEvidences`
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

**Endpoint:** `GET /getEvidences`

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
curl "http://localhost:8080/getEvidences?controlId=1234"
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
curl "http://localhost:8080/getEvidences?controlId=5678"
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

2. Login to Azure Container Registry:
```bash
task azure:acr-login
```

3. Create Resource Group and Container Apps Environment:
```bash
task azure:create-rg
task azure:create-env
```

### Deploy to Azure

1. Build and push to registry:
```bash
task docker:build
task docker:push
```

2. Deploy to Azure Container Apps:
```bash
task azure:deploy
```

Or do everything in one command:
```bash
task deploy:all
```

### Configuration

Update the following variables in `Taskfile.yml` before deployment:

```yaml
vars:
  DOCKER_REGISTRY: your-registry.azurecr.io
  RESOURCE_GROUP: evidence-service-rg
  LOCATION: eastus
  CONTAINER_APP_NAME: evidence-service
  CONTAINER_APP_ENV: evidence-env
```

Also, make sure to set your Azure subscription ID in the `azure:login` task.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

[Add your license here]
