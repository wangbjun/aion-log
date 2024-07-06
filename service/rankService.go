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
