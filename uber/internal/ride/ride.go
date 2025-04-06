package ride

import (
    "uber/internal/user"
    "uber/internal/location"
    "uber/internal/utils"
    "uber/pkg/database"
)

type Ride struct {
    ID              string
    RiderID         string
    DriverID        string
    PickupLat       float64
    PickupLng       float64
    DropoffLat      float64
    DropoffLng      float64
    Status          string
}

type RideService struct {
    userSvc     *user.UserService
    locationSvc *location.LocationService
    db          *database.DB
    logger      *utils.Logger
}

func NewRideService(us *user.UserService, ls *location.LocationService, db *database.DB, logger *utils.Logger) *RideService {
    return &RideService{userSvc: us, locationSvc: ls, db: db, logger: logger}
}

func (rs *RideService) RequestRide(riderID string, pickupLat, pickupLng, dropoffLat, dropoffLng float64) (string, error) {
    _, err := rs.userSvc.GetUser(riderID)
    if err != nil {
        rs.logger.Error("Rider not found:", err)
        return "", fmt.Errorf("rider not found")
    }
    driverID, err := rs.locationSvc.FindNearestDriver(pickupLat, pickupLng)
    if err != nil {
        rs.logger.Error("No drivers available:", err)
        return "", err
    }
    rideID := fmt.Sprintf("ride_%d", time.Now().UnixNano())
    ride := Ride{
        ID:             rideID,
        RiderID:        riderID,
        DriverID:       driverID,
        PickupLat:      pickupLat,
        PickupLng:      pickupLng,
        DropoffLat:     dropoffLat,
        DropoffLng:     dropoffLng,
        Status:         "requested",
    }
    err = rs.db.SaveRide(ride)
    if err != nil {
        rs.logger.Error("Failed to save ride:", err)
        return "", err
    }
    rs.logger.Info("Ride requested:", rideID)
    return fmt.Sprintf("Ride %s requested. Driver %s assigned.", rideID, driverID), nil
}