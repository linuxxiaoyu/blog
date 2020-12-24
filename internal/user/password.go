package user

import (
	"encoding/base64"
	"log"

	"golang.org/x/crypto/scrypt"
)

var (
	salt = []byte{
		0x11, 0x22, 0x33, 0x44, 0x33, 0x77, 0x0d, 0x0a,
	}
)

func addSalt(password string) string {
	dk, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(dk)
}
