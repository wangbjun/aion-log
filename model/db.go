package model

import (
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "modernc.org/sqlite"
)

var defaultDB *gorm.DB

const BatchInsertSize = 2000

func Init() {
	config := &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 500 * time.Millisecond,
			LogLevel:      logger.Error,
			Colorful:      true,
		}),
	}
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        "aion.db?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=temp_store(MEMORY)&_pragma=cache_size(-64000)",
	}, config)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database connection: " + err.Error())
	}
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(5)
	defaultDB = db
	if err := createTables(); err != nil {
		panic("failed to create tables: " + err.Error())
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
		if len(split) != 2 {
			continue
		}
		err := defaultDB.Exec("INSERT INTO aion_player_skill (skill, class) values (?, ?)", split[0], split[1]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func createTables() error {
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
            skill_count INTEGER DEFAULT 0,
            kill_count INTEGER DEFAULT 0,
            death_count INTEGER DEFAULT 0,
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
		err := defaultDB.Exec(sql).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func ResetData() error {
	tables := []string{
		"aion_player_skill",
		"aion_chat_log",
		"aion_player_info",
		"aion_player_rank",
		"aion_timeline",
	}
	for _, table := range tables {
		if err := defaultDB.Exec("DELETE FROM " + table).Error; err != nil {
			return err
		}
	}
	for _, table := range tables {
		if err := defaultDB.Exec("DELETE FROM sqlite_sequence WHERE name = ?", table).Error; err != nil {
			return err
		}
	}
	err := createTables()
	if err != nil {
		panic("failed to create tables: " + err.Error())
	}
	err = importSkill()
	if err != nil {
		panic("failed to import skill: " + err.Error())
	}
	return nil
}
