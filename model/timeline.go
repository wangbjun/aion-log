package model

import (
	"fmt"
	"strings"
	"time"
)

type Timeline struct {
	Time  time.Time `json:"time"`
	Value int       `json:"value"`
}

func (r Timeline) TableName() string {
	return "aion_timeline"
}

func (r Timeline) BatchInsert(items []Timeline) error {
	sql := "INSERT INTO `aion_timeline` (`time`,`value`) VALUES "
	for _, v := range items {
		sql += fmt.Sprintf("('%s',%d),", v.Time.Format(time.DateTime), v.Value)
	}
	sql = strings.TrimRight(sql, ",")
	return DB().Exec(sql).Error
}

func (r Timeline) GetAll(st, et string) ([]Timeline, error) {
	var results []Timeline
	err := DB().Where("time >= ? and time <= ?", st, et).Find(&results).Error
	return results, err
}
