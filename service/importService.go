package service

import (
	"aion/model"
	"io"
)

func ImportChatLog(reader io.Reader, sourceName string) error {
	cacheService := NewCacheService()

	UpdateImportStatus("reset", "正在清空旧统计数据", 10)
	if err := model.ResetData(); err != nil {
		return err
	}
	if err := cacheService.Load(); err != nil {
		return nil
	}

	UpdateImportStatus("parse", "正在解析战斗日志", 25)
	parser := NewParseService()
	if err := parser.RunReader(reader, sourceName); err != nil {
		return err
	}
	UpdateImportStatus("classify", "正在识别玩家职业和阵营", 65)
	if err := NewClassifyService().Run(); err != nil {
		return err
	}
	UpdateImportStatus("rank", "正在生成异常分析排行", 78)
	if err := NewRankService().Run(); err != nil {
		return err
	}
	UpdateImportStatus("timeline", "正在生成战斗时间线", 88)
	if err := NewTimelineService().Run(); err != nil {
		return err
	}
	UpdateImportStatus("cache", "正在刷新缓存", 96)
	return cacheService.Load()
}
