
# Caching Optimization for High-Performance Systems

This project demonstrates various caching strategies in the context of a social media backend system. It is the implementation for the thesis _"Caching Optimization for High-Performance Systems"_, using Go, PostgreSQL, and Redis to evaluate performance under different conditions.

## 🔧 Features

- JWT authentication
- Write-through cache and fan-out on write (timeline)
- Redis ZSet for trending
- Adaptive TTL cache control
- Full benchmarking pipeline using `wrk` and Lua scripting

## 📁 Technologies Used

- Go
- PostgreSQL
- Redis
- wrk (HTTP benchmarking)
- Lua (for scripted requests)

## .env file
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=thesis_db
REDIS_ADDR=localhost:6379
JWT_SECRET=myStrongSecretKeyHere


## 🚀 How to Run

### 1. Start the Server
```bash
go run cmd/server/main.go
```

### 2. Make Scripts Executable
```bash
chmod +x scripts/benchmarkTimeline.sh
chmod +x scripts/benchmarkTrending.sh
chmod +x scripts/benchmarkProfile.sh
```

### 3. Run Benchmarks
```bash
./scripts/benchmarkTimeline.sh
./scripts/benchmarkTrending.sh
./scripts/benchmarkProfile.sh
```

### 4. User Setup
```bash
curl -X POST http://localhost:8080/signup -H "Content-Type: application/json" -d '{"username":"user2", "password":"password2"}'
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username":"user1", "password":"password1"}'
```

### 5. Benchmarking Individual Endpoints
- Set JWT token in `.lua` script for authorization headers.
- Use `wrk` to test endpoints.

## 📊 Benchmark Results Summary

| Endpoint         | Strategy                    | Requests/sec | Avg Latency | Max Latency | Transfer/sec | Key Observations                                      |
|------------------|-----------------------------|--------------|-------------|-------------|---------------|-------------------------------------------------------|
| `/posts`         | DB Only                     | 54,176       | 87.67 ms    | 1.98 s      | 8.99 MB/s     | No caching; high throughput.                         |
| `/posts`         | Write-Through + Fan-Out     | 47,302       | 118.18 ms   | 1.98 s      | 7.85 MB/s     | Slightly slower due to Redis and fan-out writes.     |
| `/timeline`      | DB Only                     | 2,987        | 70.58 ms    | 2.00 s      | 0.60 MB/s     | Low performance due to join-heavy DB operations.     |
| `/timeline`      | Fan-Out Caching             | 23,327       | 27.68 ms    | 2.00 s      | 4.35 MB/s     | ~8× faster; Redis list accelerates feed retrieval.   |
| `/trending`      | DB Only                     | 6,528        | 265.49 ms   | 2.00 s      | 1.23 MB/s     | Slow due to aggregation/sorting logic.               |
| `/trending`      | ZSet Cache                  | 17,753       | 192.84 ms   | 1.86 s      | 2.05 MB/s     | ~2.7× speedup; Redis ZSet enables fast ranking.      |
| `/trending/ttl`  | Adaptive TTL (Conventional) | 33,149       | 293.75 ms   | 399.51 ms   | 5.93 MB/s     | Best throughput overall; stable TTL for hot content. |
| `/trending/ttl`  | Adaptive TTL (Inverted)     | 17,742       | 278.87 ms   | 675.44 ms   | 3.11 MB/s     | Lower throughput; frequent evictions.                |
| `/profile/:id`   | Cold Cache (DB only)        | 32,306       | 208.94 ms   | 1.92 s      | 6.28 MB/s     | Moderate latency with many misses.                   |
| `/profile/:id`   | Warm Cache (Redis hit)      | 40,529       | 232.29 ms   | 300.73 ms   | 7.89 MB/s     | Improved throughput and lower DB load.               |

