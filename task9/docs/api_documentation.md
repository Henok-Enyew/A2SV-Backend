# Task Management REST API Documentation

## Overview

This is a RESTful API for managing tasks built with Go and Gin Framework with MongoDB as the persistent data storage. The API provides endpoints for creating, reading, updating, and deleting tasks with JWT-based authentication and role-based authorization. All data is persisted in MongoDB, ensuring data availability across API restarts.

**Base URL**: `http://localhost:8080`

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Most endpoints require a valid JWT token in the Authorization header.

### Authentication Flow

1. Register a new user account using `POST /auth/register`
2. Login with credentials using `POST /auth/login` to receive a JWT token
3. Include the token in subsequent requests: `Authorization: Bearer <token>`

### User Roles

- **admin**: Can perform all operations including creating, updating, and deleting tasks, and promoting users
- **user**: Can only view tasks (GET operations)

**Note**: The first user registered in the system automatically becomes an admin.

## MongoDB Integration

This API uses MongoDB for persistent data storage. Tasks are stored in a MongoDB collection with automatic ObjectID generation for unique identifiers.

### Prerequisites

- MongoDB instance running (local or cloud)
- Go 1.19 or higher
- MongoDB Go Driver (automatically installed via `go mod tidy`)

### MongoDB Setup

#### Local MongoDB Setup

