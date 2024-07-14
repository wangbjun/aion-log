package service

import "aion/model"

type CacheService struct {
	cachePlayers  map[string][]*model.Player
	cacheRank     map[string][]model.RankResult
	cachePlayer   map[string]*model.Player
	cacheClass    map[string][]*model.SkillDamage
	cacheCritical map[string]model.SkillDamage
}

var defaultCacheService *CacheService

func NewCacheService() *CacheService {
	if defaultCacheService == nil {
		defaultCacheService = &CacheService{
			cachePlayers:  make(map[string][]*model.Player),
			cacheRank:     make(map[string][]model.RankResult),
			cachePlayer:   make(map[string]*model.Player),
			cacheClass:    make(map[string][]*model.SkillDamage),
			cacheCritical: make(map[string]model.SkillDamage),
		}
	}

	return defaultCacheService
}

func (s *CacheService) Load() error {
	players, err := model.Player{}.GetAll()
	if err != nil {
		return err
	}
	for _, player := range players {
		s.cachePlayer[player.Name] = player
	}

	criticalRatio, err := model.ChatLog{}.GetCriticalRatio("")
	if err != nil {
		return err
	}
	for _, critical := range criticalRatio {
		s.cacheCritical[critical.Skill] = critical
	}

	return nil
}

func (s *CacheService) GetPlayer(name string) (*model.Player, bool) {
	player, ok := s.cachePlayer[name]
	return player, ok
}

func (s *CacheService) GetPlayers(key string) ([]*model.Player, bool) {
	players, ok := s.cachePlayers[key]
	return players, ok
}

func (s *CacheService) SetPlayers(name string, data []*model.Player) {
	s.cachePlayers[name] = data
}

func (s *CacheService) GetRank(key string) ([]model.RankResult, bool) {
	cached, ok := s.cacheRank[key]
	return cached, ok
}
func (s *CacheService) SetRank(key string, data []model.RankResult) {
	s.cacheRank[key] = data
}

func (s *CacheService) GetClassTop(key string) ([]*model.SkillDamage, bool) {
	cached, ok := s.cacheClass[key]
	return cached, ok
}

func (s *CacheService) SetClassTop(key string, data []*model.SkillDamage) {
	s.cacheClass[key] = data
}

func (s *CacheService) GetSkillCritical(skill string) (float64, bool) {
	ratio, ok := s.cacheCritical[skill]
	return ratio.Critical, ok
}
