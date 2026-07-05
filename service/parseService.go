package service

import (
	"aion/model"
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

var (
	kill     = regexp.MustCompile("(\\S+)把(\\S+)打倒了。")
	death    = regexp.MustCompile("(\\S+)倒下了。")
	killedBy = regexp.MustCompile("(\\S+)受到(\\S+)的攻击而终结。")

	attackSkill  = regexp.MustCompile("(\\S+)使用(.+)技能，对(\\S+)造成了(\\S+)的伤害")
	attackNormal = regexp.MustCompile("(\\S+)给(\\S+)造成了(\\S+)的伤害")
	useSkill     = regexp.MustCompile("(\\S+)使用(.+)技能，")
)

type parsedLine struct {
	Time time.Time
	Msg  string
	Raw  string
}

type logRule struct {
	name   string
	expr   *regexp.Regexp
	handle func(*Parser, parsedLine, []string) error
}

var logRules = []logRule{
	{name: "attack_skill", expr: attackSkill, handle: (*Parser).attackSkill},
	{name: "attack_normal", expr: attackNormal, handle: (*Parser).attackNormal},
	{name: "use_skill", expr: useSkill, handle: (*Parser).useSkill},

	{name: "kill", expr: kill, handle: (*Parser).kill},
	{name: "death", expr: death, handle: (*Parser).death},
	{name: "killed_by", expr: killedBy, handle: (*Parser).killedBy},
}

type Parser struct {
	resultLog    chan model.ChatLog
	resultPlayer chan model.Player
	uniquePlayer map[string]model.Player
	skill2Class  map[string]model.PlayerSkill
}

func NewParseService() Parser {
	return Parser{
		resultLog:    make(chan model.ChatLog, 1000),
		resultPlayer: make(chan model.Player, 1000),
		uniquePlayer: make(map[string]model.Player, 1000),
		skill2Class:  defaultCacheService.cacheSkill,
	}
}

func (r *Parser) Run(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	return r.RunReader(file, fileName)
}

func (r *Parser) RunReader(reader io.Reader, sourceName string) error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	errCh := make(chan error, 2)
	go r.processPlayer(&wg, errCh)
	go r.processLog(&wg, errCh)

	st := time.Now()
	log.Printf("begin process: %s", sourceName)

	decoder := simplifiedchinese.GBK.NewDecoder()
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		b, err := decoder.Bytes(scanner.Bytes())
		if err != nil {
			log.Printf("decoder error: %s", err)
			continue
		}
		if len(b) == 0 {
			continue
		}
		line, ok := splitLogLine(string(b))
		if !ok {
			log.Printf("skip invalid line: %s", string(b))
			continue
		}

		err = r.parseLine(line)
		if err != nil {
			log.Printf("parse line error: %s", err)
		}
	}
	if err := scanner.Err(); err != nil {
		close(r.resultLog)
		close(r.resultPlayer)
		wg.Wait()
		close(errCh)
		return fmt.Errorf("scan file failed: %w", err)
	}
	close(r.resultLog)
	close(r.resultPlayer)
	wg.Wait()
	close(errCh)

	log.Printf("finish process: %s, cost: %.2fs\n", sourceName, time.Since(st).Seconds())
	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Parser) processLog(wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	var firstErr error
	logItems := make([]model.ChatLog, 0, model.BatchInsertSize)
	for {
		chatLog, ok := <-r.resultLog
		if !ok {
			break
		}
		logItems = append(logItems, chatLog)
		if len(logItems) >= model.BatchInsertSize {
			err := model.ChatLog{}.BatchInsert(logItems)
			if err != nil {
				log.Printf("batch insert log error: %s", err)
				if firstErr == nil {
					firstErr = fmt.Errorf("batch insert log failed: %w", err)
				}
			}
			logItems = make([]model.ChatLog, 0, model.BatchInsertSize)
		}
	}
	if len(logItems) > 0 {
		err := model.ChatLog{}.BatchInsert(logItems)
		if err != nil {
			log.Printf("batch insert log error: %s", err)
			if firstErr == nil {
				firstErr = fmt.Errorf("batch insert log failed: %w", err)
			}
		}
	}
	errCh <- firstErr
}

func (r *Parser) processPlayer(wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	var firstErr error

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
		if len(result) >= model.BatchInsertSize {
			err := model.Player{}.BatchInsert(result)
			if err != nil {
				log.Printf("batch insert player error: %s", err)
				if firstErr == nil {
					firstErr = fmt.Errorf("batch insert player failed: %w", err)
				}
			}
			result = make([]model.Player, 0, model.BatchInsertSize)
		}
	}
	if len(result) > 0 {
		err := model.Player{}.BatchInsert(result)
		if err != nil {
			log.Printf("batch insert player error: %s", err)
			if firstErr == nil {
				firstErr = fmt.Errorf("batch insert player failed: %w", err)
			}
		}
	}
	errCh <- firstErr
}

