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
  - Adaptive TTL for timelines (toggle via `ENABLE_ADAPTIVE_TTL`)

---

## 🧪 Benchmark Results

| API Endpoint | Mode | Requests/sec | Avg Latency | Max Latency | Comments |
|:-------------|:----:|:------------:|:-----------:|:-----------:|:---------|
| `/timeline/db` | DB Only | ~289 | ~120.92ms | ~2.0s | DB overloads under load |
| `/timeline/cache` | Redis Cache | ~32,435 | ~3.48ms | ~54ms | 110x faster with caching |
| `/trending/db` | DB Only | ~5,445 | ~38.63ms | ~281ms | Slows as data grows |
| `/trending/cache` | Redis Cache | ~45,353 | ~2.19ms | ~16ms | 8x faster with cache |
| `/profile/:id` Cold | DB Only | ~43,301 | ~4.82ms | ~1.64s | Max latency spikes without cache |
| `/profile/:id` Warm | Redis Cache | ~45,851 | ~2.17ms | ~14ms | Stable latency with caching |
| `/posts` Create | Write + Cache | ~2,616 | ~3.81ms | ~80ms | Write-Through and Fan-out impact |

✅ Caching improved read throughput by **8x–110x** and stabilized system latency dramatically.

---


