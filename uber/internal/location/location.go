package location

import (
    "uber/internal/utils"
    "uber/pkg/database"
)

type Location struct {
    DriverID string
    Lat      float64
    Lng      float64
}

type LocationService struct {
    db     *database.DB
    logger *utils.Logger
}

func NewLocationService(db *database.DB, logger *utils.Logger) *LocationService {
    return &LocationService{db: db, logger: logger}
}

func (ls *LocationService) UpdateDriverLocation(driverID string, lat, lng float64) (string, error) {
    location := Location{DriverID: driverID, Lat: lat, Lng: lng}
    err := ls.db.SaveLocation(location)
    if err != nil {
        ls.logger.Error("Failed to update location:", err)
        return "", err
    }
    ls.logger.Info("Driver location updated:", driverID)
    return fmt.Sprintf("Driver %s location updated to (%f, %f)", driverID, lat, lng), nil
}

func (ls *LocationService) FindNearestDriver(lat, lng float64) (string, error) {
    locations, err := ls.db.GetAllLocations()
    if err != nil {
        ls.logger.Error("Failed to get locations:", err)
        return "", err
    }
    if len(locations) == 0 {
        return "", fmt.Errorf("no drivers available")
    }
    // Simplified: return first driver; real system uses distance calculation
    return locations[0].DriverID, nil
}