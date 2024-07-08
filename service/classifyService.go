package service

import (
	"aion/model"
	"fmt"
	"strings"
)

type ClassifyService struct{}

func NewClassifyService() *ClassifyService {
	return &ClassifyService{}
}

var updateSql = []string{
	"UPDATE aion_player_info SET type = 1 WHERE name like '%永恒之巅'",
	"UPDATE aion_player_info SET type = 1 WHERE name like '%谁与争锋'",
	"UPDATE aion_player_info SET type = 1 WHERE name like '%永恒之岛'",
	"UPDATE aion_player_info SET type = 2 WHERE name like '%火之神殿'",
	"UPDATE aion_player_info SET type = 2 WHERE name like '%傲世八星'",
	"UPDATE aion_player_info SET type = 1 WHERE name in (select distinct player from aion_player_chat_log where " +
		"target in('魔族上级守护神将','魔族结界膜生成师','魔族城门')) and class != 0",
	"UPDATE aion_player_info SET type = 2 WHERE name in (select distinct player from aion_player_chat_log where " +
		"target in('天族上级守护神将','天族结界膜生成师','天族城门')) and class != 0",
}

func (r ClassifyService) Run() error {
	for _, sql := range updateSql {
		err := model.DB().Exec(sql).Error
		if err != nil {
			return err
		}
	}

	err := r.updateBright()
	if err != nil {
		return err
	}

	err = r.updateDark()
	if err != nil {
		return err
	}

	return nil
}

func (r ClassifyService) updateBright() error {
	var result []struct {
		Id int
	}
	sql := "select id from aion_player_info where name in (select distinct(player) name from aion_player_chat_log " +
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
	sql := "select id from aion_player_info where name in (select distinct(player) name from aion_player_chat_log " +
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
