#!/bin/bash

# InkStudio Backend Setup Script
# This script helps set up the development environment

set -e

echo "🎨 InkStudio Backend Setup"
echo "=========================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    echo "   Visit: https://golang.org/doc/install"
    exit 1
fi

echo "✅ Go version: $(go version)"
echo ""

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo "⚠️  PostgreSQL client not found. Make sure PostgreSQL is installed and running."
    echo "   Linux: sudo apt-get install postgresql"
    echo "   Mac: brew install postgresql"
else
    echo "✅ PostgreSQL client found"
fi
echo ""

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod download
go mod tidy
echo "✅ Dependencies installed"
echo ""

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "✅ .env file created. Please update it with your configuration."
else
    echo "✅ .env file already exists"
fi
echo ""

# Create bin directory
mkdir -p bin
echo "✅ Created bin directory"
echo ""

# Check if database exists
echo "🗄️  Checking database connection..."
if command -v psql &> /dev/null; then
    # Try to connect to the database
    if psql -lqt postgres://postgres:postgres@localhost:5432/postgres 2>/dev/null | cut -d \| -f 1 | grep -qw inkstudio; then
        echo "✅ Database 'inkstudio' exists"
    else
        echo "⚠️  Database 'inkstudio' not found. Creating it..."
        psql postgres://postgres:postgres@localhost:5432/postgres -c "CREATE DATABASE inkstudio;" 2>/dev/null || echo "❌ Failed to create database. Please create it manually."
    fi
fi
echo ""

# Build the project
echo "🔨 Building the project..."
go build -o bin/server ./cmd/server
echo "✅ Build successful"
echo ""

echo "🎉 Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Update .env with your database credentials"
echo "  2. Run 'make run' or './bin/server' to start the server"
echo "  3. Run 'make dev' for development mode with hot reload"
echo ""
echo "Available commands:"
echo "  make help   - Show all available commands"
echo "  make run    - Build and run the server"
echo "  make dev    - Run in development mode"
echo "  make test   - Run tests"
echo ""
