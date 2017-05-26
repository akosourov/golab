package api

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type API struct {
	echo     *echo.Echo
	bindAddr string
}

func New(bindAddr string) *API {
	a := &API{
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

func (a *API) addDriver(c echo.Context) error {
	p := &Payload{}
	if err := c.Bind(p); err != nil {
		return c.JSON(http.StatusUnsupportedMediaType, &DefaultResponse{
			Success: false,
			Message: "Set content-type application/json or check payload data",
		})
	}
	return c.JSON(http.StatusOK, DriverResponse{
		Success: false,
		Message: "Added",
	})
}

func (a *API) getDriver(c echo.Context) error {
	driverID := c.Param("id")
	id, err := strconv.Atoi(driverID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, DefaultResponse{
			Success: false,
			Message: "could not convert string to integer",
		})
	}
	return c.JSON(http.StatusOK, DriverResponse{
		Success: true,
		Message: "Driver was found",
		Driver:  id,
	})
}

func (a *API) deleteDriver(c echo.Context) error {
	driverID := c.Param("id")
	_, err := strconv.Atoi(driverID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, DefaultResponse{
			Success: false,
			Message: "Driver s was deleted",
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
			Message: "Empty coordinates",
		})
	}
	_, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "Bad latitude",
		})
	}
	_, err = strconv.ParseFloat(lon, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &DefaultResponse{
			Success: false,
			Message: "Bad longitude",
		})
	}
	// TODO add nearest
	return c.JSON(http.StatusOK, &NearestDriverResponse{
		Success: true,
		Message: "Nearest drivers was found",
	})
}

func (a *API) Start() error {
	return a.echo.Start(a.bindAddr)
}


