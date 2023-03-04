package model

import (
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	LocationName string `gorm:"column:name;type:varchar(60)"`
	Address      string `gorm:"column:address;type:varchar(500)"`
	PhotoURL     string `gorm:"column:photo_url;type:varchar(500)"`
}

func (Location) TableName() string {
	return "locations"
}
