package main

import (
	"aion/config"
	"aion/model"
	"aion/zlog"
	"bufio"
	"errors"
	"flag"
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
	regDeathA  = regexp.MustCompile("(.*?)把(.*?)打倒了")
	regDeathB  = regexp.MustCompile("(.*?)倒下了")
	regDeathC  = regexp.MustCompile("(.*?)受到(.*?)的攻击而死亡")
	regAttackA = regexp.MustCompile("(.*?)给(.*?)造成了(.*)的伤害")
	regAttackB = regexp.MustCompile("(.*?)使用(.*?)技能，[对给](.*?)造成了(.*)的伤害")
)

var classSkills = map[int][]string{
	model.JX: {"飞刀", "飞刀连射", "破灭一击", "杀气破裂", "吸血波", "回旋一击"},
	model.SH: {"主神的惩罚", "处决一击", "天雷斩", "审判", "幻影摄捕", "捕获"},
	model.SX: {"交叉斩", "反击", "影子下坠", "暗破", "暗袭", "灭杀", "猛兽之牙", "背后重击", "进击斩"},
	model.GX: {"利锥箭", "夺取之箭", "套索箭", "强袭箭", "沉默箭", "狂风箭", "百发百中", "破灭箭"},
	model.ZY: {"大地之怒", "大地的惩罚", "惩罚之电", "放电", "断罪一击"},
	model.HF: {"共鸣烟雾", "击破连锁", "必灭重击", "暗击锁", "流星一击", "灭火", "贯穿连锁"},
	model.JL: {"吸引", "幽冥之苦痛", "愤怒之漩涡", "真空爆炸", "精灵弱化", "诅咒之云"},
	model.MD: {"冬季的束缚", "冰河重击", "冷气召唤", "暴风重击", "火焰乱舞", "火焰叉", "结冰"},
}

const WorkerNum = 5

type Parser struct {
	lineChan     chan string
	resultLog    chan model.Log
	resultPlayer chan model.Player
	uniquePlayer map[string]model.Player
	skill2Class  map[string]model.Class
}

func main() {
	var file string
	var conf string
	flag.StringVar(&conf, "c", "../app.ini", "conf file path")
	flag.StringVar(&file, "f", "", "chat log file")
	flag.Parse()
	if file == "" {
		panic("file is empty")
	}
	flag.Parse()

	config.Init(conf)
	zlog.Init()
	model.Init()

	parser := NewParser()
	err := parser.Run(file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("finished process file: %s\n", file)
}

func NewParser() Parser {
	skill2Class := make(map[string]model.Class)
	for class, skills := range classSkills {
		for _, skill := range skills {
			skill2Class[skill] = class
		}
	}
	return Parser{
		lineChan:     make(chan string, 1000),
		resultLog:    make(chan model.Log, 1000),
		resultPlayer: make(chan model.Player, 1000),
		uniquePlayer: make(map[string]model.Player, 1000),
		skill2Class:  skill2Class,
	}
}

func (r *Parser) Run(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//启动worker
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup

	for i := 0; i < WorkerNum; i++ {
		wg.Add(1)
		go r.worker(&wg)
	}
	wg2.Add(1)
	go r.processResult(&wg2)

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
		r.lineChan <- text
		time.Sleep(time.Microsecond * 100)
	}
	close(r.lineChan)
	wg.Wait()
	close(r.resultLog)
	close(r.resultPlayer)
	wg2.Wait()

	for _, player := range r.uniquePlayer {
		fmt.Printf("%v\n", player)
	}

	return nil
}

func (r *Parser) worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range r.lineChan {
		var err error
		if regAttackA.Match([]byte(line)) {
			err = r.parseAttackA(line)
		} else if regAttackB.Match([]byte(line)) {
			err = r.parseAttackB(line)
		} else if regDeathA.Match([]byte(line)) {
			err = r.parseDeathA(line)
		} else if regDeathB.Match([]byte(line)) {
			err = r.parseDeathB(line)
		} else if regDeathC.Match([]byte(line)) {
			err = r.parseDeathC(line)
		}
		if err != nil {
			fmt.Printf("parse line error: %s\n", err)
		}
	}
}

