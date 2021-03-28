package model

import (
	"aion/util"
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

var CachedData = cache.New(6*time.Hour, 30*time.Minute)

type Log struct {
	Player       string    `gorm:"player" json:"player"`
	Skill        string    `gorm:"skill" json:"skill"`
	TargetPlayer string    `gorm:"target_player" json:"target_player"`
	Damage       int       `gorm:"damage" json:"damage"`
	Time         time.Time `gorm:"time" json:"time"`
	OriginDesc   string    `gorm:"origin_desc" json:"origin_desc"`
}

type BattleLog struct {
	Id int `gorm:"primaryKey" json:"id"`
	Log
}

func (r BattleLog) TableName() string {
	return "aion_player_battle_log"
}

func (r BattleLog) Insert() error {
	return DB().Create(&r).Error
}

func (r BattleLog) BacthInsert(items []Log) error {
	sql := "INSERT INTO `aion_player_battle_log` (`player`,`skill`,`target_player`,`damage`,`time`,`origin_desc`) VALUES "
	for k, v := range items {
		if len(items)-1 == k {
			sql += fmt.Sprintf("('%s','%s','%s',%d,'%s','%s')", v.Player,
				v.Skill, v.TargetPlayer, v.Damage, v.Time.Format(util.TimeFormat), v.OriginDesc)
		} else {
			sql += fmt.Sprintf("('%s','%s','%s',%d,'%s','%s'),", v.Player,
				v.Skill, v.TargetPlayer, v.Damage, v.Time.Format(util.TimeFormat), v.OriginDesc)
		}
	}
	return DB().Exec(sql).Error
}

func (r BattleLog) SavePlayer() {
	p1 := Player{
		Name: r.Player,
		Type: TypeOther,
		Time: r.Time,
	}
	_ = p1.Save()
	p2 := Player{
		Name: r.TargetPlayer,
		Type: TypeOther,
		Time: r.Time,
	}
	_ = p2.Save()
}

func (r BattleLog) GetAll(st, et string, page, pageSize int, player, skill, sort string) ([]BattleLog, int, error) {
	var results []BattleLog
	query := DB().Model(&BattleLog{})
	if st != "" {
		query = query.Where("time >= ?", st)
	}
	if et != "" {
		query = query.Where("time <= ?", et)
	}
	if player != "" {
		query = query.Where("player = ? or target_player = ?", player, player)
	}
	if skill != "" {
		query = query.Where("skill = ?", skill)
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

func (r BattleLog) GetLastTime() *time.Time {
	var result BattleLog
	err := DB().Order("time desc").Limit(1).First(&result).Error
	if err != nil {
		return nil
	}
	return &result.Time
}

type Count struct {
	Count int
}

func (r BattleLog) GetSkillCount(st, et, player string) int {
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

type weekly struct {
	Day string `json:"day"`
}

func (r BattleLog) GetWeekly() []weekly {
	var results []weekly
	sql := "SELECT date_format(time, '%Y-%m-%d') day FROM aion_player_battle_log group by day order by day desc"
	DB().Raw(sql).Find(&results)
	return results
}

func (r BattleLog) DeleteByDay(day string) error {
	sql := "delete from aion_player_battle_log where date_format(time, '%Y-%m-%d') <= '" + day + "'"
	return DB().Exec(sql).Error
}
