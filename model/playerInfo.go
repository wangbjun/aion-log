package model

import (
	"fmt"
	"strings"
	"time"
)

type Player struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Type       int       `json:"type"`
	Class      int       `json:"class"`
	SkillCount int       `json:"skill_count"`
	KillCount  int       `json:"kill_count"`
	DeathCount int       `json:"death_count"`
	Time       time.Time `json:"time"`
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
	sql := "INSERT INTO `aion_player_info` (`name`,`type`,`class`,`time`) VALUES "
	for _, v := range items {
		sql += fmt.Sprintf("('%s',%d,'%d','%s'),", v.Name, v.Type, v.Class, v.Time.Format(time.DateTime))
	}
	sql = strings.TrimRight(sql, ",")
	return DB().Exec(sql).Error
}

func (r Player) GetAll() ([]*Player, error) {
	var results []*Player
	err := DB().Find(&results).Error
	return results, err
}

func (r Player) GetByTime(st, et string) ([]*Player, error) {
	var results []*Player
	sql := fmt.Sprintf("select distinct name from (select player as name from aion_player_chat_log where "+
		"time >= '%s' AND time <= '%s' union select target as name from aion_player_chat_log where "+
		"target != '' and time >= '%s' AND time <= '%s') as result", st, et, st, et)
	err := DB().Raw(sql).Find(&results).Error
	return results, err
}
