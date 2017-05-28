package storage

import (
	"errors"
	"time"
	"sync"

	"github.com/dhconnelly/rtreego"
	"github.com/akosourov/golab/storage/lru"
)

type (
	Location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}
	Driver struct {
		ID           int      `json:"id"`
		LastLocation Location `json:"location"`
		Expiration   int64    `json:"-"`
		Locations    *lru.LRU `json:"-"`
	}
	DriverStorage struct {
		mu        *sync.RWMutex
		drivers   map[int]*Driver
		locations *rtreego.Rtree
		lruSize   int
	}
)

func New(lruSize int) *DriverStorage {
	return &DriverStorage{
		mu: new(sync.RWMutex),
		drivers: make(map[int]*Driver),
		locations: rtreego.NewTree(2, 25, 50),
		lruSize: lruSize,
	}
}

func (d *Driver) Expire() bool {
	if d.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > d.Expiration
}

func (d *Driver) Bounds() *rtreego.Rect {
	return rtreego.Point{d.LastLocation.Lon, d.LastLocation.Lat}.ToRect(0.01)
}

func (s *DriverStorage) Set(driver *Driver) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	d, ok := s.drivers[driver.ID]
	if !ok {
		cache, err := lru.New(s.lruSize)
		if err != nil {
			return err
		}
		driver.Locations = cache
		s.locations.Insert(driver)
		d = driver
	}
	d.Expiration = time.Now().UnixNano()
	d.LastLocation = driver.LastLocation
	d.Locations.Add(time.Now().UnixNano(), driver.LastLocation)

	s.drivers[driver.ID] = driver
	return nil
}

func (s *DriverStorage) Get(key int) (*Driver, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	driver, ok := s.drivers[key]
	if !ok {
		return nil, errors.New("Driver does not exist")
	}
	return driver, nil
}

func (s *DriverStorage) Delete(key int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

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
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.drivers)
}

// Nearest возвращает ближайших водителей к указанным координатам.
// count - максимальное количество найденных
func (s *DriverStorage) Nearest(point rtreego.Point, count int) []*Driver {
	s.mu.Lock()
	defer s.mu.Unlock()

	foundItems := s.locations.NearestNeighbors(count, point)
	nearest := []*Driver{}
	for _, item := range foundItems {
		if item == nil {
			continue
		}
		nearest = append(nearest, item.(*Driver))
	}
	return nearest
}

// DeleteExpired удаляет всех водителей из базы, валидность данных
// которых истекла
func (s *DriverStorage) DeleteExpired() {
	now := time.Now().UnixNano()
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, d := range s.drivers {
		if d.Expiration > 0 && now > d.Expiration {
			if s.locations.Delete(d) {
				delete(s.drivers, k)
			}
		}
	}
}