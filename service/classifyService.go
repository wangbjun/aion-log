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
