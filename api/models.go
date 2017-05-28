package api

import "github.com/akosourov/golab/storage"

// Структуры для получения запросов
type (
	Location struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
	}
	Payload struct {
		Timestamp int64    `json:"timestamp"`
		DriverID  int      `json:"driver_id"`
		Location  Location `json:"location"`
	}
)

// Структуры для формирования ответов
type (
	DefaultResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	DriverResponse struct {
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Driver  *storage.Driver `json:"driver"`
	}
	NearestDriverResponse struct {
		Success bool              `json:"success"`
		Message string            `json:"message"`
		Drivers []*storage.Driver `json:"drivers"`
	}
)

