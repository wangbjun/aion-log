package model

import (
	"aion/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
)

var dbConnections = make(map[string]*gorm.DB)

func Init() {
	for k, v := range config.DBConfig {
		db, err := openConnection(v)
		if err != nil {
			log.Fatalf("init mysql pool [%s] failed，error： %s\n", k, err.Error())
		} else {
			dbConnections[k] = db
		}
	}
}

func DB() *gorm.DB {
	return GetDB("default")
}

func GetDB(name string) *gorm.DB {
	conn, ok := dbConnections[name]
	if !ok {
		return nil
	}
	return conn
}

func openConnection(conf map[string]string) (*gorm.DB, error) {
	db, err := gorm.Open(conf["dialect"], conf["dsn"])
	if err != nil {
		log.Fatalf("open connection failed,error: %s", err.Error())
		return nil, err
	}
	idle, _ := strconv.Atoi(conf["maxIdleConns"])
	open, _ := strconv.Atoi(conf["maxOpenConns"])
	db.DB().SetMaxIdleConns(idle)
	db.DB().SetMaxOpenConns(open)
	if config.GetAPP("DEBUG").String() == "true" {
		db.LogMode(true)
	}
	return db, nil
}
