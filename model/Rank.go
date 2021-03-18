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
	Player   string `json:"player"`
	Type     int    `json:"type"`
	Count    int    `json:"count"`
	AllCount int    `json:"all_count"`
	Times    string `json:"times"`
}

func (r Rank) GetAll(st, et string) ([]RankResult, error) {
	var results []RankResult
	sql := "select player,SUBSTRING_INDEX(GROUP_CONCAT(time ORDER BY time desc),',',30) times, count(1) count from aion_player_rank where 1=1"
	if st != "" {
		sql += fmt.Sprintf(" and time >= '%s'", st)
	}
	if et != "" {
		sql += fmt.Sprintf(" and time <= '%s'", et)
	}
	sql += " group by player HAVING count >= 5"
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
	sql := "select player,count(DISTINCT(skill)) count,time from aion_player_battle_log where skill != '普通攻击'"
	lastTime := Rank{}.GetLastTime()
	if lastTime != nil {
		sql += " and time > '" + lastTime.Format(util.TimeFormat) + "'"
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
