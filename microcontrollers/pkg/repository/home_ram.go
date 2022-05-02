package repository

import (
	mc "microcontrollers"
	"sync"
)

type HomeRamRepository struct {
	mutex  *sync.Mutex
	home   map[string]*mc.Home
	client map[string]string
}

func NewHomeRamRepository() *HomeRamRepository {
	return &HomeRamRepository{
		mutex:  new(sync.Mutex),
		home:   make(map[string]*mc.Home),
		client: make(map[string]string),
	}
}

func (r *HomeRamRepository) GetHome(id string) (*mc.Home, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	h, ok := r.home[id]

	return h, ok
}

func (r *HomeRamRepository) GetHomeTG(id string) (*mc.Home, bool) {
	//r.mutex.Lock()
	//defer r.mutex.Unlock()
	c, ok := r.client[id]

	if !ok {
		return nil, ok
	}

	h, ok := r.home[c]

	return h, ok
}
func (r *HomeRamRepository) CreateHome(id, clientId string) bool {
	//r.mutex.Lock()
	//defer r.mutex.Unlock()

	_, exh := r.home[id]
	_, exc := r.client[clientId]

	if exh || exc {
		return false
	}

	r.home[id] = &mc.Home{
		ID:       id,
		ClientId: clientId,
	}

	r.client[clientId] = id

	return true
}

func (r *HomeRamRepository) UpdateHome(id string, input mc.UpdateHomeInput) bool {
	//r.mutex.Lock()
	//defer r.mutex.Unlock()

	_, ex := r.home[id]

	if !ex {
		return false
	}

	//if input.IsLedTurned != nil {
	//	r.home[id].IsLedTurned = *input.IsLedTurned
	//}

	if input.Temperature != nil {
		r.home[id].Temperature = *input.Temperature
	}

	if input.IsRobbery != nil {
		r.home[id].IsRobbery = *input.IsRobbery
	}

	return true
}

func (r *HomeRamRepository) UpdateHomeInfo(id string, input mc.UpdateHomeCommandInput) bool {
	//r.mutex.Lock()
	//defer r.mutex.Unlock()

	homeId, ex := r.client[id]

	if !ex {
		return false
	}

	if input.SecureMode != nil {
		r.home[homeId].SecureMode = *input.SecureMode
	}
	if input.LedTurn != nil {
		r.home[homeId].IsLedTurned = *input.LedTurn
	}
	if input.OpenGate != nil {
		r.home[homeId].IsGateOpened = *input.OpenGate
	}

	if input.NewClientId != nil {
		r.home[homeId].ClientId = *input.NewClientId
		r.client[*input.NewClientId] = homeId
		delete(r.client, id)
	}

	return true
}
