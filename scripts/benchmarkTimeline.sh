#!/bin/bash

# Configurable parameters
THREADS=4
CONNECTIONS=100
DURATION=30
TARGET_DB="http://localhost:8080/timeline/db"
TARGET_CACHE="http://localhost:8080/timeline/cache"
SCRIPT_PATH="scripts/wrk_timeline.lua"

# Benchmark Timeline DB
echo "===== Benchmarking /timeline/db (Database only) ====="
wrk -t$THREADS -c$CONNECTIONS -d${DURATION}s -s $SCRIPT_PATH $TARGET_DB > db_result.txt

# Extract metrics for DB
REQ_DB=$(grep "Requests/sec" db_result.txt | awk '{print $2}')
LAT_DB=$(grep -m1 "Latency" db_result.txt | awk '{print $2}')
MAX_LAT_DB=$(grep -m1 "Latency" db_result.txt | awk '{print $3}')
TRANSFER_DB=$(grep "Transfer/sec" db_result.txt | awk '{print $2}')
NON2XX_DB=$(grep "Non-2xx" db_result.txt | awk '{print $4}')

# Benchmark Timeline Cache
echo "===== Benchmarking /timeline/cache (Redis Cache) ====="
wrk -t$THREADS -c$CONNECTIONS -d${DURATION}s -s $SCRIPT_PATH $TARGET_CACHE > cache_result.txt

# Extract metrics for Cache
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
echo "Target API (DB): $TARGET_DB"
echo "Target API (Cache): $TARGET_CACHE"
echo "==============================================================="
echo ""
echo "=================== Benchmark Results ==================="
printf "%-20s %-15s %-15s %-15s %-15s %-15s\n" "Test" "Req/sec" "Avg Latency" "Max Latency" "Transfer/sec" "Non-2xx Errors"
printf "%-20s %-15s %-15s %-15s %-15s %-15s\n" "Timeline DB" "$REQ_DB" "$LAT_DB" "$MAX_LAT_DB" "$TRANSFER_DB" "${NON2XX_DB:-0}"
printf "%-20s %-15s %-15s %-15s %-15s %-15s\n" "Timeline Cache" "$REQ_CACHE" "$LAT_CACHE" "$MAX_LAT_CACHE" "$TRANSFER_CACHE" "${NON2XX_CACHE:-0}"
echo "==============================================================="
