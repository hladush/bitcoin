#!/bin/bash

# Bitcoin Tracker API Test Script

SERVER_URL="http://localhost:8080"
TEST_ADDRESS="bc1q0sg9rdst255gtldsmcf8rk0764avqy2h2ksqs5"

echo "🧪 Bitcoin Tracker API Test Suite"
echo "================================="

# Test 1: Health Check
echo "1. Testing health check..."
curl -s $SERVER_URL/health | jq '.' || echo "Health check failed"
echo ""

# Test 2: Add Address
echo "2. Adding test address..."
curl -s -X POST $SERVER_URL/addresses \
  -H "Content-Type: application/json" \
  -d "{\"address\": \"$TEST_ADDRESS\", \"label\": \"Test Address\"}" | jq '.' || echo "Add address failed"
echo ""

# Test 3: Get All Addresses
echo "3. Getting all addresses..."
curl -s $SERVER_URL/addresses | jq '.' || echo "Get addresses failed"
echo ""

# Test 4: Get Specific Address
echo "4. Getting specific address..."
curl -s $SERVER_URL/addresses/$TEST_ADDRESS | jq '.' || echo "Get specific address failed"
echo ""

# Test 5: Get Balance
echo "5. Getting address balance..."
curl -s $SERVER_URL/addresses/$TEST_ADDRESS/balance | jq '.' || echo "Get balance failed"
echo ""

# Test 6: Sync Address
echo "6. Syncing address..."
curl -s -X POST $SERVER_URL/addresses/$TEST_ADDRESS/sync | jq '.' || echo "Sync failed"
echo ""

# Test 7: Get Transactions
echo "7. Getting transactions..."
curl -s "$SERVER_URL/addresses/$TEST_ADDRESS/transactions?limit=5" | jq '.' || echo "Get transactions failed"
echo ""

echo "✅ Test suite completed!"
