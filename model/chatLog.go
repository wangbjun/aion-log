package model

import (
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
	if len(items) == 0 {
		return nil
	}
	for i := range items {
		items[i].RawMsg = strings.TrimSpace(items[i].RawMsg)
	}
	return DB().CreateInBatches(items, BatchInsertSize).Error
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
		query = query.Where("player = ? AND target = ?", player, target)
	} else if player != "" {
		query = query.Where("player = ?", player)
	} else if target != "" {
		if strings.HasPrefix(target, "-") {
			query = query.Where("target != ?", strings.TrimPrefix(target, "-"))
		} else {
			query = query.Where("target = ?", target)
		}
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
	sort = allowedChatLogSort(sort)
	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Order(sort + " desc").Find(&results).Error
	return results, count, err
}

func allowedChatLogSort(sort string) string {
	allowed := map[string]string{
		"id":     "id",
		"time":   "time",
		"value":  "value",
		"skill":  "skill",
		"player": "player",
		"target": "target",
	}
	if column, ok := allowed[sort]; ok {
		return column
	}
	return "id"
}

func (r ChatLog) GetRanks() ([]Rank, error) {
	sql := "select player,count(DISTINCT(skill)) count,time from aion_chat_log where skill not in ('','kill','killed') " +
		"and value > 0 group by player,time HAVING count >= 3"
	var results []Rank
	err := DB().Raw(sql).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

type SkillDamage struct {
	Skill    string  `json:"skill"`
	Count    int     `json:"count"`
	Damage   int     `json:"damage"`
	Average  float64 `json:"average"`
	Critical float64 `json:"critical"`
}

func (r ChatLog) GetClassTop(class, player string) ([]*SkillDamage, error) {
	classInt, err := strconv.Atoi(class)
	if err != nil {
		return nil, err
	}
	query := DB().Table("aion_chat_log").
		Select("skill, count(1) count, max(value) damage, avg(value) average").
		Where("value > 0").
		Where("target != ''").
		Where("skill in (select skill from aion_player_skill where class = ?)", classInt).
		Group("skill").
		Order("damage desc")
	if player != "" {
		query = query.Where("player = ?", player)
	}
	var results []*SkillDamage
	err = query.Find(&results).Error
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}

func (r ChatLog) GetCriticalRatio(player string) ([]SkillDamage, error) {
	criticalCondition := "target != '' and raw_msg LIKE ?"
	criticalArgs := []interface{}{"致命一击%"}
	totalCondition := "target != ''"
	var totalArgs []interface{}
	if player != "" {
		criticalCondition += " and player = ?"
		criticalArgs = append(criticalArgs, player)
		totalCondition += " and player = ?"
		totalArgs = append(totalArgs, player)
	}

	sql := "SELECT a.skill, (a.count * 1.0) / b.total critical FROM " +
		"(SELECT skill, count(1) count FROM aion_chat_log WHERE " + criticalCondition + " GROUP BY skill) a " +
		"JOIN (SELECT skill, count(1) total FROM aion_chat_log WHERE " + totalCondition + " GROUP BY skill) b ON a.skill = b.skill"
	args := append(criticalArgs, totalArgs...)
	var results []SkillDamage
	err := DB().Raw(sql, args...).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
