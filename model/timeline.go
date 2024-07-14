package model

import (
	"fmt"
	"strings"
	"time"
)

type Timeline struct {
	Time  time.Time `json:"time"`
	Value int       `json:"value"`
	Type  int       `json:"type"`
}

func (r Timeline) TableName() string {
	return "aion_timeline"
}

func (r Timeline) BatchInsert(items []Timeline) error {
	sql := "INSERT INTO `aion_timeline` (`time`,`value`, `type`) VALUES "
	for _, v := range items {
		sql += fmt.Sprintf("('%s',%d,%d),", v.Time.Format(time.DateTime), v.Value, v.Type)
	}
	sql = strings.TrimRight(sql, ",")
	return DB().Exec(sql).Error
}

func (r Timeline) GetAll(st, et string, tp int) ([]Timeline, error) {
	var results []Timeline
	err := DB().Where("time >= ? and time <= ? and type = ?", st, et, tp).Find(&results).Error
	return results, err
}
