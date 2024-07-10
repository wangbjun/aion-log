package service

import (
	"aion/model"
	"fmt"
	"time"
)

type TimelineService struct{}

func NewTimelineService() *TimelineService {
	return &TimelineService{}
}

func (r TimelineService) Run() error {
	var result []struct {
		Tb    time.Time
		Value int
	}
	err := model.DB().Exec("truncate table aion_timeline").Error
	if err != nil {
		return fmt.Errorf("clean table error:" + err.Error())
	}
	sql := "select from_unixtime(floor(unix_timestamp(time) / 10) * 10) as tb, sum(count) as value from (" +
		"select time, count(1) as count from ( select time from aion_chat_log where target != '' and skill != 'attack' " +
		"and value > 0 group by time, skill) t1 group by time) t2 group by tb order by tb"
	err = model.DB().Raw(sql).Scan(&result).Error
	if err != nil {
		return err
	}
	var items []model.Timeline
	for _, item := range result {
		items = append(items, model.Timeline{
			Time:  item.Tb,
			Value: item.Value,
		})
	}
	return model.Timeline{}.BatchInsert(items)
}
