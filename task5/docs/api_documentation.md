# Task Management REST API Documentation

## Overview

This is a RESTful API for managing tasks built with Go and Gin Framework. The API provides endpoints for creating, reading, updating, and deleting tasks.

**Base URL**: `http://localhost:8080/api/v1`

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
      "id": 1,
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

---

### 2. Get Task by ID

Retrieve details of a specific task.

**Endpoint**: `GET /tasks/:id`

**Parameters**:
- `id` (path parameter): Task ID (integer)

**Request**: No request body required

**Response**:
```json
{
  "status": "success",
  "data": {
    "id": 1,
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
- `400 Bad Request`: Invalid task ID format
- `404 Not Found`: Task not found

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
    "id": 1,
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
- `id` (path parameter): Task ID (integer)

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
    "id": 1,
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
- `400 Bad Request`: Invalid request body or task ID
- `404 Not Found`: Task not found

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
- `id` (path parameter): Task ID (integer)

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
- `400 Bad Request`: Invalid task ID format
- `404 Not Found`: Task not found

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

1. Install dependencies:
```bash
go mod tidy
```

2. Run the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Testing with Postman

### Collection Setup

1. Create a new Postman collection named "Task Management API"
2. Set the base URL variable: `{{base_url}}` = `http://localhost:8080/api/v1`

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
- URL: `{{base_url}}/tasks/1`

#### Update Task
- Method: PUT
- URL: `{{base_url}}/tasks/1`
- Headers: `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "status": "in_progress"
}
```

#### Delete Task
- Method: DELETE
- URL: `{{base_url}}/tasks/1`

## Notes

- The API uses in-memory storage. All data will be lost when the server restarts.
- Task IDs are auto-incremented starting from 1.
- The `created_at` and `updated_at` fields are automatically managed by the system.
- All endpoints return JSON responses.

