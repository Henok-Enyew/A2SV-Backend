# Task Management REST API Documentation

## Overview

This is a RESTful API for managing tasks built with Go and Gin Framework with MongoDB as the persistent data storage. The API provides endpoints for creating, reading, updating, and deleting tasks. All data is persisted in MongoDB, ensuring data availability across API restarts.

**Base URL**: `http://localhost:8080`

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

Example:
```bash
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DB="task_manager"
```

Or for MongoDB Atlas:
```bash
export MONGODB_URI="mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority"
export MONGODB_DB="task_manager"
```

## Endpoints

### 1. Get All Tasks

Retrieve a list of all tasks.

**Endpoint**: `GET /tasks`

**Request**: No request body required

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
- `500 Internal Server Error`: Database error occurred

---

### 2. Get Task by ID

Retrieve details of a specific task.

**Endpoint**: `GET /tasks/:id`

**Parameters**:
- `id` (path parameter): Task ID (MongoDB ObjectID as string, e.g., "507f1f77bcf86cd799439011")

**Request**: No request body required

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

### 3. Create Task

Create a new task.

**Endpoint**: `POST /tasks`

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

### 4. Update Task

Update an existing task.

**Endpoint**: `PUT /tasks/:id`

**Parameters**:
- `id` (path parameter): Task ID (MongoDB ObjectID as string)

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

### 5. Delete Task

Delete a task.

**Endpoint**: `DELETE /tasks/:id`

**Parameters**:
- `id` (path parameter): Task ID (MongoDB ObjectID as string)

**Request**: No request body required

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

#### Create Task
- Method: POST
- URL: `{{base_url}}/tasks`
- Headers: `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "title": "Learn Go",
  "description": "Complete Go programming course",
  "due_date": "2024-12-31T00:00:00Z",
  "status": "pending"
}
```

#### Get All Tasks
- Method: GET
- URL: `{{base_url}}/tasks`

#### Get Task by ID
- Method: GET
- URL: `{{base_url}}/tasks/507f1f77bcf86cd799439011`
  - Note: Use the actual ObjectID string returned when creating a task

#### Update Task
- Method: PUT
- URL: `{{base_url}}/tasks/507f1f77bcf86cd799439011`
  - Note: Use the actual ObjectID string returned when creating a task
- Headers: `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "status": "in_progress"
}
```

#### Delete Task
- Method: DELETE
- URL: `{{base_url}}/tasks/507f1f77bcf86cd799439011`
  - Note: Use the actual ObjectID string returned when creating a task

## Notes

- **Persistent Storage**: The API uses MongoDB for persistent data storage. All data persists across server restarts.
- **Task IDs**: Task IDs are MongoDB ObjectIDs (24-character hexadecimal strings) automatically generated by MongoDB.
- **Timestamps**: The `created_at` and `updated_at` fields are automatically managed by the system.
- **Data Persistence**: All tasks are stored in the `tasks` collection in the configured MongoDB database.
- **Backward Compatibility**: The API maintains the same endpoint structure and response format as the in-memory version, with the only change being the ID format (ObjectID string instead of integer).
- **Error Handling**: All MongoDB operations include proper error handling for network errors, database errors, and validation errors.
- **All endpoints return JSON responses**.

## MongoDB Integration Details

### Database Structure

- **Database**: `task_manager` (configurable via `MONGODB_DB` environment variable)
- **Collection**: `tasks`
- **Document Structure**: Each task is stored as a document with the following fields:
  - `_id`: MongoDB ObjectID (primary key)
  - `title`: String
  - `description`: String
  - `due_date`: ISODate
  - `status`: String (pending, in_progress, completed)
  - `created_at`: ISODate
  - `updated_at`: ISODate

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

