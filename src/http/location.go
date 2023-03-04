package http

type GetLocation struct {
	ID           int    `json:"id"`
	LocationName string `json:"location_name"`
	Address      string `json:"address"`
}

type CreateLocation struct {
    LocationName string `json:"location_name" validate:"required"`
	Address      string `json:"address" validate:"required"`
}

type UpdateLocation struct {
    LocationName string `json:"location_name"`
	Address      string `json:"address"`
}
