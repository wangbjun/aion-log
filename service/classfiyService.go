package service

import (
	"aion/zlog"
	"fmt"
	"time"
)

type ClassfiyService struct{}

func (r ClassfiyService) Start() {
	r.run()
	var t = time.NewTicker(time.Hour * 8)
	for {
		select {
		case <-t.C:
			r.run()
		}
	}
}

func (r ClassfiyService) run() {
	defer func() {
		if err := recover(); err != nil {
			zlog.Logger.Error(fmt.Sprintf("RankService Run Error: %s", err))
		}
	}()
}
