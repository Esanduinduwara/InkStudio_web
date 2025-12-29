# InkStudio Backend - Authentication API

A secure Go backend for user authentication with registration, login, and JWT-based authorization.

## Features

- ✅ User registration with email and username
- ✅ Secure password hashing using bcrypt (with automatic salting)
- ✅ User login with email and password
- ✅ JWT token generation and validation
- ✅ Protected routes with authentication middleware
- ✅ PostgreSQL database integration
- ✅ CORS support for frontend integration
- ✅ Request logging middleware
- ✅ Input validation and error handling

## Security Features

### Password Security

- **Bcrypt Hashing**: Passwords are hashed using bcrypt with a cost factor of 12
- **Automatic Salting**: Bcrypt automatically generates a unique salt for each password
- **No Plain Text Storage**: Passwords are never stored in plain text
- **Hash Comparison**: Uses constant-time comparison to prevent timing attacks

### JWT Authentication

- **Token Expiration**: Tokens expire after 24 hours
- **HMAC Signing**: Tokens are signed using HMAC-SHA256
- **Claims Validation**: Validates expiration, issued at, and not before claims

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher

## Installation

1. **Clone the repository** (if not already done)

2. **Install dependencies**:

```bash
cd backend
go mod download
```

3. **Setup PostgreSQL Database**:

```bash
# Create database
createdb inkstudio

# Or using psql
psql -U postgres
CREATE DATABASE inkstudio;
```

4. **Configure environment variables**:

```bash
cp .env.example .env
# Edit .env with your configuration
```

## Running the Server

```bash
# Development mode
cd cmd/server
go run .

# Or from backend root
go run ./cmd/server

# Build and run
go build -o bin/server ./cmd/server
./bin/server
```

The server will start on `http://localhost:8080`

## API Endpoints

### 1. Health Check

```bash
GET /api/health
```

**Response**:

```json
{
  "status": "healthy",
  "time": "2025-12-10T12:00:00Z"
}
```

### 2. Register New User

```bash
POST /api/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "securepassword123"
}
```

**Response** (201 Created):

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "johndoe",
    "created_at": "2025-12-10T12:00:00Z",
    "updated_at": "2025-12-10T12:00:00Z"
  },
  "message": "User registered successfully"
}
```

### 3. Login

```bash
POST /api/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response** (200 OK):

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "johndoe",
    "created_at": "2025-12-10T12:00:00Z",
    "updated_at": "2025-12-10T12:00:00Z"
  },
  "message": "Login successful"
}
```

### 4. Get User Profile (Protected)

```bash
GET /api/profile
Authorization: Bearer <your_jwt_token>
```

**Response** (200 OK):

```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "johndoe",
  "created_at": "2025-12-10T12:00:00Z",
  "updated_at": "2025-12-10T12:00:00Z"
}
```

## Testing with cURL

### Register a new user:

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "password123"
  }'
```

### Login:

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Access protected route:

```bash
# Save token from login response
TOKEN="your_jwt_token_here"

curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN"
```

## Database Schema

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
```

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── auth/
│   │   ├── jwt.go              # JWT token generation and validation
│   │   └── password.go         # Password hashing with bcrypt (includes salt)
│   ├── database/
│   │   ├── db.go               # Database connection and initialization
│   │   └── user.go             # User database operations
│   ├── handlers/
│   │   ├── handler.go          # Handler struct and constructor
│   │   ├── auth.go             # Registration and login handlers
│   │   └── health.go           # Health check handler
│   ├── middleware/
│   │   ├── auth.go             # JWT authentication middleware
│   │   ├── cors.go             # CORS middleware
│   │   └── logging.go          # Request logging middleware
│   └── models/
│       └── user.go             # User models and request/response structs
├── pkg/
│   └── response/
│       └── response.go         # JSON response helpers
├── config/
│   └── config.go               # Configuration management
├── go.mod                      # Go module dependencies
├── .env.example                # Example environment configuration
├── .gitignore                  # Git ignore rules
└── README.md                   # This file
```

## Environment Variables

| Variable       | Description                  | Default                                                                 |
| -------------- | ---------------------------- | ----------------------------------------------------------------------- |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/inkstudio?sslmode=disable` |
| `JWT_SECRET`   | Secret key for JWT signing   | `your-secret-key-change-in-production`                                  |
| `PORT`         | Server port                  | `8080`                                                                  |

## Error Handling

The API returns consistent error responses:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:

- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid credentials or token)
- `404` - Not Found
- `409` - Conflict (duplicate email/username)
- `500` - Internal Server Error

## Production Deployment

1. **Set strong JWT secret**:

   - Minimum 32 characters
   - Use random string generator

2. **Configure CORS**:

   - Change `Access-Control-Allow-Origin` from `*` to your frontend domain

3. **Use environment variables**:

   - Never commit `.env` file
   - Use secrets manager in production

4. **Enable SSL/TLS**:

   - Use HTTPS in production
   - Set `sslmode=require` in DATABASE_URL

5. **Database security**:
   - Use strong database password
   - Enable connection pooling
   - Regular backups

## License

MIT License
