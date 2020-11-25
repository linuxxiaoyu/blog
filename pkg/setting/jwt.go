package setting

import (
	"log"
)

var (
	signingKey string
)

func initJWT() {
	jwt, err := cfg.GetSection("jwt")
	if err != nil {
		log.Fatalf("Fail to get section 'jwt': %v", err)
	}

	signingKey = jwt.Key("SigningKey").MustString("yourSigningKey")
}

// SigningKey is a conf
func SigningKey() string {
	return signingKey
}
