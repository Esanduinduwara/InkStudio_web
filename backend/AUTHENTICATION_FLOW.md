# InkStudio Backend - Authentication Flow Diagram

## 🔐 Registration Flow with Salt & Hashing

```
┌──────────────────────────────────────────────────────────────────────────┐
│                          USER REGISTRATION                                │
└──────────────────────────────────────────────────────────────────────────┘

CLIENT                    BACKEND                       DATABASE
  │                          │                              │
  │  POST /api/register      │                              │
  ├─────────────────────────>│                              │
  │  {                       │                              │
  │    email: "user@x.com"   │  1. Validate Input          │
  │    username: "john"      │     ✓ Email format          │
  │    password: "Pass123!"  │     ✓ Password length       │
  │  }                       │     ✓ Required fields       │
  │                          │                              │
  │                          │  2. Hash Password            │
  │                          │     ┌──────────────────┐    │
  │                          │     │  bcrypt.Generate │    │
  │                          │     │  - Cost: 12      │    │
  │                          │     │  - Gen Salt      │    │
  │                          │     │  - Hash + Salt   │    │
  │                          │     └──────────────────┘    │
  │                          │     Result:                  │
  │                          │     $2a$12$N9qo8uLOickgx... │
  │                          │      │  │  └─Salt─┘└─Hash─┘ │
  │                          │      │  └─Cost               │
  │                          │      └─Algo                  │
  │                          │                              │
  │                          │  3. Store User              │
  │                          ├─────────────────────────────>│
  │                          │  INSERT INTO users          │
  │                          │  (email, username,          │
  │                          │   password_hash)            │
  │                          │                              │
  │                          │<─────────────────────────────┤
  │                          │  User Created (ID: 1)       │
  │                          │                              │
  │                          │  4. Generate JWT            │
  │                          │     ┌──────────────────┐    │
  │                          │     │ Create Token     │    │
  │                          │     │ - UserID         │    │
  │                          │     │ - Email          │    │
  │                          │     │ - Expires: 24h   │    │
  │                          │     │ - Sign with key  │    │
  │                          │     └──────────────────┘    │
  │                          │                              │
  │<─────────────────────────┤                              │
  │  201 Created             │                              │
  │  {                       │                              │
  │    token: "eyJhbG..."    │                              │
  │    user: {               │                              │
  │      id: 1               │                              │
  │      email: "user@x.com" │                              │
  │      username: "john"    │                              │
  │    },                    │                              │
  │    message: "Success"    │                              │
  │  }                       │                              │
  │                          │                              │
```

## 🔑 Login Flow with Password Verification

```
┌──────────────────────────────────────────────────────────────────────────┐
│                             USER LOGIN                                    │
└──────────────────────────────────────────────────────────────────────────┘

CLIENT                    BACKEND                       DATABASE
  │                          │                              │
  │  POST /api/login         │                              │
  ├─────────────────────────>│                              │
  │  {                       │                              │
  │    email: "user@x.com"   │  1. Validate Input          │
  │    password: "Pass123!"  │     ✓ Email present         │
  │  }                       │     ✓ Password present       │
  │                          │                              │
  │                          │  2. Fetch User              │
  │                          ├─────────────────────────────>│
  │                          │  SELECT id, email, username,│
  │                          │  password_hash              │
  │                          │  WHERE email = ?            │
  │                          │                              │
  │                          │<─────────────────────────────┤
  │                          │  User Found:                │
  │                          │  {                          │
  │                          │    id: 1                    │
  │                          │    password_hash:           │
  │                          │    "$2a$12$N9qo8u..."      │
  │                          │  }                          │
  │                          │                              │
  │                          │  3. Verify Password         │
  │                          │     ┌──────────────────┐    │
  │                          │     │ bcrypt.Compare   │    │
  │                          │     │                  │    │
  │                          │     │ Extract salt     │    │
  │                          │     │ from stored hash │    │
  │                          │     │      ↓           │    │
  │                          │     │ Hash input with  │    │
  │                          │     │ extracted salt   │    │
  │                          │     │      ↓           │    │
  │                          │     │ Compare hashes   │    │
  │                          │     │ (constant time)  │    │
  │                          │     └──────────────────┘    │
  │                          │     Result: ✓ Match         │
  │                          │                              │
  │                          │  4. Generate JWT            │
  │                          │     (Same as registration)  │
  │                          │                              │
  │<─────────────────────────┤                              │
  │  200 OK                  │                              │
  │  {                       │                              │
  │    token: "eyJhbG..."    │                              │
  │    user: {...}           │                              │
  │    message: "Success"    │                              │
  │  }                       │                              │
  │                          │                              │
```

## 🛡️ Protected Route Access

```
┌──────────────────────────────────────────────────────────────────────────┐
│                        PROTECTED ROUTE ACCESS                             │
└──────────────────────────────────────────────────────────────────────────┘

CLIENT                    BACKEND                       DATABASE
  │                          │                              │
  │  GET /api/profile        │                              │
  ├─────────────────────────>│                              │
  │  Headers:                │                              │
  │    Authorization:        │  1. Auth Middleware         │
  │    Bearer eyJhbG...      │     ┌──────────────────┐    │
  │                          │     │ Extract Token    │    │
  │                          │     │ Verify Signature │    │
  │                          │     │ Check Expiration │    │
  │                          │     │ Parse Claims     │    │
  │                          │     └──────────────────┘    │
  │                          │     ✓ Valid Token           │
  │                          │     UserID: 1                │
  │                          │                              │
  │                          │  2. Fetch User Data         │
  │                          ├─────────────────────────────>│
  │                          │  SELECT id, email, username │
  │                          │  WHERE id = ?               │
  │                          │                              │
  │                          │<─────────────────────────────┤
  │                          │  User Data                  │
  │                          │                              │
  │<─────────────────────────┤                              │
  │  200 OK                  │                              │
  │  {                       │                              │
  │    id: 1,                │                              │
  │    email: "user@x.com",  │                              │
  │    username: "john",     │                              │
  │    created_at: "..."     │                              │
  │  }                       │                              │
  │                          │                              │
```

