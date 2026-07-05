package service

import (
	"aion/model"
	"sync"
)

type CacheService struct {
	mu          sync.RWMutex
	cachePlayer map[string]*model.Player
	cacheClass  map[string][]*model.SkillDamage
	cacheSkill  map[string]model.PlayerSkill
}

var (
	defaultCacheService *CacheService
	defaultCacheOnce    sync.Once
)

func NewCacheService() *CacheService {
	defaultCacheOnce.Do(func() {
		defaultCacheService = &CacheService{
			cachePlayer: make(map[string]*model.Player),
			cacheClass:  make(map[string][]*model.SkillDamage),
			cacheSkill:  make(map[string]model.PlayerSkill),
		}
	})

	return defaultCacheService
}

func (s *CacheService) Load() error {
	players, err := model.Player{}.GetAll()
	if err != nil {
		return err
	}
	cachePlayer := make(map[string]*model.Player, len(players))
	for _, player := range players {
		cachePlayer[player.Name] = player
	}

	skills, err := model.PlayerSkill{}.GetAll()
	if err != nil {
		return err
	}
	cacheSkill := make(map[string]model.PlayerSkill, len(skills))
	for _, skill := range skills {
		cacheSkill[skill.Skill] = skill
	}

	s.mu.Lock()
	s.cachePlayer = cachePlayer
	s.cacheClass = make(map[string][]*model.SkillDamage)
	s.cacheSkill = cacheSkill
	s.mu.Unlock()

	return nil
}

func (s *CacheService) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cachePlayer = make(map[string]*model.Player)
	s.cacheClass = make(map[string][]*model.SkillDamage)
	s.cacheSkill = make(map[string]model.PlayerSkill)
}

func (s *CacheService) GetPlayer(name string) (*model.Player, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	player, ok := s.cachePlayer[name]
	return player, ok
}

func (s *CacheService) GetSkill(skill string) (model.PlayerSkill, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	player, ok := s.cacheSkill[skill]
	return player, ok
}

func (s *CacheService) GetClassTop(key string) ([]*model.SkillDamage, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cached, ok := s.cacheClass[key]
	return cached, ok
}

func (s *CacheService) SetClassTop(key string, data []*model.SkillDamage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cacheClass[key] = data
}
