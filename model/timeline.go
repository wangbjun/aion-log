package model

import "time"

type Timeline struct {
	Time  time.Time `json:"time"`
	Value int       `json:"value"`
	Type  int       `json:"type"`
}

func (r Timeline) TableName() string {
	return "aion_timeline"
}

func (r Timeline) BatchInsert(items []Timeline) error {
	if len(items) == 0 {
		return nil
	}
	return DB().CreateInBatches(items, BatchInsertSize).Error
}

func (r Timeline) GetAll(st, et string, tp int) ([]Timeline, error) {
	var results []Timeline
	query := DB().Where("type = ?", tp)
	if st != "" && et != "" {
		query = query.Where("time >= ? and time <= ?", st, et)
	}
	err := query.Find(&results).Error
	return results, err
}
