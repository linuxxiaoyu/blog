package setting

import (
	"fmt"
	"log"
	"strings"
)

var (
	releaseMode = true
	port        uint
)

func initServer() {
	initCfg()
	server, err := cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	runMode := server.Key("RunMode").MustString("release")
	if strings.EqualFold(runMode, "debug") {
		releaseMode = false
	}

	port = server.Key("Port").MustUint(8080)
}

// IsReleaseMode return RunMode from conf/server.ini
func IsReleaseMode() bool {
	return releaseMode
}

// Port return Port from conf/server.ini [server] Port
func Port() string {
	return fmt.Sprintf(":%d", port)
}
