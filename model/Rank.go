package model

import (
	"aion/util"
	"fmt"
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

func (r Rank) Save() error {
	return DB().Create(&r).Error
}

type RankResult struct {
	Player    string `json:"player"`
	Type      int    `json:"type"`
	Counts    int    `json:"counts"`
	AllCounts int    `json:"all_counts"`
	Times     string `json:"times"`
}

func (r Rank) GetAll(st, et, level string) ([]RankResult, error) {
	var results []RankResult
	sql := "select player,SUBSTRING_INDEX(GROUP_CONCAT(time ORDER BY time desc),',',30) times, count(1) counts " +
		"from aion_player_rank where count = " + level
	if st != "" {
		sql += fmt.Sprintf(" and time >= '%s'", st)
	}
	if et != "" {
		sql += fmt.Sprintf(" and time <= '%s'", et)
	}
	sql += " group by player HAVING counts >= 3"
	err := DB().Raw(sql).Find(&results).Error
	return results, err
}

func (r Rank) GetLastTime() *time.Time {
	var result Rank
	err := DB().Order("time desc").First(&result).Error
	if err != nil {
		return nil
	} else {
		return &result.Time
	}
}

func (r Rank) GetRanks() ([]Rank, error) {
	sql := "select player,count(DISTINCT(skill)) count,time from aion_player_battle_log"
	lastTime := Rank{}.GetLastTime()
	if lastTime != nil {
		sql += " where time > '" + lastTime.Format(util.TimeFormat) + "'"
	}
	sql += " group by player,time HAVING count(DISTINCT(skill)) >= 3"

	var results []Rank
	err := DB().Raw(sql).Find(&results).Error
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}

func (r Rank) DeleteByDay(day string) error {
	sql := "delete from aion_player_rank where date_format(time, '%Y-%m-%d') <= '" + day + "'"
	return DB().Exec(sql).Error
}
