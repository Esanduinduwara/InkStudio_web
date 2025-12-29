#!/bin/bash

# MySQL Setup Script for InkStudio Backend

echo "🗄️  Setting up MySQL database for InkStudio"
echo "==========================================="
echo ""

# Default values
DB_USER="root"
DB_PASS="password"
DB_NAME="inkstudio"
DB_HOST="localhost"
DB_PORT="3306"

# Check if MySQL is installed
if ! command -v mysql &> /dev/null; then
    echo "❌ MySQL client not found."
    echo ""
    echo "Please install MySQL:"
    echo "  Ubuntu/Debian: sudo apt-get install mysql-server mysql-client"
    echo "  macOS: brew install mysql"
    echo "  Or use Docker: docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password mysql:8.0"
    exit 1
fi

echo "✅ MySQL client found"
echo ""

# Prompt for credentials
read -p "Enter MySQL root password (default: password): " input_pass
if [ ! -z "$input_pass" ]; then
    DB_PASS="$input_pass"
fi

read -p "Enter database name (default: inkstudio): " input_db
if [ ! -z "$input_db" ]; then
    DB_NAME="$input_db"
fi

echo ""
echo "Creating database '$DB_NAME'..."

# Create database
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS <<EOF
CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE $DB_NAME;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_email (email),
    INDEX idx_users_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SHOW TABLES;
DESCRIBE users;
EOF

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Database created successfully!"
    echo ""
    echo "Connection string:"
    echo "  $DB_USER:$DB_PASS@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true&charset=utf8mb4"
    echo ""
    echo "Update your .env file with:"
    echo "  DATABASE_URL=$DB_USER:$DB_PASS@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true&charset=utf8mb4"
else
    echo ""
    echo "❌ Failed to create database. Please check your credentials."
    exit 1
fi
