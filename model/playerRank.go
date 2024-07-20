package model

import (
	"fmt"
	"strings"
	"time"
)

type Rank struct {
	Id     int       `json:"id"`
	Player string    `json:"player"`
	Count  int       `json:"count"`
	Time   time.Time `json:"time"`
}

func (r Rank) TableName() string {
	return "aion_player_rank"
}

func (r Rank) BatchInsert(items []Rank) error {
	sql := "INSERT INTO `aion_player_rank` (`player`,`count`,`time`) VALUES "
	for _, v := range items {
		sql += fmt.Sprintf("('%s',%d,'%s'),", v.Player, v.Count, v.Time.Format(time.DateTime))
	}
	sql = strings.TrimRight(sql, ",")
	return DB().Exec(sql).Error
}

type RankResult struct {
	Player string `json:"player"`
	Type   int    `json:"type"`
	Class  int    `json:"class"`
	Counts int    `json:"counts"`
	Times  string `json:"times"`
}

func (r Rank) GetAll(level string) ([]RankResult, error) {
	var results []RankResult
	sql := fmt.Sprintf("SELECT player, GROUP_CONCAT(time ORDER BY time DESC) AS times, COUNT(time) AS counts FROM"+
		" (SELECT player, time FROM aion_player_rank WHERE count = %s ORDER BY player, time DESC ) GROUP BY player HAVING counts >= 5", level)
	err := DB().Raw(sql).Find(&results).Error
	return results, err
}
