package model

type PlayerSkill struct {
	Skill string `json:"skill"`
	Class int    `json:"class"`
}

func (r PlayerSkill) TableName() string {
	return "aion_player_skill"
}

func (r PlayerSkill) GetAll() ([]PlayerSkill, error) {
	var results []PlayerSkill
	err := DB().Find(&results).Error
	return results, err
}
