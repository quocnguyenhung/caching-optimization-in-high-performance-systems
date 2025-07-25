#!/bin/bash

# CONFIGURATION
POST_IDS=(1 2 3 4 5 6 7 8 9 10)   # Replace with your actual post IDs
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTMxNTMxODcsInVzZXJfaWQiOjN9.cncYY59GMrP7ZszCURqlmKSojmWzF27XAXzudUX_hZU"            # Replace with a valid JWT token
LIKE_COUNT=20                     # Number of random likes
LIMIT=10                          # Trending limit
SERVER_URL="http://localhost:8080"

# 1. Like random posts
echo "Liking $LIKE_COUNT random posts..."
for i in $(seq 1 $LIKE_COUNT); do
  # Pick a random post ID
  POST_ID=${POST_IDS[$RANDOM % ${#POST_IDS[@]}]}
  echo "Liking post $POST_ID"
  curl -s -X POST "$SERVER_URL/posts/$POST_ID/like" -H "Authorization: Bearer $TOKEN" > /dev/null
done
echo "Done liking."

# 2. Clear the trending cache for this limit
echo "Clearing trending cache key in Redis..."
redis-cli DEL trending:top:$LIMIT

# 3. Access trending with adaptive TTL (conventional)
echo "Accessing /trending/ttl?mode=conventional&limit=$LIMIT"
curl "$SERVER_URL/trending/ttl?mode=conventional&limit=$LIMIT"
echo

# 4. Access trending with adaptive TTL (inverted)
echo "Accessing /trending/ttl?mode=inverted&limit=$LIMIT"
curl "$SERVER_URL/trending/ttl?mode=inverted&limit=$LIMIT"
echo

# 5. Access trending multiple times to check for cache hits
echo "Accessing trending endpoints multiple times to check for cache hits..."
for i in {1..3}; do
  curl -s "$SERVER_URL/trending/ttl?mode=conventional&limit=$LIMIT" > /dev/null
  curl -s "$SERVER_URL/trending/ttl?mode=inverted&limit=$LIMIT" > /dev/null
done

# 6. Optionally, access random posts' trending data (simulate user activity)
echo "Accessing trending for random limits..."
for i in {1..5}; do
  RAND_LIMIT=$(( (RANDOM % 10) + 1 ))
  curl -s "$SERVER_URL/trending/ttl?mode=conventional&limit=$RAND_LIMIT" > /dev/null
  curl -s "$SERVER_URL/trending/ttl?mode=inverted&limit=$RAND_LIMIT" > /dev/null
done

echo "Done. Check your server logs for [CACHE SET], [CACHE HIT], [CACHE MISS], and TTL values."