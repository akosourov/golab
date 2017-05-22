package storage

import "errors"

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

func (ds *DriverStorage) Set(key int, driver Driver) {
	return
}

func (ds *DriverStorage) Get(key int) (*Driver, error) {
	driver, ok := ds.drivers[key]
	if !ok {
		return nil, errors.New("Driver does not exist")
	}
	return driver, nil
}

func (ds *DriverStorage) Delete(key int) error {
	_, ok := ds.drivers[key]
	if !ok {
		return errors.New("Driver does not exist")
	}
	delete(ds.drivers, key)
	return nil
}

func (ds *DriverStorage) Nearest() []*Driver {
	return nil
}