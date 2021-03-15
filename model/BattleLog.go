package model

import (
	"aion/util"
	"fmt"
	"time"
)

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

func (r BattleLog) GetRank(st, et, level string) ([]BattleLog, error) {
	var results []BattleLog
	sql := "SELECT time,player FROM aion_player_battle_log where skill != '普通攻击'"
	if st != "" {
		sql += fmt.Sprintf(" and time >= '%s'", st)
	}
	if et != "" {
		sql += fmt.Sprintf(" and time <= '%s'", et)
	}
	sql += fmt.Sprintf(" group by time,player having count(distinct(skill)) = %s order by time desc", level)
	err := DB().Raw(sql).Find(&results).Error
	return results, err
}

type CountResult struct {
	Type  int `gorm:"column:type" json:"type"`
	Count int
}

func (r BattleLog) GetStat(st, et string) (int, error) {
	var result CountResult
	sql := "SELECT count(distinct(player)) count FROM aion_player_battle_log where 1 = 1"
	if st != "" {
		sql += fmt.Sprintf(" and time >= '%s'", st)
	}
	if et != "" {
		sql += fmt.Sprintf(" and time <= '%s'", et)
	}
	err := DB().Raw(sql).First(&result).Error
	return result.Count, err
}

func (r BattleLog) GetLastTime() time.Time {
	var result BattleLog
	err := DB().Order("time desc").Limit(1).First(&result).Error
	if err != nil {
		return time.Now()
	}
	return result.Time
}
