package util

import (
	"kn-assignment/property"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func GenerateAccessToken(id, username, role string) (string, error) {
	expirationTime := time.Now().Add(property.Get().Server.AccessTokenExpiry)
	claims := &Claims{
		Id:       id,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(property.Get().Secret.JWTSecretKey))
}

func GenerateRefreshToken(id, username, role string) (string, error) {
	expirationTime := time.Now().Add(property.Get().Server.RefreshTokenExpiry)
	claims := &Claims{
		Id:       id,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(property.Get().Secret.JWTSecretKey))
}

func ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(property.Get().Secret.JWTSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}
