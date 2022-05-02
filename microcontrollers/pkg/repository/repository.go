package repository

import mc "microcontrollers"

type Home interface {
	GetHome(id string) (*mc.Home, bool)
	GetHomeTG(clientId string) (*mc.Home, bool)
	CreateHome(id, clientId string) bool
	UpdateHome(id string, input mc.UpdateHomeInput) bool
	UpdateHomeInfo(id string, input mc.UpdateHomeCommandInput) bool
}

type Repository struct {
	Home
}

func NewRepository() *Repository {
	return &Repository{
		Home: NewHomeRamRepository(),
	}
}
