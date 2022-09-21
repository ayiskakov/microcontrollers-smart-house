package service

import (
	"context"

	"microcontrollers/internal/entity"
	"microcontrollers/internal/storage/memstorage"
)

func (s *Service) GetHome(ctx context.Context, homeId string) (entity.Home, error) {
	home := entity.Home{}

	err := s.storage.ExecTX(ctx, func(q memstorage.Querier) error {
		var err error
		home, err = q.GetHome(ctx, homeId)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.Home{}, err
	}

	return home, nil
}

func (s *Service) GetHomeTG(ctx context.Context, clientId string) (entity.Home, error) {
	home := entity.Home{}

	err := s.storage.ExecTX(ctx, func(q memstorage.Querier) error {
		var err error
		home, err = q.GetHomeTG(ctx, clientId)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.Home{}, err
	}

	return home, nil
}

func (s *Service) CreateHome(ctx context.Context, homeId, clientId string) (entity.Home, error) {
	home := entity.Home{}

	err := s.storage.ExecTX(ctx, func(q memstorage.Querier) error {
		var err error
		home, err = q.CreateHome(ctx, homeId, clientId)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.Home{}, err
	}

	return home, nil
}

func (s *Service) UpdateHome(ctx context.Context, id string, input entity.UpdateHomeInput) (entity.Home, error) {
	home := entity.Home{}

	err := s.storage.ExecTX(ctx, func(q memstorage.Querier) error {
		var err error
		home, err = q.UpdateHome(ctx, id, input)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.Home{}, err
	}

	return home, nil
}

func (s *Service) UpdateHomeInfo(ctx context.Context, id string, input entity.UpdateHomeCommandInput) (entity.Home, error) {
	home := entity.Home{}

	err := s.storage.ExecTX(ctx, func(q memstorage.Querier) error {
		var err error
		home, err = q.UpdateHomeInfo(ctx, id, input)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return entity.Home{}, err
	}

	return home, nil
}
