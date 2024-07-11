package model

import (
	"fmt"
	"strconv"
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
	return "aion_chat_log"
}

func (r ChatLog) BatchInsert(items []ChatLog) error {
	sql := "INSERT INTO `aion_chat_log` (`player`,`skill`,`target`,`value`,`time`,`raw_msg`) VALUES "
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

	if player != "" && target != "" {
		query = query.Where("player = ? or target = ?", player, player)
	} else if player != "" {
		query = query.Where("player = ?", player)
	} else if target != "" {
		if strings.HasPrefix(target, "-") {
			query = query.Where("target != ?", target)
		} else {
			query = query.Where("target = ?", target)
		}
		query = query.Where("target = ?", target)
	}
	if skill != "" {
		query = query.Where("skill like ?", skill+"%")
	}
	if value != "" {
		seg := strings.Split(value, "-")
		if len(seg) == 2 {
			ge, _ := strconv.Atoi(seg[0])
			le, _ := strconv.Atoi(seg[1])
			if ge == le {
				query = query.Where("value = ?", ge)
			} else if le > ge {
				query = query.Where("value >= ? and value <= ?", ge, le)
			} else if ge > 0 && le == 0 {
				query = query.Where("value >= ?", ge)
			} else if le > 0 && ge == 0 {
				query = query.Where("value <= ?", le)
			}
		} else {
			valueInt, _ := strconv.Atoi(value)
			query = query.Where("value > ?", valueInt)
		}
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

func (r ChatLog) GetRanks() ([]Rank, error) {
	sql := "select player,count(DISTINCT(skill)) count,time from aion_chat_log where skill not in ('','kill','killed') " +
		"and value > 0 group by player,time HAVING count >= 3"
	var results []Rank
	err := DB().Raw(sql).Find(&results).Error
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}

func (r ChatLog) AddIndex() error {
	return DB().Exec("ALTER TABLE `aion_chat_log` " +
		"ADD KEY `idx_player_skill_time` (`player`,`skill`,`time`)," +
		"ADD KEY `idx_skill` (`skill`)," +
		"ADD KEY `idx_time` (`time`)," +
		"ADD KEY `idx_target` (`target`)," +
		"ADD KEY `idx_value` (`value`)").Error
}

func (r ChatLog) RemoveIndex() error {
	DB().Exec("drop index idx_player_skill_time on aion_chat_log")
	DB().Exec("drop index idx_skill on aion_chat_log")
	DB().Exec("drop index idx_target on aion_chat_log")
	DB().Exec("drop index idx_time on aion_chat_log")
	return DB().Exec("drop index idx_value on aion_chat_log").Error
}
