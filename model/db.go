package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

var defaultDB *gorm.DB

func Init(initTable bool) {
	config := &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 500 * time.Millisecond,
			LogLevel:      logger.Warn,
			Colorful:      true,
		}),
	}
	dbFile := "aion.db"
	if initTable {
		err := os.Remove(dbFile)
		if err != nil && !os.IsNotExist(err) {
			panic("failed to remove old database file: " + err.Error())
		}
	}
	db, err := gorm.Open(sqlite.Open("aion.db"), config)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	if initTable {
		defaultDB = db
		err = createTables(db)
		if err != nil {
			panic("failed to create tables: " + err.Error())
		}
		err = importSkill()
		if err != nil {
			panic("failed to import skill: " + err.Error())
		}
	} else {
		defaultDB = db.Debug()
	}
}

func DB() *gorm.DB {
	return defaultDB
}

func importSkill() error {
	file, err := os.ReadFile("./storage/skill.txt")
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(file), "\n") {
		split := strings.Split(line, ",")
		if len(split) != 3 {
			continue
		}
		err := defaultDB.Exec("INSERT INTO aion_player_skill (skill, class) values (?, ?)", split[0], split[2]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func createTables(db *gorm.DB) error {
	createSql := []string{
		`CREATE TABLE IF NOT EXISTS aion_chat_log (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            player TEXT DEFAULT NULL,
            skill TEXT DEFAULT NULL,
            target TEXT DEFAULT NULL,
            value INTEGER DEFAULT NULL,
            time DATETIME DEFAULT NULL,
            raw_msg TEXT DEFAULT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS aion_player_info (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT DEFAULT NULL,
            type INTEGER DEFAULT NULL,
            class INTEGER DEFAULT NULL,
            time DATETIME DEFAULT NULL,
            critical_ratio REAL DEFAULT NULL,
            UNIQUE (name, type)
        );`,
		`CREATE TABLE IF NOT EXISTS aion_player_rank (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            player TEXT DEFAULT NULL,
            count INTEGER DEFAULT NULL,
            time DATETIME DEFAULT NULL,
            UNIQUE (player, count, time)
        );`,
		`CREATE TABLE IF NOT EXISTS aion_player_skill (
            skill TEXT NOT NULL,
            critical_ratio REAL DEFAULT NULL,
            class INTEGER NOT NULL,
            UNIQUE (class, skill)
        );`,
		`CREATE TABLE IF NOT EXISTS aion_timeline (
            time DATETIME NOT NULL,
            value INTEGER NOT NULL,
            type INTEGER NOT NULL DEFAULT 0
        );`,
	}

	for _, sql := range createSql {
		err := db.Exec(sql).Error
		if err != nil {
			return err
		}
	}
	return nil
}
