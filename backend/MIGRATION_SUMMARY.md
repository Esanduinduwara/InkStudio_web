# InkStudio Backend - Migration Summary

## Changes Made: PostgreSQL → MySQL

### ✅ What Was Changed

#### 1. Database Driver

- **Before:** `github.com/lib/pq` (PostgreSQL)
- **After:** `github.com/go-sql-driver/mysql` (MySQL)
- **File:** `go.mod`

#### 2. Database Connection

- **Before:** `sql.Open("postgres", databaseURL)`
- **After:** `sql.Open("mysql", databaseURL)`
- **File:** `internal/database/db.go`

#### 3. SQL Syntax for Table Creation

**Before (PostgreSQL):**

```sql
id SERIAL PRIMARY KEY
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

**After (MySQL):**

```sql
id INT AUTO_INCREMENT PRIMARY KEY
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
```

- **File:** `internal/database/db.go`

#### 4. Query Parameter Placeholders

**Before (PostgreSQL):**

```sql
SELECT * FROM users WHERE email = $1 AND id = $2
```

**After (MySQL):**

```sql
SELECT * FROM users WHERE email = ? AND id = ?
```

- **Files:**
  - `internal/database/user.go`
  - All database query functions

#### 5. INSERT with RETURNING → LastInsertId()

**Before (PostgreSQL):**

```go
db.QueryRow("INSERT ... RETURNING id, email, ...").Scan(&user.ID, ...)
```

**After (MySQL):**

```go
result, err := db.Exec("INSERT ...")
userID, err := result.LastInsertId()
// Then fetch the created user
db.QueryRow("SELECT ... WHERE id = ?", userID).Scan(...)
```

- **File:** `internal/database/user.go`

#### 6. Connection String Format

**Before:**

```
postgres://user:pass@localhost:5432/dbname?sslmode=disable
```

**After:**

```
user:pass@tcp(localhost:3306)/dbname?parseTime=true&charset=utf8mb4
```

- **Files:**
  - `config/config.go`
  - `.env.example`

#### 7. Docker Configuration

- Changed from PostgreSQL to MySQL 8.0
- Updated health checks
- Changed port from 5432 to 3306
- **File:** `docker-compose.yml`

#### 8. Error Detection

**Before:**

```go
if strings.Contains(err.Error(), "duplicate key")
```

**After:**

```go
if strings.Contains(err.Error(), "Duplicate entry") ||
   strings.Contains(err.Error(), "duplicate key")
```

- **File:** `internal/handlers/auth.go`

### 🔒 Security Features (Already Implemented)

The password hashing and salting were **already properly implemented** using bcrypt:

#### Bcrypt Implementation

```go
// internal/auth/password.go
func HashPassword(password string) (string, error) {
    // Cost factor 12 = 4,096 iterations
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword, password string) error {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err
}
```

#### How Bcrypt Handles Salt:

1. **Salt Generation:** Bcrypt automatically generates a unique 22-character random salt for each password
2. **Salt Storage:** The salt is embedded in the hash output (no separate column needed)
3. **Hash Format:** `$2a$12$[22-char-salt][31-char-hash]`
4. **Verification:** Bcrypt extracts the salt from the hash when verifying

**Example Hash:**

```
$2a$12$N9qo8uLOickgx2ZMRZoMye.ijfuNVqFn0rgCs9b0S7QV1K6R5wTT6
 │   │  └─────────────┬──────────────┘ └────────┬────────┘
 │   │              Salt                      Hash
 │   └─ Cost (2^12 iterations)
 └─ Algorithm version
