package service

import "aion/model"

type CacheService struct {
	cachePlayers map[string][]*model.Player
	cacheRank    map[string][]model.RankResult
	cachePlayer  map[string]*model.Player
}

var defaultCacheService *CacheService

func NewCacheService() *CacheService {
	if defaultCacheService == nil {
		defaultCacheService = &CacheService{
			cachePlayers: make(map[string][]*model.Player),
			cacheRank:    make(map[string][]model.RankResult),
			cachePlayer:  make(map[string]*model.Player),
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
