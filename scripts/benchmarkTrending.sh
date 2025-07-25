#!/bin/bash

# Configurable parameters
THREADS=4
CONNECTIONS=100000
DURATION=30
TARGET_DB="http://localhost:8080/trending/db"
TARGET_CACHE="http://localhost:8080/trending/cache"
SCRIPT_PATH="scripts/wrk_trending.lua"

# Benchmark Trending DB
echo "===== Benchmarking /trending/db (Database Query) ====="
wrk -t$THREADS -c$CONNECTIONS -d${DURATION}s -s $SCRIPT_PATH $TARGET_DB > db_result.txt

# Extract metrics for DB
REQ_DB=$(grep "Requests/sec" db_result.txt | awk '{print $2}')
LAT_DB=$(grep -m1 "Latency" db_result.txt | awk '{print $2}')
MAX_LAT_DB=$(grep -m1 "Latency" db_result.txt | awk '{print $4}')
TRANSFER_DB=$(grep "Transfer/sec" db_result.txt | awk '{print $2}')

# Benchmark Trending Cache
echo "===== Benchmarking /trending/cache (Redis Cache) ====="
wrk -t$THREADS -c$CONNECTIONS -d${DURATION}s -s $SCRIPT_PATH $TARGET_CACHE > cache_result.txt

# Extract metrics for Cache
REQ_CACHE=$(grep "Requests/sec" cache_result.txt | awk '{print $2}')
LAT_CACHE=$(grep -m1 "Latency" cache_result.txt | awk '{print $2}')
MAX_LAT_CACHE=$(grep -m1 "Latency" cache_result.txt | awk '{print $4}')
TRANSFER_CACHE=$(grep "Transfer/sec" cache_result.txt | awk '{print $2}')

# Pretty Print Full Comparison
echo ""
echo "=================== Benchmark Configuration ==================="
echo "Threads: $THREADS"
echo "Connections: $CONNECTIONS"
echo "Duration: ${DURATION}s"
echo "Target API (DB): $TARGET_DB"
echo "Target API (Cache): $TARGET_CACHE"
echo "==============================================================="
echo ""
echo "=================== Benchmark Results ==================="
printf "%-20s %-15s %-15s %-15s %-15s\n" "Test" "Req/sec" "Avg Latency" "Max Latency" "Transfer/sec"
printf "%-20s %-15s %-15s %-15s %-15s\n" "Trending DB" "$REQ_DB" "$LAT_DB" "$MAX_LAT_DB" "$TRANSFER_DB"
printf "%-20s %-15s %-15s %-15s %-15s\n" "Trending Cache" "$REQ_CACHE" "$LAT_CACHE" "$MAX_LAT_CACHE" "$TRANSFER_CACHE"
echo "==============================================================="
