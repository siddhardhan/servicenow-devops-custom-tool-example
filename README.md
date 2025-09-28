# Evidence Service

A microservice that provides evidence information based on control IDs. Built with Go and Gin framework.

## Features

- RESTful API endpoint `/getEvidences`
- Query evidences by control ID
- Randomly generated system IDs for each evidence (32 characters)
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
task all

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

**Example Request:**
```bash
curl "http://localhost:8080/getEvidences?controlId=1234"
```

**Example Response:**
```json
[
  {
    "evidenceId": "sys_1a2b3c4d5e6f7890abcdef0123456789",
    "evidenceType": "document",
    "controlId": "1234",
    "evidenceStatus": "SUCCESS"
  },
  {
    "evidenceId": "sys_0987654321fedcba0123456789abcdef",
    "evidenceType": "audit",
    "controlId": "1234",
    "evidenceStatus": "SUCCESS"
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

Default configuration:
- Port: 8080
- Docker image name: evidence-service
- Docker tag: latest

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

[Add your license here]
