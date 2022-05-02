package service

import (
	mc "microcontrollers"
	"microcontrollers/pkg/repository"
)

type Home interface {
	GetHome(id string) (*mc.Home, bool)
	GetHomeTG(id string) (*mc.Home, bool)
	CreateHome(id, clientId string) bool
	UpdateHome(id string, input mc.UpdateHomeInput) bool
	UpdateHomeInfo(id string, input mc.UpdateHomeCommandInput) bool
}

type Service struct {
	Home
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Home: NewHomeService(repos.Home),
	}
}
