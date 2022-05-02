package service

import (
	mc "microcontrollers"
	"microcontrollers/pkg/repository"
)

type HomeService struct {
	repo repository.Home
}

func NewHomeService(repo repository.Home) *HomeService {
	return &HomeService{repo: repo}
}

func (s *HomeService) GetHome(homeId string) (*mc.Home, bool) {
	return s.repo.GetHome(homeId)
}

func (s *HomeService) GetHomeTG(clientId string) (*mc.Home, bool) {
	return s.repo.GetHomeTG(clientId)
}

func (s *HomeService) CreateHome(homeId, clientId string) bool {
	return s.repo.CreateHome(homeId, clientId)
}

func (s *HomeService) UpdateHome(id string, input mc.UpdateHomeInput) bool {
	return s.repo.UpdateHome(id, input)
}

func (s *HomeService) UpdateHomeInfo(id string, input mc.UpdateHomeCommandInput) bool {
	return s.repo.UpdateHomeInfo(id, input)
}
