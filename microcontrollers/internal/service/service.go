package service

import (
	"context"

	"microcontrollers/internal/storage/memstorage"
)

type Storage interface {
	ExecTX(ctx context.Context, fn func(q memstorage.Querier) error) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}
