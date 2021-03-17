package service

import (
	"aion/model"
	"aion/util"
	"aion/zlog"
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	regKillA   = regexp.MustCompile("(.*?)把(.*?)打倒了")
	regKillB   = regexp.MustCompile("(.*?)倒下了")
	regKillC   = regexp.MustCompile("(.*?)受到(.*?)的攻击而死亡")
	regDamageA = regexp.MustCompile("(.*?)给(.*?)造成了(.*)的伤害")
	regDamageB = regexp.MustCompile("(.*?)使用(.*?)技能，[对|给](.*?)造成了(.*)的伤害")
)

var asiaShanghai, _ = time.LoadLocation("Asia/Shanghai")

type Parser struct {
	fileChan   chan string
	saveChan   chan model.Log
	playerType map[string]model.Player
	mapLocker  *sync.Mutex
	isRuning   bool
}

func NewParser() Parser {
	return Parser{
		fileChan:   make(chan string, 1),
		saveChan:   make(chan model.Log, 1),
		playerType: make(map[string]model.Player),
		mapLocker:  new(sync.Mutex),
		isRuning:   false,
	}
}

func (r Parser) IsRuning() bool {
	return r.isRuning
}

func (r *Parser) Start() {
	go r.saveBatch()
	go r.savePlayer()
	for {
		select {
		case fileName := <-r.fileChan:
			if fileName == "" {
				continue
			}
			r.isRuning = true
			zlog.Logger.Sugar().Infof("开始解析日志文件： %s", fileName)
			err := r.Run(fileName)
			if err != nil {
				zlog.Logger.Sugar().Errorf("日志解析运行失败： %s", err)
			}
			r.isRuning = false
			zlog.Logger.Sugar().Infof("结束解析日志文件： %s", fileName)
		}
	}
}

func (r Parser) Add(fileName string) {
	r.fileChan <- fileName
}

func (r Parser) Run(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := simplifiedchinese.GBK.NewDecoder()
	buff := bufio.NewReader(file)
	for {
		a, _, err := buff.ReadLine()
		if err == io.EOF {
			break
		}
		line, err := decoder.Bytes(a)
		if err != nil {
			continue
		}
		if len(line) == 0 {
			continue
		}
		text := string(line)
		if strings.Contains(text, "反弹了攻击") {
			continue
		}
		if regDamageB.Match(line) {
			r.parseDamageB(text)
		} else if regDamageA.Match(line) {
			r.parseDamageA(text)
		} else if regKillA.Match(line) {
			r.parseKillA(text)
		} else if regKillB.Match(line) {
			r.parseKillB(text)
		} else if regKillB.Match(line) {
			r.parseKillC(text)
		}
		time.Sleep(time.Microsecond * 500)
	}
	return nil
}

func (r Parser) parseDamageA(line string) {
	if strings.Count(line[22:], "给") >= 2 {
		return
	}
	match := regDamageA.FindStringSubmatch(line[22:])
	if len(match) != 4 {
		return
	}
	player := strings.ReplaceAll(match[1], "致命一击！", "")
	if player == "" {
		player = "我"
	}
	var (
		damage, _ = strconv.Atoi(strings.ReplaceAll(match[3], ",", ""))
		t, _      = time.ParseInLocation(util.TimeFormat,
			strings.ReplaceAll(line[0:19], ".", "-"), asiaShanghai)
	)
	r.saveChan <- model.Log{
		Player:       player,
		Skill:        "普通攻击",
		TargetPlayer: match[2],
		Damage:       damage,
		Time:         t,
		OriginDesc:   line[22:],
	}
	r.mapLocker.Lock()
	r.playerType[player] = model.Player{
		Name: player,
		Time: t,
	}
	r.playerType[match[2]] = model.Player{
		Name: match[2],
		Time: t,
	}
	r.mapLocker.Unlock()
}

func (r Parser) parseDamageB(line string) {
	match := regDamageB.FindStringSubmatch(line[22:])
	if len(match) != 5 {
		return
	}
	player := strings.ReplaceAll(match[1], "致命一击！", "")
	if player == "" {
		player = "我"
	}
	var (
		damage, _ = strconv.Atoi(strings.ReplaceAll(match[4], ",", ""))
		t, _      = time.ParseInLocation(util.TimeFormat,
			strings.ReplaceAll(line[0:19], ".", "-"), asiaShanghai)
	)
	r.saveChan <- model.Log{
		Player:       player,
		Skill:        match[2],
		TargetPlayer: match[3],
		Damage:       damage,
		Time:         t,
		OriginDesc:   line[22:],
	}
	r.mapLocker.Lock()
	r.playerType[player] = model.Player{
		Name: player,
		Time: t,
	}
	r.playerType[match[3]] = model.Player{
		Name: match[3],
		Time: t,
	}
	r.mapLocker.Unlock()
}

func (r Parser) parseKillA(line string) {
	if strings.Count(line[22:], "把") >= 2 {
		return
	}
	match := regKillA.FindStringSubmatch(line[22:])
	if len(match) != 3 {
		return
	}
	t, _ := time.ParseInLocation(util.TimeFormat,
		strings.ReplaceAll(line[0:19], ".", "-"), asiaShanghai)
	r.mapLocker.Lock()
	r.playerType[match[1]] = model.Player{
		Name: match[1],
		Type: model.TypeTian,
		Time: t,
	}
	r.playerType[match[2]] = model.Player{
		Name: match[2],
		Type: model.TypeMo,
		Time: t,
	}
	r.mapLocker.Unlock()
}

func (r Parser) parseKillB(line string) {
	match := regKillB.FindStringSubmatch(line[22:])
	if len(match) != 2 {
		return
	}
	t, _ := time.ParseInLocation(util.TimeFormat,
		strings.ReplaceAll(line[0:19], ".", "-"), asiaShanghai)
	r.mapLocker.Lock()
	r.playerType[match[1]] = model.Player{
		Name: match[1],
		Type: model.TypeMo,
		Time: t,
	}
	r.mapLocker.Unlock()
}

func (r Parser) parseKillC(line string) {
	match := regKillC.FindStringSubmatch(line[22:])
	if len(match) != 3 {
		return
	}
	t, _ := time.ParseInLocation(util.TimeFormat,
		strings.ReplaceAll(line[0:19], ".", "-"), asiaShanghai)
	r.mapLocker.Lock()
	r.playerType[match[2]] = model.Player{
		Name: match[2],
		Type: model.TypeTian,
		Time: t,
	}
	r.playerType[match[1]] = model.Player{
		Name: match[1],
		Type: model.TypeMo,
		Time: t,
	}
	r.mapLocker.Unlock()
}

func (r Parser) saveBatch() {
	var cached []model.Log
	for {
		if len(cached) >= 500 {
			model.BattleLog{}.BacthInsert(cached)
			cached = []model.Log{}
		}
		item := <-r.saveChan
		cached = append(cached, item)
	}
}

func (r Parser) savePlayer() {
	var t = time.NewTicker(time.Second * 60)
	for {
		select {
		case <-t.C:
			if len(r.playerType) == 0 {
				continue
			}
			r.mapLocker.Lock()
			tmpItems := make(map[string]model.Player)
			for k, v := range r.playerType {
				tmpItems[k] = v
			}
			r.playerType = make(map[string]model.Player)
			r.mapLocker.Unlock()
			zlog.Logger.Info(fmt.Sprintf("begin save player: %d", len(tmpItems)))
			for _, v := range tmpItems {
				err := v.SaveType()
				if err != nil {
					zlog.Logger.Error("save player error: " + err.Error())
				}
			}
			zlog.Logger.Info("after save player")
		}
	}
}