func (r *Parser) processResult(wg *sync.WaitGroup) {
	defer wg.Done()
	doneLog := false
	donePlayer := false

	var logItems []model.Log
	for {
		select {
		case log, ok := <-r.resultLog:
			if !ok {
				doneLog = true
			} else {
				logItems = append(logItems, log)
				if len(logItems) >= 500 {
					model.BattleLog{}.BatchInsert(logItems)
					logItems = []model.Log{}
				}
			}
		case player, ok := <-r.resultPlayer:
			if !ok {
				donePlayer = true
			} else {
				if exsited, ok := r.uniquePlayer[player.Name]; ok {
					if player.Type != 0 {
						exsited.Type = player.Type
					}
					if player.Class != 0 {
						exsited.Class = player.Class
					}
					r.uniquePlayer[player.Name] = exsited
				} else {
					r.uniquePlayer[player.Name] = player
				}
			}
		}
		if doneLog && donePlayer {
			model.BattleLog{}.BatchInsert(logItems)
			return
		}
	}
}

func (r *Parser) parseAttackA(line string) error {
	match := regAttackA.FindStringSubmatch(line[22:])
	if len(match) != 4 {
		return errors.New("ParseAttackA matches fail:" + line)
	}
	player := strings.ReplaceAll(match[1], "致命一击！", "")
	if player == "" {
		player = "我"
	}
	var (
		damage, _ = strconv.Atoi(strings.ReplaceAll(match[3], ",", ""))
		timeStr   = strings.ReplaceAll(line[0:19], ".", "-")
	)
	r.resultLog <- model.Log{
		Player: player,
		Skill:  "普通攻击",
		Target: match[2],
		Value:  damage,
		Time:   timeStr,
		RawMsg: line[22:],
	}
	r.resultPlayer <- model.Player{
		Name: player,
		Time: timeStr,
	}
	r.resultPlayer <- model.Player{
		Name: match[2],
		Time: timeStr,
	}
	return nil
}

func (r *Parser) parseAttackB(line string) error {
	match := regAttackB.FindStringSubmatch(line[22:])
	if len(match) != 5 {
		return errors.New("parseAttackB matches fail:" + line)
	}
	player := strings.ReplaceAll(match[1], "致命一击！", "")
	if player == "" {
		player = "我"
	}
	var (
		damage, _ = strconv.Atoi(strings.ReplaceAll(match[4], ",", ""))
		timeStr   = strings.ReplaceAll(line[0:19], ".", "-")
	)
	r.resultLog <- model.Log{
		Player: player,
		Skill:  match[2],
		Target: match[3],
		Value:  damage,
		Time:   timeStr,
		RawMsg: line[22:],
	}
	r.resultPlayer <- model.Player{
		Name:  player,
		Class: r.skill2Class[removeSkillLevel(match[2])],
		Time:  timeStr,
	}
	r.resultPlayer <- model.Player{
		Name: match[3],
		Time: timeStr,
	}
	return nil
}

// A把b打倒了
func (r *Parser) parseDeathA(line string) error {
	match := regDeathA.FindStringSubmatch(line[22:])
	if len(match) != 3 {
		return errors.New("parseDeathA matches fail:" + line)
	}
	timeStr := strings.ReplaceAll(line[0:19], ".", "-")
	r.resultPlayer <- model.Player{
		Name: match[1],
		Type: model.TypeTian,
		Time: timeStr,
	}
	r.resultPlayer <- model.Player{
		Name: match[2],
		Type: model.TypeMo,
		Time: timeStr,
	}
	return nil
}

func (r *Parser) parseDeathB(line string) error {
	match := regDeathB.FindStringSubmatch(line[22:])
	if len(match) != 2 {
		return errors.New("parseDeathB matches fail:" + line)
	}
	timeStr := strings.ReplaceAll(line[0:19], ".", "-")
	r.resultPlayer <- model.Player{
		Name: match[1],
		Type: model.TypeMo,
		Time: timeStr,
	}
	return nil
}

func (r *Parser) parseDeathC(line string) error {
	match := regDeathC.FindStringSubmatch(line[22:])
	if len(match) != 3 {
		return errors.New("parseDeathC matches fail:" + line)
	}
	timeStr := strings.ReplaceAll(line[0:19], ".", "-")
	if match[1] != "" {
		r.resultPlayer <- model.Player{
			Name: match[1],
			Type: model.TypeTian,
			Time: timeStr,
		}
	}
	r.resultPlayer <- model.Player{
		Name: match[2],
		Type: model.TypeMo,
		Time: timeStr,
	}
	return nil
}

func removeSkillLevel(skill string) string {
	split := strings.Split(skill, " ")
	if len(split) == 2 {
		return split[0]
	} else {
		return skill
	}
}
