package constant

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

const (
	ContentTypeApplicationJson = "application/json"
	DOBLayout                  = `01/01/2023`
	DBServerMaster             = "master"
	FilterByDay                = "day"
	FilterByWeek               = "week"
	FilterByMonth              = "month"
	StatusCheckIn              = "check-in"
	StatusCheckOut             = "check-out"
)

var (
	// db connection
	_ = godotenv.Load()
	ServiceName = os.Getenv("SERVICE_NAME")

	// jwt
	SampleSecretKey = []byte("SecretYouShouldHide")

	// error
	ErrInvalidAddress           = errors.New("invalid address")
	ErrInvalidID                = errors.New("invalid id")
	ErrInvalidFormat            = errors.New("invalid format")
	ErrInvalidLocationName      = errors.New("invalid location")
	ErrInvalidPassword          = errors.New("invalid password")
	ErrAccountExist             = errors.New("account already exist")
	ErrAccountNotRegistered     = errors.New("account not registered")
	ErrEmailAlreadyExist        = errors.New("email already exist")
	ErrLocationAlreadyExist     = errors.New("location already exist")
	ErrLocationNameAlreadyExist = errors.New("location name already exist")
	ErrLocationNotExist         = errors.New("location is not exist")
	ErrKTPNumberAlreadyExist    = errors.New("ktp number already exist")
	ErrPasswordCannotBeEmpty    = errors.New("password cannot be empty")
	ErrPhoneNumberAlreadyExist  = errors.New("phone number already exist")
	ErrUsernameAlreadyExist     = errors.New("username already exist")
)
