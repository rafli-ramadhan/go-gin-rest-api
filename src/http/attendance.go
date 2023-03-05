package http

type GetAttendance struct {
	LocationID   int     `json:"location_id"`
	LocationName string  `json:"location_name"`
	Status       string  `json:"status"`
	Time         string  `json:"time"`
	Description  string  `json:"description"`
}

type GetAttendanceByLocation struct {
	LocationID   int     `json:"location_id"`
	LocationName string  `json:"location_name"`
	Status       string  `json:"status"`
	Address      string  `json:"address"`
}

type AddAttendance struct {
	LocationID int     `json:"location_id" validate:"required"`
	Status     string  `json:"status" validate:"required"`
}
