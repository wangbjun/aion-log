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
)

type Parser struct {
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

	wg := sync.WaitGroup{}
	wg.Add(2)
	go r.processLog(&wg)
	go r.processPlayer(&wg)

	st := time.Now()
	log.Printf("begin proccess: %s", fileName)

	decoder := simplifiedchinese.GBK.NewDecoder()
	buff := bufio.NewReader(file)
	for {
		a, _, err := buff.ReadLine()
		if err == io.EOF {
			break
		}
		b, err := decoder.Bytes(a)
		if err != nil {
			log.Printf("decoder error: %s", err)
			continue
		}
		if len(b) == 0 {
			continue
		}
		line := string(b)
		if regAttackA.MatchString(line) {
			err = r.parseAttackA(line)
		} else if regAttackB.MatchString(line) {
			err = r.parseAttackB(line)
		} else if regAttackC.MatchString(line) {
			err = r.parseAttackC(line)
		} else if regDeathA.MatchString(line) {
			err = r.parseDeathA(line)
		} else if regDeathB.MatchString(line) {
			err = r.parseDeathB(line)
		} else if regDeathC.MatchString(line) {
			err = r.parseDeathC(line)
		}
		if err != nil {
			fmt.Printf("parse line error: %s\n", err)
		}
	}
	close(r.resultLog)
	close(r.resultPlayer)
	wg.Wait()

	log.Printf("finish proccess: %s, cost: %.2fs\n", fileName, time.Since(st).Seconds())
	return nil
}

func (r *Parser) processLog(wg *sync.WaitGroup) {
	defer wg.Done()
	logItems := make([]model.ChatLog, 0, 500)
	for {
		chatLog, ok := <-r.resultLog
		if !ok {
			break
		}
		logItems = append(logItems, chatLog)
		if len(logItems) >= 500 {
			err := model.ChatLog{}.BatchInsert(logItems)
			if err != nil {
				log.Printf("batch insert log error: %s", err)
			}
			logItems = make([]model.ChatLog, 0, 500)
		}
	}
	model.ChatLog{}.BatchInsert(logItems)
}

func (r *Parser) processPlayer(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		player, ok := <-r.resultPlayer
		if !ok {
			break
		}
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
	var result []model.Player
	for _, player := range r.uniquePlayer {
		result = append(result, player)
		if len(result) >= 500 {
			err := model.Player{}.BatchInsert(result)
			if err != nil {
				log.Printf("batch insert player error: %s", err)
			}
			result = make([]model.Player, 0, 500)
		}
	}
	model.Player{}.BatchInsert(result)
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
	if _, ok := invalidTarget[name]; ok {
		return false
	}
	return true
}

var invalidPlayer = map[string]int{
	"太古气息": 1, "地之气息": 1, "水之气息": 1, "旋风之气息": 1, "风之气息": 1, "高洁气息": 1, "神圣的气息": 1, "治愈之气息": 1,
	"生命之气息": 1, "火之气息": 1, "深渊的气息": 1, "水之精灵": 1, "火之精灵": 1, "风之精灵": 1, "台风之精灵": 1, "地之精灵": 1,
	"熔岩精灵": 1, "冰柱": 1, "召唤台风": 1, "高级攻城兵器": 1, "超大型连射炮": 1, "大型连射炮": 1, "飞行祝福": 1,
}

var invalidTarget = map[string]int{
	"训练用稻草人": 1, "变异的RA-98c": 1,
}
