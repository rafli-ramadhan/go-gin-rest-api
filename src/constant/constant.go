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
)

var (
	// db connection
	_ = godotenv.Load()
	ServiceName = os.Getenv("SERVICE_NAME")

	// jwt
	SampleSecretKey = []byte("SecretYouShouldHide")

	// error	
	ErrInvalidID               = errors.New("invalid id")
	ErrInvalidFormat           = errors.New("invalid format")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrAccountExist            = errors.New("account already exist")
	ErrAccountNotRegistered    = errors.New("account not registered")
	ErrEmailAlreadyExist       = errors.New("email already exist")
	ErrKTPNumberAlreadyExist   = errors.New("ktp number already exist")
	ErrPasswordCannotBeEmpty   = errors.New("password cannot be empty")
	ErrPhoneNumberAlreadyExist = errors.New("phone number already exist")
	ErrUsernameAlreadyExist    = errors.New("username already exist")
)
