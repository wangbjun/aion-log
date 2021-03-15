package controller

import (
	"aion/config"
	"aion/model"
	"aion/service"
	"aion/util"
	"bufio"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type battleController struct {
	Controller
}

var (
	BattleController = battleController{}
	taskManager      = service.NewParser()
)

func init() {
	go taskManager.Start()
}

type BattleLogResult struct {
	model.BattleLog
	PlayerType       int `json:"player_type"`
	TargetPlayerType int `json:"target_player_type"`
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
		pageSize = 20
	}
	data, count, err := model.BattleLog{}.GetAll(st, et, page, pageSize, queryPlayer, querySkill, sort)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	playerMap := make(map[string]int)
	players, err := model.Player{}.GetPlayers()
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
			BattleLog:        v,
			PlayerType:       playerMap[v.Player],
			TargetPlayerType: playerMap[v.TargetPlayer],
		})
	}
	r.Success(ctx, "ok", map[string]interface{}{"list": results, "total": count})
}

type RankResult struct {
	Player string `json:"player"`
	Type   int    `json:"type"`
	Count  int    `json:"count"`
	Times  string `json:"times"`
}

func (r battleController) GetRank(ctx *gin.Context) {
	var (
		st, _    = ctx.GetQuery("st")
		et, _    = ctx.GetQuery("et")
		level, _ = ctx.GetQuery("level")
	)
	if level == "" {
		level = "4"
	}
	data, err := model.BattleLog{}.GetRank(st, et, level)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	var (
		countMap  = make(map[string]int)
		timeMap   = make(map[string]string)
		playerMap = make(map[string]int)
	)
	for _, v := range data {
		countMap[v.Player]++
		if existed, ok := timeMap[v.Player]; ok {
			timeMap[v.Player] = existed + "," + v.Time.Format(util.TimeFormat)
		} else {
			timeMap[v.Player] = v.Time.Format(util.TimeFormat)
		}
	}
	players, err := model.Player{}.GetPlayers()
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	for _, v := range players {
		playerMap[v.Name] = v.Type
	}
	var result []RankResult
	for k, v := range countMap {
		result = append(result, RankResult{
			Player: k,
			Count:  v,
			Times:  timeMap[k],
			Type:   playerMap[k],
		})
	}
	r.Success(ctx, "ok", map[string]interface{}{"list": result})
}

func (r battleController) GetStat(ctx *gin.Context) {
	var (
		st, _ = ctx.GetQuery("st")
		et, _ = ctx.GetQuery("et")
	)
	data, err := model.BattleLog{}.GetStat(st, et)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	r.Success(ctx, "ok", map[string]interface{}{"total": data})
}

func (r battleController) GetPlayers(ctx *gin.Context) {
	counts, err := model.Player{}.GetPlayers()
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	r.Success(ctx, "ok", map[string]interface{}{"list": counts})
}

func (r battleController) ChangePlayerType(ctx *gin.Context) {
	var (
		id, _ = ctx.GetQuery("id")
		t, _  = ctx.GetQuery("type")
	)
	err := model.Player{}.ChangeType(id, t)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	r.Success(ctx, "ok", nil)
}

func (r battleController) GetTask(ctx *gin.Context) {
	r.Success(ctx, "ok", taskManager.IsRuning())
}

func (r battleController) AddTask(ctx *gin.Context) {
	if taskManager.IsRuning() {
		r.Failed(ctx, Failed, "当前有任务正在运行，请稍后再试！")
		return
	}
	file, err := ctx.FormFile("file")
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	if file.Size == 0 {
		r.Failed(ctx, ParamError, "文件不能为空")
		return
	}
	openFile, err := file.Open()
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	defer openFile.Close()
	firstLine, err := bufio.NewReader(openFile).ReadString('\n')
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	var asiaShanghai, _ = time.LoadLocation("Asia/Shanghai")
	t, err := time.ParseInLocation(util.TimeFormat, strings.ReplaceAll(firstLine[0:19], ".", "-"), asiaShanghai)
	if err != nil {
		r.Failed(ctx, Failed, "此文件格式错误")
		return
	}
	lastTime := model.BattleLog{}.GetLastTime()
	if t.Before(lastTime) {
		r.Failed(ctx, ParamError, "此日志已经解析过！")
		return
	}
	savePath := config.Conf.Section("APP").Key("UPLOAD_DIR").String() + "/" + file.Filename
	err = ctx.SaveUploadedFile(file, savePath)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
		return
	}
	if taskManager.IsRuning() {
		r.Failed(ctx, Failed, "当前有任务正在运行，请稍后再试！")
		return
	}
	go taskManager.Add(savePath)
	r.Success(ctx, "ok", nil)
}
