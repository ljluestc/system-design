package api

import (
    "encoding/json"
    "net/http"
    "uber/internal/auth"
    "uber/internal/location"
    "uber/internal/ride"
    "uber/internal/user"
    "uber/internal/payment"
    "uber/internal/notification"
    "uber/internal/utils"
    "golang.org/x/crypto/bcrypt"
)

// Handlers contains all API handlers
type Handlers struct {
    rideSvc         *ride.RideService
    userSvc         *user.UserService
    locationSvc     *location.LocationService
    paymentSvc      *payment.PaymentService
    authSvc         *auth.AuthService
    notificationSvc *notification.NotificationService
    logger          *utils.Logger
}

// NewHandlers initializes the handlers with dependency injection
func NewHandlers(rs *ride.RideService, us *user.UserService, ls *location.LocationService, ps *payment.PaymentService, as *auth.AuthService, ns *notification.NotificationService, logger *utils.Logger) *Handlers {
    return &Handlers{
        rideSvc:         rs,
        userSvc:         us,
        locationSvc:     ls,
        paymentSvc:      ps,
        authSvc:         as,
        notificationSvc: ns,
        logger:          logger,
    }
}

// RegisterUserHandler handles user registration
func (h *Handlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        ID       string `json:"id"`
        Name     string `json:"name"`
        Email    string `json:"email"`
        Role     string `json:"role"`
        Password string `json:"password"` // Added missing password field
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Error("Failed to decode request:", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Validate input
    if req.ID == "" || req.Name == "" || req.Email == "" || req.Role == "" || req.Password == "" {
        h.logger.Error("Missing required fields")
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    // Check if user already exists
    _, err := h.userSvc.GetUser(req.ID)
    if err == nil {
        h.logger.Error("User already exists:", req.ID)
        http.Error(w, "User already exists", http.StatusConflict)
        return
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        h.logger.Error("Failed to hash password:", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Register user
    resp, err := h.userSvc.RegisterUser(req.ID, req.Name, req.Email, req.Role, string(hashedPassword))
    if err != nil {
        h.logger.Error("Failed to register user:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(resp))
}

// RequestRideHandler handles ride requests
func (h *Handlers) RequestRideHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        RiderID      string  `json:"rider_id"`
        PickupLat    float64 `json:"pickup_lat"`
        PickupLng    float64 `json:"pickup_lng"`
        DropoffLat   float64 `json:"dropoff_lat"`
        DropoffLng   float64 `json:"dropoff_lng"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Error("Failed to decode ride request:", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Authenticate user
    userID := r.Header.Get("UserID")
    if userID != req.RiderID {
        h.logger.Error("Unauthorized ride request")
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    resp, err := h.rideSvc.RequestRide(req.RiderID, req.PickupLat, req.PickupLng, req.DropoffLat, req.DropoffLng)
    if err != nil {
        h.logger.Error("Failed to request ride:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte(resp))
}

// AcceptRideHandler handles ride acceptance
func (h *Handlers) AcceptRideHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        RideID   string `json:"ride_id"`
        DriverID string `json:"driver_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Error("Failed to decode accept request:", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Authenticate driver
    userID := r.Header.Get("UserID")
    if userID != req.DriverID {
        h.logger.Error("Unauthorized ride acceptance")
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    resp, err := h.rideSvc.AcceptRide(req.RideID, req.DriverID)
    if err != nil {
        h.logger.Error("Failed to accept ride:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte(resp))
}

// CompleteRideHandler handles ride completion
func (h *Handlers) CompleteRideHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        RideID string `json:"ride_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Error("Failed to decode complete request:", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Authenticate driver
    userID := r.Header.Get("UserID")
    ride, err := h.rideSvc.GetRide(req.RideID)
    if err != nil || ride.DriverID != userID {
        h.logger.Error("Unauthorized ride completion")
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    resp, err := h.rideSvc.CompleteRide(req.RideID)
    if err != nil {
        h.logger.Error("Failed to complete ride:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte(resp))
}

// UpdateLocationHandler handles driver location updates
func (h *Handlers) UpdateLocationHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        DriverID string  `json:"driver_id"`
        Lat      float64 `json:"lat"`
        Lng      float64 `json:"lng"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Error("Failed to decode location update:", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Authenticate driver
    userID := r.Header.Get("UserID")
    if userID != req.DriverID {
        h.logger.Error("Unauthorized location update")
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    resp, err := h.locationSvc.UpdateDriverLocation(req.DriverID, req.Lat, req.Lng)
    if err != nil {
        h.logger.Error("Failed to update location:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte(resp))
}

// LoginHandler handles user login
func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Error("Failed to decode login request:", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    token, err := h.authSvc.Login(req.Email, req.Password)
    if err != nil {
        h.logger.Error("Login failed:", err)
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    w.Write([]byte(token))
}