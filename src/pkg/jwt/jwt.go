package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/forkyid/go-utils/v1/aes"
	"github.com/golang-jwt/jwt"
	"go-rest-api/src/constant"
)

func GenerateJWT(accountID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["accountID"] = accountID
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(constant.SampleSecretKey)
	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(bearerToken string) (claims jwt.MapClaims, err error) {
	tokenStr := strings.Replace(bearerToken, "Bearer ", "", -1)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return constant.SampleSecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if token == nil {
		err = errors.New("token error")
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("invalid JWT Token")
		return nil, err
	}

	exp := claims["exp"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return claims, err
}

func ExtractID(bearerToken string) (int, error) {
	claimsMap, err := ValidateToken(bearerToken)
	if err != nil {
		return -1, fmt.Errorf("failed on claiming token")
	}
	id := aes.Decrypt(claimsMap["accountID"].(string))
	if id == -1 {
		return -1, fmt.Errorf("invalid ID")
	}
	return id, nil
}