func (r *Parser) parseLine(line parsedLine) error {
	for _, rule := range logRules {
		match := rule.expr.FindStringSubmatch(line.Msg)
		if len(match) == 0 {
			continue
		}
		if err := rule.handle(r, line, match); err != nil {
			return fmt.Errorf("%s: %w", rule.name, err)
		}
		return nil
	}
	return nil
}

// (.*?)使用(.*?)技能，对(.*?)造成了(.*)的伤害
func (r *Parser) attackSkill(line parsedLine, match []string) error {
	msg, tm := line.Msg, line.Time
	if len(match) != 5 {
		return errors.New("parseAttackA matches fail:" + msg)
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
		Time:   tm,
		RawMsg: msg,
	}
	r.resultPlayer <- model.Player{
		Name:  player,
		Class: r.skill2Class[skill].Class,
		Time:  tm,
	}
	r.resultPlayer <- model.Player{
		Name: target,
		Time: tm,
	}
	return nil
}

// (.*?)给(.*?)造成了(.*)的伤害
func (r *Parser) attackNormal(line parsedLine, match []string) error {
	msg, tm := line.Msg, line.Time
	if strings.Contains(msg, "反弹了攻击") {
		return nil
	}
	if strings.Count(msg, "给") >= 2 {
		return nil
	}
	if len(match) != 4 {
		return errors.New("parseAttackB matches fail:" + msg)
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
		Time:   tm,
		RawMsg: msg,
	}
	r.resultPlayer <- model.Player{
		Name: player,
		Time: tm,
	}
	r.resultPlayer <- model.Player{
		Name: target,
		Time: tm,
	}
	return nil
}

// (.*?)使用(.*?)技能
func (r *Parser) useSkill(line parsedLine, match []string) error {
	msg, tm := line.Msg, line.Time
	if len(match) != 3 {
		return errors.New("parseAttackC matches fail:" + msg)
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
		Class: r.skill2Class[skill].Class,
		Time:  tm,
	}
	return nil
}

// (.*?)把(.*?)打倒了
func (r *Parser) kill(line parsedLine, match []string) error {
	msg, tm := line.Msg, line.Time
	if len(match) != 3 {
		return errors.New("parseDeathA matches fail:" + msg)
	}

	var (
		player = match[1]
		target = match[2]
	)
	r.resultLog <- model.ChatLog{
		Player: player,
		Target: target,
		Skill:  "kill",
		Time:   tm,
		RawMsg: msg,
	}
	r.resultPlayer <- model.Player{
		Name: player,
		Type: model.TypeBright,
		Time: tm,
	}
	r.resultPlayer <- model.Player{
		Name: target,
		Type: model.TypeDark,
		Time: tm,
	}
	return nil
}

// (.*?)倒下了
func (r *Parser) death(line parsedLine, match []string) error {
	msg, tm := line.Msg, line.Time
	if len(match) != 2 {
		return errors.New("parseDeathB matches fail:" + msg)
	}

	r.resultPlayer <- model.Player{
		Name: match[1],
		Type: model.TypeDark,
		Time: tm,
	}
	return nil
}

// (.*?)受到(.*?)的攻击而死亡
func (r *Parser) killedBy(line parsedLine, match []string) error {
	msg, tm := line.Msg, line.Time
	if len(match) != 3 {
		return errors.New("parseDeathC matches fail:" + msg)
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
		Time:   tm,
		RawMsg: msg,
	}
	r.resultPlayer <- model.Player{
		Name: player,
		Type: model.TypeBright,
		Time: tm,
	}
	r.resultPlayer <- model.Player{
		Name: target,
		Type: model.TypeDark,
		Time: tm,
	}
	return nil
}

func splitLogLine(raw string) (parsedLine, bool) {
	if len(raw) <= 22 {
		return parsedLine{}, false
	}
	tm := formatTime(raw)
	if tm.IsZero() {
		return parsedLine{}, false
	}
	return parsedLine{
		Time: tm,
		Msg:  raw[22:],
		Raw:  raw,
	}, true
}

func formatTime(ts string) time.Time {
	if len(ts) < 19 {
		return time.Time{}
	}
	ts = strings.ReplaceAll(ts[:19], ".", "-")
	tm, _ := time.ParseInLocation(time.DateTime, ts, time.Local)
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
