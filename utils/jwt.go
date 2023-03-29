package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"plumelog/config"
	"plumelog/log"
	"time"
)

type JwtCustomClaims struct {
	UserName string
	Password string
	*jwt.RegisteredClaims
}

func GenerateToken(username, password string) string {
	claims := JwtCustomClaims{
		UserName: username,
		Password: password,
	}
	registeredClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Conf.ExpireTime) * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   "token",
	}
	claims.RegisteredClaims = &registeredClaims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte("jwt_token"))
	if err != nil {
		log.Error.Println(err)
	}
	return signedString
}

func ParseToken(tokenStr string) (*JwtCustomClaims, error) {
	customClaims := JwtCustomClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, &customClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte("jwt_token"), nil
	})
	if err != nil {
		log.Error.Println("parse token err:", err)
		return nil, err
	}
	return &customClaims, nil
}

func IsInvalidToken(tokenStr string) bool {
	claims, err := ParseToken(tokenStr)
	if err != nil || claims.UserName != config.Conf.Gin.Username || claims.Password != config.Conf.Gin.Password {
		return false
	}
	return true
}
