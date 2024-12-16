# Task Management Application

This is a Task Management Application built with Go, Gin, and PostgreSQL. The application provides APIs for managing tasks, including creating, updating, and retrieving tasks. It also includes authentication and authorization features using JWT.

## Table of Contents

- [Introduction](#introduction)
  - [Technologies Used](#technologies-used)
  - [Folder Structure](#folder-structure)
  - [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Running the Application](#running-the-application)
  - [Accessing Swagger Documentation](#accessing-swagger-documentation)
- [API Documentation](#api-documentation)
  - [Authentication](#authentication)
  - [Tasks](#tasks)
  - [Task Summary](#task-summary)
- [Example Requests](#example-requests)
  - [Create a Task](#create-a-task)
  - [Get a Task](#get-a-task)

## Introduction

### Technologies Used

- **Go**: The main programming language used for the backend.
- **Gin**: A web framework written in Go.
- **PostgreSQL**: The database used for storing task and user information.
- **PGX**: A PostgreSQL driver and toolkit for Go.
- **JWT**: Used for authentication and authorization.
- **Docker**: Used for containerizing the application.

### Folder Structure

The project follows the Hexagonal Architecture principles. Below is an overview of the folder structure:

**go-task-management**

```
├── cmd
│   ├── docs
│   │   └── main.go
│   └── server
│       └── main.go
├── configs
├── infrastructure
│   ├── postgres.go
│   └── query-builder.go
├── internal
│   ├── adapter
│   ├── core
│   │   ├── domain
│   │   ├── error
│   │   ├── port
│   │   └── service
│   │       ├── auth-svc
│   │       └── task-svc
│   ├── handler
│   │   ├── auth-hdl
│   │   └── task-hdl
│   ├── middleware
│   ├── repository
│   │   └── postgres
│   │       ├── auth-repo
│   │       ├── task-repo
│   │       └── user-repo
│   ├── router
│   └── util
├── migrations
├── pkg
├── property
└── server
    └── server.go
```

### Features

- **Role-Based Access Control**: Two types of users - Employer and Employee.
- **Task Management**: Create, update, and retrieve tasks.
- **Authentication**: JWT-based authentication.
- **Swagger Documentation**: API documentation using Swaggo.

## Getting Started

### Prerequisites

- **Go**: Ensure you have Go installed. You can download it from [here](https://golang.org/dl/).
- **Docker**: Ensure you have Docker installed. You can download it from [here](https://www.docker.com/products/docker-desktop).

### Installation

1. **Clone the repository**:

   ```sh
   git clone https://github.com/vuongthanh148/go-task-management
   cd go-task-management
   ```

2. **Set up environment variables**:
   Create a `.env` file in the root directory and add the following:

   ```env
    SERVICE_NAME="go-task-management"
    HOST="localhost"
    PORT="8080"
    RUN_LOCAL="true"
    API_DOCS="true"
    API_DOCS_SCHEMA="http"
    API_DOCS_VERSION="v0.0.1"

    DEBUG_MODE="true"
    CONSOLE_FORMAT="true"
    SHUTDOWN_TIMEOUT="30"

    POSTGRES_PASSWORD_SECRET="password"
    JWT_SECRET_KEY="newsecret"

    POSTGRES_HOST="127.0.0.1"
    POSTGRES_PORT="5432"
    POSTGRES_USER="user"
    POSTGRES_DATABASE="taskdb"

    POSTGRES_CONN_URI="host=%s port=%s database=%s user=%s password=%s"
    POSTGRES_MAX_CONN_LIFETIME="1h"
    POSTGRES_MAX_CONN_IDLE_TIME="30h"
    POSTGRES_MAX_CONNS="4"
    POSTGRES_MIN_CONNS="0"
   ```

3. **Run Docker Compose**:

   ```sh
   docker-compose up -d --build
   ```

4. **Generate Migration File**:

   ```sh
   make new-migration
   ```

5. **Run Migrations**:
   By default, migrations are automatically run on app start, but you can run them manually with the command below:

   ```sh
   make migrate-up
   ```

6. **Linting**:

   ```sh
   make lint
   ```

### Accessing Swagger Documentation

Once the application is running, you can access the Swagger documentation by navigating to the following URL in your web browser:

```
http://localhost:8080/swagger/index.html
```

This will open the Swagger UI where you can explore and test the API endpoints.

### API Documentation

The API documentation is available at `/docs/swagger/index.html` when the server is running.

#### Authentication

- **POST /api/v1/auth/register**: Register a new user
- **POST /api/v1/auth/login**: Login and obtain a JWT
- **POST /api/v1/auth/refresh-token**: Refresh the access token

#### Tasks

- **GET /api/v1/tasks**: Retrieve all tasks (requires authentication)
- **GET /api/v1/tasks/assignee/:assigneeID**: Retrieve tasks by assignee (requires authentication)
- **GET /api/v1/tasks/:taskID**: Retrieve a task by ID (requires authentication)
- **POST /api/v1/tasks**: Create a new task (requires authentication)
- **PATCH /api/v1/tasks/:taskID**: Update a task (requires authentication)
- **DELETE /api/v1/tasks/:taskID**: Delete a task (requires authentication)

#### Task Summary

- **GET /api/v1/tasks/summary**: Retrieve task summary for employees (requires authentication)

### Example Requests

#### Create a Task

```sh
curl -X POST http://localhost:8080/api/v1/tasks \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <your-jwt-token>" \
-d '{
  "title": "New Task",
  "description": "Task description",
  "assignee_id": "12345"
}'
```

#### Get a Task

```sh
curl -X GET http://localhost:8080/api/v1/tasks \
-H "Authorization: Bearer <your-jwt-token>"
```
