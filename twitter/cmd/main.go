package main

import (
    "fmt"
    "net/http"
    "twitter/internal/config"
    "twitter/internal/logger"
    "twitter/internal/notification"
    "twitter/internal/search"
    "twitter/internal/timeline"
    "twitter/internal/tweet"
    "twitter/internal/user"
    "twitter/pkg/cache"
    "twitter/pkg/loadbalancer"
    "twitter/pkg/messaging"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        fmt.Println("Failed to load config:", err)
        return
    }

    // Initialize logger
    log := logger.NewLogger()

    // Initialize cache and message queue
    cacheInstance := cache.NewCache()
    queue := messaging.NewMessageQueue()

    // Initialize database (mock for simplicity)
    // In a real system, this would connect to SQL/NoSQL databases
    db := pkgdatabase.NewMockDB()

    // Initialize service repositories
    userRepo := user.NewUserRepository(db)
    tweetRepo := tweet.NewTweetRepository(db)
    timelineRepo := timeline.NewTimelineRepository()
    searchRepo := search.NewSearchRepository()
    notificationRepo := notification.NewNotificationRepository()

    // Initialize services with dependencies
    userService := user.NewUserService(userRepo, log, cacheInstance)
    tweetService := tweet.NewTweetService(tweetRepo, log, cacheInstance, queue)
    timelineService := timeline.NewTimelineService(timelineRepo, log, cacheInstance, queue)
    searchService := search.NewSearchService(searchRepo, log, cacheInstance)
    notificationService := notification.NewNotificationService(notificationRepo, log, queue)

    // Initialize load balancer with service instances
    lb := loadbalancer.NewLoadBalancer(cfg.ServiceInstances)

    // Start health checks in a goroutine
    go lb.StartHealthChecks()

    // Initialize controllers with services and load balancer
    userController := user.NewUserController(userService, lb)
    tweetController := tweet.NewTweetController(tweetService, lb)
    timelineController := timeline.NewTimelineController(timelineService, lb)
    searchController := search.NewSearchController(searchService, lb)
    notificationController := notification.NewNotificationController(notificationService, lb)

    // Register HTTP handlers
    http.HandleFunc("/users", userController.HandleUserRequests)
    http.HandleFunc("/tweets", tweetController.HandleTweetRequests)
    http.HandleFunc("/timeline", timelineController.HandleTimelineRequests)
    http.HandleFunc("/search", searchController.HandleSearchRequests)
    http.HandleFunc("/notifications", notificationController.HandleNotificationRequests)

    // Start the server
    fmt.Println("Server starting on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Error("Server failed to start:", err)
    }
}