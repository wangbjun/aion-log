package controller

import (
	"aion/model"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"strconv"
	"time"
)

type battleController struct {
	Controller
}

var (
	BattleController = battleController{}
)

type LogResult struct {
	model.Log
	PlayerType  int `json:"player_type"`
	PlayerClass int `json:"player_class"`
	TargetType  int `json:"target_type"`
	TargetClass int `json:"target_class"`
}

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
	if err != nil || pageSize < 0 || pageSize > 1000 {
		pageSize = 1000
	}
	data, count, err := model.Log{}.GetAll(st, et, page, pageSize, queryPlayer, queryTarget, querySkill, sort, value)
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
			Log:         v,
			PlayerType:  playerMap[v.Player].Type,
			PlayerClass: playerMap[v.Player].Class,
		}
		if playerMap[v.Target] != nil {
			result.TargetType = playerMap[v.Target].Type
			result.TargetClass = playerMap[v.Target].Class
		}
		results = append(results, result)
	}
	r.Success(ctx, "ok", map[string]interface{}{"list": results, "total": count})
}

func (r battleController) GetRank(ctx *gin.Context) {
	var (
		st, _    = ctx.GetQuery("st")
		et, _    = ctx.GetQuery("et")
		level, _ = ctx.GetQuery("level")
	)
	data, err := model.Rank{}.GetAll(st, et, level)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	var playerMap = make(map[string]*model.Player)
	players, err := model.Player{}.GetAll()
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	for _, v := range players {
		playerMap[v.Name] = v
	}
	Log := model.Log{}
	for k, v := range data {
		data[k].Type = playerMap[v.Player].Type
		data[k].Class = playerMap[v.Player].Class
		data[k].AllCounts = Log.GetSkillCount(st, et, v.Player)
	}
	r.Success(ctx, "ok", map[string]interface{}{"list": data})
}

type skillCount struct {
	Player string
	Target string
	Count  int
}

var CachedData = cache.New(6*time.Hour, 30*time.Minute)

func (r battleController) GetPlayers(ctx *gin.Context) {
	var key = "all_player_info"
	if cached, found := CachedData.Get(key); found {
		r.Success(ctx, "ok", map[string]interface{}{"list": cached.([]*model.Player)})
		return
	}
	players, err := model.Player{}.GetAll()
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	var result []skillCount
	skillCountSql := "select player,count(1) count from (select player from aion_player_battle_log " +
		"where skill != '' and skill != '普通攻击' group by player,skill,time) t1 group by t1.player"
	err = model.DB().Raw(skillCountSql).Scan(&result).Error
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	var playerSkillCount = make(map[string]int)
	for _, v := range result {
		playerSkillCount[v.Player] = v.Count
	}

	err = model.DB().Raw("select player,target,count(1) count from aion_player_battle_log " +
		"where skill = '' and raw_msg like '%打倒了。' group by player,target").Scan(&result).Error
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

	err = model.DB().Raw("select player,target,count(1) count from aion_player_battle_log " +
		"where skill = '' and raw_msg like '%终结。' group by player,target").Scan(&result).Error
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	for _, v := range result {
		playerKillCount[v.Player] += v.Count
		playerDeathCount[v.Target] += v.Count
	}

	for _, player := range players {
		player.SkillCount = playerSkillCount[player.Name]
		player.KillCount = playerKillCount[player.Name]
		player.DeathCount = playerDeathCount[player.Name]
	}

	CachedData.Set(key, players, cache.DefaultExpiration)

	r.Success(ctx, "ok", map[string]interface{}{"list": players})
}
