package service

import (
	"aion/model"
	"fmt"
	"strings"
)

type ClassifyService struct {
	CacheService
}

func NewClassifyService() *ClassifyService {
	return &ClassifyService{CacheService: *NewCacheService()}
}

var updateSql = []string{
	"UPDATE aion_player_info SET type = 1 WHERE name like '%永恒之巅'",
	"UPDATE aion_player_info SET type = 1 WHERE name like '%谁与争锋'",
	"UPDATE aion_player_info SET type = 1 WHERE name like '%永恒之岛'",
	"UPDATE aion_player_info SET type = 2 WHERE name like '%火之神殿'",
	"UPDATE aion_player_info SET type = 2 WHERE name like '%傲世八星'",
	"UPDATE aion_player_info SET type = 1 WHERE name in (select distinct player from aion_chat_log where " +
		"target in('魔族中级守护神将','魔族上级守护神将','魔族结界膜生成师','魔族城门')) and class != 0",
	"UPDATE aion_player_info SET type = 2 WHERE name in (select distinct player from aion_chat_log where " +
		"target in('天族中级守护神将','天族上级守护神将','天族结界膜生成师','天族城门')) and class != 0",
	"UPDATE aion_player_info SET type = 1 WHERE name in (select distinct target from aion_chat_log where " +
		"player in('魔族中级守护神将','魔族上级守护神将','魔族结界膜生成师','魔族城门')) and class != 0",
	"UPDATE aion_player_info SET type = 2 WHERE name in (select distinct target from aion_chat_log where " +
		"player in('天族中级守护神将','天族上级守护神将','天族结界膜生成师','天族城门')) and class != 0",
}

func (r ClassifyService) Run() error {
	err := r.CacheService.Load()
	if err != nil {
		return err
	}

	for _, sql := range updateSql {
		err := model.DB().Exec(sql).Error
		if err != nil {
			return err
		}
	}

	for i := 0; i < 2; i++ {
		err = r.updateUnknown()
		if err != nil {
			return err
		}
		err = r.updateBright()
		if err != nil {
			return err
		}
		err = r.updateDark()
		if err != nil {
			return err
		}
		err = r.updateBright()
		if err != nil {
			return err
		}
	}

	err = r.updatePlayerCritical()
	if err != nil {
		return err
	}

	err = r.updateSkillCritical()

	return err
}

func (r ClassifyService) updateBright() error {
	var result []struct {
		Id int
	}
	sql := "select id from aion_player_info where name in (select distinct(player) name from aion_chat_log " +
		"where target in (select name from aion_player_info where type = 2)) and class != 0 and type = 0"

	err := model.DB().Raw(sql).Scan(&result).Error
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return nil
	}

	var idStr string
	for _, res := range result {
		idStr += fmt.Sprintf("%d,", res.Id)
	}
	sql = fmt.Sprintf("update aion_player_info set type = 1 where id in (%s)", strings.TrimRight(idStr, ","))
	err = model.DB().Exec(sql).Error

	if err != nil {
		return err
	}

	return nil
}

func (r ClassifyService) updateDark() error {
	var result []struct {
		Id int
	}
	sql := "select id from aion_player_info where name in (select distinct(player) name from aion_chat_log " +
		"where target in (select name from aion_player_info where type = 1)) and class != 0 and type = 0"

	err := model.DB().Raw(sql).Scan(&result).Error
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return nil
	}

	var idStr string
	for _, res := range result {
		idStr += fmt.Sprintf("%d,", res.Id)
	}
	sql = fmt.Sprintf("update aion_player_info set type = 2 where id in (%s)", strings.TrimRight(idStr, ","))
	err = model.DB().Exec(sql).Error

	if err != nil {
		return err
	}

	return nil
}

func (r ClassifyService) updateUnknown() error {
	var unknown []model.Player
	err := model.DB().Raw("select name from aion_player_info where type =0 and class != 0").Scan(&unknown).Error
	if err != nil {
		return err
	}
	knownPlayers := make(map[string]*model.Player)
	for _, player := range r.CacheService.cachePlayer {
		seg := strings.Split(player.Name, "-")
		if len(seg) == 2 {
			knownPlayers[seg[0]] = player
		}
	}

	for _, player := range unknown {
		if existed, ok := knownPlayers[player.Name]; ok {
			err = model.DB().Exec("update aion_player_info set type = ? where name = ?", existed.Type, player.Name).Error
			if err != nil {
				return err
			}
		}
	}

	return err
}

func (r ClassifyService) updatePlayerCritical() error {
	var result []struct {
		Player string
		Ratio  float64
	}
	playerCriticalSql := "SELECT a.player, (a.count * 1.0)/b.total as ratio FROM " +
		"(SELECT player, count(1) count FROM aion_chat_log WHERE target != '' and " +
		"skill not in ('attack','kill','killed') and raw_msg LIKE '致命一击%' GROUP BY player) a JOIN " +
		"(SELECT player, count(1) total FROM aion_chat_log where target != '' and " +
		"skill not in ('attack','kill','killed') GROUP BY player) b ON a.player = b.player"
	err := model.DB().Raw(playerCriticalSql).Find(&result).Error
	if err != nil {
		return err
	}
	for _, res := range result {
		err = model.DB().Exec("update aion_player_info set critical_ratio = ? where name = ?", res.Ratio, res.Player).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r ClassifyService) updateSkillCritical() error {
	var result []struct {
		Skill string
		Ratio float64
	}
	skillCriticalSql := "SELECT a.skill, (a.count * 1.0) / b.total ratio FROM (SELECT skill, count(1) count FROM aion_chat_log " +
		"WHERE target != '' and raw_msg LIKE '致命一击%' group BY skill) a JOIN (SELECT skill, count(1) total FROM " +
		"aion_chat_log where target != '' GROUP BY skill) b ON a.skill = b.skill"
	err := model.DB().Raw(skillCriticalSql).Find(&result).Error
	if err != nil {
		return err
	}
	for _, res := range result {
		err = model.DB().Exec("update aion_player_skill set critical_ratio = ? where skill = ?", res.Ratio, res.Skill).Error
		if err != nil {
			return err
		}
	}
	return nil
}
