# 🚀 Quick Start Guide - InkStudio Backend

## Prerequisites Check

```bash
# Check Go installation
go version  # Should be 1.21+

# Check MySQL installation
mysql --version  # Should be 8.0+

# Or use Docker
docker --version
docker-compose --version
```

## Option 1: Docker (Easiest) 🐳

```bash
# Navigate to backend directory
cd backend

# Start MySQL and backend together
docker-compose up -d

# Check if running
docker-compose ps

# View logs
docker-compose logs -f

# Test the API
curl http://localhost:8080/api/health

# Stop services
docker-compose down
```

**That's it! Your backend is running on http://localhost:8080**

---

## Option 2: Local Development 💻

### Step 1: Install Dependencies

```bash
cd backend

# Download Go modules
go mod download
go mod tidy
```

### Step 2: Setup MySQL Database

**Option A: Using the setup script (recommended)**

```bash
chmod +x setup-mysql.sh
./setup-mysql.sh
```

**Option B: Manual setup**

```bash
# Login to MySQL
mysql -u root -p

# Create database
CREATE DATABASE inkstudio CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# Exit MySQL
exit
```

### Step 3: Configure Environment

```bash
# Copy example env file
cp .env.example .env

# Edit .env with your database credentials
nano .env  # or vim, code, etc.
```

**Example .env:**

```env
DATABASE_URL=root:yourpassword@tcp(localhost:3306)/inkstudio?parseTime=true&charset=utf8mb4
JWT_SECRET=your-super-secret-jwt-key-minimum-32-characters-long
PORT=8080
```

### Step 4: Run the Server

**Option A: Development mode (recommended)**

```bash
go run ./cmd/server
```

**Option B: Build and run**

```bash
go build -o bin/server ./cmd/server
./bin/server
```

**Option C: Using Makefile**

```bash
make run  # Build and run
# or
make dev  # Development mode with hot reload (requires air)
```

### Step 5: Verify Installation

```bash
# Check health endpoint
curl http://localhost:8080/api/health

# Expected response:
# {"status":"healthy","time":"2025-12-29T..."}
```

---

## 🧪 Test the Authentication

### 1. Register a New User

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "mypassword123"
  }'
```

**Expected Response (201 Created):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "test@example.com",
    "username": "testuser",
    "created_at": "2025-12-29T10:30:00Z",
    "updated_at": "2025-12-29T10:30:00Z"
  },
  "message": "User registered successfully"
}
```

### 2. Login with Existing User

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "mypassword123"
  }'
```

**Expected Response (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {...},
  "message": "Login successful"
}
```

### 3. Access Protected Route

```bash
# Save the token from login/register response
TOKEN="paste_your_token_here"

curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN"
```

**Expected Response (200 OK):**

```json
{
  "id": 1,
  "email": "test@example.com",
  "username": "testuser",
  "created_at": "2025-12-29T10:30:00Z",
  "updated_at": "2025-12-29T10:30:00Z"
}
```

---

## 📝 Verify Database

```bash
# Connect to MySQL
mysql -u root -p inkstudio

# View users table
DESCRIBE users;

# Check registered users (without showing password hashes)
SELECT id, email, username, created_at FROM users;

# See how password is stored (hash with embedded salt)
SELECT id, email, LEFT(password_hash, 30) AS hash_preview FROM users;

# Example output:
# +----+-------------------+--------------------------------+
# | id | email             | hash_preview                   |
# +----+-------------------+--------------------------------+
# |  1 | test@example.com  | $2a$12$N9qo8uLOickgx2ZMRZoMye |
# +----+-------------------+--------------------------------+
```

---

## 🎯 Available Endpoints

| Method | Endpoint        | Auth Required | Description            |
| ------ | --------------- | ------------- | ---------------------- |
| GET    | `/api/health`   | ❌ No         | Server health check    |
| POST   | `/api/register` | ❌ No         | Register new user      |
| POST   | `/api/login`    | ❌ No         | Login with credentials |
| GET    | `/api/profile`  | ✅ Yes        | Get user profile       |

