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

func (r ChatLog) GetAll(st, et string, page, pageSize int, player, target, skill, sort, value, banPlayer string) ([]ChatLog, int64, error) {
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
		query = query.Where("target != ''")
	}
	if banPlayer != "" {
		banPlayers := strings.Split(banPlayer, ",")
		query = query.Where("player not in (?)", banPlayers)
	}

	var count int64
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

type SkillDamage struct {
	Skill    string  `json:"skill"`
	Count    int     `json:"count"`
	Damage   int     `json:"damage"`
	Average  float64 `json:"average"`
	Critical float64 `json:"critical"`
}

func (r ChatLog) GetClassTop(class, player string) ([]*SkillDamage, error) {
	sql := "select skill,count(1) count,max(value) damage,avg(value)  average from aion_chat_log where value > 0 and target != ''"
	if player != "" {
		sql += " and player = '" + player + "'"
	}
	sql += fmt.Sprintf(" and skill in (select skill from aion_player_skill where class = %s) group by skill order by damage desc", class)
	var results []*SkillDamage
	err := DB().Raw(sql).Find(&results).Error
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}

func (r ChatLog) GetCriticalRatio(player string) ([]SkillDamage, error) {
	condition := "target != '' and raw_msg LIKE '致命一击%'"
	if player != "" {
		condition += " and player = '" + player + "'"
	}
	sql := fmt.Sprintf("SELECT a.skill, a.count / b.total critical FROM (SELECT skill, count(1) count FROM aion_chat_log "+
		"WHERE %s GROUP BY skill) a JOIN (SELECT skill, count(1) total FROM aion_chat_log where target != ''", condition)
	if player != "" {
		sql += " and player = '" + player + "'"
	}
	sql += " GROUP BY skill) b ON a.skill = b.skill"
	var results []SkillDamage
	err := DB().Raw(sql).Find(&results).Error
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}
