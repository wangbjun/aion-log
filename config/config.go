package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

var Conf *ini.File

var DBConfig map[string]map[string]string

func Init(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("conf file [%s]  not found!", file)
	}
	conf, err := ini.Load(file)
	if err != nil {
		return fmt.Errorf("parse conf file [%s] failed, err: %s", file, err.Error())
	}
	Conf = conf
	DBConfig = map[string]map[string]string{
		"default": {
			"dialect":      Conf.Section("DB").Key("Dialect").String(),
			"dsn":          Conf.Section("DB").Key("DSN").String(),
			"maxIdleConns": Conf.Section("DB").Key("MAX_IDLE_CONN").String(),
			"maxOpenConns": Conf.Section("DB").Key("MAX_OPEN_CONN").String(),
		},
	}
	return nil
}

func GetAPP(name string) *ini.Key {
	return Conf.Section("APP").Key(name)
}
