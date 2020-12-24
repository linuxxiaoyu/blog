package setting

import (
	"log"

	"github.com/go-ini/ini"
)

var (
	cfg *ini.File
)

// Init settings from conf/server.ini
func Init() {
	var err error
	cfg, err = ini.Load("../../configs/server.ini")
	if err != nil {
		log.Fatalf("Fail to parse '../../configs/server.ini': %v", err)
	}

	initServer()
	initDB()
	initJWT()
	initCache()
}
