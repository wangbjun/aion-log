package service

import (
	"aion/model"
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
	"UPDATE aion_player_info SET type = 1 WHERE name in (select distinct player from aion_player_battle_log where " +
		"target in('魔族上级守护神将','魔族结界膜生成师','魔族城门')) and class != 0",
}

func (r ClassifyService) Run() error {
	for _, sql := range updateSql {
		err := model.DB().Exec(sql).Error
		if err != nil {
			return err
		}
	}
	return nil
}
