package service

import (
	"aion/model"
	"fmt"
)

type RankService struct{}

func NewRankService() *RankService {
	return &RankService{}
}

func (r RankService) Run() error {
	ranks, err := model.ChatLog{}.GetRanks()
	if err != nil {
		return fmt.Errorf("GetRanks Failed: " + err.Error())
	}

	var items []model.Rank
	for i, v := range ranks {
		items = append(items, v)
		if len(items) >= 500 || i == len(items)-1 {
			err = model.Rank{}.BatchInsert(items)
			if err != nil {
				return fmt.Errorf("BatchInsert Failed: " + err.Error())
			}
			items = make([]model.Rank, 0)
		}
	}
	return nil
}
