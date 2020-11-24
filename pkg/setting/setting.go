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
	cfg, err = ini.Load("conf/server.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/server.ini': %v", err)
	}

	initServer()
	initDB()
}