1. Install MongoDB Community Edition from [MongoDB Download Center](https://www.mongodb.com/try/download/community)
2. Start MongoDB service:
   ```bash
   # On Linux/Mac
   mongod
   
   # On Windows (as Administrator)
   net start MongoDB
   ```
3. MongoDB will run on `mongodb://localhost:27017` by default

#### MongoDB Atlas (Cloud) Setup

1. Create a free account at [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)
2. Create a new cluster
3. Get your connection string from the Atlas dashboard
4. Set the `MONGODB_URI` environment variable with your connection string

### Configuration

The API can be configured using environment variables:

- `MONGODB_URI`: MongoDB connection URI (default: `mongodb://localhost:27017`)
- `MONGODB_DB`: Database name (default: `task_manager`)
- `JWT_SECRET`: Secret key for JWT token signing (default: `your-secret-key-change-in-production`)

Example:
```bash
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DB="task_manager"
export JWT_SECRET="your-secure-secret-key-here"
```

Or for MongoDB Atlas:
```bash
export MONGODB_URI="mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority"
export MONGODB_DB="task_manager"
export JWT_SECRET="your-secure-secret-key-here"
```

## Authentication Endpoints

### 1. Register User

Create a new user account.

**Endpoint**: `POST /auth/register`

**Request Body**:
```json
{
  "username": "john_doe",
  "password": "securepassword123"
}
```

**Fields**:
- `username` (required): Unique username (string)
- `password` (required): Password with minimum 6 characters (string)

**Response**:
```json
{
  "status": "success",
  "message": "user registered successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "username": "john_doe",
    "role": "admin"
  }
}
```

**Status Codes**:
- `201 Created`: User registered successfully
- `400 Bad Request`: Invalid request body or validation error
- `409 Conflict`: Username already exists
- `500 Internal Server Error`: Server error

**Note**: The first user registered automatically becomes an admin.

---

### 2. Login

Authenticate user and receive JWT token.

**Endpoint**: `POST /auth/login`

**Request Body**:
```json
{
  "username": "john_doe",
  "password": "securepassword123"
}
```

**Response**:
```json
{
  "status": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "507f1f77bcf86cd799439011",
      "username": "john_doe",
      "role": "admin"
    }
  }
}
```

**Status Codes**:
- `200 OK`: Login successful
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Invalid credentials
- `500 Internal Server Error`: Server error

**Usage**: Include the token in subsequent requests:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## Task Endpoints

### 3. Get All Tasks

Retrieve a list of all tasks.

**Endpoint**: `GET /tasks`

**Authentication**: Required (Bearer token)

**Authorization**: All authenticated users (admin and user roles)

**Request**: No request body required

**Headers**:
```
Authorization: Bearer <your-jwt-token>
```

**Response**:
```json
{
  "status": "success",
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "title": "Complete project",
      "description": "Finish the task management API",
      "due_date": "2024-12-31T00:00:00Z",
      "status": "pending",
      "created_at": "2024-01-01T10:00:00Z",
      "updated_at": "2024-01-01T10:00:00Z"
    }
  ],
  "count": 1
}
```

**Status Codes**:
- `200 OK`: Successfully retrieved tasks
- `401 Unauthorized`: Missing or invalid token
- `500 Internal Server Error`: Database error occurred

---

### 4. Get Task by ID

Retrieve details of a specific task.

**Endpoint**: `GET /tasks/:id`

**Authentication**: Required (Bearer token)

**Authorization**: All authenticated users (admin and user roles)

**Parameters**:
- `id` (path parameter): Task ID (MongoDB ObjectID as string, e.g., "507f1f77bcf86cd799439011")

**Request**: No request body required

**Headers**:
```
Authorization: Bearer <your-jwt-token>
```

**Response**:
```json
{
  "status": "success",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "title": "Complete project",
    "description": "Finish the task management API",
    "due_date": "2024-12-31T00:00:00Z",
    "status": "pending",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

**Status Codes**:
- `200 OK`: Task found
- `400 Bad Request`: Invalid task ID format (not a valid ObjectID)
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: Task not found
- `500 Internal Server Error`: Database error occurred

**Error Response**:
```json
{
  "status": "error",
  "message": "task not found"
}
```

---

### 5. Create Task

Create a new task.

**Endpoint**: `POST /tasks`

**Authentication**: Required (Bearer token)

**Authorization**: Admin only

**Headers**:
```
Authorization: Bearer <your-jwt-token>
```

**Request Body**:
```json
{
  "title": "Complete project",
  "description": "Finish the task management API",
  "due_date": "2024-12-31T00:00:00Z",
  "status": "pending"
}
```

**Fields**:
- `title` (required): Task title (string)
- `description` (optional): Task description (string)
- `due_date` (required): Due date in ISO 8601 format (string)
- `status` (optional): Task status - "pending", "in_progress", or "completed" (default: "pending")

**Response**:
```json
{
  "status": "success",
  "message": "task created successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "title": "Complete project",
    "description": "Finish the task management API",
    "due_date": "2024-12-31T00:00:00Z",
    "status": "pending",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
}
```

**Status Codes**:
- `201 Created`: Task created successfully
- `400 Bad Request`: Invalid request body or validation error
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Admin access required
- `500 Internal Server Error`: Database error occurred

**Error Response**:
```json
{
  "status": "error",
  "message": "invalid request body",
  "error": "Key: 'CreateTaskRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag"
}
```

---

### 6. Update Task

Update an existing task.

**Endpoint**: `PUT /tasks/:id`

**Authentication**: Required (Bearer token)

**Authorization**: Admin only

**Parameters**:
- `id` (path parameter): Task ID (MongoDB ObjectID as string)

**Headers**:
```
Authorization: Bearer <your-jwt-token>
```

**Request Body**:
```json
{
  "title": "Updated task title",
  "description": "Updated description",
  "due_date": "2024-12-31T00:00:00Z",
  "status": "in_progress"
}
```

**Fields** (all optional):
- `title`: Task title (string)
- `description`: Task description (string)
- `due_date`: Due date in ISO 8601 format (string)
- `status`: Task status - "pending", "in_progress", or "completed" (string)

**Response**:
```json
{
  "status": "success",
  "message": "task updated successfully",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "title": "Updated task title",
    "description": "Updated description",
    "due_date": "2024-12-31T00:00:00Z",
    "status": "in_progress",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T11:00:00Z"
  }
}
```

**Status Codes**:
- `200 OK`: Task updated successfully
- `400 Bad Request`: Invalid request body or task ID format
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Admin access required
- `404 Not Found`: Task not found
- `500 Internal Server Error`: Database error occurred

**Error Response**:
```json
{
  "status": "error",
  "message": "task not found"
}
```

---

### 7. Delete Task

Delete a task.

**Endpoint**: `DELETE /tasks/:id`

**Authentication**: Required (Bearer token)

**Authorization**: Admin only

**Parameters**:
- `id` (path parameter): Task ID (MongoDB ObjectID as string)

**Request**: No request body required

**Headers**:
```
Authorization: Bearer <your-jwt-token>
```

**Response**:
```json
{
  "status": "success",
  "message": "task deleted successfully"
}
```

**Status Codes**:
- `200 OK`: Task deleted successfully
- `400 Bad Request`: Invalid task ID format (not a valid ObjectID)
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Admin access required
- `404 Not Found`: Task not found
- `500 Internal Server Error`: Database error occurred
```

