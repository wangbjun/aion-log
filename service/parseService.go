package service

import (
	"aion/model"
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	regDeathA  = regexp.MustCompile("(\\S+)把(\\S+)打倒了。")
	regDeathB  = regexp.MustCompile("(\\S+)倒下了。")
	regDeathC  = regexp.MustCompile("(\\S+)受到(\\S+)的攻击而终结。")
	regAttackA = regexp.MustCompile("(\\S+)使用(.+)技能，对(\\S+)造成了(\\S+)的伤害")
	regAttackB = regexp.MustCompile("(\\S+)给(\\S+)造成了(\\S+)的伤害")
	regAttackC = regexp.MustCompile("(\\S+)使用(.+)技能，")
	regValue   = regexp.MustCompile("(\\d{1,3}(,\\d{3})*)")
)

const WorkerNum = 8

type Parser struct {
	lineChan     chan string
	resultLog    chan model.ChatLog
	resultPlayer chan model.Player
	uniquePlayer map[string]model.Player
	skill2Class  map[string]model.Class
}

func NewParseService() Parser {
	playerSkill, err := model.PlayerSkill{}.GetAll()
	if err != nil {
		return Parser{}
	}
	skill2Class := make(map[string]model.Class)
	for _, skill := range playerSkill {
		skill2Class[skill.Skill] = skill.Class
	}
	return Parser{
		lineChan:     make(chan string, 1000),
		resultLog:    make(chan model.ChatLog, 1000),
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
			log.Printf("decoder error: %s", err)
			continue
		}
		text := string(line)
		if len(line) == 0 {
			continue
		}
		r.lineChan <- text
	}
	close(r.lineChan)
	wg.Wait()
	close(r.resultLog)
	close(r.resultPlayer)
	wg2.Wait()
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
		} else if regAttackC.Match([]byte(line)) {
			err = r.parseAttackC(line)
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

	var logItems []model.ChatLog
	for {
		select {
		case chatLog, ok := <-r.resultLog:
			if !ok {
				doneLog = true
			} else {
				logItems = append(logItems, chatLog)
				if len(logItems) > 500 {
					model.ChatLog{}.BatchInsert(logItems)
					logItems = []model.ChatLog{}
				}
			}
		case player, ok := <-r.resultPlayer:
			if !ok {
				donePlayer = true
			} else {
				if existed, ok := r.uniquePlayer[player.Name]; ok {
					if existed.Type == 0 {
						existed.Type = player.Type
					}
					if existed.Class == 0 {
						existed.Class = player.Class
					}
					existed.Time = player.Time
					r.uniquePlayer[player.Name] = existed
				} else {
					r.uniquePlayer[player.Name] = player
				}
			}
		}
		if doneLog && donePlayer {
			model.ChatLog{}.BatchInsert(logItems)
			var result []model.Player
			for _, player := range r.uniquePlayer {
				result = append(result, player)
				if len(result) > 500 {
					model.Player{}.BatchInsert(result)
					result = make([]model.Player, 0)
				}
			}
			model.Player{}.BatchInsert(result)
			return
		}
	}
}

// (.*?)使用(.*?)技能，对(.*?)造成了(.*)的伤害
func (r *Parser) parseAttackA(line string) error {
	match := regAttackA.FindStringSubmatch(line[22:])
	if len(match) != 5 {
		return errors.New("parseAttackB matches fail:" + line)
	}
	var (
		player = strings.ReplaceAll(match[1], "致命一击！", "")
		target = match[3]
		skill  = match[2]
	)
	if !isPlayerValid(player) || !isTargetValid(target) {
		return nil
	}

	r.resultLog <- model.ChatLog{
		Player: player,
		Skill:  skill,
		Target: target,
		Value:  formatValue(match[4]),
		Time:   formatTime(line),
		RawMsg: line[22:],
	}
	r.resultPlayer <- model.Player{
		Name:  player,
		Class: r.skill2Class[skill],
		Time:  formatTime(line),
	}
	r.resultPlayer <- model.Player{
		Name: target,
		Time: formatTime(line),
	}
	return nil
}

// (.*?)给(.*?)造成了(.*)的伤害
func (r *Parser) parseAttackB(line string) error {
	if strings.Contains(line, "反弹了攻击") {
		return nil
	}
	if strings.Count(line[22:], "给") >= 2 {
		return nil
	}
	match := regAttackB.FindStringSubmatch(line[22:])
	if len(match) != 4 {
		return errors.New("ParseAttackA matches fail:" + line)
	}

	var (
		player = strings.ReplaceAll(match[1], "致命一击！", "")
		target = match[2]
	)

	if !isPlayerValid(player) || !isTargetValid(target) {
		return nil
	}

	r.resultLog <- model.ChatLog{
		Player: player,
		Skill:  "attack",
		Target: target,
		Value:  formatValue(match[3]),
		Time:   formatTime(line),
		RawMsg: line[22:],
	}
	r.resultPlayer <- model.Player{
		Name: player,
		Time: formatTime(line),
	}
	r.resultPlayer <- model.Player{
		Name: target,
		Time: formatTime(line),
	}
	return nil
}