---

## 🔧 Useful Commands

### Development

```bash
# Run with hot reload (requires air)
make dev

# Format code
make format

# Run tests
make test

# Build binary
make build

# Clean build artifacts
make clean
```

### Database

```bash
# Run MySQL setup script
./setup-mysql.sh

# Connect to database
mysql -u root -p inkstudio

# Backup database
mysqldump -u root -p inkstudio > backup.sql

# Restore database
mysql -u root -p inkstudio < backup.sql
```

### Docker

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f backend

# Restart backend only
docker-compose restart backend

# Remove all data (including database)
docker-compose down -v
```

---

## 🐛 Troubleshooting

### Problem: "Failed to connect to database"

**Solution:**

```bash
# Check if MySQL is running
sudo systemctl status mysql

# Start MySQL
sudo systemctl start mysql

# Or check Docker container
docker-compose ps
```

### Problem: "Access denied for user"

**Solution:**

```bash
# Update .env with correct credentials
DATABASE_URL=your_username:your_password@tcp(localhost:3306)/inkstudio?parseTime=true&charset=utf8mb4

# Or reset MySQL root password
sudo mysql
ALTER USER 'root'@'localhost' IDENTIFIED BY 'newpassword';
FLUSH PRIVILEGES;
```

### Problem: "Database 'inkstudio' doesn't exist"

**Solution:**

```bash
# Run setup script
./setup-mysql.sh

# Or create manually
mysql -u root -p -e "CREATE DATABASE inkstudio CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### Problem: "Port 8080 already in use"

**Solution:**

```bash
# Change port in .env
PORT=8081

# Or kill process using port 8080
lsof -ti:8080 | xargs kill -9
```

### Problem: "Package not found" errors

**Solution:**

```bash
# Install dependencies
go mod download
go mod tidy

# Clean module cache if needed
go clean -modcache
go mod download
```

---

## 📚 Documentation Files

- **README.md** - Main documentation
- **SECURITY.md** - Security implementation details
- **MIGRATION_SUMMARY.md** - PostgreSQL to MySQL migration details
- **AUTHENTICATION_FLOW.md** - Visual authentication flow diagrams
- **MYSQL_REFERENCE.md** - MySQL commands and queries
- **This file** - Quick start guide

---

## 🎓 Understanding the Security

Your passwords are secured using **bcrypt with automatic salting**:

1. **Registration:** Password → Bcrypt (generates unique salt) → Hash stored in DB
2. **Login:** Password → Bcrypt (uses stored salt) → Compare hashes → Grant access
3. **Security:** Each password gets a unique 22-character random salt embedded in the hash

**Example:**

```
User enters: "mypassword123"
Stored in DB: "$2a$12$N9qo8uLOickgx2ZMRZoMye.ijfuNVqFn0rgCs9b0S7QV1K6R5wTT6"
                      │   │  └─────────────┬──────────────┘ └────────┬────────┘
                      │   │              Salt (22 chars)       Hash (31 chars)
                      │   └─ Cost factor (2^12 = 4096 iterations)
                      └─ Algorithm (bcrypt version 2a)
```

Read **SECURITY.md** for detailed explanation!

---

## ✅ Next Steps

1. ✅ Test all endpoints with cURL or Postman
2. ✅ Read SECURITY.md to understand the security implementation
3. ✅ Integrate with your frontend (Next.js in /frontend directory)
4. ✅ Customize JWT expiration time if needed
5. ✅ Set up production environment variables
6. ✅ Add rate limiting for production
7. ✅ Consider adding 2FA for extra security

---

## 🆘 Need Help?

- Check **README.md** for detailed documentation
- Review **TROUBLESHOOTING** section above
- Check **MYSQL_REFERENCE.md** for database commands
- Review **SECURITY.md** for security questions

---

**Congratulations! Your secure authentication backend is ready! 🎉**
