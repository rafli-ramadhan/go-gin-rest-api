package constant

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

const (
	DBServerMaster             = "master"
	ContentTypeApplicationJson = "application/json"
	DOBLayout                  = `01/01/2023`
)

var (
	// jwt
	SampleSecretKey = []byte("SecretYouShouldHide")

	// error	
	ErrInvalidID               = errors.New("invalid id")
	ErrInvalidFormat           = errors.New("invalid format")
	ErrAccountExist            = errors.New("user already exist")
	ErrAccountNotRegistered    = errors.New("user not registered")
	ErrEmailAlreadyExist       = errors.New("email already exist")
	ErrKTPNumberAlreadyExist   = errors.New("ktp number already exist")
	ErrPhoneNumberAlreadyExist = errors.New("phone number already exist")
	ErrUsernameAlreadyExist    = errors.New("username already exist")

	// db connection
	_ = godotenv.Load()
	ServiceName = os.Getenv("SERVICE_NAME")
)
