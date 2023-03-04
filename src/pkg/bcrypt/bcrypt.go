package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(str string) (hashedStr string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), 14)
	if err != nil {
		return "", err
	}
	hashedStr = string(bytes)
	return hashedStr, nil
}
func CheckHash(hashedStr string, str string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(str))
	return
}