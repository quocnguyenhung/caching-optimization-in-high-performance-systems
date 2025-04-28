package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/config"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/handler"
	"github.com/quocnguyenhung/caching-optimization-in-high-performance-systems/internal/middleware"
)

func main() {
	// Load environment
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	err = config.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	config.ConnectRedis()

	// Test Redis
	ctx := context.Background()
	err = config.RedisClient.Set(ctx, "test-key", "hello world", 0).Err()
	if err != nil {
		log.Fatalf("Failed to set test key in Redis: %v", err)
	}
	val, err := config.RedisClient.Get(ctx, "test-key").Result()
	if err != nil {
		log.Fatalf("Failed to get test key from Redis: %v", err)
	}
	fmt.Println("Test Redis Key:", val)

	// Public routes
	http.HandleFunc("/signup", handler.Signup)
	http.HandleFunc("/login", handler.Login)

	// Protected routes
	http.Handle("/posts", middleware.AuthMiddleware(http.HandlerFunc(handler.CreatePost)))
	http.Handle("/posts/", middleware.AuthMiddleware(http.HandlerFunc(handler.LikePost)))
	http.Handle("/follow/", middleware.AuthMiddleware(http.HandlerFunc(handler.FollowUser)))
	http.Handle("/timeline/db", middleware.AuthMiddleware(http.HandlerFunc(handler.GetTimelineFromDB)))
	http.Handle("/timeline/cache", middleware.AuthMiddleware(http.HandlerFunc(handler.GetTimelineWithCache)))
	http.Handle("/trending/cache", http.HandlerFunc(handler.GetTrendingWithCache))
	http.Handle("/trending/db", http.HandlerFunc(handler.GetTrendingFromDB))
	http.Handle("/profile/", http.HandlerFunc(handler.GetProfile))

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
