package service

import (
	"aion/model"
	"aion/zlog"
	"fmt"
	"time"
)

type CleanService struct{}

func (r CleanService) Start() {
	var t = time.NewTicker(time.Hour * 1)
	for {
		select {
		case <-t.C:
			r.run()
		}
	}
}

func (r CleanService) run() {
	defer func() {
		if err := recover(); err != nil {
			zlog.Logger.Error(fmt.Sprintf("RankService Run Error: %s", err))
		}
	}()
	days := model.BattleLog{}.GetWeekly()
	if len(days) <= 7 {
		return
	}
	err := model.BattleLog{}.DeleteByDay(days[7].Day)
	if err != nil {
		zlog.Logger.Error("delete battle log error")
	}
	err = model.Player{}.DeleteByDay(days[7].Day)
	if err != nil {
		zlog.Logger.Error("delete player error")
	}
	err = model.Rank{}.DeleteByDay(days[7].Day)
	if err != nil {
		zlog.Logger.Error("delete rank error")
	}
}
