# Caching Optimization for High-Performance Systems

## 📎 Resources

- 🔗 [Demo Video](https://www.youtube.com/watch?v=UdKiaBQyl5A)
- 📄 [Thesis Report & Slides](https://drive.google.com/drive/folders/1zzhBtGxGyeFITczU0Efz6mBtrRfZf2im)

---

This project explores caching strategies in a social media backend. Implemented in Go, PostgreSQL, and Redis, it evaluates static TTL, adaptive TTL, and write-through caching using real benchmarks.

## 🔧 Features

- JWT authentication
- Write-through caching with fan-out on write (timeline)
- Redis ZSet-based trending cache
- Adaptive TTL strategy for dynamic content
- Benchmarking pipeline with `wrk` and Lua scripting

## 🧪 Tech Stack

- Go
- PostgreSQL
- Redis
- wrk (HTTP benchmarking)
- Lua (scripting)

## ⚙️ Environment (.env)
```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=thesis_db
REDIS_ADDR=localhost:6379
JWT_SECRET=myStrongSecretKeyHere
```


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
chmod +x benchmark_trending_ttl.sh
```

### 3. Run Benchmarks
```bash
./scripts/benchmarkTimeline.sh
./scripts/benchmarkTrending.sh
./scripts/benchmarkProfile.sh
./benchmark_trending_ttl.sh # tests conventional vs inverted TTL
```

### 4. User Setup
```bash
curl -X POST http://localhost:8080/signup -H "Content-Type: application/json" -d '{"username":"user2", "password":"password2"}'
curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"username":"user1", "password":"password1"}'
```

### 5. Benchmarking Individual Endpoints
Set your JWT in the Lua scripts under `scripts/` before running the following
benchmark commands.

- **Posts**
  ```bash
  wrk -t4 -c1000 -d30s -s scripts/wrk_create_post.lua http://localhost:8080/posts
  ```
- **Timeline (DB vs Cache)**
  ```bash
  ./scripts/benchmarkTimeline.sh
  ```
- **Profile (Cold vs Warm)**
  ```bash
  ./scripts/benchmarkProfile.sh
  ```
- **Trending (DB vs Static TTL Cache)**
  ```bash
  ./scripts/benchmarkTrending.sh
  ```
- **Adaptive TTL – Conventional & Inverted**
  ```bash
  ./benchmark_trending_ttl.sh
  ```

## 🧪 API Testing Guide
Use the JWT returned from `/login` in the `Authorization` header for protected routes.

### Creating Posts
```bash
curl -X POST http://localhost:8080/posts \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"content":"hello world"}'
```

### Following Users
```bash
curl -X POST http://localhost:8080/follow/2 \
  -H "Authorization: Bearer <TOKEN>"
```

### Timeline Endpoints
```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/timeline/db
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/timeline/cache
```

### Trending Endpoints
```bash
curl http://localhost:8080/trending/db?limit=10
curl http://localhost:8080/trending/cache?limit=10
```

### Adaptive TTL (Conventional vs. Inverted)
```bash
curl http://localhost:8080/trending/ttl?mode=conventional&limit=10
curl http://localhost:8080/trending/ttl?mode=inverted&limit=10
```

### TTL Metrics
```bash
curl http://localhost:8080/metrics/ttl?strategy=conventional
curl http://localhost:8080/metrics/ttl?strategy=inverted
```

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

