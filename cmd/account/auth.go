package main

import (
	"time"

	"github.com/linuxxiaoyu/blog/internal/setting"

	"github.com/dgrijalva/jwt-go"
)

type userClaims struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

func generateToken(uid uint32, username string) (string, error) {
	claims := userClaims{
		uid,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
			Issuer:    "blog",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(setting.SigningKey()))
}

func parseToken(tokenString string) (*userClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(setting.SigningKey()), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*userClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
