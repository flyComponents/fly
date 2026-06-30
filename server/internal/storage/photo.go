package storage

import (
	"context"
	"fly_server/internal/domain"
	"io"
)

type PhotoRepository interface {
	Save(ctx context.Context, photo domain.Photo) error
	Get(ctx context.Context, id string) (*domain.Photo, error)
	GetByOwnerId(ctx context.Context, ownerId string) ([]*domain.Photo, error)
}

type FileStorage interface {
	Save(ctx context.Context, id string, r io.Reader) (string, error)
	Delete(ctx context.Context, path string) error
}
