package tool

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(claims map[string] any, jwtSecret string, expiration time.Duration) (token string, err error) {
	now := time.Now()
	claims["iat"]  = now.Unix()                 // 签发时间
	claims["exp"]  = now.Add(expiration).Unix() // 过期时间
	claims["nbf"]  = now.Unix()                 // 生效时间

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims)).SignedString([] byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseJwtToken(token string, jwtSecret string) (claims map[string] any, err error) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return [] byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tokenClaims.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token claims is not map claims")
	}
	return claims, nil
}

