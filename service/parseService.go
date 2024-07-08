package service

import (
	"aion/model"
	"aion/util"
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
	regDeathA  = regexp.MustCompile("(.*?)把(.*?)打倒了")
	regDeathB  = regexp.MustCompile("(.*?)倒下了")
	regDeathC  = regexp.MustCompile("(.*?)受到(.*?)的攻击而终结")
	regAttackA = regexp.MustCompile("(.*?)使用(.*?)技能，对(.*?)造成了(.*)的伤害")
	regAttackB = regexp.MustCompile("(.*?)给(.*?)造成了(.*)的伤害")
	regAttackC = regexp.MustCompile("(.*?)使用(.*?)技能")
)

const WorkerNum = 5

type Parser struct {
	lineChan     chan string
	resultLog    chan model.Log
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

	var logItems []model.Log
	for {
		select {
		case battleLog, ok := <-r.resultLog:
			if !ok {
				doneLog = true
			} else {
				if strings.Contains(battleLog.Target, "训练用稻草人") {
					continue
				}
				logItems = append(logItems, battleLog)
				if len(logItems) >= 500 {
					model.Log{}.BatchInsert(logItems)
					logItems = []model.Log{}
				}
			}
		case player, ok := <-r.resultPlayer:
			if !ok {
				donePlayer = true
			} else {
				if player.Name == "" {
					continue
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
		}
		if doneLog && donePlayer {
			model.Log{}.BatchInsert(logItems)
			var result []model.Player
			for _, player := range r.uniquePlayer {
				result = append(result, player)
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
	player := strings.ReplaceAll(match[1], "致命一击！", "")
	if player == "" {
		player = "我"
	}
	skill := match[2]
	r.resultLog <- model.Log{
		Player: player,
		Skill:  skill,
		Target: match[3],
		Value:  formatDamage(match[4]),
		Time:   formatTime(line),
		RawMsg: line[22:],
	}
	r.resultPlayer <- model.Player{
		Name:  player,
		Class: r.skill2Class[util.RemoveRomanNumber(skill)],
		Time:  formatTime(line),
	}
	r.resultPlayer <- model.Player{
		Name: match[3],
		Time: formatTime(line),
	}
	return nil
}

// (.*?)给(.*?)造成了(.*)的伤害
func (r *Parser) parseAttackB(line string) error {
	if strings.Contains(line, "反弹了攻击") {
		return nil
	}
	match := regAttackB.FindStringSubmatch(line[22:])
	if len(match) != 4 {
		return errors.New("ParseAttackA matches fail:" + line)
	}
	player := strings.ReplaceAll(match[1], "致命一击！", "")
	if player == "" {
		player = "我"
	}
	r.resultLog <- model.Log{
		Player: player,
		Skill:  "普通攻击",
		Target: match[2],
		Value:  formatDamage(match[3]),
		Time:   formatTime(line),
		RawMsg: line[22:],
	}
	r.resultPlayer <- model.Player{
		Name: player,
		Time: formatTime(line),
	}
	r.resultPlayer <- model.Player{
		Name: match[2],
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
	player := strings.ReplaceAll(match[1], "致命一击！", "")
	if player == "" {
		player = "我"
	}
	skill := match[2]
	r.resultLog <- model.Log{
		Player: player,
		Skill:  skill,
		Time:   formatTime(line),
		RawMsg: line[22:],
	}
	r.resultPlayer <- model.Player{
		Name:  player,
		Class: r.skill2Class[util.RemoveRomanNumber(skill)],
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
	r.resultLog <- model.Log{
		Player: match[1],
		Target: match[2],
		Time:   formatTime(line),
		RawMsg: line[22:],
	}
	r.resultPlayer <- model.Player{
		Name: match[1],
		Type: model.TypeBright,
		Time: formatTime(line),
	}
	r.resultPlayer <- model.Player{
		Name: match[2],
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
	r.resultLog <- model.Log{
		Player: match[2],
		Target: match[1],
		Time:   formatTime(line),
		RawMsg: line[22:],
	}
	r.resultPlayer <- model.Player{
		Name: match[1],
		Type: model.TypeBright,
		Time: formatTime(line),
	}
	r.resultPlayer <- model.Player{
		Name: match[2],
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

func formatDamage(ds string) int {
	d, _ := strconv.Atoi(strings.ReplaceAll(ds, ",", ""))
	return d
}
