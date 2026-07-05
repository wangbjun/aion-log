package service

import (
	"aion/model"
)

func ImportChatLog(fileName string) error {
	cacheService := NewCacheService()

	if err := model.ResetData(); err != nil {
		return err
	}

	parser := NewParseService()
	if err := parser.Run(fileName); err != nil {
		return err
	}
	if err := NewClassifyService().Run(); err != nil {
		return err
	}
	if err := NewRankService().Run(); err != nil {
		return err
	}
	if err := NewTimelineService().Run(); err != nil {
		return err
	}
	return cacheService.Load()
}
