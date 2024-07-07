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

var classSkills = map[int][]string{
	model.JX: {"飞刀", "飞刀连射", "猛烈一击", "身体重击", "愤怒一击", "杀气破裂", "吸血波", "地震波动", "冲击波动", "激怒爆炸",
		"剑气波动", "剑气破裂", "杀气波动", "强制捆绑", "下刺", "最后一击", "腾空斩击", "生存姿态", "集中格挡", "翅膀强化", "突击姿态", "愤怒波动"},
	model.SH: {"主神的惩罚", "脚踝重击", "盾牌反击", "处决一击", "天雷斩", "主神惩罚", "审判", "集中捕获", "捕获", "精神破坏", "激昂",
		"闪光斩", "保护之盾", "庇护之盔甲", "阻断之甲", "破坏之气合", "主神盔甲", "双重盔甲", "俘虏"},
	model.SX: {"暗杀者步伐", "猛兽的咆哮", "绝魂斩", "反击", "影子下坠", "暗破", "暗袭", "灭杀", "猛兽之牙", "背后重击", "进击斩",
		"短剑投掷", "神速契约", "命中之契约", "回避契约", "烟幕弹", "六感最大化", "强袭姿态", "影子步行", "进击斩"},
	model.GX: {"利锥箭", "狙击", "疾风箭", "暴走箭", "连射", "套索箭", "强袭箭", "沉默箭", "狂风箭", "百发百中", "破灭箭", "祝福之弓",
		"攻击之眼", "风之疾走", "猎人的决心"},
	model.ZY: {"大地之怒", "大地的惩罚", "惩戒之电", "放电", "断罪一击", "惩戒", "闪电", "审判之电", "破灭之诉说", "灿烂之佑护",
		"净化之光辉", "治疗保护膜", "不死之帐幕", "集中祈祷", "神速之祈祷", "免罪", "痊愈之闪光", "治疗之风", "苦行", "尤斯迪埃之光"},
	model.HF: {"共鸣烟雾", "打击锁链", "必灭重击", "白热一击", "暗击锁", "流星一击", "灭火", "贯穿连锁", "铁壁之咒语",
		"神速之咒语", "激怒之咒语", "鼓吹之咒语", "保护阵", "守护之祝福", "阻断之幕", "疾走之咒语"},
	model.JL: {"吸引", "幽冥之苦痛", "愤怒之漩涡", "真空爆炸", "精灵弱化", "诅咒之云", "大地之锁链", "范围侵蚀", "魔法解除", "侵蚀",
		"精神协调", "召唤:风之精灵"},
	model.MD: {"火焰熔解", "混乱的咒语", "冬季的束缚", "冰河重击", "流星重击", "火焰爆发", "元气吸收", "台风重击", "冷气召唤", "暴风重击",
		"火焰乱舞", "火焰叉", "结冰", "时空扭曲", "白杰尔的智慧", "神速的恩惠", "元素结界", "冰雪甲冑", "铁甲之恩惠", "魔力倍增"},
	model.ZXZ: {"攻", "断", "击", "雷", "电击突袭", "爆雷打", "雷电大爆炸", "斩", "电流波动", "呼啸雷电", "电雷重击", "命运之执行者"},
}

const WorkerNum = 5

type Parser struct {
	lineChan     chan string
	resultLog    chan model.Log
	resultPlayer chan model.Player
	uniquePlayer map[string]model.Player
	skill2Class  map[string]model.Class
}

func NewParseService() Parser {
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
	r.resultLog <- model.Log{
		Player: player,
		Skill:  match[2],
		Target: match[3],
		Value:  formatDamage(match[4]),
		Time:   formatTime(line),
		RawMsg: line[22:],
	}
	r.resultPlayer <- model.Player{
		Name:  player,
		Class: r.skill2Class[util.RemoveRomanNumber(match[2])],
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
	r.resultPlayer <- model.Player{
		Name:  player,
		Class: r.skill2Class[util.RemoveRomanNumber(match[2])],
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
