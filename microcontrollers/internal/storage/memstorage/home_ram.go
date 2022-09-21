package memstorage

import (
	"context"
	"sync"

	"microcontrollers/internal/entity"
	"microcontrollers/internal/pkg/database"

	"github.com/pkg/errors"
)

type Queries struct {
	mutex  *sync.Mutex
	home   map[string]*entity.Home
	client map[string]string
}

func NewQueries() *Queries {
	return &Queries{
		mutex:  new(sync.Mutex),
		home:   make(map[string]*entity.Home),
		client: make(map[string]string),
	}
}

func (q *Queries) GetHome(ctx context.Context, id string) (entity.Home, error) {
	h, ok := q.home[id]

	if !ok {
		return entity.Home{}, errors.Wrapf(database.ErrNotFound, "no such home with [id]: %s", id)
	}

	select {
	case <-ctx.Done():
		return entity.Home{}, errors.Wrap(database.ErrCtxDone, "GetHome: ctx done")
	default:

	}

	return *h, nil
}

func (q *Queries) GetHomeTG(ctx context.Context, id string) (entity.Home, error) {
	c, ok := q.client[id]

	if !ok {
		return entity.Home{}, errors.Wrapf(database.ErrNotFound, "no such home with [id]: %s", id)
	}

	select {
	case <-ctx.Done():
		return entity.Home{}, errors.Wrap(database.ErrCtxDone, "GetHomeTG: ctx done")
	default:

	}

	h, ok := q.home[c]

	if !ok {
		return entity.Home{}, errors.Wrapf(database.ErrNotFound, "no such home with [id]: %s", id)
	}

	return *h, nil
}

func (q *Queries) CreateHome(ctx context.Context, id, clientId string) (entity.Home, error) {
	_, ok := q.home[id]

	if ok {
		return entity.Home{}, errors.Wrapf(database.ErrDuplicateKey, "home with [id]: %s already exists", id)
	}

	select {
	case <-ctx.Done():
		return entity.Home{}, errors.Wrap(database.ErrCtxDone, "CreateHome: ctx done")
	default:

	}

	_, ok = q.client[clientId]

	if ok {
		return entity.Home{}, errors.Wrapf(database.ErrDuplicateKey, "client with [client_id]: %s already exists", clientId)
	}

	q.home[id] = &entity.Home{
		ID:       id,
		ClientId: clientId,
	}

	q.client[clientId] = id

	return *q.home[id], nil
}

func (q *Queries) UpdateHome(ctx context.Context, id string, input entity.UpdateHomeInput) (entity.Home, error) {
	_, ok := q.home[id]

	if !ok {
		return entity.Home{}, errors.Wrapf(database.ErrNotFound, "no such home with [id]: %s", id)
	}

	//if input.IsLedTurned != nil {
	//	q.home[id].IsLedTurned = *input.IsLedTurned
	//}

	select {
	case <-ctx.Done():
		return entity.Home{}, errors.Wrap(database.ErrCtxDone, "UpdateHome: ctx done")
	default:

	}

	if input.Temperature != nil {
		q.home[id].Temperature = *input.Temperature
	}

	if input.IsRobbery != nil {
		q.home[id].IsRobbery = *input.IsRobbery
	}

	return *q.home[id], nil
}

func (q *Queries) UpdateHomeInfo(ctx context.Context, id string, input entity.UpdateHomeCommandInput) (entity.Home, error) {
	homeId, ok := q.client[id]

	if !ok {
		return entity.Home{}, errors.Wrapf(database.ErrNotFound, "no such home with [id]: %s", id)
	}

	select {
	case <-ctx.Done():
		return entity.Home{}, errors.Wrap(database.ErrCtxDone, "UpdateHomeInfo: ctx done")
	default:

	}

	if input.SecureMode != nil {
		q.home[homeId].SecureMode = *input.SecureMode
	}
	if input.LedTurn != nil {
		q.home[homeId].IsLedTurned = *input.LedTurn
	}
	if input.OpenGate != nil {
		q.home[homeId].IsGateOpened = *input.OpenGate
	}

	if input.NewClientId != nil {
		q.home[homeId].ClientId = *input.NewClientId
		q.client[*input.NewClientId] = homeId
		delete(q.client, id)
	}

	return *q.home[homeId], nil
}
