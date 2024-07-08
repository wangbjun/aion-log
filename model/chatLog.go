package model

import (
	"fmt"
	"strings"
	"time"
)

type ChatLog struct {
	Id     int       `gorm:"primaryKey" json:"id"`
	Player string    `gorm:"player" json:"player"`
	Skill  string    `gorm:"skill" json:"skill"`
	Target string    `gorm:"target" json:"target"`
	Value  int       `gorm:"value" json:"value"`
	Time   time.Time `gorm:"time" json:"time"`
	RawMsg string    `gorm:"raw_msg" json:"raw_msg"`
}

func (r ChatLog) TableName() string {
	return "aion_player_chat_log"
}

func (r ChatLog) BatchInsert(items []ChatLog) error {
	sql := "INSERT INTO `aion_player_chat_log` (`player`,`skill`,`target`,`value`,`time`,`raw_msg`) VALUES "
	for _, v := range items {
		sql += fmt.Sprintf("('%s','%s','%s',%d,'%s','%s'),", v.Player, v.Skill, v.Target, v.Value, v.Time.Format(time.DateTime), strings.TrimSpace(v.RawMsg))
	}
	sql = strings.TrimRight(sql, ",")
	return DB().Exec(sql).Error
}

func (r ChatLog) GetAll(st, et string, page, pageSize int, player, target, skill, sort, value string) ([]ChatLog, int, error) {
	var results []ChatLog
	query := DB().Model(&ChatLog{})
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
		query = query.Where("skill like ?", skill+"%")
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

var cacheData = make(map[string]int)

func (r ChatLog) GetSkillCount(player string) int {
	if cached, ok := cacheData[player]; ok {
		return cached
	}
	var result []struct{ skill string }
	sql := "select skill from aion_player_chat_log where skill not in ('','kill','killed') " +
		"and value > 0 and player = '" + player + "'"
	DB().Raw(sql).Find(&result)
	count := len(result)
	if count > 0 {
		cacheData[player] = count
	}
	return count
}

func (r ChatLog) GetRanks() ([]Rank, error) {
	sql := "select player,count(DISTINCT(skill)) count,time from aion_player_chat_log where skill not in ('','kill','killed') " +
		"and value > 0 group by player,time HAVING count >= 3"
	var results []Rank
	err := DB().Raw(sql).Find(&results).Error
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}
