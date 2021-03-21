package service

import (
	"aion/model"
	"aion/zlog"
	"fmt"
	"time"
)

type ClassfiyService struct{}

func (r ClassfiyService) Start() {
	r.run()
	var t = time.NewTicker(time.Hour * 8)
	for {
		select {
		case <-t.C:
			r.run()
		}
	}
}

var skill2Pro = map[int][]string{
	model.JX: {"飞刀", "飞刀连射", "幻影弹", "猛烈一击", "会心一击", "会心一击", "破灭一击", "杀气破裂", "吸血波", "回旋一击"},
	model.SH: {"主神的惩罚", "剑风", "吸血惩罚", "处决一击", "天雷斩", "审判", "幻影摄捕", "捕获", "血剑斩", "连续乱打"},
	model.SX: {"交叉斩", "反击", "后方打击", "影子下坠", "暗破", "暗袭", "灭杀", "猛兽之牙", "背后重击", "进击斩"},
	model.GX: {"利锥箭", "夺取之箭", "套索箭", "强袭箭", "格力普尼克斯之箭", "沉默箭", "狂风箭", "百发百中", "破灭箭", "降魔之箭"},
	model.ZY: {"大地之怒", "大地的惩罚", "大地的报应", "天罚", "审判套索", "惩罚之电", "放电", "断罪一击", "权能爆发", "闪电"},
	model.HF: {"共鸣烟雾", "击破连锁", "强力击", "必灭重击", "暗击锁", "流星一击", "灭火", "白热一击", "破裂击", "贯穿连锁"},
	model.JL: {"元素打击", "吸引", "幽冥之苦痛", "愤怒之漩涡", "灵魂抢夺", "真空爆炸", "精灵弱化", "诅咒之云"},
	model.MD: {"冬季的束缚", "冰河重击", "冷气召唤", "暴风重击", "火焰乱舞", "火焰叉", "灵魂冻结", "神圣的咒语", "结冰", "魔力烈焰"},
}

func (r ClassfiyService) run() {
	defer func() {
		if err := recover(); err != nil {
			zlog.Logger.Error(fmt.Sprintf("RankService Run Error: %s", err))
		}
	}()
	for k, v := range skill2Pro {
		model.Player{}.UpdateBySkills(k, v)
	}
}
