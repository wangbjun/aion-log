package service

import (
	"aion/model"
	"fmt"
	"time"
)

type RankService struct{}

func NewRankService() *RankService {
	return &RankService{}
}

func (r RankService) Run() error {
	err := model.DB().Exec("truncate table aion_player_rank").Error
	if err != nil {
		return fmt.Errorf("clean table error:" + err.Error())
	}
	ranks, err := model.Rank{}.GetRanks()
	if err != nil {
		return fmt.Errorf("GetRanks Failed: " + err.Error())
	}
	for _, v := range ranks {
		_ = v.Save()
		time.Sleep(time.Millisecond * 100)
	}
	return nil
}
