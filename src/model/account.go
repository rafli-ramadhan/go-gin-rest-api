package model

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	gorm.Model
	Username          string    `gorm:"column:username;type:varchar(50)"`
	FullName          string    `gorm:"column:full_name;type:varchar(150)"`
	Email             *string   `gorm:"column:email;type:varchar(150)"`
	Password          string    `gorm:"column:password;type:varchar(64)"`
	Address           *string   `gorm:"column:address;type:varchar(50)"`
	EmployeeNumber    *string   `gorm:"column:employee_number;type:varchar(50)"`
	JobPosition       *string   `gorm:"column:job_position;type:varchar(50)"`
	KTPNumber         *string   `gorm:"column:ktp_number;type:varchar(50)"`
	PhoneNumber       *string   `gorm:"column:phone_number;type:varchar(20)"`
	PhotoURL          string    `gorm:"column:photo_url;type:varchar(200)"`
	Gender            string    `gorm:"column:gender"`
	DateOfBirth       time.Time `gorm:"column:date_of_birth;type:date"`
	IsVerified        bool      `gorm:"column:is_verified;type:bool"`
}

func (Account) TableName() string {
	return "accounts"
}
