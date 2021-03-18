package service

import (
	"aion/model"
	"aion/zlog"
	"time"
)

type Ranking struct{}

func (r Ranking) Start() {
	var t = time.NewTicker(time.Millisecond * 10)
	for {
		select {
		case <-t.C:
			r.run()
		}
	}
}

func (r Ranking) run() {
	ranks, err := model.Rank{}.GetRanks()
	if err != nil {
		zlog.Logger.Error("GetRanks Failed: " + err.Error())
		return
	}
	for _, v := range ranks {
		_ = v.Save()
		time.Sleep(time.Millisecond * 200)
	}
}
