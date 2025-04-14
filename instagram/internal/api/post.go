package api

import (
    "encoding/json"
    "github.com/calelin/instagram/internal/auth"
    "github.com/calelin/instagram/internal/cache"
    "github.com/calelin/instagram/internal/db"
    "github.com/calelin/instagram/internal/queue"
    "github.com/calelin/instagram/internal/storage"
    "github.com/calelin/instagram/pkg/model"
    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
    "net/http"
)

// Handler manages API endpoints
type Handler struct {
    authService   *auth.AuthService
    postgresDB    *db.PostgresDB
    cassandraDB   *db.CassandraDB
    redisClient   *cache.RedisClient
    s3Storage     *storage.S3Storage
    kafkaProducer *queue.KafkaProducer
    log           *logrus.Logger
}

// NewHandler creates a new Handler
func NewHandler(authService *auth.AuthService, postgresDB *db.PostgresDB, cassandraDB *db.CassandraDB, redisClient *cache.RedisClient, s3Storage *storage.S3Storage, kafkaProducer *queue.KafkaProducer, log *logrus.Logger) *Handler {
    return &Handler{
        authService:   authService,
        postgresDB:    postgresDB,
        cassandraDB:   cassandraDB,
        redisClient:   redisClient,
        s3Storage:     s3Storage,
        kafkaProducer: kafkaProducer,
        log:           log,
    }
}

// PostMediaHandler handles POST /postMedia
func (h *Handler) PostMediaHandler(w http.ResponseWriter, r *http.Request) {
    var req model.PostMediaRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.log.Errorf("Invalid request: %v", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request")
        return
    }

    // Validate user
    userID := r.Context().Value("userID").(string)
    if err := h.authService.ValidateUser(userID); err != nil {
        h.log.Errorf("Unauthorized: %v", err)
        respondWithError(w, http.StatusUnauthorized, "Unauthorized")
        return
    }

    // Upload media to S3
    mediaURL, err := h.s3Storage.UploadMedia(req.MediaFile, req.MediaType, userID)
    if err != nil {
        h.log.Errorf("Failed to upload media: %v", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to upload media")
        return
    }

    // Save post metadata
    postID, err := h.postgresDB.SavePost(userID, mediaURL, req.MediaType, req.Caption, req.Hashtags)
    if err != nil {
        h.log.Errorf("Failed to save post: %v", err)
        respondWithError(w, http.StatusInternalServerError, "Failed to save post")
        return
    }

    // Update feed in Cassandra
    if err := h.cassandraDB.UpdateFeed(userID, postID, mediaURL, req.Caption); err != nil {
        h.log.Warnf("Failed to update feed: %v", err)
    }

    // Notify followers via Kafka
    h.kafkaProducer.NotifyFollowers(userID, postID)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"post_id": postID})
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
    resp := model.ErrorResponse{Error: message}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(resp)
}