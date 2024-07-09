package controller

import (
	"aion/model"
	"aion/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
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
		st, _            = ctx.GetQuery("st")
		et, _            = ctx.GetQuery("et")
		queryPage, _     = ctx.GetQuery("page")
		queryPageSize, _ = ctx.GetQuery("pageSize")
		queryPlayer, _   = ctx.GetQuery("player")
		queryTarget, _   = ctx.GetQuery("target")
		querySkill, _    = ctx.GetQuery("skill")
		sort, _          = ctx.GetQuery("sort")
		value, _         = ctx.GetQuery("value")
	)
	page, err := strconv.Atoi(queryPage)
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(queryPageSize)
	if err != nil || pageSize < 0 || pageSize > 500 {
		pageSize = 500
	}
	data, count, err := model.ChatLog{}.GetAll(st, et, page, pageSize, queryPlayer, queryTarget, querySkill, sort, value)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	playerMap := make(map[string]*model.Player)
	players, err := model.Player{}.GetAll()
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	for _, v := range players {
		playerMap[v.Name] = v
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
	if cached, ok := r.cache.GetRank(level); ok {
		r.Success(ctx, "ok", map[string]interface{}{"list": cached})
		return
	}
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
	r.Success(ctx, "ok", map[string]interface{}{"list": data})
}

// GetPlayers 获取所有玩家
func (r battleController) GetPlayers(ctx *gin.Context) {
	var (
		st, _ = ctx.GetQuery("st")
		et, _ = ctx.GetQuery("et")
	)
	if st == "" || et == "" {
		r.Failed(ctx, ParamError, "请选择时间范围")
		return
	}
	var key = "player_" + st + "_" + et
	if cached, ok := r.cache.GetPlayers(key); ok {
		r.Success(ctx, "ok", map[string]interface{}{"list": cached})
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
			v.Time = existed.Time
		}
	}
	var result []struct {
		Player string
		Target string
		Count  int
	}
	skillCountSql := fmt.Sprintf("select player,count(DISTINCT skill, time) as count from aion_player_chat_log " +
		"where skill not in ('','kill','killed') and time >= '" + st + "' AND time <= '" + et + "' group by player")
	err = model.DB().Raw(skillCountSql).Find(&result).Error
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	var playerSkillCount = make(map[string]int)
	for _, v := range result {
		playerSkillCount[v.Player] = v.Count
	}

	err = model.DB().Raw("select player,target,count(1) count from aion_player_chat_log " +
		"where skill = 'kill' and time >= '" + st + "' AND time <= '" + et + "' group by player,target").Find(&result).Error
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	var playerKillCount = make(map[string]int)
	var playerDeathCount = make(map[string]int)
	for _, v := range result {
		playerKillCount[v.Player] += v.Count
		playerDeathCount[v.Target] += v.Count
	}

	err = model.DB().Raw("select player,target,count(1) count from aion_player_chat_log " +
		"where skill = 'killed' and time >= '" + st + "' AND time <= '" + et + "' group by player,target").Find(&result).Error
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	for _, v := range result {
		playerKillCount[v.Target] += v.Count
		playerDeathCount[v.Player] += v.Count
	}

	for _, player := range players {
		player.SkillCount = playerSkillCount[player.Name]
		player.KillCount = playerKillCount[player.Name]
		player.DeathCount = playerDeathCount[player.Name]
	}

	if len(players) > 0 {
		r.cache.SetPlayers(key, players)
	}
	r.Success(ctx, "ok", map[string]interface{}{"list": players})
}
