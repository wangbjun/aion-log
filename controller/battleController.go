package controller

import (
	"aion/model"
	"aion/service"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type battleController struct {
	Controller
	cache *service.CacheService
}

var (
	BattleController = battleController{
		cache: service.NewCacheService(),
	}
)

type LogResult struct {
	model.ChatLog
	PlayerType  int `json:"player_type"`
	PlayerClass int `json:"player_class"`
	TargetType  int `json:"target_type"`
	TargetClass int `json:"target_class"`
}

// GetAll 获取所有日志
func (r battleController) GetAll(ctx *gin.Context) {
	var (
		st, _        = ctx.GetQuery("st")
		et, _        = ctx.GetQuery("et")
		page, _      = ctx.GetQuery("page")
		pageSize, _  = ctx.GetQuery("pageSize")
		player, _    = ctx.GetQuery("player")
		target, _    = ctx.GetQuery("target")
		skill, _     = ctx.GetQuery("skill")
		sorter, _    = ctx.GetQuery("sort")
		value, _     = ctx.GetQuery("value")
		banPlayer, _ = ctx.GetQuery("banPlayer")
	)
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt < 0 || pageSizeInt > 500 {
		pageSizeInt = 500
	}
	data, count, err := model.ChatLog{}.GetAll(st, et, pageInt, pageSizeInt, player, target, skill, sorter, value, banPlayer)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	var results []LogResult
	for _, v := range data {
		result := LogResult{
			ChatLog: v,
		}
		if cached, ok := r.cache.GetPlayer(v.Player); ok {
			result.PlayerType = cached.Type
			result.PlayerClass = cached.Class
		}
		if cached, ok := r.cache.GetPlayer(v.Target); ok {
			result.TargetType = cached.Type
			result.TargetClass = cached.Class
		}
		results = append(results, result)
	}
	r.Success(ctx, "ok", map[string]interface{}{"list": results, "total": count})
}

// GetRank 获取封神榜
func (r battleController) GetRank(ctx *gin.Context) {
	var (
		level, _ = ctx.GetQuery("level")
	)
	data, err := model.Rank{}.GetAll(level)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}

	for k, v := range data {
		if cached, ok := r.cache.GetPlayer(v.Player); ok {
			data[k].Type = cached.Type
			data[k].Class = cached.Class
		}
	}
	r.Success(ctx, "ok", data)
}

// GetPlayers 获取时间段内的玩家
func (r battleController) GetPlayers(ctx *gin.Context) {
	var (
		st, _ = ctx.GetQuery("st")
		et, _ = ctx.GetQuery("et")
	)

	if st == "" || et == "" {
		players, err := model.Player{}.GetAll()
		if err != nil {
			r.Failed(ctx, Failed, err.Error())
			return
		}
		r.Success(ctx, "ok", players)
		return
	}

	players, err := model.Player{}.GetByTime(st, et)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}

	for _, v := range players {
		if existed, ok := r.cache.GetPlayer(v.Name); ok {
			v.Id = existed.Id
			v.Type = existed.Type
			v.Class = existed.Class
			v.CriticalRatio = existed.CriticalRatio
			v.Time = existed.Time
		}
	}

	if err := (model.Player{}).FillStats(players, st, et); err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	r.Success(ctx, "ok", players)
}

func (r battleController) GetTimeline(ctx *gin.Context) {
	var (
		st, _ = ctx.GetQuery("st")
		et, _ = ctx.GetQuery("et")
	)
	killTime, err := model.Timeline{}.GetAll(st, et, 1)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	killedTime, err := model.Timeline{}.GetAll(st, et, 2)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	killTimes := mergeTimes(killTime, killedTime)

	var times []string
	var killValue []int
	var killedValue []int
	for _, t := range killTimes {
		times = append(times, t.Format(time.DateTime))
		killValue = append(killValue, getValue(killTime, t))
		killedValue = append(killedValue, getValue(killedTime, t))
	}
	r.Success(ctx, "ok", map[string]interface{}{
		"timeData":    times,
		"killValue":   killValue,
		"killedValue": killedValue,
	})
}

func (r battleController) GetClassTop(ctx *gin.Context) {
	var (
		class, _  = ctx.GetQuery("class")
		player, _ = ctx.GetQuery("player")
	)
	if class == "" {
		r.Failed(ctx, ParamError, "请选择职业")
		return
	}
	key := fmt.Sprintf("%s_%s", class, player)
	if cached, ok := r.cache.GetClassTop(key); ok {
		r.Success(ctx, "ok", cached)
		return
	}
	result, err := model.ChatLog{}.GetClassTop(class, player)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	for _, res := range result {
		if skill, ok := r.cache.GetSkill(res.Skill); ok {
			res.Critical = skill.CriticalRatio
		}
	}
	if len(result) != 0 {
		r.cache.SetClassTop(key, result)
	}
	r.Success(ctx, "ok", result)
}

func (r battleController) ImportLog(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		r.Failed(ctx, ParamError, "missing log file")
		return
	}

	uploadDir := filepath.Join("storage", "uploads")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}

	saveName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	savePath := filepath.Join(uploadDir, saveName)
	if err := ctx.SaveUploadedFile(file, savePath); err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}

	err = service.ImportChatLog(savePath)
	if err != nil {
		r.Failed(ctx, Failed, "importChatLog failed: "+err.Error())
		return
	}

	r.Success(ctx, "ok", nil)
}

func mergeTimes(data1, data2 []model.Timeline) []time.Time {
	timeSet := make(map[time.Time]struct{})
	for _, dp := range data1 {
		timeSet[dp.Time] = struct{}{}
	}
	for _, dp := range data2 {
		timeSet[dp.Time] = struct{}{}
	}
	mergedTimes := make([]time.Time, 0, len(timeSet))
	for t := range timeSet {
		mergedTimes = append(mergedTimes, t)
	}
	sort.Slice(mergedTimes, func(i, j int) bool {
		return mergedTimes[i].Before(mergedTimes[j])
	})
	return mergedTimes
}

func getValue(data []model.Timeline, targetTime time.Time) int {
	for _, dp := range data {
		if dp.Time.Equal(targetTime) {
			return dp.Value
		}
	}
	return 0
}