---

### 8. Promote User to Admin

Promote a user to admin role.

**Endpoint**: `POST /promote`

**Authentication**: Required (Bearer token)

**Authorization**: Admin only

**Request Body**:
```json
{
  "username": "jane_doe"
}
```

**Headers**:
```
Authorization: Bearer <your-jwt-token>
```

**Response**:
```json
{
  "status": "success",
  "message": "user promoted to admin successfully"
}
```

**Status Codes**:
- `200 OK`: User promoted successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Admin access required
- `404 Not Found`: User not found
- `500 Internal Server Error`: Database error occurred

**Error Response**:
```json
{
  "status": "error",
  "message": "task not found"
}
```

---

## Access Control Summary

| Endpoint | Method | Authentication | Authorization |
|----------|--------|----------------|---------------|
| `/auth/register` | POST | Not required | Public |
| `/auth/login` | POST | Not required | Public |
| `/tasks` | GET | Required | All users |
| `/tasks/:id` | GET | Required | All users |
| `/tasks` | POST | Required | Admin only |
| `/tasks/:id` | PUT | Required | Admin only |
| `/tasks/:id` | DELETE | Required | Admin only |
| `/promote` | POST | Required | Admin only |

## Task Status Values

Valid status values:
- `pending`: Task is pending
- `in_progress`: Task is in progress
- `completed`: Task is completed

## Date Format

All dates should be in ISO 8601 format: `YYYY-MM-DDTHH:MM:SSZ`

Example: `2024-12-31T23:59:59Z`

## Error Handling

All error responses follow this format:
```json
{
  "status": "error",
  "message": "error description"
}
```

Additional error details may be included in the `error` field for validation errors.

### Common Error Responses

- **401 Unauthorized**: Missing or invalid JWT token
- **403 Forbidden**: Insufficient permissions (admin access required)
- **404 Not Found**: Resource not found
- **409 Conflict**: Username already exists

## Running the API

### Prerequisites

1. Ensure MongoDB is running (local or cloud)
2. Install Go dependencies:
```bash
go mod tidy
```

### Running the Server

1. Set environment variables (optional, defaults provided):
```bash
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DB="task_manager"
```

2. Run the server:
```bash
go run main.go
```

The server will:
- Connect to MongoDB
- Start on `http://localhost:8080`
- Display connection status

### Verifying MongoDB Connection

The API will attempt to connect to MongoDB on startup. If connection fails, the application will exit with an error message. Ensure MongoDB is running before starting the API.

### Verifying Data in MongoDB

You can verify data stored in MongoDB using:

1. **MongoDB Compass** (GUI):
   - Connect to your MongoDB instance
   - Navigate to `task_manager` database
   - View the `tasks` collection

2. **MongoDB Shell**:
```bash
mongosh
use task_manager
db.tasks.find().pretty()
```

3. **Direct Query**:
```javascript
db.tasks.find({ status: "pending" })
db.tasks.findOne({ _id: ObjectId("507f1f77bcf86cd799439011") })
```

## Testing with Postman

### Collection Setup

1. Create a new Postman collection named "Task Management API"
2. Set the base URL variable: `{{base_url}}` = `http://localhost:8080`

### Example Requests

#### Register User
- Method: POST
- URL: `{{base_url}}/auth/register`
- Headers: `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "username": "john_doe",
  "password": "securepassword123"
}
```

#### Login
- Method: POST
- URL: `{{base_url}}/auth/login`
- Headers: `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "username": "john_doe",
  "password": "securepassword123"
}
```
- Save the token from the response for subsequent requests

#### Get All Tasks
- Method: GET
- URL: `{{base_url}}/tasks`
- Headers: 
  - `Authorization: Bearer <your-jwt-token>`

