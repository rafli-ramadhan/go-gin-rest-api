package http

type GetLocation struct {
	ID           int    `json:"id"`
	LocationName string `json:"location_name"`
	Address      string `json:"address"`
	PhotoURL     string `json:"photo_url"`
}

type CreateLocation struct {
    LocationName string `json:"location_name" validate:"required"`
	Address      string `json:"address" validate:"required"`
	PhotoURL     string `json:"photo_url"`
}

type UpdateLocation struct {
    LocationName string `json:"location_name"`
	Address      string `json:"address"`
}
