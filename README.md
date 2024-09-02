# Authorization gRPC SSO-Service

## Overview

This project is an Authorization gRPC Service, designed as a microservice for a large To-Do application. The service is responsible for managing user authentication and authorization, handling JWT token creation and verification, and storing user credentials in a PostgreSQL database. It is written in Go (Golang) version 1.23 and uses several third-party libraries and tools to ensure robustness, security, and ease of deployment.

## Features

- **User Authentication**: Handles user login, registration, and password management.
- **JWT Tokens**: Issues and validates JSON Web Tokens (JWT) for secure communication between services.
- **PostgreSQL Integration**: Stores user credentials and related data in a PostgreSQL database.
- **Configuration Management**: Uses `cleanenv` for configuration management.
- **Database Migrations**: Uses `goose` for managing database schema migrations.
- **Containerization**: Fully containerized using Docker and Docker Compose for easy deployment and scaling.

## Tech Stack

- **Golang**: v1.23
- **gRPC**: Communication protocol for the service.
- **PostgreSQL**: Database for storing user data.
- **JWT**: JSON Web Tokens for secure authentication.
- **Cleanenv**: Library for configuration management.
- **Goose**: Tool for database schema migrations.
- **Docker & Docker Compose**: Containerization tools for deployment.

## Getting Started

### Prerequisites

- **Go**: v1.23 or higher
- **Docker**: Ensure Docker and Docker Compose are installed.
- **PostgreSQL**: Running instance of PostgreSQL (can be handled via Docker).

### Configuration

The application configuration is managed using `cleanenv`. You can configure the application via the `config.yaml` file located in the `config/` directory.

Example `config.yaml`:

```yaml
server:
  host: "0.0.0.0"
  port: 50051

database:
  host: "localhost"
  port: 5432
  user: "user"
  password: "password"
  name: "auth_db"

jwt:
  secret_key: "your-secret-key"
  token_expiry: "24h"
```

## Migrations
```bash
goose -dir ./migrations postgres "user=youruser password=yourpassword dbname=yourdb sslmode=disable" up
```

## gRPC API

The service exposes several gRPC endpoints for managing user authentication and authorization. The `.proto` files defining the service and message structures are located in the `proto/` directory.

## Examples

You can interact with the service using gRPC clients like `grpcurl` or by integrating it with other services in your To-Do application.

## Deployment

For production deployments, ensure that the `config.yaml` file is correctly configured with environment-specific values. Consider using Docker orchestration tools like Kubernetes for managing containers at scale.

## Contributing

Feel free to submit issues and pull requests. Please ensure that your code follows the project’s coding standards and is well-documented.
