package model

import (
	"aion/zlog"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type Player struct {
	Id   int       `json:"id"`
	Name string    `json:"name"`
	Type int       `json:"type"`
	Pro  int       `json:"pro"`
	Time time.Time `json:"time"`
}

const (
	TypeOther = iota
	TypeTian
	TypeMo
)

type pro = int

const (
	Unknown pro = iota
	JX
	SH
	SX
	GX
	ZY
	HF
	JL
	MD
)

func (r Player) TableName() string {
	return "aion_player_info"
}

func (r Player) Save() error {
	var existed Player
	err := DB().First(&existed, "name = ?", r.Name).Error
	if err == nil {
		return DB().Model(&existed).UpdateColumn("time", r.Time).Error
	} else if err == gorm.ErrRecordNotFound {
		return DB().Create(&r).Error
	}
	return nil
}

func (r Player) GetAll() ([]Player, error) {
	var results []Player
	err := DB().Order("name asc").Find(&results).Error
	return results, err
}

func (r Player) SaveType() error {
	var existed Player
	err := DB().First(&existed, "name = ?", r.Name).Error
	if err == nil {
		var update = Player{Time: r.Time}
		if r.Type != 0 {
			update = Player{Time: r.Time, Type: r.Type}
		}
		return DB().Model(&existed).UpdateColumns(update).Error
	} else if err == gorm.ErrRecordNotFound {
		return DB().Create(&r).Error
	}
	return nil
}

func (r Player) ChangeType(id, t string) error {
	sql := fmt.Sprintf("update %s set type = %s where id = %s", r.TableName(), t, id)
	return DB().Exec(sql).Error
}

func (r Player) DeleteByDay(day string) error {
	sql := "delete from aion_player_info where date_format(time, '%Y-%m-%d') <= '" + day + "'"
	return DB().Exec(sql).Error
}

func (r Player) UpdateBySkills(p pro, skills []string) {
	var results []BattleLog
	sql := "select distinct(player) from aion_player_battle_log where skill in ('" + strings.Join(skills, "','") + "')"
	err := DB().Raw(sql).Find(&results).Error
	if err != nil {
		return
	}
	var players []string
	for _, v := range results {
		players = append(players, v.Player)
	}
	updateSql := fmt.Sprintf("update aion_player_info set pro = %d where name in ('%s')",
		p, strings.Join(players, "','"))
	err = DB().Exec(updateSql).Error
	if err != nil {
		zlog.Logger.Error("update player pro failed")
		return
	}
}
