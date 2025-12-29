# InkStudio Backend - Go + MySQL

A RESTful API backend built with Go and MySQL, featuring user authentication with JWT tokens and bcrypt password hashing.

## Features

- ✅ User Registration with email validation
- ✅ User Login with JWT authentication
- ✅ Password hashing with bcrypt (includes automatic salt generation)
- ✅ MySQL database
- ✅ Docker and Docker Compose setup
- ✅ CORS middleware
- ✅ Logging middleware
- ✅ Protected routes with JWT validation
- ✅ Health check endpoint

## Tech Stack

- **Language**: Go 1.21
- **Database**: MySQL 8.0
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Containerization**: Docker & Docker Compose

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── config/
│   └── config.go             # Configuration management
├── internal/
│   ├── auth/
│   │   ├── jwt.go           # JWT token generation and validation
│   │   └── password.go       # Password hashing with bcrypt
│   ├── database/
│   │   ├── db.go            # Database connection
│   │   └── user.go          # User database operations
│   ├── handlers/
│   │   ├── auth.go          # Authentication handlers
│   │   ├── handler.go       # Base handler
│   │   └── health.go        # Health check handler
│   ├── middleware/
│   │   ├── auth.go          # JWT authentication middleware
│   │   ├── cors.go          # CORS middleware
│   │   └── logging.go       # Request logging middleware
│   └── models/
│       └── user.go          # User models and DTOs
├── pkg/
│   └── response/
│       └── response.go       # HTTP response helpers
├── docker-compose.yml        # Docker Compose configuration
├── Dockerfile               # Docker image definition
├── init.sql                 # Database initialization script
├── go.mod                   # Go module dependencies
├── go.sum                   # Go module checksums
├── Makefile                 # Build automation
└── README.md
```

## Getting Started

### Prerequisites

- Docker and Docker Compose installed
- (Optional) Go 1.21+ for local development

### Running with Docker (Recommended)

1. **Start the services**:
   ```bash
   cd backend
   make docker-up
   ```
   Or manually:
   ```bash
   docker-compose up -d
   ```

2. **Check the logs**:
   ```bash
   make docker-logs
   ```
   Or:
   ```bash
   docker-compose logs -f
   ```

3. **Stop the services**:
   ```bash
   make docker-down
   ```

The backend will be available at `http://localhost:8080` and MySQL at `localhost:3306`.

### Local Development (Without Docker)

1. **Start MySQL** (you'll need MySQL running locally):
   ```bash
   mysql -u root -p
   CREATE DATABASE inkstudio;
   ```

2. **Set environment variables**:
   ```bash
   export DB_HOST=localhost
   export DB_PORT=3306
   export DB_USER=inkstudio
   export DB_PASSWORD=inkstudio_password
   export DB_NAME=inkstudio
   export JWT_SECRET=your-super-secret-jwt-key
   ```

3. **Run the application**:
   ```bash
   make run
   ```

## API Endpoints

### Health Check
- **GET** `/health`
  - Returns service health status
  - Response: `200 OK`
    ```json
    {
      "status": "healthy",
      "database": "connected"
    }
    ```

### Authentication

#### Register
- **POST** `/api/auth/register`
  - Register a new user
  - Request body:
    ```json
    {
      "email": "user@example.com",
      "username": "username",
      "password": "password123"
    }
    ```
  - Response: `201 Created`
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIs...",
      "user": {
        "id": 1,
        "email": "user@example.com",
        "username": "username",
        "created_at": "2025-12-29T10:00:00Z",
        "updated_at": "2025-12-29T10:00:00Z"
      },
      "message": "User registered successfully"
    }
    ```

#### Login
- **POST** `/api/auth/login`
  - Login with existing credentials
  - Request body:
    ```json
    {
      "email": "user@example.com",
      "password": "password123"
    }
    ```
  - Response: `200 OK`
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIs...",
      "user": {
        "id": 1,
        "email": "user@example.com",
        "username": "username",
        "created_at": "2025-12-29T10:00:00Z",
        "updated_at": "2025-12-29T10:00:00Z"
      },
      "message": "Login successful"
    }
    ```

#### Get Current User (Protected)
- **GET** `/api/auth/me`
  - Get current user information
  - Headers: `Authorization: Bearer <token>`
  - Response: `200 OK`
    ```json
    {
      "id": 1,
      "email": "user@example.com",
      "username": "username",
      "created_at": "2025-12-29T10:00:00Z",
      "updated_at": "2025-12-29T10:00:00Z"
    }
    ```

## Testing the API

### Using curl

**Register a new user**:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "password123"
  }'
```

**Login**:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Get current user** (replace TOKEN with the JWT from login):
```bash
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer TOKEN"
```

## Security Features

### Password Security
- **Bcrypt hashing**: All passwords are hashed using bcrypt with a cost factor of 12
- **Automatic salt generation**: Bcrypt automatically generates a unique salt for each password
- **Password validation**: Minimum password length of 6 characters

### JWT Authentication
- **Token expiration**: Tokens expire after 24 hours
- **Secure signing**: Tokens are signed with HS256 algorithm
- **Protected routes**: Sensitive endpoints require valid JWT tokens

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | MySQL host | `localhost` |
| `DB_PORT` | MySQL port | `3306` |
| `DB_USER` | MySQL user | `inkstudio` |
| `DB_PASSWORD` | MySQL password | `inkstudio_password` |
| `DB_NAME` | MySQL database name | `inkstudio` |
| `JWT_SECRET` | JWT signing secret | (change in production!) |
| `PORT` | Application port | `8080` |

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_username (username)
);
```

## Makefile Commands

```bash
make help           # Show available commands
make build          # Build the Go application
make run            # Run the application locally
make test           # Run tests
make clean          # Clean build artifacts
make docker-build   # Build Docker images
make docker-up      # Start Docker containers
make docker-down    # Stop Docker containers
make docker-logs    # View Docker logs
make docker-restart # Restart Docker containers
make docker-clean   # Stop containers and remove volumes
```

## Troubleshooting

### Can't connect to MySQL
- Wait a few seconds after `docker-compose up` for MySQL to initialize
- Check logs: `docker-compose logs mysql`
- Verify MySQL is healthy: `docker-compose ps`

### Port already in use
- Change ports in `docker-compose.yml`
- Or stop conflicting services

### Database connection errors
- Ensure MySQL container is running: `docker ps`
- Check environment variables are set correctly
- Verify database credentials

## Production Considerations

⚠️ **Before deploying to production**:

1. Change `JWT_SECRET` to a strong, random value
2. Use environment-specific configuration
3. Enable HTTPS/TLS
4. Implement rate limiting
5. Add input validation and sanitization
6. Set up proper logging and monitoring
7. Use secrets management (e.g., Docker secrets, HashiCorp Vault)
8. Restrict CORS origins to your frontend domain
9. Regular security updates and dependency scanning
10. Implement database backups

## License

MIT
