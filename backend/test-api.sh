#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080"

echo -e "${BLUE}=== InkStudio Backend API Tests ===${NC}\n"

# Test 1: Health Check
echo -e "${BLUE}1. Testing Health Check...${NC}"
HEALTH=$(curl -s -X GET "$BASE_URL/health")
echo "Response: $HEALTH"
if echo "$HEALTH" | grep -q "healthy"; then
    echo -e "${GREEN}✓ Health check passed${NC}\n"
else
    echo -e "${RED}✗ Health check failed${NC}\n"
fi

# Test 2: Register New User
echo -e "${BLUE}2. Testing User Registration...${NC}"
TIMESTAMP=$(date +%s)
EMAIL="user${TIMESTAMP}@example.com"
USERNAME="user${TIMESTAMP}"
PASSWORD="testpass123"

REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"username\": \"$USERNAME\",
    \"password\": \"$PASSWORD\"
  }")

echo "Response: $REGISTER_RESPONSE"
TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ ! -z "$TOKEN" ]; then
    echo -e "${GREEN}✓ Registration successful${NC}"
    echo "Token: ${TOKEN:0:50}..."
else
    echo -e "${RED}✗ Registration failed${NC}"
fi
echo ""

# Test 3: Login
echo -e "${BLUE}3. Testing User Login...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

echo "Response: $LOGIN_RESPONSE"
LOGIN_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ ! -z "$LOGIN_TOKEN" ]; then
    echo -e "${GREEN}✓ Login successful${NC}"
else
    echo -e "${RED}✗ Login failed${NC}"
fi
echo ""

# Test 4: Get Current User (Protected Route)
echo -e "${BLUE}4. Testing Protected Route (/api/auth/me)...${NC}"
ME_RESPONSE=$(curl -s -X GET "$BASE_URL/api/auth/me" \
  -H "Authorization: Bearer $LOGIN_TOKEN")

echo "Response: $ME_RESPONSE"
if echo "$ME_RESPONSE" | grep -q "$EMAIL"; then
    echo -e "${GREEN}✓ Protected route access successful${NC}"
else
    echo -e "${RED}✗ Protected route access failed${NC}"
fi
echo ""

# Test 5: Wrong Password
echo -e "${BLUE}5. Testing Login with Wrong Password...${NC}"
WRONG_LOGIN=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"wrongpassword\"
  }")

echo "Response: $WRONG_LOGIN"
if echo "$WRONG_LOGIN" | grep -q "Invalid email or password"; then
    echo -e "${GREEN}✓ Wrong password rejected correctly${NC}"
else
    echo -e "${RED}✗ Wrong password test failed${NC}"
fi
echo ""

# Test 6: Duplicate Registration
echo -e "${BLUE}6. Testing Duplicate Registration...${NC}"
DUP_REGISTER=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"username\": \"$USERNAME\",
    \"password\": \"$PASSWORD\"
  }")

echo "Response: $DUP_REGISTER"
if echo "$DUP_REGISTER" | grep -q "already exists"; then
    echo -e "${GREEN}✓ Duplicate registration rejected correctly${NC}"
else
    echo -e "${RED}✗ Duplicate registration test failed${NC}"
fi
echo ""

# Test 7: Protected Route Without Token
echo -e "${BLUE}7. Testing Protected Route without Token...${NC}"
NO_TOKEN=$(curl -s -X GET "$BASE_URL/api/auth/me")

echo "Response: $NO_TOKEN"
if echo "$NO_TOKEN" | grep -q "Missing authorization header"; then
    echo -e "${GREEN}✓ Unauthorized access blocked correctly${NC}"
else
    echo -e "${RED}✗ Unauthorized access test failed${NC}"
fi
echo ""

echo -e "${BLUE}=== All Tests Completed ===${NC}"
