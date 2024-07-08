package model

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"strings"
	"time"
)

var CachedData = cache.New(6*time.Hour, 30*time.Minute)

type Log struct {
	Id     int       `gorm:"primaryKey" json:"id"`
	Player string    `gorm:"player" json:"player"`
	Skill  string    `gorm:"skill" json:"skill"`
	Target string    `gorm:"target" json:"target"`
	Value  int       `gorm:"value" json:"value"`
	Time   time.Time `gorm:"time" json:"time"`
	RawMsg string    `gorm:"raw_msg" json:"raw_msg"`
}

func (r Log) TableName() string {
	return "aion_player_battle_log"
}

func (r Log) BatchInsert(items []Log) error {
	sql := "INSERT INTO `aion_player_battle_log` (`player`,`skill`,`target`,`value`,`time`,`raw_msg`) VALUES "
	for _, v := range items {
		sql += fmt.Sprintf("('%s','%s','%s',%d,'%s','%s'),", v.Player, v.Skill, v.Target, v.Value, v.Time.Format(time.DateTime), strings.TrimSpace(v.RawMsg))
	}
	sql = strings.TrimRight(sql, ",")
	return DB().Exec(sql).Error
}

func (r Log) GetAll(st, et string, page, pageSize int, player, target, skill, sort, value string) ([]Log, int, error) {
	var results []Log
	query := DB().Model(&Log{})
	if st != "" {
		query = query.Where("time >= ?", st)
	}
	if et != "" {
		query = query.Where("time <= ?", et)
	}
	if player != "" {
		query = query.Where("player = ?", player)
	}
	if target != "" {
		query = query.Where("target = ?", target)
	}
	if skill != "" {
		query = query.Where("skill = ?", skill)
	}
	if value != "" {
		query = query.Where("value > ?", value)
	}
	var count int
	err := query.Count(&count).Error
	if err != nil {
		return results, 0, err
	}
	if sort == "" {
		sort = "id"
	}
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Order(sort + " desc").Find(&results).Error
	return results, count, err
}

type Count struct {
	Count int
}

func (r Log) GetSkillCount(st, et, player string) int {
	var key = st + et + player
	if cached, found := CachedData.Get(key); found {
		return cached.(int)
	}
	var result Count
	sql := "select count(distinct(time)) count from aion_player_battle_log where player = '" + player + "'"
	if st != "" {
		sql += " and time >= '" + st + "'"
	}
	if et != "" {
		sql += " and time <= '" + et + "'"
	}
	DB().Raw(sql).Find(&result)
	CachedData.Set(key, result.Count, cache.DefaultExpiration)
	return result.Count
}
