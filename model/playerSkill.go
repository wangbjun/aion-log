package model

type PlayerSkill struct {
	Skill         string  `json:"skill"`
	Class         int     `json:"class"`
	CriticalRatio float64 `json:"critical_ratio"`
}

func (r PlayerSkill) TableName() string {
	return "aion_player_skill"
}

func (r PlayerSkill) GetAll() ([]PlayerSkill, error) {
	var results []PlayerSkill
	err := DB().Find(&results).Error
	return results, err
}

func (r PlayerSkill) GetBySkill(skill string) (*PlayerSkill, error) {
	var result PlayerSkill
	err := DB().Raw("select * from aion_player_skill where skill = ?", skill).Find(&result).Error
	return &result, err
}
