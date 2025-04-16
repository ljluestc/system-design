package routing

type GeoRouter struct{}

func NewGeoRouter() *GeoRouter {
    return &GeoRouter{}
}

func (gr *GeoRouter) FindNearest(location string) int {
    return 0 // Mock geo-based routing
}