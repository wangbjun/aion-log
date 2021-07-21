package service

import (
	"aion/model"
	"aion/zlog"
	"fmt"
	"time"
)

type RankService struct{}

func (r RankService) Start() {
	r.run()
	var t = time.NewTicker(time.Hour * 6)
	for {
		select {
		case <-t.C:
			r.run()
		}
	}
}

func (r RankService) run() {
	defer func() {
		if err := recover(); err != nil {
			zlog.Logger.Error(fmt.Sprintf("RankService Run Error: %s", err))
		}
	}()
	ranks, err := model.Rank{}.GetRanks()
	if err != nil {
		zlog.Logger.Error("GetRanks Failed: " + err.Error())
		return
	}
	zlog.Logger.Error(fmt.Sprintf("GetRanks: %d个", len(ranks)))
	for _, v := range ranks {
		_ = v.Save()
		time.Sleep(time.Millisecond * 100)
	}
}
