package model

type Player struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Type  int    `json:"type"`
	Class int    `json:"class"`
	Time  string `json:"time"`
}

const (
	TypeOther = iota
	TypeTian
	TypeMo
)

type Class = int

const (
	Unknown Class = iota
	JX
	SH
	SX
	GX
	ZY
	HF
	JL
	MD
)

func (r Player) TableName() string {
	return "aion_player_info"
}

func (r Player) GetAll() ([]Player, error) {
	var results []Player
	err := DB().Order("name asc").Find(&results).Error
	return results, err
}
