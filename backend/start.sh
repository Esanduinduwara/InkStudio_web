#!/bin/bash

# InkStudio Backend Quick Start Script

echo "🚀 InkStudio Backend - Quick Start"
echo "=================================="
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

echo "✓ Docker is running"
echo ""

# Navigate to backend directory
cd "$(dirname "$0")"

# Stop and remove existing containers
echo "🧹 Cleaning up existing containers..."
docker compose down -v 2>/dev/null || true
echo ""

# Build and start containers
echo "🔨 Building and starting containers..."
docker compose up --build -d

# Wait for services to be healthy
echo ""
echo "⏳ Waiting for services to be ready..."
sleep 5

# Check if MySQL is healthy
MYSQL_HEALTHY=false
for i in {1..30}; do
    if docker exec inkstudio-mysql mysqladmin ping -h localhost -u root -prootpassword --silent 2>/dev/null; then
        MYSQL_HEALTHY=true
        break
    fi
    echo -n "."
    sleep 1
done

echo ""
if [ "$MYSQL_HEALTHY" = true ]; then
    echo "✓ MySQL is healthy"
else
    echo "❌ MySQL failed to start"
    docker compose logs mysql
    exit 1
fi

# Check if backend is responding
BACKEND_HEALTHY=false
for i in {1..10}; do
    if curl -s http://localhost:8080/health > /dev/null; then
        BACKEND_HEALTHY=true
        break
    fi
    echo -n "."
    sleep 1
done

echo ""
if [ "$BACKEND_HEALTHY" = true ]; then
    echo "✓ Backend is healthy"
else
    echo "❌ Backend failed to start"
    docker compose logs backend
    exit 1
fi

echo ""
echo "✅ All services are running!"
echo ""
echo "📍 Services:"
echo "   - Backend API: http://localhost:8080"
echo "   - MySQL: localhost:3307"
echo ""
echo "📋 Useful commands:"
echo "   - View logs:        docker compose logs -f"
echo "   - Stop services:    docker compose down"
echo "   - Restart:          docker compose restart"
echo "   - Run tests:        ./test-api.sh"
echo ""
echo "🧪 Running quick API test..."
echo ""

# Quick test
HEALTH=$(curl -s http://localhost:8080/health)
echo "Health Check: $HEALTH"
echo ""
echo "✨ Ready to go! Try running ./test-api.sh for comprehensive tests."
