# Task Management API - Clean Architecture

This is a refactored version of the Task Management API following Clean Architecture principles.

## Architecture Overview

The codebase is organized into distinct layers with clear separation of concerns:

```
task8/
├── domain/          # Core business entities and interfaces
├── usecase/         # Business logic and use cases
├── repository/      # Data access layer (implements domain interfaces)
├── infrastructure/  # External dependencies (MongoDB, JWT, Password hashing)
├── delivery/        # HTTP handlers and routing
└── main.go          # Application entry point
```

## Layer Dependencies

Dependencies flow inward following Clean Architecture principles:

- **Delivery** → **UseCase** → **Domain**
- **Repository** → **Domain** (implements domain interfaces)
- **Infrastructure** → **Domain** (implements domain interfaces)

## Layer Descriptions

### Domain Layer (`domain/`)

Core business entities and interfaces. This layer is independent of external frameworks and libraries.

- `entities.go`: Core business entities (Task, User, Request DTOs)
- `interfaces.go`: Repository and service interfaces

### Use Case Layer (`usecase/`)

Contains business logic and orchestrates interactions between layers.

- `task_usecase.go`: Task-related business logic
- `auth_usecase.go`: Authentication and authorization business logic

### Repository Layer (`repository/`)

Implements data access interfaces defined in the domain layer.

- `task_repository.go`: MongoDB implementation of TaskRepository
- `user_repository.go`: MongoDB implementation of UserRepository

### Infrastructure Layer (`infrastructure/`)

External dependencies and implementations of domain interfaces.

- `database.go`: MongoDB connection management
- `jwt.go`: JWT token generation and validation
- `password.go`: Password hashing using bcrypt

### Delivery Layer (`delivery/`)

HTTP handlers, middleware, and routing.

- `http/`: HTTP handlers (task_handler.go, auth_handler.go, dto.go)
- `middleware/`: Authentication and authorization middleware
- `router.go`: Route configuration

## Benefits of Clean Architecture

1. **Independence**: Business logic is independent of frameworks, databases, and external libraries
2. **Testability**: Easy to test business logic with mock repositories
3. **Flexibility**: Easy to swap implementations (e.g., change database, add new delivery mechanism)
4. **Maintainability**: Clear separation of concerns makes code easier to understand and maintain
5. **Scalability**: Easy to add new features without affecting existing code

## Running the Application

```bash
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DB="task_manager"
export JWT_SECRET="your-secret-key"

go run main.go
```

## API Endpoints

Same as previous version:
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login and get JWT token
- `GET /tasks` - Get all tasks (authenticated)
- `GET /tasks/:id` - Get task by ID (authenticated)
- `POST /tasks` - Create task (admin only)
- `PUT /tasks/:id` - Update task (admin only)
- `DELETE /tasks/:id` - Delete task (admin only)
- `POST /promote` - Promote user to admin (admin only)

