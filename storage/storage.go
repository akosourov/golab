package storage

import (
	"errors"

	"github.com/dhconnelly/rtreego"
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
	DriverStorage struct {
		drivers   map[int]*Driver
		locations *rtreego.Rtree
	}
)

func (d *Driver) Bounds() *rtreego.Rect {
	return rtreego.Point{d.LastLocation.Lon, d.LastLocation.Lat}.ToRect(0.01)
}

func New() *DriverStorage {
	return &DriverStorage{
		drivers: make(map[int]*Driver),
		locations: rtreego.NewTree(2, 25, 50),
	}
}

func (s *DriverStorage) Set(key int, driver *Driver) {
	_, ok := s.drivers[key]
	if !ok {
		s.locations.Insert(driver)
	}
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
	d, ok := s.drivers[key]
	if !ok {
		return errors.New("Driver does not exist")
	}
	if s.locations.Delete(d) {
		delete(s.drivers, key)
		return nil
	}
	return errors.New("Could not remove")
}

func (s *DriverStorage) Len() int {
	return len(s.drivers)
}

func (s *DriverStorage) Nearest(count int, lat, lon float64) []*Driver {
	p := rtreego.Point{lat, lon}
	foundItems := s.locations.NearestNeighbors(count, p)
	nearest := []*Driver{}
	for _, item := range foundItems {
		if item == nil {
			continue
		}
		nearest = append(nearest, item.(*Driver))
	}
	return nearest
}

