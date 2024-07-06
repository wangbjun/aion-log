package model

import (
	"fmt"
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

func (r Player) Insert(player Player) error {
	value := fmt.Sprintf("('%s', %d, %d, '%s')", player.Name, player.Type, player.Class, player.Time.Format(time.DateTime))
	sql := "INSERT INTO aion_player_info (name, type, class, time) VALUES " + value + " ON DUPLICATE KEY UPDATE time = VALUES(time)"
	return DB().Exec(sql).Error
}

func (r Player) GetAll() ([]*Player, error) {
	var results []*Player
	err := DB().Find(&results).Error
	return results, err
}
