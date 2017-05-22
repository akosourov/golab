package api

// Структуры для получения запросов
type (
	Location struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
	}
	Payload struct {
		Timestamp int64    `json:"timestamp"`
		DriverID  int      `json:"driver_id"`
		//Location  Location `json:"location"`
	}
)

// Структуры для формирования ответов
type (
	DefaultResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	DriverResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Driver  int    `json:"driver"`
	}
	NearestDriverResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Drivers []int  `json:"drivers"`
	}
)

