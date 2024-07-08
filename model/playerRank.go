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
	Player    string `json:"player"`
	Type      int    `json:"type"`
	Class     int    `json:"class"`
	Counts    int    `json:"counts"`
	AllCounts int    `json:"all_counts"`
	Times     string `json:"times"`
}

func (r Rank) GetAll(level string) ([]RankResult, error) {
	var results []RankResult
	sql := "select player,SUBSTRING_INDEX(GROUP_CONCAT(time ORDER BY time desc),',',30) times, count(1) counts " +
		"from aion_player_rank where count = " + level
	sql += " group by player HAVING counts >= 5"
	err := DB().Raw(sql).Find(&results).Error
	return results, err
}
