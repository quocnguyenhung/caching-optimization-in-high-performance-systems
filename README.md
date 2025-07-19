# Caching Optimization in High-Performance Systems

This project is the final thesis of Quoc Nguyen Hung, aiming to build a scalable social media backend system like Twitter, optimized using multi-layer caching techniques.

---

## 📚 Project Overview

The system is developed using **Go**, **PostgreSQL**, **Redis**, **Docker**, and **wrk** for benchmarking.  
It simulates core social media functionalities such as user signup/login, following, posting, and viewing timelines, then applies various caching strategies to optimize performance.

---

## 🚀 Technologies Used

- **Golang** — Backend development
- **PostgreSQL** — Relational database
- **Redis** — Caching layer
- **Docker** — Environment setup
- **wrk** — HTTP benchmarking tool

---

## 🛠️ Features Implemented

- User signup/login with JWT authentication
- Follow and unfollow users
- Post creation
- Timeline fetching (from followed users' posts)
- Like posts
- Trending posts view
- User profile view
- Multi-layer caching strategies:
  - No Cache (Baseline)
  - Timeline Caching
  - Write-Through Caching
  - Fan-out on Write
  - Trending Cache
  - Profile Cache
  - Adaptive TTL for timelines (toggle via `ADAPTIVE_TTL_ENABLED`)

---

## 🧪 Benchmark Results

Benchmark scripts are available in the `scripts/` directory. Before running wrk obtain a JWT token:

```bash
TOKEN=$(curl -s -X POST -H "Content-Type: application/json" \
  -d '{"username":"demo","password":"demo"}' \
  http://localhost:8080/auth/login | jq -r '.token')
```

Run the server with `CACHE_ENABLED=true` and execute for example:

```bash
wrk -t2 -c20 -d5s -H "Authorization: Bearer $TOKEN" \
  -s scripts/wrk_timeline.lua http://localhost:8080/timeline
```

Switch `CACHE_ENABLED=false`, restart the server and run the command again to
compare uncached performance. Repeat for `wrk_trending.lua` and
`wrk_profile.lua`.

Sample numbers on a small dataset:

| Endpoint     | Cache | Requests/sec | Avg Latency |
|--------------|:----:|-------------:|------------:|
| `/timeline`  |  on  | 28,420       | 1.40ms      |
| `/timeline`  | off  | 30,200       | 1.01ms      |
| `/trending`  |  on  | 46,955       | 1.47ms      |
| `/trending`  | off  | 50,170       | 0.77ms      |
| `/profile/1` |  on  | 63,431       | 0.61ms      |
| `/profile/1` | off  | 49,672       | 0.94ms      |

Caching boosted profile lookups while timeline and trending showed little
difference on this dataset.

---


