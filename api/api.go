package api

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/akosourov/golab/storage"
	"github.com/labstack/echo"
	"github.com/dhconnelly/rtreego"
)

const MAX_NEAREST_DRIVERS = 10

type API struct {
	database  *storage.DriverStorage
	waitGroup sync.WaitGroup
	echo      *echo.Echo
	bindAddr  string
}

func New(bindAddr string, lruSize int) *API {
	a := &API{
		database: storage.New(lruSize),
		echo: echo.New(),
		bindAddr: bindAddr,
	}
	g := a.echo.Group("/api")
	g.POST("/driver/", a.addDriver)
	g.GET("/driver/:id", a.getDriver)
	g.DELETE("/driver/:id", a.deleteDriver)
	g.GET("/driver/:lat/:lon/nearest", a.nearestDrivers)
	return a
}


func (a *API) Start() {
	a.waitGroup.Add(1)
	go func() {
		a.echo.Start(a.bindAddr)
		a.waitGroup.Done()
	}()
	go a.deleteExpired()
}

func (a *API) WaitStop() {
	a.waitGroup.Wait()
}

func (a *API) deleteExpired() {
	for range time.Tick(1) {
		a.database.DeleteExpired()
	}
}

func (a *API) addDriver(c echo.Context) error {
	p := &Payload{}
	if err := c.Bind(p); err != nil {
		return c.JSON(http.StatusUnsupportedMediaType, &DefaultResponse{
			Success: false,
			Message: "Set content-type application/json or check payload data",
		})
	}
	driver := &storage.Driver{
		ID: p.DriverID,
		LastLocation: storage.Location{
			Lat: p.Location.Latitude,
			Lon: p.Location.Longitude,
		},
	}
	if err := a.database.Set(driver); err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: fmt.Sprintf("Could not add driver - %s", err.Error()),
		})
	}
	return c.JSON(http.StatusOK, DriverResponse{
		Success: true,
		Message: "Driver was added",
	})
}

func (a *API) getDriver(c echo.Context) error {
	driverID := c.Param("id")
	id, err := strconv.Atoi(driverID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "Could not convert string to integer",
		})
	}
	driver, err := a.database.Get(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: fmt.Sprintf("Could not get driver - %s", err.Error()),
		})
	}
	return c.JSON(http.StatusOK, &DriverResponse{
		Success: true,
		Message: "Driver was found",
		Driver:  driver,
	})
}

func (a *API) deleteDriver(c echo.Context) error {
	driverID := c.Param("id")
	id, err := strconv.Atoi(driverID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, DefaultResponse{
			Success: false,
			Message: "Could not convert string to integer",
		})
	}
	if err := a.database.Delete(id); err != nil {
		return c.JSON(http.StatusBadRequest, DefaultResponse{
			Success: false,
			Message: fmt.Sprintf("Could not delete driver - %s", err.Error()),
		})
	}
	return c.JSON(http.StatusOK, &DefaultResponse{
		Success: true,
		Message: "Driver was deleted",
	})
}

func (a *API) nearestDrivers(c echo.Context) error {
	lat := c.Param("lat")
	lon := c.Param("lon")
	if lat == "" || lon == "" {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "Bad coordinates",
		})
	}
	lt, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "Bad latitude",
		})
	}
	ln, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "Bad longitude",
		})
	}
	drivers := a.database.Nearest(rtreego.Point{lt, ln}, MAX_NEAREST_DRIVERS)
	return c.JSON(http.StatusOK, &NearestDriverResponse{
		Success: true,
		Message: "Nearest drivers was found",
		Drivers: drivers,
	})
}



