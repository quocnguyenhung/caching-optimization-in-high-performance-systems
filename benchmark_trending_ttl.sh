#!/bin/bash

# CONFIGURATION
POST_IDS=(1 2 3 4 5 6 7 8 9 10)   # Replace with your actual post IDs
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTMxNTMxODcsInVzZXJfaWQiOjN9.cncYY59GMrP7ZszCURqlmKSojmWzF27XAXzudUX_hZU"            # Replace with a valid JWT token
LIKE_COUNT=20                     # Number of random likes
LIMIT=10                          # Trending limit
SERVER_URL="http://localhost:8080"
WRK_THREADS=4
WRK_CONNECTIONS=100000
WRK_DURATION=30
WRK_SCRIPT="scripts/wrk_trending.lua"

# 1. Like random posts
echo "Liking $LIKE_COUNT random posts..."
for i in $(seq 1 $LIKE_COUNT); do
  POST_ID=${POST_IDS[$RANDOM % ${#POST_IDS[@]}]}
  echo "Liking post $POST_ID"
  curl -s -X POST "$SERVER_URL/posts/$POST_ID/like" -H "Authorization: Bearer $TOKEN" > /dev/null
done
echo "Done liking."


# 3. Benchmark adaptive TTL (conventional)
echo "Benchmarking /trending/ttl?mode=conventional&limit=$LIMIT"
wrk -t$WRK_THREADS -c$WRK_CONNECTIONS -d${WRK_DURATION}s -s $WRK_SCRIPT "$SERVER_URL/trending/ttl?mode=conventional&limit=$LIMIT" > wrk_conventional.txt

# 4. Benchmark adaptive TTL (inverted)
echo "Benchmarking /trending/ttl?mode=inverted&limit=$LIMIT"
wrk -t$WRK_THREADS -c$WRK_CONNECTIONS -d${WRK_DURATION}s -s $WRK_SCRIPT "$SERVER_URL/trending/ttl?mode=inverted&limit=$LIMIT" > wrk_inverted.txt

# 5. Show summary
echo "==== Conventional TTL Results ===="
grep "Requests/sec\\|Latency\\|Transfer/sec" wrk_conventional.txt

echo "==== Inverted TTL Results ===="
grep "Requests/sec\\|Latency\\|Transfer/sec" wrk_inverted.txt

echo "Check your server logs for [CACHE SET], [CACHE HIT], [CACHE MISS], and TTL values."