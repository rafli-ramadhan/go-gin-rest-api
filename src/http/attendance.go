package http

type GetAttendance struct {
	AccountID  string  `json:"account_id"`
	LocationID string  `json:"location_id"`
	Status     string  `json:"status"`
}

type AddAttendance struct {
	LocationID int     `json:"location_id" validate:"required"`
	Status     string  `json:"status" validate:"required"`
}
