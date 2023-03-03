package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHarsh(password string) (hashedPassword string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	hashedPassword = string(bytes)
	return hashedPassword, nil
}
func CheckPasswordHarsh(hashedPassword string, password string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return
}