#### Get Task by ID
- Method: GET
- URL: `{{base_url}}/tasks/507f1f77bcf86cd799439011`
- Headers: 
  - `Authorization: Bearer <your-jwt-token>`

#### Create Task (Admin Only)
- Method: POST
- URL: `{{base_url}}/tasks`
- Headers: 
  - `Content-Type: application/json`
  - `Authorization: Bearer <your-admin-jwt-token>`
- Body (raw JSON):
```json
{
  "title": "Learn Go",
  "description": "Complete Go programming course",
  "due_date": "2024-12-31T00:00:00Z",
  "status": "pending"
}
```

#### Update Task (Admin Only)
- Method: PUT
- URL: `{{base_url}}/tasks/507f1f77bcf86cd799439011`
- Headers: 
  - `Content-Type: application/json`
  - `Authorization: Bearer <your-admin-jwt-token>`
- Body (raw JSON):
```json
{
  "status": "in_progress"
}
```

#### Delete Task (Admin Only)
- Method: DELETE
- URL: `{{base_url}}/tasks/507f1f77bcf86cd799439011`
- Headers: 
  - `Authorization: Bearer <your-admin-jwt-token>`

#### Promote User (Admin Only)
- Method: POST
- URL: `{{base_url}}/promote`
- Headers: 
  - `Content-Type: application/json`
  - `Authorization: Bearer <your-admin-jwt-token>`
- Body (raw JSON):
```json
{
  "username": "jane_doe"
}
```

## Notes

- **Persistent Storage**: The API uses MongoDB for persistent data storage. All data persists across server restarts.
- **Task IDs**: Task IDs are MongoDB ObjectIDs (24-character hexadecimal strings) automatically generated by MongoDB.
- **Timestamps**: The `created_at` and `updated_at` fields are automatically managed by the system.
- **Data Persistence**: All tasks are stored in the `tasks` collection in the configured MongoDB database.
- **Backward Compatibility**: The API maintains the same endpoint structure and response format as the in-memory version, with the only change being the ID format (ObjectID string instead of integer).
- **Error Handling**: All MongoDB operations include proper error handling for network errors, database errors, and validation errors.
- **All endpoints return JSON responses**.

## Security Features

### Password Hashing

- Passwords are hashed using bcrypt before storage
- Original passwords are never stored in the database
- Password comparison is done using bcrypt's secure comparison

### JWT Tokens

- Tokens are signed using HS256 algorithm
- Token expiration: 24 hours
- Tokens contain user ID, username, and role
- Secret key configurable via `JWT_SECRET` environment variable

### Authorization

- Role-based access control (RBAC) implemented
- Middleware validates JWT tokens on protected routes
- Admin middleware enforces admin-only access
- First registered user automatically becomes admin

## MongoDB Integration Details

### Database Structure

- **Database**: `task_manager` (configurable via `MONGODB_DB` environment variable)
- **Collections**: 
  - `tasks`: Stores task documents
  - `users`: Stores user documents

#### Tasks Collection
Each task is stored as a document with the following fields:
  - `_id`: MongoDB ObjectID (primary key)
  - `title`: String
  - `description`: String
  - `due_date`: ISODate
  - `status`: String (pending, in_progress, completed)
  - `created_at`: ISODate
  - `updated_at`: ISODate

#### Users Collection
Each user is stored as a document with the following fields:
  - `_id`: MongoDB ObjectID (primary key)
  - `username`: String (unique)
  - `password`: String (bcrypt hashed)
  - `role`: String (admin or user)

### Connection Management

- The API establishes a connection to MongoDB on startup
- Connection is maintained throughout the application lifecycle
- Proper connection cleanup on application shutdown
- Connection timeout: 10 seconds

### Error Handling

The API handles various MongoDB-related errors:
- **Connection Errors**: Fails fast on startup if MongoDB is unavailable
- **Query Errors**: Returns appropriate HTTP status codes (404 for not found, 500 for server errors)
- **Validation Errors**: Returns 400 for invalid ObjectID formats
- **Network Errors**: Proper timeout handling for all database operations