// (.*?)使用(.*?)技能
func (r *Parser) parseAttackC(line string) error {
	match := regAttackC.FindStringSubmatch(line[22:])
	if len(match) != 3 {
		return errors.New("parseAttackC matches fail:" + line)
	}

	var (
		player = strings.ReplaceAll(match[1], "致命一击！", "")
		skill  = match[2]
	)
	if !isPlayerValid(player) {
		return nil
	}

	value := 0
	matchValue := regValue.FindStringSubmatch(line[22:])
	if len(matchValue) > 1 && matchValue[1] != "1" {
		value = formatValue(matchValue[1])
	}

	r.resultLog <- model.ChatLog{
		Player: player,
		Skill:  skill,
		Value:  value,
		Time:   formatTime(line),
		RawMsg: line[22:],
	}
	r.resultPlayer <- model.Player{
		Name:  player,
		Class: r.skill2Class[skill],
		Time:  formatTime(line),
	}
	return nil
}

// (.*?)把(.*?)打倒了
func (r *Parser) parseDeathA(line string) error {
	match := regDeathA.FindStringSubmatch(line[22:])
	if len(match) != 3 {
		return errors.New("parseDeathA matches fail:" + line)
	}

	var (
		player = match[1]
		target = match[2]
	)
	r.resultLog <- model.ChatLog{
		Player: player,
		Target: target,
		Skill:  "kill",
		Time:   formatTime(line),
		RawMsg: line[22:],
	}
	r.resultPlayer <- model.Player{
		Name: player,
		Type: model.TypeBright,
		Time: formatTime(line),
	}
	r.resultPlayer <- model.Player{
		Name: target,
		Type: model.TypeDark,
		Time: formatTime(line),
	}
	return nil
}

// (.*?)倒下了
func (r *Parser) parseDeathB(line string) error {
	match := regDeathB.FindStringSubmatch(line[22:])
	if len(match) != 2 {
		return errors.New("parseDeathB matches fail:" + line)
	}

	r.resultPlayer <- model.Player{
		Name: match[1],
		Type: model.TypeDark,
		Time: formatTime(line),
	}
	return nil
}

// (.*?)受到(.*?)的攻击而死亡
func (r *Parser) parseDeathC(line string) error {
	match := regDeathC.FindStringSubmatch(line[22:])
	if len(match) != 3 {
		return errors.New("parseDeathC matches fail:" + line)
	}

	var (
		player = match[1]
		target = match[2]
	)

	if !isPlayerValid(player) {
		return nil
	}

	r.resultLog <- model.ChatLog{
		Player: player,
		Target: target,
		Skill:  "killed",
		Time:   formatTime(line),
		RawMsg: line[22:],
	}
	r.resultPlayer <- model.Player{
		Name: player,
		Type: model.TypeBright,
		Time: formatTime(line),
	}
	r.resultPlayer <- model.Player{
		Name: target,
		Type: model.TypeDark,
		Time: formatTime(line),
	}
	return nil
}

func formatTime(ts string) time.Time {
	ts = strings.ReplaceAll(ts[:19], ".", "-")
	tm, _ := time.Parse(time.DateTime, ts)
	return tm
}

func formatValue(ds string) int {
	d, _ := strconv.Atoi(strings.ReplaceAll(ds, ",", ""))
	return d
}

var invalidPlayer = map[string]int{
	"太古气息": 1, "地之气息": 1, "水之气息": 1, "旋风之气息": 1, "风之气息": 1, "高洁气息": 1, "神圣的气息": 1, "治愈之气息": 1,
	"生命之气息": 1, "火之气息": 1, "深渊的气息": 1, "水之精灵": 1, "火之精灵": 1, "风之精灵": 1, "台风之精灵": 1, "地之精灵": 1,
	"熔岩精灵": 1, "冰柱": 1, "召唤台风": 1, "高级攻城兵器": 1, "超大型连射炮": 1, "大型连射炮": 1,
}

func isPlayerValid(name string) bool {
	if name == "" {
		return false
	}
	if _, ok := invalidPlayer[name]; ok {
		return false
	}
	return true
}

func isTargetValid(name string) bool {
	if name == "" {
		return false
	}
	if name == "训练用稻草人" {
		return false
	}
	return true
}
