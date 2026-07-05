package model

import (
	"fmt"
	"time"
)

type Player struct {
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	Type          int       `json:"type"`
	Class         int       `json:"class"`
	SkillCount    int       `json:"skill_count"`
	KillCount     int       `json:"kill_count"`
	DeathCount    int       `json:"death_count"`
	CriticalRatio float64   `json:"critical_ratio"`
	Time          time.Time `json:"time"`
}

const (
	TypeOther = iota
	TypeBright
	TypeDark
)

type Class = int

const (
	Unknown Class = iota
	JX
	SH
	SX
	GX
	ZY
	HF
	JL
	MD
	ZXZ
)

func (r Player) TableName() string {
	return "aion_player_info"
}

func (r Player) BatchInsert(items []Player) error {
	if len(items) == 0 {
		return nil
	}
	return DB().CreateInBatches(items, BatchInsertSize).Error
}

func (r Player) GetAll() ([]*Player, error) {
	var results []*Player
	err := DB().Find(&results).Error
	return results, err
}

func (r Player) GetByTime(st, et string) ([]*Player, error) {
	var results []*Player
	sql1, args1 := chatLogNameQuery("player", st, et)
	sql2, args2 := chatLogNameQuery("target", st, et)
	sql := fmt.Sprintf("select name from (%s union %s) as results", sql1, sql2)
	args := append(args1, args2...)
	err := DB().Raw(sql, args...).Find(&results).Error
	return results, err
}

func (r Player) FillStats(players []*Player, st, et string) error {
	if len(players) == 0 {
		return nil
	}

	playerSkillCount, err := r.getSkillCount(st, et)
	if err != nil {
		return err
	}
	playerKillCount, playerDeathCount, err := r.getKillAndDeathCount(st, et)
	if err != nil {
		return err
	}

	for _, player := range players {
		player.SkillCount = playerSkillCount[player.Name]
		player.KillCount = playerKillCount[player.Name]
		player.DeathCount = playerDeathCount[player.Name]
	}
	return nil
}

func chatLogNameQuery(column, st, et string) (string, []interface{}) {
	sql := fmt.Sprintf("select %s as name from aion_chat_log", column)
	if st == "" || et == "" {
		return sql, nil
	}
	return sql + " where time >= ? and time <= ?", []interface{}{st, et}
}

func chatLogTimeFilter(sql, st, et string) (string, []interface{}) {
	if st == "" || et == "" {
		return sql, nil
	}
	return sql + " and time >= ? and time <= ?", []interface{}{st, et}
}

func (r Player) getSkillCount(st, et string) (map[string]int, error) {
	var result []struct {
		Player string
		Count  int
	}
	sql := "select player,count(1) as count from aion_chat_log " +
		"where target != '' and skill not in ('attack','kill','killed')"
	sql, args := chatLogTimeFilter(sql, st, et)
	sql += " group by player"
	err := DB().Raw(sql, args...).Find(&result).Error
	if err != nil {
		return nil, err
	}

	playerSkillCount := make(map[string]int, len(result))
	for _, v := range result {
		playerSkillCount[v.Player] = v.Count
	}
	return playerSkillCount, nil
}

func (r Player) getKillAndDeathCount(st, et string) (map[string]int, map[string]int, error) {
	playerKillCount := make(map[string]int)
	playerDeathCount := make(map[string]int)

	if err := r.addKillAndDeathCount("kill", st, et, playerKillCount, playerDeathCount); err != nil {
		return nil, nil, err
	}
	if err := r.addKillAndDeathCount("killed", st, et, playerKillCount, playerDeathCount); err != nil {
		return nil, nil, err
	}

	return playerKillCount, playerDeathCount, nil
}

func (r Player) addKillAndDeathCount(skill, st, et string, playerKillCount, playerDeathCount map[string]int) error {
	var result []struct {
		Player string
		Target string
		Count  int
	}
	sql := "select player,target,count(1) count from aion_chat_log where skill = ?"
	args := []interface{}{skill}
	if st != "" && et != "" {
		sql += " and time >= ? and time <= ?"
		args = append(args, st, et)
	}
	sql += " group by player,target"
	err := DB().Raw(sql, args...).Find(&result).Error
	if err != nil {
		return err
	}

	for _, v := range result {
		if skill == "kill" {
			playerKillCount[v.Player] += v.Count
			playerDeathCount[v.Target] += v.Count
			continue
		}
		playerKillCount[v.Target] += v.Count
		playerDeathCount[v.Player] += v.Count
	}
	return nil
}
