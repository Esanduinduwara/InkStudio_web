# InkStudio Backend - Security Documentation

## Password Security Implementation

This backend implements industry-standard security practices for user authentication.

### 🔐 Password Hashing with Bcrypt

**Bcrypt** is used for password hashing, which automatically includes several security features:

#### 1. **Automatic Salt Generation**

- Each password gets a **unique, random salt** automatically
- Salt is generated using cryptographically secure random number generator
- The salt is stored **within the hash itself** (no separate storage needed)
- Format: `$2a$[cost]$[22-character salt][31-character hash]`

#### 2. **Adaptive Hashing (Cost Factor: 12)**

```go
// From internal/auth/password.go
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    // Cost factor 12 = 2^12 = 4,096 iterations
    // This provides strong security while maintaining reasonable performance
}
```

**Why Cost Factor 12?**

- Cost factor determines how computationally expensive hashing is
- Higher cost = more secure but slower
- Cost 12 is the recommended minimum for 2024+
- Can be increased in the future as hardware improves

#### 3. **How Salting Works**

**Without Salt (Vulnerable):**

```
User 1: password123 → hash1234567890
User 2: password123 → hash1234567890  ❌ Same hash!
```

Attackers can use rainbow tables (pre-computed hashes) to crack passwords.

**With Salt (Secure):**

```
User 1: password123 + salt_abc → hash_xyz123
User 2: password123 + salt_def → hash_abc789  ✅ Different hashes!
```

Each user gets a unique hash even with the same password.

### 🔑 Complete Authentication Flow

#### Registration Process:

```
1. User sends: { email, username, password: "MyPass123!" }
2. Backend validates input
3. Bcrypt generates random salt (e.g., "$2a$12$N9qo8uLOickgx2ZMRZoMye")
4. Bcrypt hashes password with salt
5. Result: "$2a$12$N9qo8uLOickgx2ZMRZoMye.ijfuNVqFn0rgCs9b0S7QV1K6R5wTT6"
                 │     │  └─────────────┬──────────────┘ └────────┬────────┘
                 │     │              Salt (22 chars)       Hash (31 chars)
                 │     └─ Cost factor
                 └─ Algorithm version
6. Stored in database: password_hash column
7. JWT token generated and returned
```

#### Login Process:

```
1. User sends: { email, password: "MyPass123!" }
2. Backend fetches user's stored hash from database
3. Bcrypt extracts salt from stored hash
4. Bcrypt hashes provided password with extracted salt
5. Compares new hash with stored hash using constant-time comparison
6. If match: Generate JWT and return success
   If no match: Return "Invalid credentials" error
```

### 🛡️ Security Features

#### 1. **Protection Against Common Attacks**

| Attack Type             | How We Prevent It                                          |
| ----------------------- | ---------------------------------------------------------- |
| **Rainbow Tables**      | Unique salt per password makes pre-computed tables useless |
| **Dictionary Attacks**  | High cost factor (12) makes brute force extremely slow     |
| **Timing Attacks**      | Bcrypt uses constant-time comparison                       |
| **SQL Injection**       | Parameterized queries (?) prevent injection                |
| **Duplicate Passwords** | Even identical passwords have different hashes             |
| **Database Leaks**      | Passwords cannot be reversed from hashes                   |

#### 2. **Input Validation**

```go
// Minimum 6 characters
if len(req.Password) < 6 {
    return "Password must be at least 6 characters long"
}

// Email format validation
if !strings.Contains(req.Email, "@") {
    return "Invalid email format"
}
```

#### 3. **JWT Token Security**

- Tokens expire after 24 hours
- Signed with HMAC-SHA256
- Contains user claims (ID, email, username)
- Protected routes require valid token

### 📊 Database Schema

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,  -- Stores bcrypt hash (60 chars)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_email (email),
    INDEX idx_users_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

**Key Points:**

- `password_hash` stores the full bcrypt output (includes salt)
- No separate `salt` column needed (salt is embedded in hash)
- `VARCHAR(255)` accommodates bcrypt's 60-character output
- Indexed on email and username for fast lookups

### 🔬 Example Hash Breakdown

```
$2a$12$N9qo8uLOickgx2ZMRZoMye.ijfuNVqFn0rgCs9b0S7QV1K6R5wTT6
│││  │  └─────────────┬──────────────┘ └────────┬────────┘
│││  │              Salt                      Hash
│││  └─ Cost (2^12 = 4096 rounds)
││└─ Minor version
│└─ Major version (2a = bcrypt)
└─ Identifier
```

### 🚀 Performance Considerations

**Hash Generation Time (Cost Factor 12):**

- ~250-300ms per password on modern hardware
- Intentionally slow to prevent brute force attacks
- Fast enough for user experience
- Scales with hardware improvements (cost can be increased)

**Database Storage:**

- Each hash: 60 bytes
- 1 million users: ~60 MB for password hashes
- Negligible storage overhead

### ✅ Best Practices Implemented

1. ✅ **Never store plain-text passwords**
2. ✅ **Use strong, adaptive hashing algorithm (bcrypt)**
3. ✅ **Unique salt per password (automatic with bcrypt)**
4. ✅ **Sufficient cost factor (12 = industry standard)**
5. ✅ **Parameterized queries (SQL injection prevention)**
6. ✅ **Input validation (email format, password length)**
7. ✅ **Secure token generation (JWT with expiration)**
8. ✅ **Password hash never exposed in API responses**
9. ✅ **Constant-time comparison (timing attack prevention)**
10. ✅ **HTTPS recommended for production**

### 🔧 Configuration Options

**Increase Security (Higher Cost):**

```go
// In internal/auth/password.go
// Increase from 12 to 13 or 14 for extra security
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 13)
```

**Note:** Each increase in cost factor **doubles** the computation time.

### 📚 References

- [Bcrypt Wikipedia](https://en.wikipedia.org/wiki/Bcrypt)
- [OWASP Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- [Go Bcrypt Package Documentation](https://pkg.go.dev/golang.org/x/crypto/bcrypt)

### 🔒 Additional Security Recommendations

For production environments, also consider:

1. **Rate Limiting**: Prevent brute force attacks on login endpoint
2. **Account Lockout**: Lock accounts after N failed login attempts
3. **Password Complexity**: Enforce minimum complexity rules
4. **2FA**: Add two-factor authentication
5. **Password Reset**: Secure password recovery mechanism
6. **Audit Logging**: Log authentication events
7. **HTTPS**: Always use TLS/SSL in production
8. **Environment Variables**: Never commit secrets to version control
9. **Regular Updates**: Keep dependencies updated
10. **Security Audits**: Regular code reviews and penetration testing
