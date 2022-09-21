package memstorage

import (
	"context"

	"microcontrollers/internal/entity"
)

type Querier interface {
	GetHome(ctx context.Context, id string) (entity.Home, error)
	GetHomeTG(ctx context.Context, clientId string) (entity.Home, error)
	CreateHome(ctx context.Context, id, clientId string) (entity.Home, error)
	UpdateHome(ctx context.Context, id string, input entity.UpdateHomeInput) (entity.Home, error)
	UpdateHomeInfo(ctx context.Context, id string, input entity.UpdateHomeCommandInput) (entity.Home, error)
}

type Storage struct {
	q *Queries
}

func NewStorage() *Storage {
	return &Storage{
		q: NewQueries(),
	}
}

func (s *Storage) ExecTX(_ context.Context, fn func(Querier) error) error {
	s.q.mutex.Lock()
	defer s.q.mutex.Unlock()

	cp := make(map[string]*entity.Home)
	for id, home := range s.q.home {
		cpHome := *home
		cp[id] = &cpHome
	}

	err := fn(s.q)

	if err != nil {
		s.q.home = cp
		return err
	}

	return nil
}
