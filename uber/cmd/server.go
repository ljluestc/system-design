package main

import (
    "log"
    "net/http"
    "uber/internal/api"
    "uber/internal/auth"
    "uber/internal/location"
    "uber/internal/ride"
    "uber/internal/user"
    "uber/internal/payment"
    "uber/internal/notification"
    "uber/internal/utils"
    "uber/pkg/database"
)

func main() {
    // Initialize configuration
    config := utils.LoadConfig()

    // Initialize logger
    logger := utils.NewLogger()

    // Initialize database
    db, err := database.NewDB(config.DBPath)
    if err != nil {
        logger.Fatal("Failed to initialize database:", err)
    }

    // Initialize services
    userSvc := user.NewUserService(db, logger)
    locationSvc := location.NewLocationService(db, logger)
    rideSvc := ride.NewRideService(userSvc, locationSvc, db, logger)
    paymentSvc := payment.NewPaymentService(db, logger)
    authSvc := auth.NewAuthService(userSvc, config.JWTSecret, logger)
    notificationSvc := notification.NewNotificationService(logger)

    // Setup API
    router := api.NewRouter(rideSvc, userSvc, locationSvc, paymentSvc, authSvc, notificationSvc, logger)
    http.Handle("/", router)

    // Start server
    logger.Info("Server starting on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        logger.Fatal("Server failed:", err)
    }
}