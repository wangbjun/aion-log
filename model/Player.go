package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Player struct {
	Id   int       `json:"id"`
	Name string    `json:"name"`
	Type int       `json:"type"`
	Time time.Time `json:"time"`
}

const (
	TypeOther = 0
	TypeTian  = 1
	TypeMo    = 2
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

func (r Player) GetPlayers() ([]Player, error) {
	var results []Player
	err := DB().Order("name asc").Find(&results).Error
	return results, err
}

func (r Player) ChangeType(id, t string) error {
	sql := fmt.Sprintf("update %s set type = %s where id = %s", r.TableName(), t, id)
	return DB().Exec(sql).Error
}
