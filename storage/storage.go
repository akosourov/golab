package storage

import (
	"errors"
	"math"
)

type (
	Location struct {
		Lat float64
		Lon float64
	}
	Driver struct {
		ID int
		LastLocation Location
	}
)

type DriverStorage struct {
	drivers map[int]*Driver
}

func New() *DriverStorage {
	return &DriverStorage{
		drivers: make(map[int]*Driver),
	}
}

func (s *DriverStorage) Set(key int, driver *Driver) {
	s.drivers[key] = driver
}

func (s *DriverStorage) Get(key int) (*Driver, error) {
	driver, ok := s.drivers[key]
	if !ok {
		return nil, errors.New("Driver does not exist")
	}
	return driver, nil
}

func (s *DriverStorage) Delete(key int) error {
	_, ok := s.drivers[key]
	if !ok {
		return errors.New("Driver does not exist")
	}
	delete(s.drivers, key)
	return nil
}

func (s *DriverStorage) Len() int {
	return len(s.drivers)
}

func (s *DriverStorage) Nearest(radius, lat, lon float64) []*Driver {
	nearest := []*Driver{}
	for _, driver := range s.drivers {
		dist := Distance(lat, lon, driver.LastLocation.Lat, driver.LastLocation.Lon)
		if dist <= radius {
			nearest = append(nearest, driver)
		}
	}
	return nearest
}


// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180
	r = 6378100 // Earth radius in METERS
	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)
	return 2 * r * math.Asin(math.Sqrt(h))
}
