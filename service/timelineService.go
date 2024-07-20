package service

import (
	"aion/model"
)

type TimelineService struct{}

func NewTimelineService() *TimelineService {
	return &TimelineService{}
}

func (r TimelineService) Run() error {
	err := r.runKill()
	if err != nil {
		return err
	}
	err = r.runKilled()
	if err != nil {
		return err
	}
	return nil
}

func (r TimelineService) runKill() error {
	var result []model.Timeline
	sql := "select time, sum(count) as value, 1 as type from (select time, count(1) count from " +
		"aion_chat_log where skill = 'kill' group by time) t1 group by time order by time asc"
	err := model.DB().Raw(sql).Find(&result).Error
	if err != nil {
		return err
	}
	return model.Timeline{}.BatchInsert(result)
}

func (r TimelineService) runKilled() error {
	var result []model.Timeline
	sql := "select time, sum(count) as value, 2 as type from (select time, count(1) count from " +
		"aion_chat_log where skill = 'killed' group by time) t1 group by time order by time asc"
	err := model.DB().Raw(sql).Find(&result).Error
	if err != nil {
		return err
	}
	return model.Timeline{}.BatchInsert(result)
}