```

### 📁 New Files Created

1. **SECURITY.md** - Comprehensive security documentation

   - Explains salt and hashing in detail
   - Shows how bcrypt works
   - Security best practices

2. **MYSQL_REFERENCE.md** - MySQL command reference

   - Common SQL queries
   - Backup/restore commands
   - Troubleshooting guide

3. **setup-mysql.sh** - MySQL database setup script
   - Automated database creation
   - Table initialization
   - Connection string generation

### 📋 File Structure (Unchanged)

```
backend/
├── cmd/server/main.go          # Application entry
├── config/config.go            # Configuration (updated for MySQL)
├── internal/
│   ├── auth/
│   │   ├── jwt.go             # JWT tokens (unchanged)
│   │   └── password.go        # Bcrypt hashing (unchanged)
│   ├── database/
│   │   ├── db.go              # MySQL connection (updated)
│   │   └── user.go            # User queries (updated)
│   ├── handlers/
│   │   ├── handler.go         # Handler struct (unchanged)
│   │   ├── auth.go            # Auth handlers (minor update)
│   │   └── health.go          # Health check (unchanged)
│   ├── middleware/
│   │   ├── auth.go            # JWT middleware (unchanged)
│   │   ├── cors.go            # CORS (unchanged)
│   │   └── logging.go         # Logging (unchanged)
│   └── models/
│       └── user.go            # User models (unchanged)
├── pkg/response/
│   └── response.go            # JSON helpers (unchanged)
├── go.mod                     # Dependencies (updated)
├── .env.example               # Config template (updated)
├── docker-compose.yml         # Docker setup (updated)
├── README.md                  # Documentation (updated)
├── SECURITY.md                # New: Security docs
├── MYSQL_REFERENCE.md         # New: MySQL reference
└── setup-mysql.sh             # New: Setup script
```

### 🚀 How to Use

#### Option 1: Quick Start with Docker

```bash
cd backend
docker-compose up -d
```

#### Option 2: Local Development

```bash
cd backend

# 1. Setup MySQL database
./setup-mysql.sh

# 2. Install dependencies
go mod download

# 3. Configure .env
cp .env.example .env
# Edit .env with your MySQL credentials

# 4. Run the server
go run ./cmd/server
```

### 🔑 Connection String Examples

**Local Development:**

```bash
DATABASE_URL=root:password@tcp(localhost:3306)/inkstudio?parseTime=true&charset=utf8mb4
```

**Docker:**

```bash
DATABASE_URL=inkstudio_user:inkstudio_pass@tcp(mysql:3306)/inkstudio?parseTime=true&charset=utf8mb4
```

**Remote Server:**

```bash
DATABASE_URL=user:pass@tcp(192.168.1.100:3306)/inkstudio?parseTime=true&charset=utf8mb4
```

### ✨ What Stayed the Same

- ✅ Password hashing with bcrypt (cost factor 12)
- ✅ Automatic salt generation and embedding
- ✅ JWT token generation and validation
- ✅ All API endpoints and response formats
- ✅ Middleware (CORS, logging, authentication)
- ✅ Project structure and organization
- ✅ Error handling and validation
- ✅ Security best practices

### 📚 Documentation Files

- **README.md** - Getting started, API documentation
- **SECURITY.md** - Detailed security explanation with examples
- **MYSQL_REFERENCE.md** - MySQL commands and troubleshooting
- **This file** - Summary of changes

### 🎯 Next Steps

1. Run `go mod tidy` to install MySQL driver
2. Setup MySQL database (use `./setup-mysql.sh`)
3. Update `.env` with your database credentials
4. Test the endpoints using cURL or Postman
5. Review SECURITY.md to understand the security implementation

### 🔍 Testing the Authentication

```bash
# 1. Register a user
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","username":"testuser","password":"securepass123"}'

# 2. Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"securepass123"}'

# 3. Access protected route (use token from login response)
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

### ❓ FAQ

**Q: Do I need to change anything in the auth logic?**
A: No! The bcrypt implementation already handles salt and hashing correctly.

**Q: Where is the salt stored?**
A: The salt is embedded in the password_hash field automatically by bcrypt.

**Q: Is this secure?**
A: Yes! Bcrypt with cost factor 12 is industry-standard and recommended by OWASP.

**Q: Can I use PostgreSQL instead?**
A: Yes, the original PostgreSQL code is in the git history. You can revert the changes.

**Q: What if I need to increase security?**
A: Increase the cost factor in `internal/auth/password.go` from 12 to 13 or 14.

---

**Summary:** Successfully migrated from PostgreSQL to MySQL while maintaining all security features including proper password hashing with bcrypt (automatic salting).
