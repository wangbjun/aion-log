package controller

import (
	"aion/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

type battleController struct {
	Controller
}

var (
	BattleController = battleController{}
)

type BattleLogResult struct {
	model.BattleLog
	PlayerType int `json:"player_type"`
	TargetType int `json:"target_type"`
}

func (r battleController) GetAll(ctx *gin.Context) {
	var (
		st, _            = ctx.GetQuery("st")
		et, _            = ctx.GetQuery("et")
		queryPage, _     = ctx.GetQuery("page")
		queryPageSize, _ = ctx.GetQuery("pageSize")
		queryPlayer, _   = ctx.GetQuery("player")
		querySkill, _    = ctx.GetQuery("skill")
		sort, _          = ctx.GetQuery("sort")
	)
	page, err := strconv.Atoi(queryPage)
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(queryPageSize)
	if err != nil || pageSize < 0 || pageSize > 100 {
		pageSize = 100
	}
	data, count, err := model.BattleLog{}.GetAll(st, et, page, pageSize, queryPlayer, querySkill, sort)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	playerMap := make(map[string]int)
	players, err := model.Player{}.GetAll()
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	for _, v := range players {
		playerMap[v.Name] = v.Type
	}
	var results []BattleLogResult
	for _, v := range data {
		results = append(results, BattleLogResult{
			BattleLog:  v,
			PlayerType: playerMap[v.Player],
			TargetType: playerMap[v.Target],
		})
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
	var playerMap = make(map[string]model.Player)
	players, err := model.Player{}.GetAll()
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	for _, v := range players {
		playerMap[v.Name] = v
	}
	battleLog := model.BattleLog{}
	for k, v := range data {
		data[k].Type = playerMap[v.Player].Type
		data[k].Class = playerMap[v.Player].Class
		data[k].AllCounts = battleLog.GetSkillCount(st, et, v.Player)
	}
	r.Success(ctx, "ok", map[string]interface{}{"list": data})
}

func (r battleController) GetPlayers(ctx *gin.Context) {
	counts, err := model.Player{}.GetAll()
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	r.Success(ctx, "ok", map[string]interface{}{"list": counts})
}
