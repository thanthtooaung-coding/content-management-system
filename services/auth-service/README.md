# Auth Service Documentation

## Overview
A Go-based authentication service built with clean architecture principles, featuring REST API endpoints with Fiber web framework and Uber's Fx dependency injection.


## HTTP test are in the http_test folder 

> **Note:** GraphQL functionality is currently omitted from this implementation.

---

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker (optional)
- golangci-lint (for development)

### Development Commands
```bash
# Complete development workflow
make all              # Format → Vet → Lint → Test → Build

# Individual commands
make run              # Start development server
make test             # Run tests
make test-coverage    # Generate coverage report
make clean            # Remove build artifacts

# Docker deployment
make docker-build     # Build container image
make docker-run       # Run on localhost:8080
```

---

## 📁 Project Structure

```
auth-service/
├── 🐳 Dockerfile              # Container configuration
├── 🔧 Makefile               # Build automation
├── 📝 README.md              # Project documentation
├── 📦 go.mod/go.sum          # Go dependencies
├── ⚙️  gqlgen.yml            # GraphQL config (unused)
│
├── 🎯 cmd/
│   └── main.go               # Application entry point
│
├── 🔒 internal/              # Private application code
│   ├── config/               # Configuration management
│   ├── graph/                # GraphQL code (unused)
│   ├── handler/              # HTTP request handlers
│   │   ├── graph/            # GraphQL handlers (unused)
│   │   └── rest/             # REST API handlers
│   │       ├── dummy_handler.go
│   │       └── provider/
│   ├── model/                # Data models
│   │   ├── dto/              # Data Transfer Objects
│   │   │   └── user_dto.go
│   │   └── types/            # Internal types
│   │       ├── response.go
│   │       └── users.go
│   ├── repository/           # Data access layer (empty)
│   └── service/              # Business logic
│       └── user_service.go
│
├── 📚 pkg/                   # Reusable packages
│   ├── db/                   # Database connectivity
│   ├── fiber_app/            # Fiber framework setup
│   ├── fx_app/               # Dependency injection
│   ├── logger/               # Logging utilities
│   └── utils/                # General utilities (empty)
│
├── 📄 api/                   # API documentation (empty)
├── 📖 docs/                  # Project docs (empty)
├── 📋 logs/                  # Application logs
├── 🗂️  schema/               # GraphQL schema (unused)
└── 🗃️  tmp/                  # Temporary files
```

---

## 🏗️ Architecture

### Clean Architecture Layers
| Layer | Purpose | Location |
|-------|---------|----------|
| **Presentation** | HTTP/API endpoints | `/internal/handler` |
| **Business Logic** | Core application logic | `/internal/service` |
| **Data Access** | Database operations | `/internal/repository` |
| **Models** | Data structures | `/internal/model` |

### Dependency Injection
- **Framework:** Uber Fx
- **Pattern:** Provider pattern
- **Configuration:** Located in `/provider` subdirectories

---

## 🛠️ Development Tools

### Makefile Commands

#### **Development Workflow**
```bash
make all          # Complete CI pipeline
make build        # Compile to bin/server
make run          # Start development server
make clean        # Remove build artifacts
```

#### **Code Quality**
```bash
make fmt          # Format Go code
make vet          # Static analysis
make lint         # Comprehensive linting
make tidy         # Clean dependencies
```

#### **Testing**
```bash
make test         # Run all tests
make test-coverage # Generate HTML coverage report
```

#### **Docker Operations**
```bash
make docker-build # Build image (server:latest)
make docker-run   # Run container on :8080
```

### Configuration Variables
- `APP_NAME`: Application name (default: `server`)
- `DOCKER_IMAGE`: Docker image tag (default: `server:latest`)

---

## 🔧 Technology Stack

| Component | Technology |
|-----------|------------|
| **Language** | Go |
| **Web Framework** | Fiber |
| **Dependency Injection** | Uber Fx |
| **Architecture** | Clean Architecture |
| **Containerization** | Docker |
| **Code Quality** | golangci-lint, go vet |

---

## 📝 Development Notes

### Current State
- ✅ REST API implementation
- ✅ Clean architecture structure
- ✅ Dependency injection setup
- ✅ Docker support
- ❌ GraphQL (omitted)
- ❌ Database integration (planned)

### Key Features
- **Modular Design**: Clear separation of concerns
- **Dependency Injection**: Uber Fx for clean dependencies
- **Code Quality**: Automated formatting and linting
- **Testing**: Comprehensive test coverage support
- **Docker Ready**: Containerized deployment

### Getting Started Steps
1. **Install Dependencies:** `go mod download`
2. **Run Development:** `make run`
3. **Run Tests:** `make test`
4. **Build for Production:** `make build`
5. **Deploy with Docker:** `make docker-build && make docker-run`

---

## 📊 Project Status
- **Phase:** Development
- **API:** REST only
- **Database:** Not yet integrated
- **Documentation:** In progress
- **Testing:** Framework ready