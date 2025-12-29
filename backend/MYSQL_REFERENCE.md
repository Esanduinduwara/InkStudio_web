# MySQL Quick Reference for InkStudio Backend

## Connection

```bash
# Connect to MySQL
mysql -u root -p

# Connect to specific database
mysql -u root -p inkstudio
```

## Database Commands

```sql
-- Show all databases
SHOW DATABASES;

-- Use inkstudio database
USE inkstudio;

-- Show all tables
SHOW TABLES;

-- Describe users table
DESCRIBE users;

-- Show table creation SQL
SHOW CREATE TABLE users;
```

## User Management

```sql
-- View all users
SELECT id, email, username, created_at FROM users;

-- Count total users
SELECT COUNT(*) FROM users;

-- Find user by email
SELECT * FROM users WHERE email = 'test@example.com';

-- Delete a user (be careful!)
DELETE FROM users WHERE email = 'test@example.com';

-- Delete all users (be very careful!)
TRUNCATE TABLE users;
```

## Password Hash Information

```sql
-- View password hashes (for debugging only, never expose in production)
SELECT id, email,
       LEFT(password_hash, 7) as algorithm,
       SUBSTRING(password_hash, 5, 2) as cost_factor,
       SUBSTRING(password_hash, 8, 22) as salt,
       LENGTH(password_hash) as hash_length
FROM users;

-- Example output:
-- algorithm: $2a$12$  (bcrypt version 2a, cost 12)
-- salt: 22 random characters
-- hash_length: 60 characters total
```

## Database Maintenance

```sql
-- Check table size
SELECT
    table_name AS 'Table',
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)'
FROM information_schema.TABLES
WHERE table_schema = 'inkstudio'
    AND table_name = 'users';

-- Optimize table
OPTIMIZE TABLE users;

-- Check table status
SHOW TABLE STATUS WHERE Name = 'users';
```

## Backup and Restore

```bash
# Backup database
mysqldump -u root -p inkstudio > backup.sql

# Backup with timestamp
mysqldump -u root -p inkstudio > inkstudio_$(date +%Y%m%d_%H%M%S).sql

# Restore database
mysql -u root -p inkstudio < backup.sql

# Backup only structure (no data)
mysqldump -u root -p --no-data inkstudio > structure.sql

# Backup only data (no structure)
mysqldump -u root -p --no-create-info inkstudio > data.sql
```

## Testing Queries

```sql
-- Insert test user (password is hashed "test123")
INSERT INTO users (email, username, password_hash)
VALUES (
    'test@example.com',
    'testuser',
    '$2a$12$R9h/cIPz0gi.URNNX3kh2OPST9/PgBkqquzi.Ss7KIUgO2t0jWMUW'
);

-- View recently created users
SELECT id, email, username, created_at
FROM users
ORDER BY created_at DESC
LIMIT 10;

-- Check for duplicate emails
SELECT email, COUNT(*) as count
FROM users
GROUP BY email
HAVING count > 1;

-- Users created today
SELECT * FROM users
WHERE DATE(created_at) = CURDATE();
```

## Index Management

```sql
-- Show indexes on users table
SHOW INDEX FROM users;

-- Create new index (if needed)
CREATE INDEX idx_created_at ON users(created_at);

-- Drop index (if needed)
DROP INDEX idx_created_at ON users;

-- Analyze index usage
ANALYZE TABLE users;
```

## Security Commands

```sql
-- Create dedicated application user
CREATE USER 'inkstudio_app'@'localhost' IDENTIFIED BY 'strong_password_here';

-- Grant necessary permissions
GRANT SELECT, INSERT, UPDATE ON inkstudio.* TO 'inkstudio_app'@'localhost';

-- Flush privileges
FLUSH PRIVILEGES;

-- View user privileges
SHOW GRANTS FOR 'inkstudio_app'@'localhost';

-- Remove user (if needed)
DROP USER 'inkstudio_app'@'localhost';
```

## Monitoring

```sql
-- Show active connections
SHOW PROCESSLIST;

-- Show connection statistics
SHOW STATUS LIKE 'Threads_connected';
SHOW STATUS LIKE 'Max_used_connections';

-- Show slow queries
SHOW VARIABLES LIKE 'slow_query_log';
SHOW VARIABLES LIKE 'long_query_time';
```

## Common Issues

### Connection Issues

```bash
# Check if MySQL is running
sudo systemctl status mysql

# Start MySQL
sudo systemctl start mysql

# Check MySQL port
sudo netstat -tlnp | grep 3306
```

### Permission Issues

```sql
-- If you get "Access denied"
-- Reset password for root
ALTER USER 'root'@'localhost' IDENTIFIED BY 'new_password';
FLUSH PRIVILEGES;
```

### Character Set Issues

```sql
-- Check character set
SHOW VARIABLES LIKE 'character_set%';

-- Set database character set
ALTER DATABASE inkstudio CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Set table character set
ALTER TABLE users CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## Performance Tuning

```sql
-- Check query performance
EXPLAIN SELECT * FROM users WHERE email = 'test@example.com';

-- Enable query cache (MySQL 5.7 and earlier)
SET GLOBAL query_cache_size = 1048576;
SET GLOBAL query_cache_type = ON;

-- Check buffer pool size (InnoDB)
SHOW VARIABLES LIKE 'innodb_buffer_pool_size';
```

## Connection String Examples

```bash
# Local development
DATABASE_URL=root:password@tcp(localhost:3306)/inkstudio?parseTime=true&charset=utf8mb4

# Remote server
DATABASE_URL=user:pass@tcp(192.168.1.100:3306)/inkstudio?parseTime=true&charset=utf8mb4

# With SSL
DATABASE_URL=user:pass@tcp(host:3306)/inkstudio?tls=true&parseTime=true&charset=utf8mb4

# Docker container
DATABASE_URL=root:password@tcp(mysql:3306)/inkstudio?parseTime=true&charset=utf8mb4
```

## Useful Commands for Development

```bash
# Drop and recreate database
mysql -u root -p -e "DROP DATABASE IF EXISTS inkstudio; CREATE DATABASE inkstudio CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# Run SQL file
mysql -u root -p inkstudio < schema.sql

# Export users to CSV
mysql -u root -p inkstudio -e "SELECT id, email, username, created_at FROM users" | sed 's/\t/,/g' > users.csv

# Quick check if backend can connect
mysql -u root -p inkstudio -e "SELECT 'Connection successful!' AS status;"
```
