#!/bin/bash

# Bitcoin Tracker Demo Script
# This script demonstrates the key features of the Bitcoin Tracker API

echo "🚀 Bitcoin Tracker Demo"
echo "======================"
echo ""

# Check if server is running
echo "Checking if server is running..."
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "❌ Server is not running. Please start it with: go run cmd/server/main.go"
    exit 1
fi
echo "✅ Server is running!"
echo ""

# Sample Bitcoin addresses for demo
ADDRESSES=(
    "bc1q0sg9rdst255gtldsmcf8rk0764avqy2h2ksqs5"
    "3E8ociqZa9mZUSwGdSmAEMAoAxBK3FNDcd"
)

LABELS=(
    "Sample Address 1"
    "Sample Address 2"
)

echo "📝 Adding sample Bitcoin addresses..."
for i in "${!ADDRESSES[@]}"; do
    echo "Adding ${ADDRESSES[$i]}..."
    curl -s -X POST http://localhost:8080/addresses \
        -H "Content-Type: application/json" \
        -d "{\"address\": \"${ADDRESSES[$i]}\", \"label\": \"${LABELS[$i]}\"}" \
        | jq '.success, .message // .error' || echo "Failed to add address"
done
echo ""

echo "📋 Listing all tracked addresses..."
curl -s http://localhost:8080/addresses | jq '.data[] | {address: .address, label: .label, balance_btc: .balance.balance_btc}' || echo "Failed to get addresses"
echo ""

echo "💰 Getting balance for first address..."
curl -s http://localhost:8080/addresses/${ADDRESSES[0]}/balance | jq '.data' || echo "Failed to get balance"
echo ""

echo "🔄 Syncing first address..."
curl -s -X POST http://localhost:8080/addresses/${ADDRESSES[0]}/sync | jq '.message // .error' || echo "Failed to sync"
echo ""

echo "📊 Getting recent transactions..."
curl -s "http://localhost:8080/addresses/${ADDRESSES[0]}/transactions?limit=3" | jq '.data[] | {hash: .hash, amount: .amount, type: .type}' || echo "Failed to get transactions"
echo ""

echo "✅ Demo completed!"
echo ""
echo "💡 Try these commands manually:"
echo "   curl http://localhost:8080/addresses"
echo "   curl http://localhost:8080/addresses/${ADDRESSES[0]}/balance"
echo "   curl http://localhost:8080/addresses/${ADDRESSES[0]}/transactions"
