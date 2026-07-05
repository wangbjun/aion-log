package model

import (
	"strconv"
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
	if len(items) == 0 {
		return nil
	}
	return DB().CreateInBatches(items, BatchInsertSize).Error
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
	levelInt, err := strconv.Atoi(level)
	if err != nil {
		levelInt = 3
	}
	sql := "SELECT player, GROUP_CONCAT(time ORDER BY time DESC) AS times, COUNT(time) AS counts FROM" +
		" (SELECT player, time FROM aion_player_rank WHERE count = ? ORDER BY player, time DESC ) GROUP BY player HAVING counts >= 5"
	err = DB().Raw(sql, levelInt).Find(&results).Error
	return results, err
}