## 🔒 Password Storage in Database

```
┌─────────────────────────────────────────────────────────────────────┐
│                        USERS TABLE (MySQL)                          │
├─────┬──────────────────┬──────────┬───────────────────────────────┤
│ id  │ email            │ username │ password_hash                 │
├─────┼──────────────────┼──────────┼───────────────────────────────┤
│ 1   │ user@example.com │ john     │ $2a$12$N9qo8uLOickgx2ZMRZo...│
│     │                  │          │  │   │  └─────Salt──────┘    │
│     │                  │          │  │   └─ Cost Factor (12)     │
│     │                  │          │  └─ Bcrypt Version (2a)      │
├─────┼──────────────────┼──────────┼───────────────────────────────┤
│ 2   │ jane@example.com │ jane     │ $2a$12$xP3D9fLMnOqrst4UVWx...│
│     │                  │          │  (Different hash - diff salt) │
└─────┴──────────────────┴──────────┴───────────────────────────────┘

Note: Even if two users have the same password, the hashes are different
      because each password gets a unique random salt!

Example with same password "test123":
User A: $2a$12$AbCdEfGhIjKlMnOpQrStUv...  ← Unique salt
User B: $2a$12$XyZwVuTsRqPoNmLkJiHgFe...  ← Different salt
        └──────────────┬──────────────┘
              Different salts = Different hashes
```

## 🎯 Security Features Visualization

```
┌────────────────────────────────────────────────────────────────────┐
│                    SECURITY LAYERS                                 │
└────────────────────────────────────────────────────────────────────┘

Input Layer:
  ┌──────────────────────────────────────────┐
  │ ✓ Email format validation                │
  │ ✓ Password length (min 6 chars)          │
  │ ✓ Required fields check                  │
  │ ✓ SQL injection prevention (? params)    │
  └──────────────────────────────────────────┘
                    ↓
Hashing Layer:
  ┌──────────────────────────────────────────┐
  │ ✓ Bcrypt algorithm (industry standard)   │
  │ ✓ Cost factor 12 (4,096 iterations)      │
  │ ✓ Unique random salt per password        │
  │ ✓ Salt embedded in hash (60 chars total) │
  │ ✓ Constant-time comparison                │
  └──────────────────────────────────────────┘
                    ↓
Storage Layer:
  ┌──────────────────────────────────────────┐
  │ ✓ No plain-text passwords stored         │
  │ ✓ Password hash in VARCHAR(255)          │
  │ ✓ UTF-8 MB4 character set                │
  │ ✓ InnoDB engine with transactions        │
  │ ✓ Indexed email/username for fast lookup │
  └──────────────────────────────────────────┘
                    ↓
Authentication Layer:
  ┌──────────────────────────────────────────┐
  │ ✓ JWT tokens with 24h expiration         │
  │ ✓ HMAC-SHA256 signature                  │
  │ ✓ Bearer token authentication            │
  │ ✓ Claims validation (exp, iat, nbf)      │
  └──────────────────────────────────────────┘
                    ↓
Transport Layer:
  ┌──────────────────────────────────────────┐
  │ ✓ CORS middleware                         │
  │ ✓ Request logging                         │
  │ ✓ HTTPS recommended for production        │
  └──────────────────────────────────────────┘
```

## 🚫 What Attackers CANNOT Do

```
❌ Rainbow Table Attack
   └─> Every password has unique salt → Pre-computed tables useless

❌ Dictionary Attack
   └─> Cost factor 12 = 4096 iterations → Too slow to brute force

❌ Timing Attack
   └─> Bcrypt uses constant-time comparison → No timing leaks

❌ SQL Injection
   └─> Parameterized queries (?) → Injection attempts fail

❌ Database Leak Recovery
   └─> Bcrypt is one-way → Cannot reverse hash to password

❌ Duplicate Password Detection
   └─> Same password = Different hashes → Cannot detect patterns
```

## 📊 Performance Metrics

```
Hash Generation (Cost 12):
  ┌──────────────────────────────┐
  │ Time: ~250-300ms per password │
  │ CPU: Moderate                 │
  │ Memory: Minimal               │
  │ Security: Very High           │
  └──────────────────────────────┘

Database Storage:
  ┌──────────────────────────────┐
  │ Hash size: 60 bytes          │
  │ 1K users: ~60 KB             │
  │ 1M users: ~60 MB             │
  │ Overhead: Negligible         │
  └──────────────────────────────┘

Login Verification:
  ┌──────────────────────────────┐
  │ Time: ~250-300ms             │
  │ DB Query: <10ms              │
  │ Hash Compare: ~240-290ms     │
  │ JWT Gen: <10ms               │
  └──────────────────────────────┘
```
