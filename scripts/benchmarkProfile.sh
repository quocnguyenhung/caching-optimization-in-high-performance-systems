#!/bin/bash

# Configurable parameters
THREADS=4
CONNECTIONS=100000
DURATION=30
USER_ID=3 
TARGET_PROFILE="http://localhost:8080/profile/$USER_ID"
SCRIPT_PATH="scripts/wrk_profile.lua"

# Step 1: Clear Redis Cache for Cold Start
echo "===== Clearing Profile Cache for Cold Start ====="
redis-cli DEL profile:$USER_ID

# Benchmark Profile Cold (DB Read)
echo "===== Benchmarking /profile/:id (COLD - Database Read) ====="
wrk -t$THREADS -c$CONNECTIONS -d${DURATION}s -s $SCRIPT_PATH $TARGET_PROFILE > db_result.txt

# Extract metrics for Cold
REQ_DB=$(grep "Requests/sec" db_result.txt | awk '{print $2}')
LAT_DB=$(grep -m1 "Latency" db_result.txt | awk '{print $2}')
MAX_LAT_DB=$(grep -m1 "Latency" db_result.txt | awk '{print $3}')
TRANSFER_DB=$(grep "Transfer/sec" db_result.txt | awk '{print $2}')
NON2XX_DB=$(grep "Non-2xx" db_result.txt | awk '{print $4}')

# Step 2: Warm Up Cache
echo "===== Warming up Cache for Warm Start ====="
curl -s $TARGET_PROFILE > /dev/null

# Benchmark Profile Warm (Redis Cache)
echo "===== Benchmarking /profile/:id (WARM - Redis Cache) ====="
wrk -t$THREADS -c$CONNECTIONS -d${DURATION}s -s $SCRIPT_PATH $TARGET_PROFILE > cache_result.txt

# Extract metrics for Warm
REQ_CACHE=$(grep "Requests/sec" cache_result.txt | awk '{print $2}')
LAT_CACHE=$(grep -m1 "Latency" cache_result.txt | awk '{print $2}')
MAX_LAT_CACHE=$(grep -m1 "Latency" cache_result.txt | awk '{print $3}')
TRANSFER_CACHE=$(grep "Transfer/sec" cache_result.txt | awk '{print $2}')
NON2XX_CACHE=$(grep "Non-2xx" cache_result.txt | awk '{print $4}')

# Pretty Print Full Comparison
echo ""
echo "=================== Benchmark Configuration ==================="
echo "Threads: $THREADS"
echo "Connections: $CONNECTIONS"
echo "Duration: ${DURATION}s"
echo "Target API (Profile): $TARGET_PROFILE"
echo "==============================================================="
echo ""
echo "=================== Benchmark Results ==================="
printf "%-20s %-15s %-15s %-15s %-15s %-15s\n" "Test" "Req/sec" "Avg Latency" "Max Latency" "Transfer/sec" "Non-2xx Errors"
printf "%-20s %-15s %-15s %-15s %-15s %-15s\n" "Profile Cold" "$REQ_DB" "$LAT_DB" "$MAX_LAT_DB" "$TRANSFER_DB" "${NON2XX_DB:-0}"
printf "%-20s %-15s %-15s %-15s %-15s %-15s\n" "Profile Warm" "$REQ_CACHE" "$LAT_CACHE" "$MAX_LAT_CACHE" "$TRANSFER_CACHE" "${NON2XX_CACHE:-0}"
echo "==============================================================="
