package setting

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB for CRUD
var db *gorm.DB

func InitDB() {
	initCfg()
	database, err := cfg.GetSection("database")
	if err != nil {
		log.Fatalf("Fail to get section 'database': %v", err)
	}

	userName := database.Key("User").MustString("root")
	password := database.Key("Password").MustString("")
	host := database.Key("Host").MustString("127.0.0.1")
	port := database.Key("Port").MustUint(3306)
	dbName := database.Key("DBName").MustString("")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local",
		userName,
		password,
		host,
		port,
		dbName,
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Fail to open database: %v", err)
	}
}

func DB() *gorm.DB {
	if db == nil {
		InitDB()
	}
	return db
}
